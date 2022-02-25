package grpc

import (
	"entry/web/common"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-22

var Pool *ConnPool

// Init 初始化
func Init() {
	Pool = NewPool(&Config{
		InitCount:     500,
		MaxOpenCount:  1000,
		MaxIdleCount:  1000,
		MaxLifeTime:   time.Hour,
		MaxWaitTime:   3 * time.Second,
		RpcServerAddr: common.RpcServerAddr,
	})
}
