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

func TestCheckUsername(t *testing.T) {
	name1 := "k"
	name2 := "zzzzzzzzzzkkkkkkkkkk"
	name3 := "chen zikang"
	var err error

	err = util.CheckUsername(name1)
	assert.NotNil(t, err)
	fmt.Printf("name: %s, result: %s\n", name1, err.Error())
	err = util.CheckUsername(name2)
	assert.NotNil(t, err)
	fmt.Printf("name: %s, result: %s\n", name2, err.Error())
	err = util.CheckUsername(name3)
	assert.Nil(t, err)
	fmt.Printf("name: %s, result: %s\n", name3, "ok")
}

func TestCheckPassword(t *testing.T) {
	pass1 := "k"                      // err
	pass2 := "zzzzzzzzzzzkkkkkkkkkkk" // err
	pass3 := "123456"                 // level 1
	pass4 := "chen zikang"            // level 1
	pass5 := "czk123"                 // level 2
	pass6 := "czk123CZK"              // level 3
	pass7 := "czk123@CZK"             // level 4
	var err error

	err = util.CheckPassword(pass1)
	assert.NotNil(t, err)
	fmt.Printf("pass: %s, result: %s\n", pass1, err.Error())
	err = util.CheckPassword(pass2)
	assert.NotNil(t, err)
	fmt.Printf("pass: %s, result: %s\n", pass2, err.Error())
	err = util.CheckPassword(pass3)
	assert.NotNil(t, err)
	fmt.Printf("pass: %s, result: %s\n", pass3, err.Error())
	err = util.CheckPassword(pass4)
	assert.NotNil(t, err)
	fmt.Printf("pass: %s, result: %s\n", pass4, err.Error())
	err = util.CheckPassword(pass5)
	assert.Nil(t, err)
	fmt.Printf("pass: %s, result: %s\n", pass5, "ok")
	err = util.CheckPassword(pass6)
	assert.Nil(t, err)
	fmt.Printf("pass: %s, result: %s\n", pass6, "ok")
	err = util.CheckPassword(pass7)
	assert.Nil(t, err)
	fmt.Printf("pass: %s, result: %s\n", pass7, "ok")
}
