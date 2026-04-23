package main

import (
	"flag"
	"go-server/internal/bootstrap"
	"go-server/internal/dao"
	"go-server/internal/task"
	"go-server/pkg/config"
	"go-server/pkg/log"
	"go-server/pkg/server/scheduler"
	"go-server/pkg/sid"
)

func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	conf := config.NewConfig(*envConf)

	// 初始化组件
	logger := log.NewLog(conf)          // 初始化日志
	DB := bootstrap.NewDB(conf, logger) // 初始化 MySQL
	// RDB := bootstrap.NewRedis(conf)     // 初始化 Redis
	// RDBCache := bootstrap.NewRDBCache(RDB) // 初始化 Redis

	repositoryRepository := dao.NewRepository(logger, DB)   // 初始化 Repository/dao，注入 Logger,DB,RDB
	transaction := dao.NewTransaction(repositoryRepository) // 初始化 Transaction，注入 Repository/dao
	sidSid := sid.NewSid()

	taskTask := task.NewTask(transaction, logger, sidSid)

	// 1. 初始化 scheduler
	s := scheduler.NewScheduler()

	// 2. 注册任务
	bootstrap.RegisterTasks(s, taskTask, repositoryRepository)

	// 3. 启动
	s.Start()

	// 4. 阻塞
	select {}
}
