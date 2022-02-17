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
	FileStoragePath     = "./web/public/"
	CookieTokenKey      = "sessionId"
	CookieTokenTimeout  = 24 * time.Hour
	RpcSuccessCode      = 10000
	DefaultErrorType    = "服务繁忙"
	DefaultErrorMessage = "请稍后再试"
)
