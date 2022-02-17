package util

import (
	"entry/tcp/common"
	"entry/tcp/util/e"
	"regexp"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

// CheckUsername 校验用户名
func CheckUsername(username string) int {
	if len(username) < common.NameMinLen {
		return e.ErrorUsernameTooShort
	}
	if len(username) > common.NameMaxLen {
		return e.ErrorUsernameTooLong
	}
	return e.SUCCESS
}

// CheckPassword 校验密码
func CheckPassword(password string) int {
	if len(password) < common.PassMinLen {
		return e.ErrorPasswordTooShort
	}
	if len(password) > common.PassMaxLen {
		return e.ErrorPasswordTooLong
	}

	var level int = common.PassLevelD
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*_+]`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, password)
		if match {
			level++
		}
	}

	if level < common.PassMinLevel {
		return e.ErrorPasswordNotStrong
	}
	return e.SUCCESS
}
