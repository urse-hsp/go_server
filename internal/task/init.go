package task

import (
	"go-server/internal/bootstrap"
	"go-server/pkg/server/scheduler"
)

// 注册任务
func RegisterTasks(s *scheduler.Scheduler) {
	rdb := bootstrap.InitRedis()
	// locker := scheduler.NewMemoryLocker() // 内容锁
	locker := scheduler.NewRedisLocker(rdb) // redis锁

	// demo()
	// RegisterDemoTask(s)
	RegisterDemoTask2(s, locker)
}
