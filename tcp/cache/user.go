package cache

import "time"

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// UserTokenKeyPrefix 用户token存储的key
// UserTokenTimeout   用户token过期时间
const (
	UserTokenKeyPrefix = "entry:user:session:"
	UserTokenTimeout   = time.Hour * 24
)
