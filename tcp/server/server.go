package server

import (
	"entry/pb"
	"entry/tcp/common"
	"entry/tcp/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

// Start 启动tcp server
func Start() {
	listener, err := net.Listen("tcp", common.ServerAddr)
	if err != nil {
		log.Fatalln("Failed to start tcp server :", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &service.Server{})
	reflection.Register(s)
	log.Printf("GRPC tcp server is serving at [%s]", common.ServerAddr)

	if err = s.Serve(listener); err != nil {
		log.Fatalf("GRPC tcp server failed to serve at [%s]: %s\n", common.ServerAddr, err)
	}
}
