package common

import (
	"entry/tcp/cache"
	"entry/tcp/model"
	"gopkg.in/ini.v1"
	"log"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

// Load 初始化
func Load() {
	file, err := ini.Load("./tcp/conf/conf.ini")
	if err != nil {
		log.Fatalln("Load config file error, please check file path:", err)
	} else {
		log.Println("Loading config file ...")
	}
	loadServerConfig(file)
	model.Load(file)
	cache.Load(file)
}

// loadServerConfig 导入服务配置
func loadServerConfig(file *ini.File) {
	ServerAddr = file.Section("server").Key("ServerAddr").String()
}
