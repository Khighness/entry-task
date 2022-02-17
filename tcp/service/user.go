package service

import (
	"context"
	"entry/pb"
	"entry/tcp/cache"
	"entry/tcp/mapper"
	"entry/tcp/util"
	"entry/tcp/util/e"
	"log"
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
		status = e.ErrorOperateDatabase
		return &pb.RegisterResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
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
	hashedPassword, _ = util.EncryptPass(user.Password)

	// 保存到数据库
	err = mapper.SaveUser(user.Username, hashedPassword)
	if err != nil {
		status = e.ErrorOperateDatabase
		return &pb.RegisterResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}

	log.Printf("【用户注册】 用户名：%s \n", user.Username)
	return &pb.RegisterResponse{
		Code: int32(status),
		Msg:  e.GetMsg(status),
	}, nil
}

// Login 用户登录
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
		log.Printf("pass:%s, hash:%s, res:%v\n", user.Password, hashedPassword, util.VerifyPass(user.Password, hashedPassword))
		status = e.ErrorPasswordIncorrect
		return &pb.LoginResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}

	// 生成sessionId，并在redis中缓存用户id
	sessionId := util.GenerateSessionId()
	tokenKey := cache.UserTokenKeyPrefix + sessionId
	cache.RedisClient.HSet(tokenKey, "id", id)
	cache.RedisClient.Expire(tokenKey, cache.UserTokenTimeout)

	log.Printf("【用户登录】 用户名：%s \n", user.Username)
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

// CheckToken 检查token
func (s *Server) CheckToken(ctx context.Context, in *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	var status = e.SUCCESS
	var tokenKey = cache.UserTokenKeyPrefix + in.SessionId
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

// GetProfile 获取信息
func (s *Server) GetProfile(ctx context.Context, in *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var status = e.SUCCESS

	// 从缓存中获取用户id
	id, err := cache.RedisClient.HGet(cache.UserTokenKeyPrefix+in.SessionId, "id").Int64()
	if err != nil {
		status = e.ErrorTokenExpired
		return &pb.GetProfileResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
			User: nil,
		}, nil
	}

	// 从数据库中查询
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

// UpdateProfile 更新信息
func (s *Server) UpdateProfile(ctx context.Context, user *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	var status = e.SUCCESS

	// 从缓存中获取用户id
	id, err := cache.RedisClient.HGet(cache.UserTokenKeyPrefix+user.SessionId, "id").Int64()
	if err != nil {
		status = e.ErrorTokenExpired
		return &pb.UpdateProfileResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}

	// 用户名不为空
	if user.Username != "" {
		if status = util.CheckUsername(user.Username); status != e.SUCCESS {
			return &pb.UpdateProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
			}, nil
		}
		exist, err := mapper.CheckUserUsernameExist(user.Username)
		if err != nil {
			status = e.ErrorOperateDatabase
			return &pb.UpdateProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
			}, nil
		}
		if exist {
			status = e.ErrorUsernameAlreadyExist
			return &pb.UpdateProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
			}, nil
		}
		err = mapper.UpdateUserUsernameById(id, user.Username)
		if err != nil {
			status = e.ErrorOperateDatabase
			return &pb.UpdateProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
			}, nil
		}

		log.Printf("【用户更新】 id：%d，用户名：%s\n", id, user.Username)
	}

	// 用户头像不为空
	if user.ProfilePicture != "" {
		err = mapper.UpdateUserProfilePictureById(id, user.ProfilePicture)
		if err != nil {
			status = e.ErrorOperateDatabase
			return &pb.UpdateProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
			}, nil
		}

		log.Printf("【用户更新】 id：%d，头像：%s\n", id, user.ProfilePicture)
	}

	return &pb.UpdateProfileResponse{
		Code: int32(status),
		Msg:  e.GetMsg(status),
	}, nil
}
