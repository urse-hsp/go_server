// migration 执行器
// 每次启动都会执行所有 SQL
package main

import (
	"go-server/internal/bootstrap"
	"go-server/pkg/config"
	"go-server/pkg/log"
	"os"

	"go.uber.org/zap"
)

// 迁移数据库
func main() {
	conf := config.NewConfig("config/local.yaml")
	logger := log.NewLog(conf)

	db := bootstrap.NewDB(conf, logger)
	m := bootstrap.NewMigrateServer(db, logger)

	if err := m.Start(); err != nil {
		logger.Error("migration failed", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("migration done")

}
