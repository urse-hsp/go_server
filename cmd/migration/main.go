// migration 执行器
// 每次启动都会执行所有 SQL
package migration

import (
	"fmt"
	"go-server/internal/bootstrap"
	"go-server/pkg/config"
	"go-server/pkg/log"
	"io/ioutil"
	"path/filepath"
)

// 迁移数据库
func RunMigrations() {
	files, err := filepath.Glob("./*.sql")
	if err != nil {
		panic(err)
	}

	envConf := "config/local.yaml"
	conf := config.NewConfig(envConf)
	logger := log.NewLog(conf)          // 初始化日志
	DB := bootstrap.NewDB(conf, logger) // 初始化 MySQL

	for _, file := range files {
		fmt.Println("执行 migration:", file)

		sqlBytes, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		err = DB.Exec(string(sqlBytes)).Error
		if err != nil {
			panic("执行失败: " + file + " | " + err.Error())
		}
	}
}
