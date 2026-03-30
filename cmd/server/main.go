package main

import (
	"go-server/config"
	"go-server/internal/bootstrap"
	"go-server/internal/router"
)

// @title go-server API
// @version 1.0
// @description 接口文档
// @host localhost:8080
// @BasePath /
func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		panic(err)
	}
	// 初始化各种组件
	bootstrap.InitMysql()
	bootstrap.InitRedis()

	//  初始化日志
	bootstrap.InitLogger(config.Conf.Log)

	// 启动服务
	r := router.SetupRouter() // gin router 应用实例
	r.Run(":" + config.Conf.Server.Port)
}
