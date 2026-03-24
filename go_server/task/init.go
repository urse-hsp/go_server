package task

import (
	"go-demo-server/pkg/scheduler"

	"github.com/redis/go-redis/v9"
)

// 注册任务
func RegisterTasks(s *scheduler.Scheduler) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	// locker := scheduler.NewMemoryLocker() // 内容锁
	locker := scheduler.NewRedisLocker(rdb) // redis锁

	// demo()
	// RegisterDemoTask(s)
	RegisterDemoTask2(s, locker)
}
