package main

import (
	"context"
	"encoding/gob"
	"log"
	"sync"
	"time"

	"github.com/Khighness/entry-task/rpc"
	"github.com/Khighness/entry-task/rpc/demo/public"
)

// @Author KHighness
// @Email  zikang.chen@shopee.com
// @Since 2022-02-20

var QueryUser func(int64) (public.ResponseQueryUser, error)

func main() {
	gob.Register(public.ResponseQueryUser{})
	gob.Register(public.ResponseQueryUser{})

	ctx := context.Background()
	config := &rpc.Config{
		MaxOpenCount:  3,
		MaxIdleCount:  2,
		RpcServerAddr: "127.0.0.1:30000",
	}
	connPool := rpc.Init(ctx, config)
	permission1, _ := connPool.Achieve(ctx)
	permission2, _ := connPool.Achieve(ctx)
	permission3, _ := connPool.Achieve(ctx)
	go connPool.Achieve(ctx)

	permission1.RpcCli.Call("queryUser", &QueryUser)
	u, err := QueryUser(1)
	if err != nil {
		log.Printf("query error: %v\n", err)
	} else {
		log.Printf("query result: %v %v\n", u.Id, u.Name)
	}

	time.Sleep(2 * time.Second)
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		connPool.Release(permission1.RpcCli, ctx)
		wg.Done()
	}()
	go func() {
		connPool.Release(permission2.RpcCli, ctx)
		wg.Done()
	}()
	go func() {
		connPool.Release(permission3.RpcCli, ctx)
		wg.Done()
	}()

	wg.Wait()
}
