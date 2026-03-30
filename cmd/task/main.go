package main

import (
	"go-server/config"
	"go-server/internal/task"
	"go-server/pkg/server/scheduler"
	"log"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

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
