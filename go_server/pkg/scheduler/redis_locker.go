// redis_locker.go
// redis 锁[分布式锁,集群版]
package scheduler

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLocker struct {
	client *redis.Client
}

func NewRedisLocker(client *redis.Client) *RedisLocker {
	return &RedisLocker{client: client}
}

func (r *RedisLocker) Lock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, "1", ttl).Result()
}

func (r *RedisLocker) Unlock(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
