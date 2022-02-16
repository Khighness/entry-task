package cache

import "time"

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// UserTokenKey 用户token存储的key
const (
	UserTokenKey     = "entry:user:session:"
	UserTokenTimeout = time.Hour * 24
)
