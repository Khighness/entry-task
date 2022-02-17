package common

import "time"

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

var (
	HttpAddr string
	RpcAddr  string
)

const (
	Get                 = "GET"
	Post                = "POST"
)

const (
	FileStoragePath     = "./web/public/"
	CookieTokenKey      = "sessionId"
	CookieTokenTimeout  = 24 * time.Hour
)

const (
	RpcSuccessCode      = 10000
	DefaultErrorType    = "服务繁忙"
	DefaultErrorMessage = "请稍后再试"
	CookieErrorType     = "认证失败"
	CookieErrorMessage  = "登录状态已过期"
)

