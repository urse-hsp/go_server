package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// 定义配置结构体，方便后续使用
type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpireTime int    `mapstructure:"expire_time"` // 单位：小时 (或者秒)
}

type LogConfig struct {
	LogLevel     string `mapstructure:"log_level"`
	Encoding     string `mapstructure:"encoding"`
	LogFileName  string `mapstructure:"log_file_name"`
	MaxBackups   int    `mapstructure:"max_backups"`
	MaxAge       int    `mapstructure:"max_age"`
	MaxSize      int    `mapstructure:"max_size"`
	LogFileLevel string `mapstructure:"log_file_level"`
	Compress     bool   `mapstructure:"compress"`
}

var Conf *AppConfig

// InitConfig 初始化配置
func InitConfig() error {
	v := viper.New()

	// 1. 设置配置文件名（不含扩展名）
	v.SetConfigName("local")
	// 2. 设置配置文件类型
	v.SetConfigType("yaml")
	// 3. 设置查找路径
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	// 4. 尝试读取配置
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 5. 将配置反序列化到结构体中
	if err := v.Unmarshal(&Conf); err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}

	fmt.Println("✅ viper 配置加载成功！")
	return nil
}
