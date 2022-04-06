package cache

import (
	"entry/tcp/logging"
	"github.com/go-redis/redis"
	"gopkg.in/ini.v1"
	"time"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-15

// RedisClient
var (
	RedisClient *redis.Client
	RedisAddr   string
	RedisAuth   string
	RedisDb     int
)

// Load 初始化Redis
func Load(file *ini.File) {
	loadRedisConfig(file)
	connectRedis()
}

// loadRedisConfig 读取Redis配置
func loadRedisConfig(file *ini.File) {
	redis := file.Section("redis")
	RedisAddr = redis.Key("RedisAddr").String()
	RedisAuth = redis.Key("RedisAuth").String()
	Db, err := redis.Key("RedisDb").Int()
	if err != nil {
		logging.Log.Fatalf("Wrong configuration of [RedisDb] in config file: %s", err)
	} else {
		RedisDb = Db
	}
}

// connectRedis 连接到Redis
func connectRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:         RedisAddr,
		Password:     RedisAuth,
		DB:           RedisDb,
		PoolSize:     30000,
		MinIdleConns: 1000,
		IdleTimeout:  time.Hour,
	})

	if _, err := client.Ping().Result(); err != nil {
		logging.Log.Fatalf("Failed to connect to redis server [%s]: %s", RedisAddr, err)
	} else {
		logging.Log.Printf("Succeed to connect to redis server [%s]", RedisAddr)
	}

	RedisClient = client
}
