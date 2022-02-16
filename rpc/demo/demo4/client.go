package main

import (
	"context"
	pb "entry/pb"
	"google.golang.org/grpc"
	"log"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

func main() {
	conn, err := grpc.Dial("127.0.0.1:12345", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Connect to tcp server failed:", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	name := "KHighness"
	reply, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalln("grpc failed:", err)
	}
	log.Println("reply:", reply)

}
