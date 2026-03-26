package task

import (
	"context"
	"fmt"
	"go-demo-server/pkg/server/scheduler"
	"time"
)

// 定时任务[无锁版]
// cron 到时间 → 一定执行。可能会并发执行 数据可能冲突
func RegisterDemoTask(s *scheduler.Scheduler) {
	job := scheduler.NewJob(
		"demo_task",
		func(ctx context.Context) error {
			fmt.Println("执行 demo 任务")
			return nil
		},
		scheduler.WithRetry(3, 2*time.Second),
		scheduler.WithTimeout(10*time.Second),
	)

	// 每5秒执行
	s.AddJob("*/5 * * * * *", job)
}

// 加锁版任务
// cron 到时间 → 先抢锁 → 成功才执行
// 到点尝试执行，抢不到锁就跳过”
// 执行时有锁的权限才能执行，如果被占用了，就不能执行
func RegisterDemoTask2(s *scheduler.Scheduler, locker scheduler.Locker) {
	handler := func(ctx context.Context) error {
		fmt.Println("执行任务逻辑")
		return nil
	}

	// 👉 用内存锁包一层
	handler = scheduler.WithLock(
		locker,
		"cron:lock:demo_task", // 唯一 key
		10*time.Second,        // 锁时间
		handler,
	)

	job := scheduler.NewJob(
		"demo_task",
		handler,
		scheduler.WithTimeout(5*time.Second),
	)

	s.AddJob("*/5 * * * * *", job)
}

// 分布式集群使用方式
func RedisRrunkedRegisterDemoTask(s *scheduler.Scheduler, locker scheduler.Locker) {
	job := scheduler.NewJob(
		"demo_task",
		scheduler.WithLock(
			locker,
			"lock:demo_task", // 🔥 全局唯一key
			10*time.Second,   // 🔥 锁过期时间
			func(ctx context.Context) error {
				fmt.Println("真正执行任务")
				return nil
			},
		),
		scheduler.WithTimeout(15*time.Second),
	)

	// 每5秒执行
	s.AddJob("*/5 * * * * *", job)
}
