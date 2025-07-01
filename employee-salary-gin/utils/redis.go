package utils

import (
	"context"
	"employee-salary-gin/config"
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	cfg := config.LoadConfig()
	opt, _ := redis.ParseURL("redis://" + cfg.RedisAddr + "/0")
	RedisClient = redis.NewClient(opt)

	// 测试连接
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
}
