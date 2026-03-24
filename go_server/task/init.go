package task

import "go-demo-server/pkg/scheduler"

// 注册任务
func RegisterTasks(s *scheduler.Scheduler) {
	// demo()
	RegisterDemoTask(s)
}
