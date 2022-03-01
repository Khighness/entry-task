package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

var pass = "123456"

// BCR 加密和解密时间，都大于200ms，太恐怖了
func TestEncryptAndVerifyByBCR(t *testing.T) {
	encryptStartTime := time.Now()
	hash, err := EncryptPassByBCR(pass)
	fmt.Println("encrypt time:", time.Since(encryptStartTime))
	fmt.Printf("encrypt:%s, len:%d\n", hash, len(hash))
	assert.Nil(t, err)
	verifyStartTime := time.Now()
	result := VerifyPassByMD5(pass, hash)
	fmt.Println("verify time:", time.Since(verifyStartTime))
	assert.Equal(t, true, result)
}

// MD5 加密5us，解密600ns
func TestEncryptAndVerifyByMD5(t *testing.T) {
	encryptStartTime := time.Now()
	hash, err := EncryptPassByMd5(pass)
	fmt.Println("encrypt time:", time.Since(encryptStartTime))
	fmt.Printf("encrypt:%s, len:%d\n", hash, len(hash))
	assert.Nil(t, err)
	verifyStartTime := time.Now()
	result := VerifyPassByMD5(pass, hash)
	fmt.Println("verify time:", time.Since(verifyStartTime))
	assert.Equal(t, true, result)

	fmt.Println(EncryptPassByMd5("czk911"))
}
