package main

import (
	"github.com/Khighness/entry-task/tcp/cache"
	"github.com/Khighness/entry-task/tcp/model"
	"github.com/Khighness/entry-task/tcp/server"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

func main() {
	model.InitMySQL()
	cache.InitRedis()
	server.Start()
}
