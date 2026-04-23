package main

import (
	"flag"
	"go-server/internal/bootstrap"
	"go-server/internal/controller"
	"go-server/internal/dao"
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
	// envConf := "config/local.yaml"
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	conf := config.NewConfig(*envConf)

	// 初始化组件
	logger := log.NewLog(conf)             // 初始化日志
	DB := bootstrap.NewDB(conf, logger)    // 初始化 MySQL
	RDB := bootstrap.NewRedis(conf)        // 初始化 Redis
	RDBCache := bootstrap.NewRDBCache(RDB) // 初始化 Redis

	repositoryRepository := dao.NewRepository(logger, DB)   // 初始化 Repository/dao，注入 Logger,DB,RDB
	transaction := dao.NewTransaction(repositoryRepository) // 初始化 Transaction，注入 Repository/dao

	sidSid := sid.NewSid()
	jwtJWT := jwt.NewJwt(conf)
	if jwtJWT == nil {
		panic("jwt init failed")
	}

	serviceService := service.NewService(transaction, logger, sidSid, jwtJWT, RDBCache) // 初始化 Service，注入 Transaction、Logger、Sid 和 JWT
	Handler := controller.NewHandler(logger)                                            // 初始化 Handler

	routerDeps := router.RouterDeps{
		Logger:     logger,
		Config:     conf,
		JWT:        jwtJWT,
		Repository: repositoryRepository,
		Service:    serviceService,
		Handler:    Handler,
	}

	// 启动服务
	r := router.SetupRouter(routerDeps) // gin router 应用实例
	r.Run(":" + conf.GetString("server.port"))
}
