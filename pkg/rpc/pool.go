package rpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Khighness/entry-task/pkg/logger"
	"github.com/sirupsen/logrus"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-21

// ConnPool 连接池
type ConnPool struct {
	logger *logrus.Logger // 日志
	lock   sync.Mutex     // 锁

	connPool      []int                   // 连接池
	openCount     int                     // 当前连接数
	waitCount     int                     // 等待个数
	waitQueue     map[int]chan Permission // 等待队列
	availableConn map[int]Permission      // 连接池
	nextConnIndex NextConnIndex           // 下一个连接的ID标识

	macOpenCount  int    // 最大连接数
	maxIdleCount  int    // 最大空闲连接数
	rpcServerAddr string // 远程服务器地址
}

// Permission 权限，包装连接
type Permission struct {
	NextConnIndex
	RpcCli      *Client
	CreateAt    time.Time
	MaxLifeTime time.Duration
}

// NextConnIndex 下一个连接的标识
type NextConnIndex struct {
	Index int
}

// Stat 连接池状态
type Stat struct {
	OpenCount int
	IdleCount int
	WaitCount int
}

// Config 连接池配置
type Config struct {
	MaxOpenCount  int
	MaxIdleCount  int
	RpcServerAddr string
}

// NewPool 创建连接吃
func NewPool(config *Config) *ConnPool {
	return &ConnPool{
		logger:        logger.NewLogger(logrus.DebugLevel, "", true),
		connPool:      []int{},
		openCount:     0,
		waitCount:     0,
		waitQueue:     make(map[int]chan Permission),
		availableConn: make(map[int]Permission),
		macOpenCount:  config.MaxOpenCount,
		maxIdleCount:  config.MaxOpenCount,
		rpcServerAddr: config.RpcServerAddr,
	}
}

var nowFunc = time.Now

// Achieve 获取连接
func (pool *ConnPool) Achieve(ctx context.Context) (permission Permission, err error) {
	pool.lock.Lock()

	select {
	default:
	case <-ctx.Done():
		// context取消或者超时，退出
		pool.lock.Unlock()
		return Permission{}, errors.New("fail to create a new connection, context canceled")
	}

	// (1) 连接池不为空，直接获取连接
	if len(pool.availableConn) > 0 {
		var (
			popPermission Permission
			popReqKey     int
		)
		for popReqKey, popPermission = range pool.availableConn {
			break
		}

		delete(pool.availableConn, popReqKey)
		pool.logger.Debugf("Achieve connection[fromPool] successfully, openCount:%d, idleCount:%v", pool.openCount, len(pool.availableConn))
		pool.lock.Unlock()

		return popPermission, nil
	}

	// (2) 当前连接数大于上限，则加入等待队列
	if pool.openCount >= pool.macOpenCount {
		nextConnIndex := getNextConnIndex(pool)

		req := make(chan Permission, 1)
		pool.waitQueue[nextConnIndex] = req
		pool.waitCount++
		pool.lock.Unlock()

		select {
		case <-time.After(time.Second * time.Duration(5)):
			pool.logger.Debugf("Achieve connection failed, cause: wait timeout")
			return
		case ret, ok := <-req:
			if !ok {
				return Permission{}, errors.New("get connection failed, cause: no available connection release")
			}
			pool.logger.Debugf("Achieve connection[released] successfully, openCount:%d, idleCount:%v", pool.openCount, len(pool.availableConn))
			return ret, nil
		}
	}

	// (3) 当前连接数低于上限，创建新连接
	pool.openCount++
	pool.lock.Unlock()
	nextConnIndex := getNextConnIndex(pool)

	c, err := net.Dial("tcp", pool.rpcServerAddr)
	if err != nil {
		e := fmt.Sprintf("Failed to connect to server %s, err: %s", pool.rpcServerAddr, err)
		pool.logger.Errorf(e)
		return Permission{}, errors.New(e)
	}
	client := NewClient(c)
	permission = Permission{
		NextConnIndex: NextConnIndex{nextConnIndex},
		RpcCli:        client,
		CreateAt:      nowFunc(),
		MaxLifeTime:   0,
	}
	pool.logger.Debugf("Achieve connection[created], openCount:%d, idleCount:%v", pool.openCount, len(pool.availableConn))
	return permission, nil
}

func getNextConnIndex(conn *ConnPool) int {
	currentIndex := conn.nextConnIndex.Index
	conn.nextConnIndex.Index = currentIndex + 1
	return conn.nextConnIndex.Index
}

// Release 释放连接
func (pool *ConnPool) Release(client *Client, ctx context.Context) (result bool, err error) {
	pool.lock.Lock()

	// (1) 有任务在等待获取连接
	// 将释放的连接通过channel发送给该阻塞任务
	// 然后从等待队列中删除该任务
	if len(pool.waitQueue) > 0 {
		var req chan Permission
		var reqKey int
		for reqKey, req = range pool.waitQueue {
			break
		}

		permission := Permission{
			NextConnIndex: NextConnIndex{reqKey},
			RpcCli:        client,
			CreateAt:      nowFunc(),
			MaxLifeTime:   time.Second * 5,
		}
		req <- permission
		delete(pool.waitQueue, reqKey)
		pool.waitCount--
		pool.logger.Debugf("Release connection to wait task, openCount:%d, idleCount:%v", pool.openCount, len(pool.availableConn))
	} else {
		// (2) 没有等待任务，将连接放入连接池
		if pool.openCount > 0 {
			pool.openCount--
			if len(pool.availableConn) < pool.maxIdleCount {
				nextConnIndex := getNextConnIndex(pool)
				permission := Permission{
					NextConnIndex: NextConnIndex{nextConnIndex},
					RpcCli:        client,
					CreateAt:      nowFunc(),
					MaxLifeTime:   time.Second * 5,
				}
				pool.availableConn[nextConnIndex] = permission
				pool.logger.Debugf("Release connection to conn pool, openCount:%d, idleCount:%v", pool.openCount, len(pool.availableConn))
			}
		}
	}

	pool.lock.Unlock()
	return
}

// Stat 获取连接池状态
func (pool *ConnPool) Stat() Stat {
	return Stat{
		OpenCount: pool.openCount,
		IdleCount: len(pool.availableConn),
		WaitCount: pool.waitCount,
	}
}

// Exec 执行函数，返回结果
func (pool *ConnPool) Exec(ctx context.Context, handle func(client *Client)) error {
	permission, err := pool.Achieve(ctx)
	if err != nil {
		return err
	}
	defer pool.Release(permission.RpcCli, ctx)
	handle(permission.RpcCli)
	return nil
}
