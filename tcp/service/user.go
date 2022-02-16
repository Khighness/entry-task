package service

import (
	"context"
	"entry/pb"
	"entry/tcp/common"
	"entry/tcp/mapper"
	"entry/tcp/util"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

//type User interface {
//	Register()
//}
//
//type UserDetail struct {
//	Id             int    `json:"id"`
//	Username       string `json:"username"`
//	Password       string `json:"password"`
//	ProfilePicture string `json:"profile_picture"`
//}
//
//type UserBase struct {
//	Username string `json:"username"`
//	Password string `json:"password"`
//}
//
//type UserInfo struct {
//	Id             int    `json:"id"`
//	Username       string `json:"username"`
//	ProfilePicture string `json:"profile_picture"`
//}
//
//type UserAvatar struct {
//	Id             int    `json:"id"`
//	ProfilePicture string `json:"profile_picture"`
//}
//
//// Register 用户注册
//func (user *UserBase) Register() {
//	encrypt, err1 := util.EncryptPass(user.Password)
//	if err1 != nil {
//		return
//	}
//	res, err2 := model.DB.Exec("INSERT INTO user(username, password) VALUES(?, ?)", user.Username, encrypt)
//	if err2 != nil {
//		return
//	}
//	affected, err3 := res.RowsAffected()
//	if err3 != nil {
//		return
//	}
//	fmt.Println(affected)
//}
//
//// Login 用户登录
//func (user *UserBase) Login() {
//
//}
//
//// UpdateInfo 用户修改信息
//func UpdateInfo() {
//
//}
//
//// UploadAvatar 用户上传头像
//func UploadAvatar() {
//
//}
//
//// GetList 查询所有用户
//func GetList() {
//
//}

type server struct{}

// Register 用户注册
func (s *server) Register(ctx context.Context, user *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var err error

	// 校验用户名和密码
	if err = util.CheckUsername(user.Username); err != nil {
		return &pb.RegisterResponse{
			Code: common.ErrorCode,
			Msg:  err.Error(),
		}, nil
	}
	if err = util.CheckPassword(user.Password); err != nil {
		return &pb.RegisterResponse{
			Code: common.ErrorCode,
			Msg:  err.Error(),
		}, nil
	}

	// BCR加密，计算hash
	var hashedPassword string
	hashedPassword, err = util.EncryptPass(user.Password)
	if err != nil {
		return nil, err
	}

	// 保存到数据库
	err = mapper.SaveUser(user.Username, hashedPassword)
	if err != nil {
		return &pb.RegisterResponse{
			Code: common.ErrorCode,
			Msg:  "注册失败",
		}, nil
	}
	return &pb.RegisterResponse{
		Code: common.SuccessCode,
		Msg:  "注册成功",
	}, nil
}

func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	return nil, nil
}

func (s *server) GetProfile(ctx context.Context, in *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	return nil, nil
}

func (s *server) CheckToken(ctx context.Context, in *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	return nil, nil
}

func (s *server) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	return nil, nil
}
