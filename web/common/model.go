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

type HttpResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
}
