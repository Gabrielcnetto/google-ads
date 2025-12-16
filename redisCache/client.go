package rediscache

import "github.com/redis/go-redis/v9"

func RedisConnection() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: "localhost:6379"})
}
