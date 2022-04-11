package grpc

import (
	"github.com/Khighness/entry-task/web/common"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-03-09

var GP *GrpcPool

func Init() {
	connector := GrpcConnector{GrpcServerAddr: common.RpcServerAddr}
	GP = NewGrpcPool(connector, &GrpcPoolConfig{
		MaxOpenCount: 10000,
		MaxIdleCount: 5000,
		MaxLifeTime:  30 * time.Minute,
		MaxIdleTime:  10 * time.Minute,
	})
}
