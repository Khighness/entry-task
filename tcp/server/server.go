package server

import (
	"entry/pb"
	"entry/tcp/common"
	"entry/tcp/logging"
	"entry/tcp/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
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

	enforcementPolicy := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Minute, // 客户端两次ping的等待
		PermitWithoutStream: true,            // 即使没有活动流也允许ping
	}
	serverParameters := keepalive.ServerParameters{
		MaxConnectionIdle:     30 * time.Minute, // 如果客户端空闲30m，断连
		MaxConnectionAge:      time.Hour,        // 任何客户端存活1h，断连
		MaxConnectionAgeGrace: 5 * time.Second,  // 在强制关闭连接之前，等待5s，让rpc完成
		Time:                  1 * time.Minute,  // 如果客户端空闲1分钟，ping客户端以确保连接正常
		Timeout:               1 * time.Second,  // 如果ping请求1s内未恢复，则认为连接断开
	}
	s := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(enforcementPolicy), grpc.KeepaliveParams(serverParameters))
	pb.RegisterUserServiceServer(s, &service.Server{})
	reflection.Register(s)
	logging.Log.Printf("GRPC tcp server is serving at [%s]", common.ServerAddr)

	if err = s.Serve(listener); err != nil {
		logging.Log.Fatalf("GRPC tcp server failed to serve at [%s]: %s", common.ServerAddr, err)
	}
}
