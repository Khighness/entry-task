package main

import (
	"github.com/Khighness/entry-task/web/config"
	"github.com/Khighness/entry-task/web/grpc"
	"github.com/Khighness/entry-task/web/router"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

func main() {
	config.Load()
	grpc.InitPool()
	router.Start()
}
