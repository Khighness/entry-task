package grpc

import (
	"github.com/Khighness/entry-task/web/config"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-03-09

var GP *GrpcPool

func InitPool() {
	rpcCfg := config.AppCfg.Rpc
	connector := GrpcConnector{GrpcServerAddr: rpcCfg.Addr}
	GP = NewGrpcPool(connector, &GrpcPoolConfig{
		MaxOpenCount: rpcCfg.MaxOpen,
		MaxIdleCount: rpcCfg.MaxIdle,
		MaxLifeTime:  time.Duration(rpcCfg.MaxLiftTime) * time.Second,
		MaxIdleTime:  time.Duration(rpcCfg.MaxIdleTime) * time.Second,
	})
}
