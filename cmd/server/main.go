package main

import (
	"go-server/internal/bootstrap"
	"go-server/internal/router"
	"go-server/internal/service"
	"go-server/pkg/config"
	"go-server/pkg/jwt"
	"go-server/pkg/log"
	"go-server/pkg/sid"
)

// @title go-server API
// @version 1.0
// @description 接口文档
// @host localhost:8080
// @BasePath /
func main() {
	envConf := "config/local.yaml"
	conf := config.NewConfig(envConf)

	// 初始化组件
	logger := log.NewLog(conf)          // 初始化日志
	DB := bootstrap.NewDB(conf, logger) // 初始化 MySQL
	RDB := bootstrap.NewRedis(conf)     // 初始化 Redis

	repositoryRepository := bootstrap.NewRepository(logger, DB, RDB) // 初始化 Repository/dao，注入 Logger,DB,RDB
	transaction := bootstrap.NewTransaction(repositoryRepository)    // 初始化 Transaction，注入 Repository/dao

	sidSid := sid.NewSid()
	jwtJWT := jwt.NewJwt(conf)
	if jwtJWT == nil {
		panic("jwt init failed")
	}

	serviceService := service.NewService(transaction, logger, sidSid, jwtJWT) // 初始化 Service，注入 Transaction、Logger、Sid 和 JWT

	routerDeps := router.RouterDeps{
		Logger:     logger,
		Config:     conf,
		Repository: repositoryRepository,
		Service:    serviceService,
		JWT:        jwtJWT,
	}

	// 启动服务
	r := router.SetupRouter(routerDeps) // gin router 应用实例
	r.Run(":" + conf.GetString("server.port"))
}
