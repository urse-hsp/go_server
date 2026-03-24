package main

import (
	"go-demo-server/pkg/scheduler"
	"go-demo-server/task"
	"log"
)

func main() {
	log.Println("🚀 启动定时任务服务")

	// 1. 初始化 scheduler
	s := scheduler.NewScheduler()

	// 2. 注册任务
	task.RegisterTasks(s)

	// 3. 启动
	s.Start()

	// 4. 阻塞
	select {}
}
