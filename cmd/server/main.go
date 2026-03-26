package main

import (
	"go-demo-server/config"
	"go-demo-server/internal/bootstrap"
	"go-demo-server/internal/router"
)

// @title go-demo-server API
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
	bootstrap.InitLogger()

	// 启动服务
	r := router.SetupRouter() // gin router 应用实例
	r.Run(":" + config.Conf.Server.Port)
}
