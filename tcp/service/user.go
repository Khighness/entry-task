package service

import (
	"context"
	"entry/pb"
	"entry/tcp/cache"
	"entry/tcp/mapper"
	"entry/tcp/model"
	"entry/tcp/util"
	"entry/tcp/util/e"
	"log"
	"time"
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

	// 生成sessionId，并在redis中缓存用户信息
	sessionId := util.GenerateSessionId()
	cache.SetUserInfo(sessionId, &model.User{
		Id:             id,
		Username:       user.Username,
		Password:       "",
		ProfilePicture: profilePicture,
	})

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
	_, err := cache.GetUserId(in.SessionId)
	if err != nil {
		status = e.ErrorTokenExpired
		return &pb.CheckTokenResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	//log.Printf("【令牌校验】 sessionId：%s，id：%d\n", in.SessionId, id)
	return &pb.CheckTokenResponse{
		Code: int32(status),
		Msg:  e.GetMsg(status),
	}, nil
}

// GetProfile 获取信息
// TODO: 多节点，分布式锁
func (s *Server) GetProfile(ctx context.Context, in *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	var status = e.SUCCESS
	var id int64
	var username string
	var profilePicture string

	// 从缓存中获取用户信息
	user, err := cache.GetUserInfo(in.SessionId)
	if err != nil {
		status = e.ErrorTokenExpired
		return &pb.GetProfileResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
			User: nil,
		}, nil
	}

	// 用户信息失效，说明已更新
	// 从数据库获取，添加到缓存
	if user.Username == "" || user.ProfilePicture == "" {
		username, profilePicture, err = mapper.QueryUserById(user.Id)
		if err != nil {
			status = e.ErrorOperateDatabase
			return &pb.GetProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
				User: nil,
			}, nil
		}
		cache.SetUserField(in.SessionId, "username", username)
		cache.SetUserField(in.SessionId, "profile_picture", profilePicture)
	} else {
		username = user.Username
		profilePicture = user.ProfilePicture
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
// 更新字段，延时2S双删
func (s *Server) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	var status = e.SUCCESS

	// 从缓存中获取用户id
	user, err := cache.GetUserInfo(in.SessionId)
	if err != nil {
		status = e.ErrorTokenExpired
		return &pb.UpdateProfileResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	// 从数据库查询最新信息
	username, profilePicture, err := mapper.QueryUserById(user.Id)
	if err != nil {
		status = e.ErrorOperateDatabase
		return &pb.UpdateProfileResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	user.Username = username
	user.ProfilePicture = profilePicture

	// 用户名不为空
	if in.Username != "" {
		// 检查用户名合法性
		if status = util.CheckUsername(in.Username); status != e.SUCCESS {
			return &pb.UpdateProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
			}, nil
		}
		// 检查用户名是否变动
		if in.Username != user.Username {
			// 检查是否有其他人已使用
			exist, err := mapper.CheckUserUsernameExist(in.Username)
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

			cache.DelUserField(in.SessionId, "username")
			defer func() {
				go func() {
					time.Sleep(2 * time.Second)
					cache.DelUserField(in.SessionId, "username")
				}()
			}()
			err = mapper.UpdateUserUsernameById(user.Id, in.Username)
			if err != nil {
				status = e.ErrorOperateDatabase
				return &pb.UpdateProfileResponse{
					Code: int32(status),
					Msg:  e.GetMsg(status),
				}, nil
			}
			log.Printf("【用户更新】 id：%d，用户名：%s\n", user.Id, in.Username)
		}
	}

	// 用户头像不为空，说明已更新
	if in.ProfilePicture != "" {
		cache.DelUserField(in.SessionId, "username")
		defer func() {
			go func() {
				time.Sleep(2 * time.Second)
				cache.DelUserField(in.SessionId, "profile_picture")
			}()
		}()
		err = mapper.UpdateUserProfilePictureById(user.Id, in.ProfilePicture)
		if err != nil {
			status = e.ErrorOperateDatabase
			return &pb.UpdateProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
			}, nil
		}
		log.Printf("【用户更新】 id：%d，头像：%s\n", user.Id, in.ProfilePicture)
	}

	return &pb.UpdateProfileResponse{
		Code: int32(status),
		Msg:  e.GetMsg(status),
	}, nil
}
