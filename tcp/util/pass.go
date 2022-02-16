package util

import "golang.org/x/crypto/bcrypt"

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

const (
	PasswordCost = 12
)

// EncryptPass 加密
func EncryptPass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyPass 校验
// password 用户输入密码
// hashedPassword 数据库存储密码
func VerifyPass(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
