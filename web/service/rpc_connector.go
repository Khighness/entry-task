package service

import (
	"context"
	"errors"
	"net"

	"github.com/Khighness/entry-task/pkg/rpc"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-22

var errConnFailed error = errors.New("failed to connect to server")

// RpcConnector RPC连接器
type RpcConnector struct {
	rpcServerAddr string
}

// NewRpcConnector 创建RPC连接器
func NewRpcConnector(rpcServerAddr string) RpcConnector {
	return RpcConnector{rpcServerAddr: rpcServerAddr}
}

// Connect 连接RPC服务器，返回客户端
func (c *RpcConnector) Connect(ctx context.Context) (*rpc.Client, error) {
	conn, err := net.Dial("tcp", c.rpcServerAddr)
	if err != nil {
		return nil, errConnFailed
	}
	client := rpc.NewClient(conn)
	return client, nil
}

// Close 关闭RPC客户端对连接
func (c *RpcConnector) Close(client *rpc.Client) {
	client.GetConn().Close()
}
