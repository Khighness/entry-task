package common

import (
	"entry/web/logging"
	"gopkg.in/ini.v1"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// Load 初始化
func Load() {
	file, err := ini.Load("./web/conf/conf.ini")
	if err != nil {
		logging.Log.Fatalln("Load config file error, please check file path:", err)
	} else {
		logging.Log.Infoln("Loading config file ...")
	}
	loadServerConfig(file)
	LoadRpcConfig(file)
}

// loadServerConfig 导入服务配置
func loadServerConfig(file *ini.File) {
	server := file.Section("server")
	HttpServerAddr = server.Key("HttpServerAddr").String()
}

// LoadRpcConfig 导入rpc配置
func LoadRpcConfig(file *ini.File) {
	rpc := file.Section("rpc")
	RpcServerAddr = rpc.Key("RpcServerAddr").String()
}
