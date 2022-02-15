package conf

import (
	"entry/tcp/cache"
	"entry/tcp/model"
	web "entry/web/common"
	"gopkg.in/ini.v1"
	"log"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// Load 导入配置文件
func Load() {
	file, err := ini.Load("./conf/conf.ini")
	if err != nil {
		log.Println("Load config file error, please check file path")
		panic(err)
	} else {
		log.Println("Loading config file ...")
	}
	web.Http(file)
	model.MySQL(file)
	cache.Redis(file)
}
