package main

import (
	"log"
	"net"
	"net/rpc"
)

// @Author KHighness
// @Update 2022-02-16


type HelloService struct {}

func (p *HelloService) Hello(name string, reply *string) error {
	*reply = "Hello " + name
	return nil
}

func main() {
	rpc.RegisterName("HelloService", new(HelloService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Start tcp server failed:", err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Println("Server accept client's connection failed:", err)
	}

	rpc.ServeConn(conn)
}