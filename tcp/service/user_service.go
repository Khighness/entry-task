package service

import (
	"time"

	"github.com/Khighness/entry-task/pb"
	"github.com/Khighness/entry-task/tcp/cache"
	"github.com/Khighness/entry-task/tcp/common"
	"github.com/Khighness/entry-task/tcp/common/e"
	"github.com/Khighness/entry-task/tcp/logging"
	"github.com/Khighness/entry-task/tcp/mapper"
	"github.com/Khighness/entry-task/tcp/model"
	"github.com/Khighness/entry-task/tcp/util"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-22

// UserService 用户业务
type UserService struct {
	userMapper *mapper.UserMapper
	userCache  *cache.UserCache
}

// NewUserService 创建用户业务操作
func NewUserService(userMapper *mapper.UserMapper, userCache *cache.UserCache) *UserService {
	return &UserService{
		userMapper: userMapper,
		userCache:  userCache,
	}
}

// Register 用户注册
func (s *UserService) Register(in pb.RegisterRequest) (pb.RegisterResponse, error) {
	var status int
	var err error

	// 校验用户名和密码
	if status = util.CheckUsername(in.Username); status != e.SUCCESS {
		return pb.RegisterResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	if status = util.CheckPassword(in.Password); status != e.SUCCESS {
		return pb.RegisterResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}

	// 检查用户名的唯一性
	exist, err := s.userMapper.CheckUserUsernameExist(in.Username)
	if err != nil {
		status = e.ErrorOperateDatabase
		return pb.RegisterResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	if exist {
		status = e.ErrorUsernameAlreadyExist
		return pb.RegisterResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}

	// 加密，计算hash
	var hashedPassword string
	hashedPassword, _ = util.EncryptPassByMd5(in.Password)

	// 保存到数据库
	id, err := s.userMapper.SaveUser(&model.User{
		Username:       in.Username,
		Password:       hashedPassword,
		ProfilePicture: common.DefaultProfilePicture,
	})
	if err != nil {
		status = e.ErrorOperateDatabase
		return pb.RegisterResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}

	logging.Log.Infof("[user register] userId: %v, username：%s", id, in.Username)
	return pb.RegisterResponse{
		Code: int32(status),
		Msg:  e.GetMsg(status),
	}, nil
}

// Login 用户登录
func (s *UserService) Login(in pb.LoginRequest) (pb.LoginResponse, error) {
	var status = e.SUCCESS

	// 校验用户名是否存在
	dbUser, err := s.userMapper.QueryUserByUsername(in.Username)
	if err != nil {
		status = e.ErrorUsernameIncorrect
		return pb.LoginResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}

	// 验证密码
	if !util.VerifyPassByMD5(in.Password, dbUser.Password) {
		status = e.ErrorPasswordIncorrect
		return pb.LoginResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}

	// 生成token，并在redis中缓存用户信息
	token := util.GenerateToken()
	go s.userCache.SetUserInfo(token, &model.User{
		Id:             dbUser.Id,
		Username:       in.Username,
		Password:       dbUser.Password,
		ProfilePicture: dbUser.ProfilePicture,
	})

	logging.Log.Infof("[user login] userId: %d, username: %s", dbUser.Id, in.Username)
	return pb.LoginResponse{
		Code:  int32(status),
		Msg:   e.GetMsg(status),
		Token: token,
		User: &pb.User{
			Id:             dbUser.Id,
			Username:       in.Username,
			ProfilePicture: dbUser.ProfilePicture,
		},
	}, nil
}

// CheckToken 检查token
func (s *UserService) CheckToken(in pb.CheckTokenRequest) (pb.CheckTokenResponse, error) {
	var status = e.SUCCESS
	id, err := s.userCache.GetUserId(in.Token)
	if err != nil {
		status = e.ErrorTokenExpired
		return pb.CheckTokenResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	logging.Log.Infof("[check token] token：%s，id：%d", in.Token, id)
	return pb.CheckTokenResponse{
		Code: int32(status),
		Msg:  e.GetMsg(status),
	}, nil
}

// GetProfile 获取信息
// TODO 多节点，分布式锁
func (s *UserService) GetProfile(in pb.GetProfileRequest) (pb.GetProfileResponse, error) {
	var status = e.SUCCESS

	// 从缓存中获取用户信息
	caUser, err := s.userCache.GetUserInfo(in.Token)
	if err != nil {
		status = e.ErrorTokenExpired
		return pb.GetProfileResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
			User: pb.User{},
		}, nil
	}

	// 用户信息失效，说明已更新
	// 从数据库获取，添加到缓存
	if caUser.Username == "" || caUser.ProfilePicture == "" {
		dbUser, err := s.userMapper.QueryUserById(caUser.Id)
		if err != nil {
			status = e.ErrorOperateDatabase
			return pb.GetProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
				User: pb.User{},
			}, nil
		}
		caUser.Username = dbUser.Username
		caUser.ProfilePicture = dbUser.ProfilePicture
		s.userCache.SetUserField(in.Token, "username", caUser.Username)
		s.userCache.SetUserField(in.Token, "profile_picture", caUser.ProfilePicture)
	}

	return pb.GetProfileResponse{
		Code: int32(status),
		Msg:  e.GetMsg(status),
		User: pb.User{
			Id:             caUser.Id,
			Username:       caUser.Username,
			ProfilePicture: caUser.ProfilePicture,
		},
	}, nil
}

// UpdateProfile 更新信息
// 更新字段，延时2S双删
func (s *UserService) UpdateProfile(in pb.UpdateProfileRequest) (pb.UpdateProfileResponse, error) {
	var status = e.SUCCESS

	// 从缓存中获取用户id
	caUser, err := s.userCache.GetUserInfo(in.Token)
	if err != nil {
		status = e.ErrorTokenExpired
		return pb.UpdateProfileResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	// 从数据库查询最新信息
	dbUser, err := s.userMapper.QueryUserById(caUser.Id)
	if err != nil {
		status = e.ErrorOperateDatabase
		return pb.UpdateProfileResponse{
			Code: int32(status),
			Msg:  e.GetMsg(status),
		}, nil
	}
	caUser.Username = dbUser.Username
	caUser.ProfilePicture = dbUser.ProfilePicture

	// 用户名不为空
	if in.Username != "" {
		// 检查用户名合法性
		if status = util.CheckUsername(in.Username); status != e.SUCCESS {
			return pb.UpdateProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
			}, nil
		}
		// 检查用户名是否变动
		if in.Username != caUser.Username {
			// 检查是否有其他人已使用
			exist, err := s.userMapper.CheckUserUsernameExist(in.Username)
			if err != nil {
				status = e.ErrorOperateDatabase
				return pb.UpdateProfileResponse{
					Code: int32(status),
					Msg:  e.GetMsg(status),
				}, nil
			}
			if exist {
				status = e.ErrorUsernameAlreadyExist
				return pb.UpdateProfileResponse{
					Code: int32(status),
					Msg:  e.GetMsg(status),
				}, nil
			}

			s.userCache.DelUserField(in.Token, "username")
			defer func() {
				go func() {
					time.Sleep(2 * time.Second)
					s.userCache.DelUserField(in.Token, "username")
				}()
			}()
			err = s.userMapper.UpdateUserUsernameById(caUser.Id, in.Username)
			if err != nil {
				status = e.ErrorOperateDatabase
				return pb.UpdateProfileResponse{
					Code: int32(status),
					Msg:  e.GetMsg(status),
				}, nil
			}
			logging.Log.Infof("[user update] userId：%d，username：%s", caUser.Id, in.Username)
		}
	}

	// 用户头像不为空，说明已更新
	if in.ProfilePicture != "" {
		s.userCache.DelUserField(in.Token, "username")
		defer func() {
			go func() {
				time.Sleep(2 * time.Second)
				s.userCache.DelUserField(in.Token, "profile_picture")
			}()
		}()
		err = s.userMapper.UpdateUserProfilePictureById(caUser.Id, in.ProfilePicture)
		if err != nil {
			status = e.ErrorOperateDatabase
			return pb.UpdateProfileResponse{
				Code: int32(status),
				Msg:  e.GetMsg(status),
			}, nil
		}
		logging.Log.Infof("[user update] userId：%d，avatar：%s", caUser.Id, in.ProfilePicture)
	}

	return pb.UpdateProfileResponse{
		Code: int32(status),
		Msg:  e.GetMsg(status),
	}, nil
}

// Logout 退出登录
func (s *UserService) Logout(in pb.LogoutRequest) (pb.LogoutResponse, error) {
	logging.Log.Infof("[user logout] token：%s", in.Token)
	s.userCache.DelUserInfo(in.Token)
	return pb.LogoutResponse{
		Code: e.SUCCESS,
		Msg:  e.GetMsg(e.SUCCESS),
	}, nil
}
