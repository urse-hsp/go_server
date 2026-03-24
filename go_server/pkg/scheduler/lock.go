// lock.go（分布式锁，预留）
// 既可以是单机锁，也可以是分布式锁

package scheduler

import (
	"context"
	"log"
	"time"
)

type Locker interface {
	Lock(ctx context.Context, key string, ttl time.Duration) (bool, error)
	Unlock(ctx context.Context, key string) error
}

// 包装 Handler
func WithLock(locker Locker, key string, ttl time.Duration, fn func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		ok, err := locker.Lock(ctx, key, ttl)
		if err != nil || !ok {
			log.Println("[Job] 获取锁失败，跳过")
			return nil
		}
		defer locker.Unlock(ctx, key)

		return fn(ctx)
	}
}
