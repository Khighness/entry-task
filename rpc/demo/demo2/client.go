package main

import (
	"entry/rpc/demo/demo2/hello"
	"log"
	"net/rpc"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

type HelloServiceClient struct {
	*rpc.Client
}

var _ hello.HelloServiceInterface = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{Client: c}, nil
}

func (p *HelloServiceClient) Hello(request string, reply *string) error {
	return p.Client.Call(hello.HelloServiceName+".Hello", request, reply)
}

func main() {
	client, err := DialHelloService("tcp", "127.0.0.1:12345")
	if err != nil {
		log.Fatal("Client connect to server failed:", err)
	}

	var reply string
	err = client.Hello("KHighness", &reply)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(reply)
}
