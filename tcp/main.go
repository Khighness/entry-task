package main

import (
	"entry/tcp/cache"
	"entry/tcp/common"
	"fmt"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

func main() {
	common.Load()
	//cache.RedisClient.HSet("k", "id", "1")
	//cache.RedisClient.Expire("k", 3 * time.Second)
	fmt.Println(cache.RedisClient.HGet("k", "id").Int())
}
