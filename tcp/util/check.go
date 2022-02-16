package util

import (
	"fmt"
	"regexp"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

const (
	NameMinLen   = 3
	NameMaxLen   = 18
	PassMinLen   = 6
	PassMaxLen   = 20
	PassMinLevel = 2
)

const (
	PassLevelD = iota
	PassLevelC
	PassLevelB
	PassLevelA
	PassLevelS
)

// CheckUsername 校验用户名
func CheckUsername(username string) error {
	if len(username) < NameMinLen {
		return fmt.Errorf("用户名长度不得小于%d", NameMinLen)
	}
	if len(username) > NameMaxLen {
		return fmt.Errorf("用户名长度不得大于%d", NameMaxLen)
	}
	return nil
}

// CheckPassword 校验密码
func CheckPassword(password string) error {
	if len(password) < PassMinLen {
		return fmt.Errorf("密码长度不得小于%d", PassMinLen)
	}
	if len(password) > PassMaxLen {
		return fmt.Errorf("密码长度不得大于%d", PassMaxLen)
	}

	var level int = PassLevelD
	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&*_+]`}
	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, password)
		if match {
			level++
		}
	}

	if level < PassMinLevel {
		return fmt.Errorf("密码强度较弱，最少包含数字/字母/特殊符号中的以上两种")
	}
	return nil
}
