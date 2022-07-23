package service

import (
	"time"

	"github.com/Khighness/entry-task/web/config"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-22

var Pool *RpcPool

func InitPool() {
	rpcCfg := config.AppCfg.Rpc
	connector := NewRpcConnector(rpcCfg.Addr)
	poolCfg := &RpcPoolConfig{
		InitialCount: rpcCfg.Initial,
		MaxOpenCount: rpcCfg.MaxOpen,
		MaxIdleCount: rpcCfg.MaxIdle,
		MaxLifeTime:  time.Duration(rpcCfg.MaxLifeTime) * time.Second,
		MaxIdleTime:  time.Duration(rpcCfg.MaxIdleTime) * time.Second,
	}
	Pool = NewRpcPool(connector, poolCfg)
}
