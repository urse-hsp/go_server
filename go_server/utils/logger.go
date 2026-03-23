package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

// InitLogger 初始化日志库
// 在实际项目中，你可以从 Viper 读取配置来决定日志级别和输出文件
func InitLogger() {
	// 这里演示一个生产环境的配置示例
	config := zap.NewProductionConfig()

	// 自定义配置
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel) // 设置级别
	config.Encoding = "json"                               // 输出 JSON 格式
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式
	config.EncoderConfig.CallerKey = "file"
	config.EncoderConfig.MessageKey = "msg"

	// 构建 Logger
	logger, err := config.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		panic(err)
	}

	// 转换为 SugaredLogger 以便使用
	Logger = logger.Sugar()
}

// Sync 确保所有缓冲的日志都被写入
func Sync() {
	Logger.Sync()
}
