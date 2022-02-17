package grpc

import (
	"entry/pb"
	"entry/web/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

var Client pb.UserServiceClient

// TODO rpc连接池

func Init() {
	conn, err := grpc.Dial(common.RpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to onnect to grpc tcp server [%s]", common.RpcAddr)
	} else {
		log.Printf("Succeed to connect to grpc tcp server [%s]", common.RpcAddr)
	}

	Client = pb.NewUserServiceClient(conn)
}
