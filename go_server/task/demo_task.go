package task

import (
	"context"
	"fmt"
	"go-demo-server/pkg/scheduler"
	"time"
)

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
