package main

import (
	"github.com/Khighness/entry-task/web/config"
	"github.com/Khighness/entry-task/web/router"
	"github.com/Khighness/entry-task/web/service"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

func main() {
	config.Load()
	service.InitPool()
	router.Start()
}
