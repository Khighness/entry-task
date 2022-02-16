package web

import (
	"gopkg.in/ini.v1"
	"log"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// Load 初始化
func Load() {
	file, err := ini.Load("./doc/doc.ini")
	if err != nil {
		log.Println("Load config file error, please check file path")
		panic(err)
	} else {
		log.Println("Loading config file ...")
	}
	loadServerConfig(file)
}

func loadServerConfig(file *ini.File) {
	server := file.Section("server")
	HttpPort = server.Key("HttpPort").String()
}
