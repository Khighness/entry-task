package service

import (
	"entry/tcp/model"
	"entry/tcp/util"
	"fmt"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

type UserDetail struct {
	Id             int    `json:"id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profile_picture"`
}

type UserBase struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserData struct {
	Id             int    `json:"id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

type UserAvatar struct {
	Id             int    `json:"id"`
	ProfilePicture string `json:"profile_picture"`
}

// Register 用户注册
func (user *UserBase) Register() {
	encrypt, err1 := util.EncryptPass(user.Password)
	if err1 != nil {
		return
	}
	res, err2 := model.DB.Exec("INSERT INTO user(username, password) VALUES(?, ?)", user.Username, encrypt)
	if err2 != nil {
		return
	}
	affected, err3 := res.RowsAffected()
	if err3 != nil {
		return
	}
	fmt.Println(affected)
}

// Login 用户登录
func (user *UserBase) Login() {

}

// Update 用户修改信息
func Update() {

}

// UserUploadAvatar 用户上传头像
func UserUploadAvatar() {

}

// UserList 查询所有用户
func UserList() {

}
