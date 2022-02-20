package test

import (
	"entry/tcp/common/e"
	"entry/tcp/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

func TestCheckUsername(t *testing.T) {
	name1 := "k"
	name2 := "zzzzzzzzzzkkkkkkkkkk"
	name3 := "chen zikang"
	var status int

	status = util.CheckUsername(name1)
	assert.Equal(t, e.ErrorUsernameTooShort, status)
	status = util.CheckUsername(name2)
	assert.Equal(t, e.ErrorUsernameTooLong, status)
	status = util.CheckUsername(name3)
	assert.Equal(t, e.SUCCESS, status)
}

func TestCheckPassword(t *testing.T) {
	pass1 := "k"                      // error
	pass2 := "zzzzzzzzzzzkkkkkkkkkkk" // error
	pass3 := "123456"                 // level 1
	pass4 := "chen zikang"            // level 1
	pass5 := "czk123"                 // level 2
	pass6 := "czk123CZK"              // level 3
	pass7 := "czk123@CZK"             // level 4
	var status int

	status = util.CheckPassword(pass1)
	assert.Equal(t, e.ErrorPasswordTooShort, status)
	status = util.CheckPassword(pass2)
	assert.Equal(t, e.ErrorPasswordTooLong, status)
	status = util.CheckPassword(pass3)
	assert.Equal(t, e.ErrorPasswordNotStrong, status)
	status = util.CheckPassword(pass4)
	assert.Equal(t, e.ErrorPasswordNotStrong, status)
	status = util.CheckPassword(pass5)
	assert.Equal(t, e.SUCCESS, status)
	status = util.CheckPassword(pass6)
	assert.Equal(t, e.SUCCESS, status)
	status = util.CheckPassword(pass7)
	assert.Equal(t, e.SUCCESS, status)
}
