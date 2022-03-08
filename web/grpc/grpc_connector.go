package grpc

import (
	"context"
	"entry/pb"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-03-08

var errConnFailed error = errors.New("connect to server failed")

// GrpcConnector grpc连接器
type GrpcConnector struct {
	GrpcServerAddr string
}

// Connect 连接grpc服务器，返回客户端
func (connector *GrpcConnector) Connect(ctx context.Context) (pb.UserServiceClient, error) {
	clientParameters := keepalive.ClientParameters{
		Time:                30 * time.Second, // 客户端每空闲30s ping一下服务器
		Timeout:             1 * time.Second,  // 假设连接已死之前等待1s，等待ping的ack确认
		PermitWithoutStream: true,             // 即使没有活动流也允许ping
	}
	cc, err := grpc.Dial(connector.GrpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(clientParameters))
	if err != nil {
		return nil, errConnFailed
	}
	client := pb.NewUserServiceClient(cc)
	return client, nil
}
