package database

import (
	"github.com/Skyvko6607/go-api-example/config"
	"github.com/redis/go-redis/v9"
)

type RedisContext struct {
	Client *redis.Client
}

func NewRedisContext(config *config.AppSettings) RedisContext {
	redisConfig := config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
		Protocol: redisConfig.Protocol,
	})
	return RedisContext{Client: rdb}
}
