package main

import (
	"github.com/Khighness/entry-task/tcp/common"
	"github.com/Khighness/entry-task/tcp/server"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

func main() {
	common.Load()
	server.Start()
}
