package cache

import (
	"time"

	"github.com/go-redis/redis"

	"github.com/Khighness/entry-task/tcp/model"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

const (
	// UserTokenKeyPrefix  用户token存储的key
	UserTokenKeyPrefix = "entry:user:token:"
	// UserTokenTimeout   用户token过期时间
	UserTokenTimeout = time.Hour * 24
)

// UserCache 用户缓存操作
type UserCache struct {
	client *redis.Client
}

// NewUserCache 创建 UserCache
func NewUserCache(client *redis.Client) *UserCache {
	return &UserCache{client: client}
}

// GetUserId 获取用户id
func (userCache *UserCache) GetUserId(token string) (int64, error) {
	return userCache.client.HGet(userCache.generateUserTokenKey(token), "id").Int64()
}

// GetUserInfo 获取用户信息
func (userCache *UserCache) GetUserInfo(token string) (*model.User, error) {
	userTokenKey := userCache.generateUserTokenKey(token)
	id, err := userCache.client.HGet(userTokenKey, "id").Int64()
	if err != nil {
		return nil, err
	}

	username := userCache.client.HGet(userTokenKey, "username").Val()
	profilePicture := userCache.client.HGet(userTokenKey, "profile_picture").Val()
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
	userCache.client.HSet(userTokenKey, "id", user.Id)
	userCache.client.HSet(userTokenKey, "username", user.Username)
	userCache.client.HSet(userTokenKey, "profile_picture", user.ProfilePicture)
	userCache.client.Expire(userTokenKey, UserTokenTimeout)
}

// DelUserInfo 删除用户信息
func (userCache *UserCache) DelUserInfo(token string) {
	userCache.client.Del(userCache.generateUserTokenKey(token))
}

// SetUserField 缓存用户字段信息
func (userCache *UserCache) SetUserField(token, key, val string) {
	userCache.client.HSet(userCache.generateUserTokenKey(token), key, val)
}

// DelUserField 删除用户字段信息
func (userCache *UserCache) DelUserField(token, key string) {
	userCache.client.HDel(userCache.generateUserTokenKey(token), key)
}

// generateUserTokenKey 生成用户token的key
func (userCache *UserCache) generateUserTokenKey(token string) string {
	return UserTokenKeyPrefix + token
}
