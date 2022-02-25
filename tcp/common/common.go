package common

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// TokenBytes token字节数组长度
// DefaultProfilePicture 注册的默认图片
const (
	TokenBytes            = 16
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
