package main

import (
	"entry/tcp/common"
	"entry/tcp/server"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

func main() {
	common.Load()
	server.Start()
}
