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
	file, err := ini.Load("./conf/conf.ini")
	if err != nil {
		log.Println("Load config file error, please check file path")
		panic(err)
	} else {
		log.Println("Loading config file ...")
	}
	model.MySQL(file)
	cache.Redis(file)
}
