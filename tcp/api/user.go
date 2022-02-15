package api

import "entry/tcp/service"

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// 用户注册
func UserRegister() {
	user := &service.UserBase{
		Username: "",
		Password: "",
	}
	user.Register()
}

// 用户登录
func UserLogin() {

}

// 用户修改信息
func UserUpdate() {

}

// 用户上传头像
func UserUploadAvatar() {

}

// 查询所有用户
func UserList() {

}
