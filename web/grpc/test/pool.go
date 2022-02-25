package main

import (
	"context"
	"entry/web/common"
	"entry/web/grpc"
	"sync"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-25

func main() {
	Pool := grpc.NewPool(&grpc.Config{
		InitCount:     50,
		MaxOpenCount:  100,
		MaxIdleCount:  100,
		MaxLifeTime:   time.Hour,
		MaxWaitTime:   time.Second,
		RpcServerAddr: common.RpcServerAddr,
	})

	var wg = sync.WaitGroup{}
	wg.Add(210)

	for i := 0; i < 200; i++ {
		//wg.Add(1)
		go func() {
			permission, _ := Pool.Achieve(context.Background())
			defer Pool.Release(permission, context.Background())
			wg.Done()
		}()
	}

	wg.Wait()

}
