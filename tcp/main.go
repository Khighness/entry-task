package main

import (
	"entry/tcp/common"
	"entry/tcp/util"
	"fmt"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

func main() {
	common.Load()
	util.CheckPassword("12ka.")
	fmt.Println(util.PassLevelD)
	fmt.Println(util.PassLevelC)
	fmt.Println(util.PassLevelB)
	fmt.Println(util.PassLevelA)
	fmt.Println(util.PassLevelS)
}
