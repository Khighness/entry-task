package service

import (
	"context"
	"entry/pb"
	"entry/tcp/cache"
	"entry/tcp/common"
	"entry/tcp/mapper"
	"entry/tcp/util"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

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

	// 检查用户名的唯一性
	exist, err := mapper.CheckUserUsernameExist(user.Username)
	if err != nil {
		return nil, err
	}
	if exist {
		return &pb.RegisterResponse{
			Code: common.ErrorCode,
			Msg:  "该用户名已被其他用户使用",
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
		return nil, err
	}
	return &pb.RegisterResponse{
		Code: common.SuccessCode,
		Msg:  "注册成功",
	}, nil
}

func (s *server) Login(ctx context.Context, user *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 校验用户名是否存在
	id, hashedPassword, profilePicture, err := mapper.QueryUserByUsername(user.Username)
	if err != nil {
		return &pb.LoginResponse{
			Code: common.ErrorCode,
			Msg:  "用户名错误",
		}, nil
	}
	// 验证密码
	if !util.VerifyPass(user.Password, hashedPassword) {
		return &pb.LoginResponse{
			Code: common.ErrorCode,
			Msg:  "密码错误",
		}, nil
	}

	// 生成sessionId，并在redis中缓存用户id
	sessionId := util.GenerateSessionId()
	tokenKey := cache.UserTokenKey + sessionId
	cache.RedisClient.HSet(tokenKey, "id", id)
	cache.RedisClient.Expire(tokenKey, cache.UserTokenTimeout)

	return &pb.LoginResponse{
		Code:      common.SuccessCode,
		Msg:       "登陆成功",
		SessionId: sessionId,
		User:      &pb.User{
			Id:             id,
			Username:       user.Username,
			ProfilePicture: profilePicture,
		},
	}, nil
}

func (s *server) CheckToken(ctx context.Context, in *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	_, err := cache.RedisClient.HGet(cache.UserTokenKey+in.SessionId, "id").Int()
	if err != nil {
		return &pb.CheckTokenResponse{
			Code: common.ErrorCode,
			Msg:  "登录状态已过期",
		}, nil
	}
	return &pb.CheckTokenResponse{
		Code: common.SuccessCode,
		Msg:  "session校验成功",
	}, nil
}

func (s *server) GetProfile(ctx context.Context, in *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	return nil, nil
}

func (s *server) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	return nil, nil
}
