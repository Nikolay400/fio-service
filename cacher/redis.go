package cacher

import (
	"fio-service/env"
	"time"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func NewRedis() *Redis {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     env.GetRedisUrl(),
		Password: "",
		DB:       0,
	})
	return &Redis{redisClient}
}

func (redis *Redis) Set(key string, value interface{}, t time.Duration) error {
	return redis.client.Set(key, value, t).Err()
}

func (redis *Redis) Del(key string) error {
	return redis.client.Del(key).Err()
}

func (redis *Redis) Get(key string) (string, error) {
	return redis.client.Get(key).Result()
}

func (redis *Redis) Close() {
	redis.client.Close()
}
