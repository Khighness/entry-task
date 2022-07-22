package rpc

import (
	"github.com/Khighness/entry-task/rpc"
	"github.com/Khighness/entry-task/web/config"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-22

var Pool *rpc.ConnPool

func InitPool() {
	rpcCfg := config.AppCfg.Rpc
	poolCfg := &rpc.Config{
		MaxOpenCount:  rpcCfg.MaxOpen,
		MaxIdleCount:  rpcCfg.MaxIdle,
		RpcServerAddr: rpcCfg.Addr,
	}
	Pool = rpc.NewPool(poolCfg)
}
