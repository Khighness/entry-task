package common

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

var ServerAddr string

// SessionIdBytes sessionId字节数组长度
const (
	SessionIdBytes = 16
)

const (
	NameMinLen   = 3
	NameMaxLen   = 18
	PassMinLen   = 6
	PassMaxLen   = 20
	PassMinLevel = 2
)

const (
	PassLevelD = iota
	PassLevelC
	PassLevelB
	PassLevelA
	PassLevelS
)
