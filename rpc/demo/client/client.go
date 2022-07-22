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

	ctx := context.Background()
	config := &rpc.Config{
		MaxOpenCount:  3,
		MaxIdleCount:  2,
		RpcServerAddr: "127.0.0.1:30000",
	}
	connPool := rpc.NewPool(config)
	permission1, _ := connPool.Achieve(ctx)
	log.Printf("%+v", connPool.Stat())
	permission2, _ := connPool.Achieve(ctx)
	log.Printf("%+v", connPool.Stat())

	queryUser := func(client *rpc.Client) {
		client.Call("queryUser", &QueryUser)
	}
	connPool.Exec(ctx, queryUser)
	u, err := QueryUser(1)
	if err != nil {
		log.Printf("query error: %v\n", err)
	} else {
		log.Printf("query result: %v %v\n", u.Id, u.Name)
	}

	time.Sleep(2 * time.Second)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		connPool.Release(permission1.RpcCli, ctx)
		log.Printf("%+v", connPool.Stat())
		wg.Done()
	}()
	go func() {
		connPool.Release(permission2.RpcCli, ctx)
		log.Printf("%+v", connPool.Stat())
		wg.Done()
	}()

	wg.Wait()

	time.Sleep(2 * time.Second)
	log.Printf("%+v", connPool.Stat())
}
