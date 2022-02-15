package web

import "gopkg.in/ini.v1"

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

func Http(file *ini.File) {
	loadServerConfig(file)
}

func loadServerConfig(file *ini.File) {
	server := file.Section("server")
	HttpPort = server.Key("HttpPort").String()
}
