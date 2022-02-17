package e

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

const (
	SUCCESS = 10000 // 成功
	ERROR   = 20000 // 失败

	ErrorUsernameTooShort     = 30001
	ErrorUsernameTooLong      = 30002
	ErrorUsernameAlreadyExist = 30003
	ErrorPasswordTooShort     = 30004
	ErrorPasswordTooLong      = 30005
	ErrorPasswordNotStrong    = 30006
	ErrorUsernameIncorrect    = 30007
	ErrorPasswordIncorrect    = 30008
	ErrorTokenIncorrect       = 30009
	ErrorTokenExpired         = 30010
)
