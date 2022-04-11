package mapper

import (
	"github.com/Khighness/entry-task/tcp/model"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

type UserMapper struct {
}

// SaveUser 保存用户信息
func (userMapper *UserMapper) SaveUser(user *model.User) (int64, error) {
	result, err := model.DB.Exec("INSERT INTO user(`username`, `password`, `profile_picture`) values(?, ?, ?)", user.Username, user.Password, user.ProfilePicture)
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	return id, nil
}

// UpdateUserUsernameById 根据id更新用户名
func (userMapper *UserMapper) UpdateUserUsernameById(id int64, username string) error {
	_, err := model.DB.Exec("UPDATE `user` SET `username` = ? WHERE `id` = ?", username, id)
	return err
}

// UpdateUserProfilePictureById 根据id更新用户头像
func (userMapper *UserMapper) UpdateUserProfilePictureById(id int64, profilePicture string) error {
	_, err := model.DB.Exec("UPDATE `user` SET `profile_picture` = ? WHERE `id` = ?", profilePicture, id)
	return err
}

// CheckUserUsernameExist 检查用户名是否已存在
// MySQL比较字符串在大小写敏感的情况下，必须转binary
// 在binary下，查询username无法走index_username
// 为了保证性能，用户名的大小写不敏感
func (userMapper *UserMapper) CheckUserUsernameExist(username string) (bool, error) {
	var count int64
	err := model.DB.QueryRow("SELECT COUNT(`username`) FROM `user` WHERE `username` = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// QueryUserById 根据id查询用户信息
func (userMapper *UserMapper) QueryUserById(id int64) (user *model.User, err error) {
	user = new(model.User)
	err = model.DB.QueryRow("SELECT `username`, `password`, `profile_picture` FROM `user` WHERE id = ?", id).Scan(&user.Username, &user.Password, &user.ProfilePicture)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// QueryUserByUsername 根据username查询用户信息
func (userMapper *UserMapper) QueryUserByUsername(username string) (user *model.User, err error) {
	user = new(model.User)
	err = model.DB.QueryRow("SELECT `id`, `password`, `profile_picture` FROM `user` WHERE `username` = ?", username).Scan(&user.Id, &user.Password, &user.ProfilePicture)
	if err != nil {
		return nil, err
	}
	return user, nil
}
