package common

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

const (
	// TokenBytes Token字节数组长度
	TokenBytes = 16
	// DefaultProfilePicture 注册账户的默认头像
	DefaultProfilePicture = "http://127.0.0.1:10000/avatar/show/default.jpg"
)

// 注册要求
const (
	NameMinLen   = 3
	NameMaxLen   = 18
	PassMinLen   = 6
	PassMaxLen   = 20
	PassMinLevel = PassLevelB
)

// 密码强度
const (
	PassLevelD = iota
	PassLevelC
	PassLevelB
	PassLevelA
	PassLevelS
)
