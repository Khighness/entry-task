package common

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

const (
	// HeaderTokenKey HTTP请求Header中的令牌Key
	HeaderTokenKey = "Authorization"
	// AvatarStoragePath 头像存储相对路径
	AvatarStoragePath = "./web/public/avatar/"
)

const (
	RpcSuccessCode             = 10000
	HttpSuccessCode            = 10000
	HttpSuccessMessage         = "SUCCESS"
	HttpErrorCode              = 20000
	HttpErrorMessage           = "ERROR"
	HttpErrorServerBusyCode    = 20001
	HttpErrorServerBusyMessage = "Server is busy, please try again later"
	HttpErrorRpcRequestCode    = 20002
	HttpErrorRpcRequestMessage = "RPC failed or timeout"
)
