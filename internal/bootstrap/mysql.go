package bootstrap

import (
	"fmt"
	"go-server/config"
	"go-server/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 声明全局变量
var DB *gorm.DB

// 初始化函数
func InitMysql() {
	dbConfig := config.Conf.Database

	// dsn := "root:123456@tcp(127.0.0.1:3306)/go_demo"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("❌ 数据库连接失败: " + dsn)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("❌ 获取 DB 实例失败: " + err.Error())
	}

	if err := sqlDB.Ping(); err != nil {
		panic("❌ 数据库连接失败（Ping）: " + err.Error())
	}

	fmt.Println("✅ MySQL 连接成功")

	DB = db

	autoMigrateErr := DB.AutoMigrate(model.GetModels()...)
	if autoMigrateErr != nil {
		panic("❌ 数据库迁移失败: " + autoMigrateErr.Error())
	}
}
