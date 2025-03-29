package datasources

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	Context   context.Context
	RedisWR   *redis.Client
	RedisRead *redis.Client
}

func NewRedisConnection() *RedisConnection {

	opt, _ := redis.ParseURL(os.Getenv("REDIS_URI"))
	rdb := redis.NewClient(opt)

	optRead, _ := redis.ParseURL(os.Getenv("REDISREAD_URI"))
	rdbRead := redis.NewClient(optRead)

	return &RedisConnection{
		Context:   context.Background(),
		RedisWR:   rdb,
		RedisRead: rdbRead,
	}
}
