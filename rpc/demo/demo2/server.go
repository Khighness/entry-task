package main

import (
	"entry/rpc/demo/demo2/hello"
	"log"
	"net"
	"net/rpc"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

// RPC服务 接口规范
// 1. 服务名称
// 2. 服务要实现的详细方法列表
// 3. 注册该类型服务的函数

type HelloService struct{}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "Hello " + request
	return nil
}

func main() {
	hello.RegisterHelloService(new(HelloService))

	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatal("Start tcp server failed:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Server accept client's connection failed:", err)
		}

		go rpc.ServeConn(conn)
	}
}
