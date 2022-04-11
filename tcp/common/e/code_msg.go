package e

import (
	"github.com/Khighness/entry-task/tcp/common"
	"fmt"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

// RPC 状态码
const (
	SUCCESS = 10000
	ERROR   = 20000

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

	ErrorOperateDatabase = 40001
)

// 状态码对应信息字典
var codeMessageDic = map[int]string{
	SUCCESS: "SUCCESS",
	ERROR:   "ERROR",

	ErrorUsernameTooShort:     fmt.Sprintf("用户名长度不得小于%d", common.NameMinLen),
	ErrorUsernameTooLong:      fmt.Sprintf("用户名长度不得大于%d", common.NameMaxLen),
	ErrorUsernameAlreadyExist: "用户名已存在，请换个试试",
	ErrorPasswordTooShort:     fmt.Sprintf("密码长度不得小于%d", common.PassMinLen),
	ErrorPasswordTooLong:      fmt.Sprintf("密码长度不得大于%d", common.PassMaxLen),
	ErrorPasswordNotStrong:    "密码强度较弱，最少需要包含数字/字母/特殊符号中的以上两种",
	ErrorUsernameIncorrect:    "用户名错误",
	ErrorPasswordIncorrect:    "密码错误",
	ErrorTokenIncorrect:       "令牌非法",
	ErrorTokenExpired:         "登陆状态已过期",

	ErrorOperateDatabase: "操作数据库失败",
}

// GetMsg 根据状态码获取信息
func GetMsg(code int) string {
	msg, ok := codeMessageDic[code]
	if ok {
		return msg
	}
	return codeMessageDic[ERROR]
}
