package e

import (
	"entry/tcp/common"
	"fmt"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-17

var codeMessageDic = map[int]string{
	SUCCESS: "OK",
	ERROR:   "ERROR",

	ErrorUsernameTooShort:     fmt.Sprintf("用户名长度不得小于%d", common.NameMinLen),
	ErrorUsernameTooLong:      fmt.Sprintf("用户名长度不得小于%d", common.NameMinLen),
	ErrorUsernameAlreadyExist: "用户名已存在，请换个试试",
	ErrorPasswordTooShort:     fmt.Sprintf("密码长度不得小于%d", common.PassMinLen),
	ErrorPasswordTooLong:      fmt.Sprintf("密码长度不得小于%d", common.PassMaxLen),
	ErrorPasswordNotStrong:    "密码强度较弱，最少包含数字/字母/特殊符号中的以上两种",
	ErrorUsernameIncorrect:    "用户名错误",
	ErrorPasswordIncorrect:    "密码错误",
	ErrorTokenIncorrect:       "cookie非法",
	ErrorTokenExpired:         "登陆状态已过期",
}

// GetMsg 根据状态码获取信息
func GetMsg(code int) string {
	msg, ok := codeMessageDic[code]
	if ok {
		return msg
	}
	return codeMessageDic[ERROR]
}
