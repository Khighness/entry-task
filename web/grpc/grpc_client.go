package grpc

import (
	"entry/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"log"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

func NewClient(rpcServerAddr string) (pb.UserServiceClient, error) {
	clientParameters := keepalive.ClientParameters{
		Time:                30 * time.Second, // 客户端每空闲30s ping一下服务器
		Timeout:             1 * time.Second,  // 假设连接已死之前等待1s，等待ping的ack确认
		PermitWithoutStream: true,             // 即使没有活动流也允许ping
	}
	cc, err := grpc.Dial(rpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(clientParameters))
	if err != nil {
		log.Printf("[grpc client] Failed to connect to rpc server %s, err: %s\n", rpcServerAddr, err)
		return nil, err
	}
	return pb.NewUserServiceClient(cc), nil
}
