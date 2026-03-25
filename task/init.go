package task

import (
	"go-demo-server/bootstrap"
	"go-demo-server/pkg/scheduler"
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
