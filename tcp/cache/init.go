package cache

import (
	"github.com/go-redis/redis"
	"gopkg.in/ini.v1"
	"log"
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

// 读取Redis配置
func loadRedisConfig(file *ini.File) {
	redis := file.Section("redis")
	RedisAddr = redis.Key("RedisAddr").String()
	RedisAuth = redis.Key("RedisAuth").String()
	Db, err := redis.Key("RedisDb").Int()
	if err != nil {
		log.Fatalf("Wrong configuration of [RedisDb] in config file: %s\n", err)
	} else {
		RedisDb = Db
	}
}

// 连接到Redis
func connectRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:         RedisAddr,
		Password:     RedisAuth,
		DB:           RedisDb,
		PoolSize:     20,
		MinIdleConns: 1,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Fatalf("Failed to connect to redis server [%s]: %s\n", RedisAddr, err)
	} else {
		log.Printf("Succeed to connect to redis server [%s]\n", RedisAddr)
	}

	RedisClient = client
}
