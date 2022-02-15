package main

import (
	"entry/conf"
	"entry/web/router"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

func main() {
	conf.Load()
	router.NewRouter()
}
