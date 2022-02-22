package main

import (
	"context"
	"entry/pb"
	"entry/web/grpc"
	"log"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-22

func main() {
	config := &grpc.Config{
		InitCount:     10,
		MaxOpenCount:  100,
		MaxIdleCount:  30,
		MaxLifeTime:   time.Hour,
		MaxWaitTime:   time.Second,
		RpcServerAddr: "127.0.0.1:20000",
	}

	ctx := context.Background()
	pool := grpc.NewPool(ctx, config)

	permission, _ := pool.Achieve(ctx)
	cli := permission.RpcCli
	timeout1, cancelFunc1 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc1()
	response, err := cli.Login(timeout1, &pb.LoginRequest{
		Username: "KHighness",
		Password: "czk911",
	})

	if err != nil {
		log.Println(err)
	} else {
		log.Printf("%v\n", response)
	}

	timeout2, cancelFunc2 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFunc2()
	pool.Release(cli, timeout2)
}
