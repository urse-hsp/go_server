package main

import (
	"go-demo-server/bootstrap"
	"go-demo-server/config"
	"go-demo-server/router"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	// 初始化各种组件
	bootstrap.InitMysql()
	bootstrap.InitRedis()

	//  初始化定时任务
	bootstrap.InitScheduler()
	//  初始化日志
	bootstrap.InitLogger()

	// 初始化 http client
	// task.InitTask()

	// 启动服务
	r := router.SetupRouter() // gin router 应用实例
	r.Run(":" + config.Conf.Server.Port)
}
