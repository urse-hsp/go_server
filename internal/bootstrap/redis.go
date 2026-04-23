package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// var RDB *redis.Client

func NewRedis(conf *viper.Viper) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.GetString("data.redis.addr"),
		Password: conf.GetString("data.redis.password"),
		DB:       conf.GetInt("data.redis.db"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	}

	println("✅ Redis 连接成功")

	return rdb
}

type RDBCache struct {
	rdb *redis.Client
}

func NewRDBCache(rdb *redis.Client) *RDBCache {
	return &RDBCache{rdb: rdb}
}

func (c *RDBCache) SetJSON(ctx context.Context, key string, val any, ttl time.Duration) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, key, data, ttl).Err()
}

func GetJSON[T any](c *RDBCache, ctx context.Context, key string) (T, error) {
	var res T

	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return res, err
	}

	err = json.Unmarshal([]byte(val), &res)
	return res, err
}
