package dao

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

var RC *RedisCache

func init() {
	opt, _ := redis.ParseURL("redis://default:123456@127.0.0.1:6380/0")
	rdb := redis.NewClient(opt)
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	RC = &RedisCache{
		rdb: rdb,
	}
}


func (c *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	return c.rdb.Set(ctx, key, value, expire).Err()
}


func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := c.rdb.Get(ctx, key).Result()
	return result, err
}