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

// CheckPass 校验
func CheckPass(formPass string, dbPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(formPass), []byte(dbPass))
	return err == nil
}
