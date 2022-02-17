package service

import (
	"context"
	"entry/pb"
	"entry/tcp/cache"
	"entry/tcp/mapper"
	"entry/tcp/util"
	"entry/tcp/util/e"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

type Server struct{}

// Register 用户注册
func (s *Server) Register(ctx context.Context, user *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var status = e.SUCCESS
	var err error

	// 校验用户名和密码
	if status = util.CheckUsername(user.Username); status != e.SUCCESS {
		return &pb.RegisterResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	if status = util.CheckPassword(user.Password); status != e.SUCCESS {
		return &pb.RegisterResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}

	// 检查用户名的唯一性
	exist, err := mapper.CheckUserUsernameExist(user.Username)
	if err != nil {
		return nil, err
	}
	if exist {
		status = e.ErrorUsernameAlreadyExist
		return &pb.RegisterResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
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
		Code: int32(status),
		Msg:  e.GetMsg(status),
	}, nil
}

func (s *Server) Login(ctx context.Context, user *pb.LoginRequest) (*pb.LoginResponse, error) {
	var status = e.SUCCESS

	// 校验用户名是否存在
	id, hashedPassword, profilePicture, err := mapper.QueryUserByUsername(user.Username)
	if err != nil {
		status = e.ErrorUsernameIncorrect
		return &pb.LoginResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	// 验证密码
	if !util.VerifyPass(user.Password, hashedPassword) {
		status = e.ErrorPasswordIncorrect
		return &pb.LoginResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}

	// 生成sessionId，并在redis中缓存用户id
	sessionId := util.GenerateSessionId()
	tokenKey := cache.UserTokenKey + sessionId
	cache.RedisClient.HSet(tokenKey, "id", id)
	cache.RedisClient.Expire(tokenKey, cache.UserTokenTimeout)

	return &pb.LoginResponse{
		Code:      int32(status),
		Msg:       e.GetMsg(status),
		SessionId: sessionId,
		User: &pb.User{
			Id:             id,
			Username:       user.Username,
			ProfilePicture: profilePicture,
		},
	}, nil
}

func (s *Server) CheckToken(ctx context.Context, in *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	var status = e.SUCCESS
	var tokenKey = cache.UserTokenKey + in.SessionId
	_, err := cache.RedisClient.HGet(tokenKey, "id").Int()
	if err != nil {
		status = e.ErrorTokenExpired
		return &pb.CheckTokenResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	return &pb.CheckTokenResponse{
		Code: int32(status),
		Msg:  e.GetMsg(status),
	}, nil
}

func (s *Server) GetProfile(ctx context.Context, in *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var status = e.SUCCESS
	id, err := cache.RedisClient.HGet(cache.UserTokenKey+in.SessionId, "id").Int64()
	if err != nil {
		status = e.ErrorTokenExpired
		return &pb.GetProfileResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
			User: nil,
		}, nil
	}
	username, profilePicture, err := mapper.QueryUserById(id)
	if err != nil {
		status = e.ErrorTokenExpired
		return &pb.GetProfileResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
			User: nil,
		}, nil
	}
	return &pb.GetProfileResponse{
		Code: int32(status),
		Msg:  e.GetMsg(status),
		User: &pb.User{
			Id:             id,
			Username:       username,
			ProfilePicture: profilePicture,
		},
	}, nil
}

func (s *Server) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	var status = e.SUCCESS
	id, err := cache.RedisClient.HGet(cache.UserTokenKey+in.SessionId, "id").Int64()
	if err != nil {
		status = e.ErrorTokenExpired
		return &pb.UpdateProfileResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	if in.Username != "" {
		exist, err := mapper.CheckUserUsernameExist(in.Username)
		if err != nil {
			return nil, err
		}
		if exist {
			status = e.ErrorUsernameAlreadyExist
			return &pb.UpdateProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
			}, nil
		}
		err = mapper.UpdateUserUsernameById(id, in.Username)
		if err != nil {
			return nil, err
		}
	}
	if in.ProfilePicture != "" {
		err = mapper.UpdateUserProfilePictureById(id, in.ProfilePicture)
		if err != nil {
			return nil, err
		}
	}
	return &pb.UpdateProfileResponse{
		Code: int32(status),
		Msg:  e.GetMsg(status),
	}, nil
}
