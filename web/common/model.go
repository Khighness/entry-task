package common

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

// UserInfo 用户信息
type UserInfo struct {
	Id             int64  `json:"id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

// SuccessMsg 成功提示
type SuccessMsg struct {
	SucType   string
	Message   string
	ReturnTip string
	ReturnUrl string
}

// ErrorMsg 错误提示
type ErrorMsg struct {
	ErrType   string
	Message   string
	ReturnTip string
	ReturnUrl string
}
