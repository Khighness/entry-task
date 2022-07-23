package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Khighness/entry-task/pkg/rpc"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-03-08

// RpcPool RPC连接池
type RpcPool struct {
	connector RpcConnector // 连接器

	mu           sync.RWMutex                // 读写锁
	closed       bool                        // 连接池的状态
	numOpen      int                         // 打开的连接数
	freeConn     []*RpcConn                  // 空闲连接池
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

// RpcPoolConfig RPC连接池配置
type RpcPoolConfig struct {
	InitialCount int
	MaxOpenCount int
	MaxIdleCount int
	MaxLifeTime  time.Duration
	MaxIdleTime  time.Duration
}

// RpcPoolStats grpc连接池状态
type RpcPoolStats struct {
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

// RpcConn 封装RPC客户端
type RpcConn struct {
	createdAt  time.Time   // 当前连接创建的时间
	client     *rpc.Client // grpc客户端
	inUse      bool        // 标识当前连接的状态，是否正在使用
	returnedAt time.Time   // 当前连接创建或者归还的时间
}

// connRequest 等待获取连接的请求
type connRequest struct {
	conn *RpcConn
}

const (
	defaultMaxIdleCount        = 2
	connectionRequestQueueSize = 100000
)

var nowFunc func() time.Time = time.Now
var errBadConn error = errors.New("bad connection")
var errPoolClosed error = errors.New("grpc is closed")

// nextRequestKeyLocked 返回下一个获取连接的请求的标识
func (rp *RpcPool) nextRequestKeyLocked() uint64 {
	next := rp.nextRequest
	rp.nextRequest++
	return next
}

// archive 获取连接
// [1] 从空闲连接池获取空闲连接
// [2] 从管道中获取被释放的连接
// [3] 创建新的连接
func (rp *RpcPool) archive(ctx context.Context) (*RpcConn, error) {
	rp.mu.Lock()
	if rp.closed {
		rp.mu.Unlock()
		return nil, errPoolClosed
	}

	select {
	default:
	case <-ctx.Done():
		rp.mu.Unlock()
		return nil, ctx.Err()
	}
	lifeTime := rp.maxLifeTime

	numFree := len(rp.freeConn)
	if numFree > 0 {
		rc := rp.freeConn[0]
		copy(rp.freeConn, rp.freeConn[1:])
		rp.freeConn = rp.freeConn[:numFree-1]
		rc.inUse = true

		if rc.expired(lifeTime) {
			rp.maxLifeTimeClosed++
			rp.numOpen--
			rp.mu.Unlock()
			return nil, errBadConn
		}
		rp.mu.Unlock()

		return rc, nil
	}

	if rp.maxOpenCount > 0 && rp.numOpen >= rp.maxOpenCount {
		req := make(chan connRequest, 1)
		reqKey := rp.nextRequestKeyLocked()
		rp.connRequests[reqKey] = req
		rp.waitCount++
		rp.mu.Unlock()

		select {
		case <-ctx.Done():
			rp.mu.Lock()
			delete(rp.connRequests, reqKey)
			rp.mu.Unlock()

			select {
			default:
			case ret, ok := <-req:
				if ok && ret.conn != nil {
					rp.release(ret.conn)
				}
			}
			return nil, ctx.Err()
		case ret, ok := <-req:
			if !ok {
				return nil, errPoolClosed
			}

			if ret.conn.expired(lifeTime) {
				rp.mu.Lock()
				rp.maxLifeTimeClosed++
				rp.numOpen--
				rp.mu.Unlock()
				rp.connector.Close(ret.conn.client)
				return nil, errBadConn
			}

			return ret.conn, nil
		}
	}

	rp.numOpen++ // optimistically
	rp.mu.Unlock()
	ci, err := rp.connector.Connect(ctx)
	if err != nil {
		rp.mu.Lock()
		rp.numOpen--
		rp.mayOpenNewConnections()
		rp.mu.Unlock()
		return nil, err
	}
	rp.mu.Lock()
	rc := &RpcConn{
		createdAt:  nowFunc(),
		client:     ci,
		inUse:      true,
		returnedAt: nowFunc(),
	}
	rp.mu.Unlock()
	return rc, nil
}

// release 释放连接
func (rp *RpcPool) release(rc *RpcConn) {
	rp.mu.Lock()
	if !rc.inUse {
		rp.mu.Unlock()
		panic("connection returned that was never out")
	}
	var err error
	if rc.expired(rp.maxLifeTime) {
		rp.maxLifeTimeClosed++
		rp.numOpen--
		err = errBadConn
	}
	rc.inUse = false
	rc.returnedAt = nowFunc()

	if err == errBadConn {
		rp.mayOpenNewConnections()
		rp.mu.Unlock()
		rp.connector.Close(rc.client)
		return
	}

	added := rp.putConnLocked(rc)
	rp.mu.Unlock()
	if !added {
		rp.connector.Close(rc.client)
		return
	}
}

// connectionOpener 监听 RpcPool.openerCh ，有消息则创建连接
func (rp *RpcPool) connectionOpener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-rp.openerCh:
			rp.openNewConnection(ctx)
		}
	}
}

// openNewConnection 创建新的连接
func (rp *RpcPool) openNewConnection(ctx context.Context) {
	ci, err := rp.connector.Connect(ctx)
	rp.mu.Lock()
	defer rp.mu.Unlock()
	if rp.closed {
		if err == nil {
			rp.connector.Close(ci)
		}
		rp.numOpen--
		return
	}
	if err != nil {
		rp.numOpen--
		rp.mayOpenNewConnections()
		return
	}
	rc := &RpcConn{
		createdAt:  nowFunc(),
		client:     ci,
		inUse:      false,
		returnedAt: nowFunc(),
	}
	if !rp.putConnLocked(rc) {
		rp.numOpen--
	}
}

// putConnLocked 处理连接
// [1] 将连接传给另一个请求，返回true
// [2] 将连接放回空闲连接池，返回true
// [3] 将连接丢弃，返回false
func (rp *RpcPool) putConnLocked(rc *RpcConn) bool {
	if rp.closed {
		return false
	}
	if rp.maxOpenCount > 0 && rp.numOpen > rp.maxOpenCount {
		return false
	}
	if c := len(rp.connRequests); c > 0 {
		var req chan connRequest
		var reqKey uint64
		for reqKey, req = range rp.connRequests {
			break
		}
		delete(rp.connRequests, reqKey)
		rc.inUse = true
		req <- connRequest{conn: rc}
		return true
	} else if !rp.closed {
		if rp.maxIdleCountLocked() > len(rp.freeConn) {
			rp.freeConn = append(rp.freeConn, rc)
			rp.startCleanerLocked()
			return true
		}
		rp.connector.Close(rc.client)
		rp.maxIdleClosed++
		rp.numOpen--
	}
	return false
}

// mayOpenNewConnection 创建连接或者释放连接失败的情况下调用，保证不会阻塞获取连接的请求
func (rp *RpcPool) mayOpenNewConnections() {
	numRequests := len(rp.connRequests)
	if rp.maxOpenCount > 0 {
		numCanOpen := rp.maxOpenCount - rp.numOpen
		if numRequests > numCanOpen {
			numRequests = numCanOpen
		}
	}
	for numRequests > 0 {
		rp.numOpen++ // optimistically
		numRequests--
		rp.openerCh <- struct{}{}
	}
}

// startCleanerLocked 开启 RpcPool.cleanerCh ，清理超时连接
func (rp *RpcPool) startCleanerLocked() {
	if (rp.maxLifeTime > 0 || rp.maxIdleTime > 0) && rp.numOpen > 0 && rp.cleanerCh == nil {
		rp.cleanerCh = make(chan struct{}, 1)
		go rp.connectionCleaner(rp.shortestIdleTimeLocked())
	}
}

// connectionCleaner 监听 RpcPool.cleanerCh ，有消息则检查连接池
func (rp *RpcPool) connectionCleaner(d time.Duration) {
	const minInterval = time.Second
	if d < minInterval {
		d = minInterval
	}
	t := time.NewTimer(d)

	for {
		select {
		case <-t.C:
		}

		rp.mu.Lock()
		d = rp.shortestIdleTimeLocked()
		if rp.closed || rp.numOpen == 0 || d <= 0 {
			rp.cleanerCh = nil
			rp.mu.Unlock()
			return
		}

		closing := rp.connectionCleanerRunLocked()
		rp.mu.Unlock()
		for _, rpcConn := range closing {
			rp.connector.Close(rpcConn.client)
		}

		if d < minInterval {
			d = minInterval
		}
		t.Reset(d)
	}
}

// connectionCleanerRunLocked 执行检查逻辑，返回存活超时的连接切片
func (rp *RpcPool) connectionCleanerRunLocked() (closing []*RpcConn) {
	if rp.maxLifeTime > 0 {
		expiredSince := nowFunc().Add(-rp.maxLifeTime)
		for i := 0; i < len(rp.freeConn); i++ {
			rc := rp.freeConn[i]
			if rc.createdAt.Before(expiredSince) {
				closing = append(closing, rc)
				last := len(rp.freeConn) - 1
				rp.freeConn[i] = rp.freeConn[last]
				rp.freeConn[last] = nil
				rp.freeConn = rp.freeConn[:last]
				i--
			}
		}
		rp.maxLifeTimeClosed += int64(len(closing))
		rp.numOpen -= len(closing)
	}

	if rp.maxIdleTime > 0 {
		expiredSince := nowFunc().Add(-rp.maxIdleTime)
		var expiredCount int64
		for i := 0; i < len(rp.freeConn); i++ {
			rc := rp.freeConn[i]
			if rp.maxIdleTime > 0 && rc.returnedAt.Before(expiredSince) {
				closing = append(closing, rc)
				last := len(rp.freeConn) - 1
				rp.freeConn[i] = rp.freeConn[last]
				rp.freeConn[last] = nil
				rp.freeConn = rp.freeConn[:last]
				expiredCount++
				i--
			}
		}
		rp.maxIdleTimeClosed += expiredCount
		rp.numOpen -= int(expiredCount)
	}
	return
}

// maxIdleCountLocked 获取最大空闲连接数量
func (rp *RpcPool) maxIdleCountLocked() int {
	n := rp.maxIdleCount
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
func (rp *RpcPool) shortestIdleTimeLocked() time.Duration {
	if rp.maxIdleTime <= 0 {
		return rp.maxIdleTime
	}
	if rp.maxLifeTime <= 0 {
		return rp.maxIdleTime
	}

	min := rp.maxIdleTime
	if min > rp.maxLifeTime {
		min = rp.maxLifeTime
	}
	return min
}

// expired 检查连接是否过期
func (rc *RpcConn) expired(timeout time.Duration) bool {
	if timeout < 0 {
		return false
	}
	return rc.createdAt.Add(timeout).Before(nowFunc())
}

// NewRpcPool 构造一个grpc连接池
func NewRpcPool(connector RpcConnector, config *RpcPoolConfig) *RpcPool {
	ctx, cancel := context.WithCancel(context.Background())
	rp := &RpcPool{
		connector:    connector,
		mu:           sync.RWMutex{},
		connRequests: make(map[uint64]chan connRequest),
		openerCh:     make(chan struct{}, connectionRequestQueueSize),
		maxOpenCount: config.MaxOpenCount,
		maxIdleCount: config.MaxIdleCount,
		maxLifeTime:  config.MaxLifeTime,
		maxIdleTime:  config.MaxIdleTime,
		stop:         cancel,
	}

	if config.InitialCount > config.MaxIdleCount {
		config.InitialCount = config.MaxIdleCount
	}
	for i := 0; i < config.InitialCount; i++ {
		conn, err := rp.connector.Connect(ctx)
		if err != nil {
			panic(err)
		}
		rc := &RpcConn{
			createdAt:  nowFunc(),
			client:     conn,
			inUse:      false,
			returnedAt: nowFunc(),
		}
		rp.freeConn = append(rp.freeConn, rc)
	}
	rp.numOpen = config.InitialCount

	go rp.connectionOpener(ctx)
	rp.startCleanerLocked()
	return rp
}

// Stat 获取连接池的状态
func (rp *RpcPool) Stat() RpcPoolStats {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	stats := RpcPoolStats{
		MaxOpenCount:      rp.maxOpenCount,
		MaxIdleCount:      rp.maxIdleCount,
		OpenNum:           rp.numOpen,
		Idle:              len(rp.freeConn),
		InUse:             rp.numOpen - len(rp.freeConn),
		WaitCount:         rp.waitCount,
		MaxIdleClosed:     rp.maxIdleClosed,
		MaxIdleTimeClosed: rp.maxIdleTimeClosed,
		MaxLifeTimeClosed: rp.maxLifeTimeClosed,
	}
	return stats
}

// Close 关闭连接池
func (rp *RpcPool) Close() {
	rp.mu.Lock()
	if rp.closed {
		rp.mu.Unlock()
	}
	if rp.cleanerCh != nil {
		close(rp.cleanerCh)
	}
	for _, rpcConn := range rp.freeConn {
		rp.connector.Close(rpcConn.client)
	}
	rp.freeConn = nil
	rp.closed = true
	rp.mu.Unlock()
	rp.stop()
}

// Exec 执行函数，返回结果
// 🤣 基于golang闭包的大一统接口
func (rp *RpcPool) Exec(handle func(client *rpc.Client)) error {
	rc, err := rp.archive(context.Background())
	if err != nil {
		return err
	}
	handle(rc.client)
	go rp.release(rc)
	return nil
}
