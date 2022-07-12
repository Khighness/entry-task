package cache

import (
	"time"

	"github.com/go-redis/redis"

	"github.com/Khighness/entry-task/tcp/config"
	"github.com/Khighness/entry-task/tcp/logging"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

var (
	RedisClient *redis.Client
)

// InitRedis 初始化Redis连接池
func InitRedis() {
	RedisClient = ConnectRedis(config.AppCfg.Redis)
}

// ConnectRedis 连接到Redis
func ConnectRedis(redisCfg *config.RedisConfig) *redis.Client {
	options := &redis.Options{
		Addr:         redisCfg.Addr,
		Password:     redisCfg.Pass,
		DB:           redisCfg.Db,
		PoolSize:     redisCfg.MaxConn,
		MinIdleConns: redisCfg.MinIdle,
		MaxRetries:   redisCfg.MaxRetries,
		DialTimeout:  time.Duration(redisCfg.DialTimeout) * time.Second,
		IdleTimeout:  time.Duration(redisCfg.IdleTimeout) * time.Second,
		MaxConnAge:   time.Duration(redisCfg.MaxConnAge) * time.Second,
	}
	redisClient := redis.NewClient(options)

	if _, err := redisClient.Ping().Result(); err != nil {
		logging.Log.Fatalf("Failed to connect to redis server [%s]: %s", redisCfg.Addr, err)
	} else {
		logging.Log.Printf("Succeed to connect to redis server [%s]", redisCfg.Addr)
	}
	return redisClient
}
