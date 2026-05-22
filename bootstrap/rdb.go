package bootstrap

import (
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(env *Env) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         env.RedisAddr,
		Password:     env.RedisPass, // 没有密码，默认值
		DB:           0,             // 默认DB 0
		PoolFIFO:     false,
		PoolSize:     200,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		MaxRetries:   3,
		MinIdleConns: 30,
		MaxIdleConns: 200,
	})
	return rdb
}
