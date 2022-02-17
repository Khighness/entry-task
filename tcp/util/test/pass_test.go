package test

import (
	"entry/tcp/util"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

var pass = "czk911"

func TestEncryptAndVerify(t *testing.T) {
	hash, err := util.EncryptPass(pass)
	assert.Nil(t, err)
	fmt.Printf("encrypt:%s, len:%d\n", hash, len(hash))
	result := util.VerifyPass(pass, hash)
	assert.Equal(t, true, result)
}
