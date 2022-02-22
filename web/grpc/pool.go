package grpc

import (
	"context"
	"entry/pb"
	"entry/web/common"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-22

type ConnPool struct {
	lock sync.Mutex // 锁

	connPool      []int                   // 连接池
	openCount     int                     // 当前连接数
	waitCount     int                     // 等待个数
	waitQueue     map[int]chan Permission // 等待队列
	availableConn map[int]Permission      // 连接池
	nextConnIndex NextConnIndex           // 下一个连接的ID标识

	initCount     int           // 初始化连接数
	macOpenCount  int           // 最大连接数
	maxIdleCount  int           // 最大空闲连接数
	maxLifeTime   time.Duration // 连接最大存活时间
	maxWaitTime   time.Duration // 任务最大等待时间
	rpcServerAddr string        // 远程服务器地址
}

type Permission struct {
	NextConnIndex
	RpcCli      pb.UserServiceClient
	CreateAt    time.Time
	MaxLifeTime time.Duration
}

type NextConnIndex struct {
	Index int
}

type Config struct {
	InitCount     int
	MaxOpenCount  int
	MaxIdleCount  int
	MaxLifeTime   time.Duration
	MaxWaitTime   time.Duration
	RpcServerAddr string
}

var nowFunc = time.Now

//  load 初始化
func NewPool(ctx context.Context, config *Config) (conn *ConnPool) {
	if config.InitCount > 10000 || config.InitCount > config.MaxOpenCount {
		return nil
	}

	log.Println("rpc server:", config.RpcServerAddr)

	pool := &ConnPool{
		connPool:      []int{},
		openCount:     config.InitCount,
		waitCount:     0,
		waitQueue:     make(map[int]chan Permission),
		availableConn: make(map[int]Permission),
		macOpenCount:  config.MaxOpenCount,
		maxIdleCount:  config.MaxOpenCount,
		maxLifeTime:   config.MaxLifeTime,
		maxWaitTime:   config.MaxWaitTime,
		rpcServerAddr: config.RpcServerAddr,
	}

	for i := 0; i < config.InitCount; i++ {
		nextConnIndex := getNextConnIndex(pool)
		conn, err := grpc.Dial(pool.rpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			e := fmt.Sprintf("Failed to connect to rpc server %s, err: %s\n", pool.rpcServerAddr, err)
			log.Printf(e)
			return nil
		}
		client := pb.NewUserServiceClient(conn)
		permission := Permission{
			NextConnIndex: NextConnIndex{nextConnIndex},
			RpcCli:        client,
			CreateAt:      nowFunc(),
			MaxLifeTime:   pool.maxLifeTime,
		}
		pool.availableConn[nextConnIndex] = permission
	}

	return pool
}

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
		log.Printf("Achieve connection[fromPool] successfully, openCount:%d, idleCount:%v\n", pool.openCount, len(pool.availableConn))
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
		case <-time.After(pool.maxWaitTime):
			log.Println("Achieve connection failed, cause: wait timeout")
			return
		case ret, ok := <-req:
			if !ok {
				return Permission{}, errors.New("get connection failed, cause: no available connection release")
			}
			log.Printf("Achieve connection[released] successfully, openCount:%d, idleCount:%v\n", pool.openCount, len(pool.availableConn))
			return ret, nil
		}
	}

	// (3) 当前连接数低于上限，创建新连接
	pool.openCount++
	nextConnIndex := getNextConnIndex(pool)
	pool.lock.Unlock()

	conn, err := grpc.Dial(common.RpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		e := fmt.Sprintf("Failed to connect to server %s, err: %s\n", pool.rpcServerAddr, err)
		log.Printf(e)
		return Permission{}, errors.New(e)
	}
	client := pb.NewUserServiceClient(conn)
	permission = Permission{
		NextConnIndex: NextConnIndex{nextConnIndex},
		RpcCli:        client,
		CreateAt:      nowFunc(),
		MaxLifeTime:   pool.maxLifeTime,
	}
	log.Printf("Achieve connection[created], openCount:%d, idleCount:%v\n", pool.openCount, len(pool.availableConn))
	return permission, nil
}

func getNextConnIndex(conn *ConnPool) int {
	currentIndex := conn.nextConnIndex.Index
	conn.nextConnIndex.Index = currentIndex + 1
	return conn.nextConnIndex.Index
}

// Release 释放连接
func (pool *ConnPool) Release(client pb.UserServiceClient, ctx context.Context) (result bool, err error) {
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
		log.Printf("Release connection to wait task, openCount:%d, idleCount:%v\n", pool.openCount, len(pool.availableConn))
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
					MaxLifeTime:   pool.maxLifeTime,
				}
				pool.availableConn[nextConnIndex] = permission
				log.Printf("Release connection to conn pool, openCount:%d, idleCount:%v\n", pool.openCount, len(pool.availableConn))
			}
		}
	}

	pool.lock.Unlock()
	return
}
