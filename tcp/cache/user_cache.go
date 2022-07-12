package cache

import (
	"time"

	"github.com/Khighness/entry-task/tcp/model"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// UserTokenKeyPrefix  用户token存储的key
// UserTokenTimeout   用户token过期时间
const (
	UserTokenKeyPrefix = "entry:user:token:"
	UserTokenTimeout   = time.Hour * 24
)

// UserCache 用户缓存操作
type UserCache struct {
}

// GetUserId 获取用户id
func (userCache *UserCache) GetUserId(token string) (int64, error) {
	return RedisClient.HGet(userCache.generateUserTokenKey(token), "id").Int64()
}

// GetUserInfo 获取用户信息
func (userCache *UserCache) GetUserInfo(token string) (*model.User, error) {
	userTokenKey := userCache.generateUserTokenKey(token)
	id, err := RedisClient.HGet(userTokenKey, "id").Int64()
	if err != nil {
		return nil, err
	}

	username := RedisClient.HGet(userTokenKey, "username").Val()
	profilePicture := RedisClient.HGet(userTokenKey, "profile_picture").Val()
	return &model.User{
		Id:             id,
		Username:       username,
		Password:       "",
		ProfilePicture: profilePicture,
	}, nil
}

// SetUserInfo 缓存用户信息
func (userCache *UserCache) SetUserInfo(token string, user *model.User) {
	userTokenKey := userCache.generateUserTokenKey(token)
	RedisClient.HSet(userTokenKey, "id", user.Id)
	RedisClient.HSet(userTokenKey, "username", user.Username)
	RedisClient.HSet(userTokenKey, "profile_picture", user.ProfilePicture)
	RedisClient.Expire(userTokenKey, UserTokenTimeout)
}

// DelUserInfo 删除用户信息
func (userCache *UserCache) DelUserInfo(token string) {
	RedisClient.Del(userCache.generateUserTokenKey(token))
}

// SetUserField 缓存用户字段信息
func (userCache *UserCache) SetUserField(token, key, val string) {
	RedisClient.HSet(userCache.generateUserTokenKey(token), key, val)
}

// DelUserField 删除用户字段信息
func (userCache *UserCache) DelUserField(token, key string) {
	RedisClient.HDel(userCache.generateUserTokenKey(token), key)
}

// generateUserTokenKey 生成用户token的key
func (userCache *UserCache) generateUserTokenKey(token string) string {
	return UserTokenKeyPrefix + token
}
