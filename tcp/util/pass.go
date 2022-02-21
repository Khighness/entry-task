package util

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

const (
	PasswordCost = 12
)

// EncryptPassByBCR BCR加密
func EncryptPassByBCR(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPassByBCR BCR校验
// password 用户输入密码
// hashedPassword 数据库存储密码
func VerifyPassByBCR(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// EncryptPassByMd5 MD5加密
func EncryptPassByMd5(password string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", err
	}
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes), nil
}

// VerifyPassByMD5 MD5校验
// password 用户输入密码
// hashedPassword 数据库存储密码
func VerifyPassByMD5(password string, hashedPassword string) bool {
	inputPassword, _ := EncryptPassByMd5(password)
	return inputPassword == hashedPassword
}
