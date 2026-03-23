package main

import (
	"go-demo-server/bootstrap"
	"go-demo-server/config"
	"go-demo-server/router"
	"go-demo-server/utils"
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

	// 启动服务
	r := router.SetupRouter() // gin router 应用实例
	r.Run(":" + config.Conf.Server.Port)

	//  初始化日志
	utils.InitLogger()
	defer utils.Sync() // 程序退出前确保日志写完
	utils.Logger.Info("服务启动成功", "port", 8080)
}
