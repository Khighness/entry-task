package mapper

import "entry/tcp/model"

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-16

// SaveUser 保存用户信息
func SaveUser(username, password string) error {
	_, err := model.DB.Exec("INSERT INTO user(`username`, `password`) values(?, ?)", username, password)
	return err
}

// UpdateUserUsernameById 根据id更新用户名
func UpdateUserUsernameById(id int64, username string) error {
	_, err := model.DB.Exec("UPDATE user SET username = ? WHERE id = ?", username, id)
	return err
}

// UpdateUserProfilePictureById 根据id更新用户头像
func UpdateUserProfilePictureById(id int64, profilePicture string) error {
	_, err := model.DB.Exec("UPDATE user SET profile_picture = ? WHERE id = ?", profilePicture, id)
	return err
}

// CheckUserUsernameExist 检查用户名是否已存在
func CheckUserUsernameExist(username string) (bool, error) {
	var count int64
	err := model.DB.QueryRow("SELECT COUNT(username) FROM user WHERE username = ?", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// QueryUserById 根据id查询用户信息
func QueryUserById(id int64) (username, profilePicture string, err error) {
	err = model.DB.QueryRow("SELECT username, profile_picture FROM user WHERE id = ?", id).Scan(&username, &profilePicture)
	if err != nil {
		return "", "", nil
	}
	return username, profilePicture, nil
}

// QueryUserPasswordById 根据id查询用户密码
func QueryUserPasswordById(id int64) (password string, err error) {
	err = model.DB.QueryRow("SELECT password FROM user WHERE id = ?", id).Scan(&password)
	if err != nil {
		return "", nil
	}
	return password, nil
}
