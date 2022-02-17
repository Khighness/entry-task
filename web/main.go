package main

import (
	"entry/web/common"
	"entry/web/grpc"
	"entry/web/router"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

func main() {
	common.Load()
	grpc.Init()
	router.Start()
}
