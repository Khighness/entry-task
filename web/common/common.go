package common

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

var (
	HttpServerAddr string
	RpcServerAddr  string
)

const (
	HeaderTokenKey    = "Authorization"
	AvatarStoragePath = "./web/public/avatar/"
)

const (
	RpcSuccessCode     = 10000
	HttpSuccessCode    = 10000
	HttpSuccessMessage = "SUCCESS"
	HttpErrorCode      = 50000
	HttpErrorMessage   = "ERROR"
)
