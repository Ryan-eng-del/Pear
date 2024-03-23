package dao

import (
	"context"
	"log"
	"time"

	"cyan.com/pear-user/config"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rdb *redis.Client
}

var RC *RedisCache

func init() {
	rdb := redis.NewClient(config.C.InitRedis())
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