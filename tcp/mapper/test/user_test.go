package main

import (
	"entry/tcp/mapper"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

func TestSaveUser(t *testing.T) {
	err := mapper.SaveUser("Khighness", "123456")
	assert.Nil(t, err)
}

func TestUpdateUserUsernameById(t *testing.T) {
	err := mapper.UpdateUserUsernameById(1, "Chen Zikang")
	assert.Nil(t, err)
}

func TestUpdateUserProfilePictureById(t *testing.T) {
	err := mapper.UpdateUserProfilePictureById(1, "czk.jpg")
	assert.Nil(t, err)
}

func TestCheckUserUsernameExist(t *testing.T) {
	timeStart := time.Now()
	exist, err := mapper.CheckUserUsernameExist("Khighness")
	fmt.Println("use time:", time.Since(timeStart))
	assert.Nil(t, err)
	assert.Equal(t, true, exist)
}

func TestQueryUserById(t *testing.T) {
	username, profilePicture, err := mapper.QueryUserById(1)
	assert.Nil(t, err)
	fmt.Println("username:", username)
	fmt.Println("profilePicture:", profilePicture)
}

func TestQueryUserByUsername(t *testing.T) {
	timeStart := time.Now()
	id, password, profilePicture, err := mapper.QueryUserByUsername("user_10000010")
	fmt.Println("use time:", time.Since(timeStart)) // 6ms
	assert.Nil(t, err)
	fmt.Println("id:", id)
	fmt.Println("password:", password)
	fmt.Println("profilePicture:", profilePicture)
}
