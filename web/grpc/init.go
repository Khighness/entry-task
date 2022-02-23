package grpc

import (
	"context"
	"entry/web/common"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-22

var Pool *ConnPool

// Init 初始化
func Init() {
	Pool = NewPool(context.Background(), &Config{
		InitCount:     50,
		MaxOpenCount:  100,
		MaxIdleCount:  100,
		MaxLifeTime:   time.Hour,
		MaxWaitTime:   time.Second,
		RpcServerAddr: common.RpcServerAddr,
	})
}
