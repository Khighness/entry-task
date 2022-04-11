package common

import (
	"github.com/Khighness/entry-task/tcp/cache"
	"github.com/Khighness/entry-task/tcp/logging"
	"github.com/Khighness/entry-task/tcp/model"
	"gopkg.in/ini.v1"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

var (
	ServerAddr string
)

// Load 初始化
func Load() {
	file, err := ini.Load("./tcp/conf/conf.ini")
	if err != nil {
		logging.Log.Fatalln("Load config file error, please check file path:", err)
	} else {
		logging.Log.Infoln("Loading config file ...")
	}
	loadServerConfig(file)
	model.Load(file)
	cache.Load(file)
}

// Load 导入配置
func loadServerConfig(file *ini.File) {
	ServerAddr = file.Section("server").Key("ServerAddr").String()
}
