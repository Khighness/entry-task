package common

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

// UserInfo 用户信息
type UserInfo struct {
	Id             int64  `json:"id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profilePicture"`
}

// HttpResponse 接口返回信息
type HttpResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// LoginRequest 登陆请求
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 登陆结果
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UpdateProfileRequest 更新账户请求
type UpdateProfileRequest struct {
	Username string `json:"username"`
}
