package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func (c *Config) InitRedisConnection() *Config {
	redisSingleton.Do(func() {
		log.Println("Trying to open redis connection pool . . . .")
		opts, err := redis.ParseURL(c.RedisDSN)
		if err != nil {
			log.Fatalf("REDIS_ERROR: %s", err.Error())
		}
		RedisPool = redis.NewClient(opts)
		if err := RedisPool.Ping(context.Background()).Err(); err != nil {
			log.Fatalf("REDIS_ERROR: %s", err.Error())
		}
		log.Println("Redis connection pool created . . . .")
	})
	return c
}
