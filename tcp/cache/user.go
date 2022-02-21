package cache

import (
	"entry/tcp/model"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// UserTokenKeyPrefix 用户token存储的key
// UserTokenTimeout   用户token过期时间
const (
	UserTokenKeyPrefix = "entry:user:session:"
	UserTokenTimeout   = time.Hour * 24
)

// GetUserId 获取用户id
func GetUserId(sessionId string) (int64, error) {
	return RedisClient.HGet(generateUserTokenKey(sessionId), "id").Int64()
}

// GetUserInfo 获取用户信息
func GetUserInfo(sessionId string) (*model.User, error) {
	userTokenKey := generateUserTokenKey(sessionId)
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
func SetUserInfo(sessionId string, user *model.User) {
	userTokenKey := generateUserTokenKey(sessionId)
	RedisClient.HSet(userTokenKey, "id", user.Id)
	RedisClient.HSet(userTokenKey, "username", user.Username)
	RedisClient.HSet(userTokenKey, "profile_picture", user.ProfilePicture)
	RedisClient.Expire(userTokenKey, UserTokenTimeout)
}

// DelUserInfo 删除用户信息
func DelUserInfo(sessionId string) {
	RedisClient.Del(generateUserTokenKey(sessionId))
}

// SetUserField 缓存用户字段信息
func SetUserField(sessionId, key, val string) {
	RedisClient.HSet(generateUserTokenKey(sessionId), key, val)
}

// DelUserField 删除用户字段信息
func DelUserField(sessionId, key string) {
	RedisClient.HDel(generateUserTokenKey(sessionId), key)
}

// generateUserTokenKey 生成用户token的key
func generateUserTokenKey(sessionId string) string {
	return UserTokenKeyPrefix + sessionId
}
