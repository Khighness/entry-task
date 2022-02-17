package common

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

// UserInfo 用户信息
type UserInfo struct {
	Id             int64
	Username       string
	ProfilePicture string
}

// ErrorMsg 错误提示
type ErrorMsg struct {
	ErrType string
	Message string
}
