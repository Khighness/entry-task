package main

import (
	"context"
	"encoding/gob"
	"entry/rpc"
	"entry/rpc/demo/public"
	"log"
	"time"
)

// @Author KHighness
// @Update 2022-02-20

var QueryUser func(int64) (public.ResponseQueryUser, error)

func main() {
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
	go connPool.Release(permission1.RpcCli, ctx)
	go connPool.Release(permission2.RpcCli, ctx)
	go connPool.Release(permission3.RpcCli, ctx)
	for {

	}
}
