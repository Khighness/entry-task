package grpc

import (
	"context"
	"github.com/Khighness/entry-task/pb"
	"errors"
	"sync"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-03-08

// GrpcPool grpc连接池
type GrpcPool struct {
	connector GrpcConnector // 连接器

	mu           sync.Mutex                  // 互斥锁
	closed       bool                        // 连接池的状态
	numOpen      int                         // 打开的连接数
	freeConn     []*GrpcConn                 // 空闲连接池
	connRequests map[uint64]chan connRequest // 等待队列
	nextRequest  uint64                      // 下一个获取连接请求的唯一标识

	openerCh     chan struct{} // 通道，监听打开连接的消息
	maxOpenCount int           // 最大打开连接数量
	maxIdleCount int           // 最大空闲连接数量
	maxLifeTime  time.Duration // 连接最长存活时间
	maxIdleTime  time.Duration // 连接最长空闲时间

	cleanerCh         chan struct{} // 通道，监听关闭连接的消息
	waitCount         int64         // 等待获取连接的请求数量
	maxIdleClosed     int64         // 由于连接池已满而关闭的连接总数
	maxIdleTimeClosed int64         // 由于超过最长空闲时间而关闭的连接总数
	maxLifeTimeClosed int64         // 由于超过最长存活时间而关闭的连接总数

	stop func() // 连接池关闭，关闭 openerCh
}

// GrpcPoolConfig grpc连接池配置
type GrpcPoolConfig struct {
	MaxOpenCount int
	MaxIdleCount int
	MaxLifeTime  time.Duration
	MaxIdleTime  time.Duration
}

// GrpcPoolStats grpc连接池状态
type GrpcPoolStats struct {
	MaxOpenCount int
	MaxIdleCount int

	OpenNum int
	Idle    int
	InUse   int

	WaitCount         int64
	MaxIdleClosed     int64
	MaxIdleTimeClosed int64
	MaxLifeTimeClosed int64
}

// GrpcConn 封装grpc客户端
type GrpcConn struct {
	createdAt time.Time // 当前连接创建的时间

	// TODO: close(grpcCli)
	grpcCli    pb.UserServiceClient // grpc客户端
	inUse      bool                 // 标识当前连接的状态，是否正在使用
	returnedAt time.Time            // 当前连接创建或者归还的时间
}

// connRequest 等待获取连接的请求
type connRequest struct {
	conn *GrpcConn
}

const (
	defaultMaxIdleCount        = 2
	connectionRequestQueueSize = 100000
)

var nowFunc func() time.Time = time.Now
var errBadConn error = errors.New("bad connection")
var errPoolClosed error = errors.New("grpc is closed")

// nextRequestKeyLocked 返回下一个获取连接的请求的标识
func (gp *GrpcPool) nextRequestKeyLocked() uint64 {
	next := gp.nextRequest
	gp.nextRequest++
	return next
}

// archive 获取连接
// [1] 从空闲连接池获取空闲连接
// [2] 从管道中获取被释放的连接
// [3] 创建新的连接
func (gp *GrpcPool) archive(ctx context.Context) (*GrpcConn, error) {
	gp.mu.Lock()
	if gp.closed {
		gp.mu.Unlock()
		return nil, errPoolClosed
	}

	select {
	default:
	case <-ctx.Done():
		gp.mu.Unlock()
		return nil, ctx.Err()
	}
	lifeTime := gp.maxLifeTime

	numFree := len(gp.freeConn)
	if numFree > 0 {
		gc := gp.freeConn[0]
		copy(gp.freeConn, gp.freeConn[1:])
		gp.freeConn = gp.freeConn[:numFree-1]
		gc.inUse = true

		if gc.expired(lifeTime) {
			gp.maxLifeTimeClosed++
			gp.numOpen--
			gp.mu.Unlock()
			return nil, errBadConn
		}
		gp.mu.Unlock()

		return gc, nil
	}

	if gp.maxOpenCount > 0 && gp.numOpen >= gp.maxOpenCount {
		req := make(chan connRequest, 1)
		reqKey := gp.nextRequestKeyLocked()
		gp.connRequests[reqKey] = req
		gp.waitCount++
		gp.mu.Unlock()

		select {
		case <-ctx.Done():
			gp.mu.Lock()
			delete(gp.connRequests, reqKey)
			gp.mu.Unlock()

			select {
			default:
			case ret, ok := <-req:
				if ok && ret.conn != nil {
					gp.release(ret.conn)
				}
			}
			return nil, ctx.Err()
		case ret, ok := <-req:
			if !ok {
				return nil, errPoolClosed
			}

			if ret.conn.expired(lifeTime) {
				gp.mu.Lock()
				gp.maxLifeTimeClosed++
				gp.numOpen--
				gp.mu.Unlock()
				// close(gc)
				return nil, errBadConn
			}

			return ret.conn, nil
		}
	}

	gp.numOpen++ // optimistically
	gp.mu.Unlock()
	ci, err := gp.connector.Connect(ctx)
	if err != nil {
		gp.mu.Lock()
		gp.numOpen--
		gp.mayOpenNewConnections()
		gp.mu.Unlock()
		return nil, err
	}
	gp.mu.Lock()
	gc := &GrpcConn{
		createdAt:  nowFunc(),
		grpcCli:    ci,
		inUse:      true,
		returnedAt: nowFunc(),
	}
	gp.mu.Unlock()
	return gc, nil
}

// release 释放连接
func (gp *GrpcPool) release(gc *GrpcConn) {
	gp.mu.Lock()
	if !gc.inUse {
		gp.mu.Unlock()
		panic("connection returned that was never out")
	}
	var err error
	if gc.expired(gp.maxLifeTime) {
		gp.maxLifeTimeClosed++
		gp.numOpen--
		err = errBadConn
	}
	gc.inUse = false
	gc.returnedAt = nowFunc()

	if err == errBadConn {
		gp.mayOpenNewConnections()
		gp.mu.Unlock()
		// close(gc)
		return
	}

	added := gp.putConnLocked(gc)
	gp.mu.Unlock()
	if !added {
		// close(gc)
		return
	}
}

// connectionOpener 监听 GrpcPool.openerCh ，有消息则创建连接
func (gp *GrpcPool) connectionOpener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-gp.openerCh:
			gp.openNewConnection(ctx)
		}
	}
}

// openNewConnection 创建新的连接
func (gp *GrpcPool) openNewConnection(ctx context.Context) {
	ci, err := gp.connector.Connect(ctx)
	gp.mu.Lock()
	defer gp.mu.Unlock()
	if gp.closed {
		if err == nil {
			// close(ci)
		}
		gp.numOpen--
		return
	}
	if err != nil {
		gp.numOpen--
		gp.mayOpenNewConnections()
		return
	}
	gc := &GrpcConn{
		createdAt:  nowFunc(),
		grpcCli:    ci,
		inUse:      false,
		returnedAt: nowFunc(),
	}
	if !gp.putConnLocked(gc) {
		gp.numOpen--
	}
}

// putConnLocked 处理连接
// [1] 将连接传给另一个请求，返回true
// [2] 将连接放回空闲连接池，返回true
// [3] 将连接丢弃，返回false
func (gp *GrpcPool) putConnLocked(gc *GrpcConn) bool {
	if gp.closed {
		return false
	}
	if gp.maxOpenCount > 0 && gp.numOpen > gp.maxOpenCount {
		return false
	}
	if c := len(gp.connRequests); c > 0 {
		var req chan connRequest
		var reqKey uint64
		for reqKey, req = range gp.connRequests {
			break
		}
		delete(gp.connRequests, reqKey)
		gc.inUse = true
		req <- connRequest{conn: gc}
		return true
	} else if !gp.closed {
		if gp.maxIdleCountLocked() > len(gp.freeConn) {
			gp.freeConn = append(gp.freeConn, gc)
			gp.startCleanerLocked()
			return true
		}
		// close(gc)
		gp.maxIdleClosed++
		gp.numOpen--
	}
	return false
}

// mayOpenNewConnection 创建连接或者释放连接失败的情况下调用，保证不会阻塞获取连接的请求
func (gp *GrpcPool) mayOpenNewConnections() {
	numRequests := len(gp.connRequests)
	if gp.maxOpenCount > 0 {
		numCanOpen := gp.maxOpenCount - gp.numOpen
		if numRequests > numCanOpen {
			numRequests = numCanOpen
		}
	}
	for numRequests > 0 {
		gp.numOpen++ // optimistically
		numRequests--
		gp.openerCh <- struct{}{}
	}
}

// startCleanerLocked 开启 GrpcPool.cleanerCh ，清理超时连接
func (gp *GrpcPool) startCleanerLocked() {
	if (gp.maxLifeTime > 0 || gp.maxIdleTime > 0) && gp.numOpen > 0 && gp.cleanerCh == nil {
		gp.cleanerCh = make(chan struct{}, 1)
		go gp.connectionCleaner(gp.shortestIdleTimeLocked())
	}
}

// connectionCleaner 监听 GrpcPool.cleanerCh ，有消息则检查连接池
func (gp *GrpcPool) connectionCleaner(d time.Duration) {
	const minInterval = time.Second
	if d < minInterval {
		d = minInterval
	}
	t := time.NewTimer(d)

	for {
		select {
		case <-t.C:
		}

		gp.mu.Lock()
		d = gp.shortestIdleTimeLocked()
		if gp.closed || gp.numOpen == 0 || d <= 0 {
			gp.cleanerCh = nil
			gp.mu.Unlock()
			return
		}

		closing := gp.connectionCleanerRunLocked()
		gp.mu.Unlock()
		for _, _ = range closing {
			// close(gc)
		}

		if d < minInterval {
			d = minInterval
		}
		t.Reset(d)
	}
}

// connectionCleanerRunLocked 执行检查逻辑，返回存活超时的连接切片
func (gp *GrpcPool) connectionCleanerRunLocked() (closing []*GrpcConn) {
	if gp.maxLifeTime > 0 {
		expiredSince := nowFunc().Add(-gp.maxLifeTime)
		for i := 0; i < len(gp.freeConn); i++ {
			gc := gp.freeConn[i]
			if gc.createdAt.Before(expiredSince) {
				closing = append(closing, gc)
				last := len(gp.freeConn) - 1
				gp.freeConn[i] = gp.freeConn[last]
				gp.freeConn[last] = nil
				gp.freeConn = gp.freeConn[:last]
				i--
			}
		}
		gp.maxLifeTimeClosed += int64(len(closing))
		gp.numOpen -= len(closing)
	}

	if gp.maxIdleTime > 0 {
		expiredSince := nowFunc().Add(-gp.maxIdleTime)
		var expiredCount int64
		for i := 0; i < len(gp.freeConn); i++ {
			gc := gp.freeConn[i]
			if gp.maxIdleTime > 0 && gc.returnedAt.Before(expiredSince) {
				closing = append(closing, gc)
				last := len(gp.freeConn) - 1
				gp.freeConn[i] = gp.freeConn[last]
				gp.freeConn[last] = nil
				gp.freeConn = gp.freeConn[:last]
				expiredCount++
				i--
			}
		}
		gp.maxIdleTimeClosed += expiredCount
		gp.numOpen -= int(expiredCount)
	}
	return
}

// maxIdleCountLocked 获取最大空闲连接数量
func (gp *GrpcPool) maxIdleCountLocked() int {
	n := gp.maxIdleCount
	switch {
	case n == 0:
		return defaultMaxIdleCount
	case n < 0:
		return 0
	default:
		return n
	}
}

// shortestIdleTimeLocked 获取连接的最大可空闲时间
func (gp *GrpcPool) shortestIdleTimeLocked() time.Duration {
	if gp.maxIdleTime <= 0 {
		return gp.maxIdleTime
	}
	if gp.maxLifeTime <= 0 {
		return gp.maxIdleTime
	}

	min := gp.maxIdleTime
	if min > gp.maxLifeTime {
		min = gp.maxLifeTime
	}
	return min
}

// expired 检查连接是否过期
func (gc *GrpcConn) expired(timeout time.Duration) bool {
	if timeout < 0 {
		return false
	}
	return gc.createdAt.Add(timeout).Before(nowFunc())
}

// NewGrpcPool 构造一个grpc连接池
func NewGrpcPool(connector GrpcConnector, config *GrpcPoolConfig) *GrpcPool {
	ctx, cancel := context.WithCancel(context.Background())
	gp := &GrpcPool{
		connector:    connector,
		mu:           sync.Mutex{},
		connRequests: make(map[uint64]chan connRequest),
		openerCh:     make(chan struct{}, connectionRequestQueueSize),
		maxOpenCount: config.MaxOpenCount,
		maxIdleCount: config.MaxIdleCount,
		maxLifeTime:  config.MaxLifeTime,
		maxIdleTime:  config.MaxIdleTime,
		stop:         cancel,
	}

	go gp.connectionOpener(ctx)
	return gp
}

// Exec 执行函数，返回结果
// 🤣 基于golang闭包的大一统接口
func (gp *GrpcPool) Exec(f func(client pb.UserServiceClient) (interface{}, error)) (interface{}, error) {
	gc, err := gp.archive(context.Background())
	if err != nil {
		return nil, err
	}
	defer gp.release(gc)
	return f(gc.grpcCli)
}

// Stat 获取连接池的状态
func (gp *GrpcPool) Stat() GrpcPoolStats {
	gp.mu.Lock()
	defer gp.mu.Unlock()
	stats := GrpcPoolStats{
		MaxOpenCount:      gp.maxOpenCount,
		MaxIdleCount:      gp.maxIdleCount,
		OpenNum:           gp.numOpen,
		Idle:              len(gp.freeConn),
		InUse:             gp.numOpen - len(gp.freeConn),
		WaitCount:         gp.waitCount,
		MaxIdleClosed:     gp.maxIdleClosed,
		MaxIdleTimeClosed: gp.maxIdleTimeClosed,
		MaxLifeTimeClosed: gp.maxLifeTimeClosed,
	}
	return stats
}

// Close 关闭连接池
func (gp *GrpcPool) Close() {
	gp.mu.Lock()
	if gp.closed {
		gp.mu.Unlock()
	}
	if gp.cleanerCh != nil {
		close(gp.cleanerCh)
	}
	// close(freeConn)
	gp.freeConn = nil
	gp.closed = true
	gp.mu.Unlock()
	gp.stop()
}
