package bootstrap

import (
	"go-demo-server/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.SugaredLogger

func InitLogger() {
	LOGConfig := config.Conf.Log

	// 1. 配置日志写入位置 (使用 lumberjack 实现自动切割)
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   LOGConfig.LogFileName, // 日志文件路径 (确保 logs 目录存在)
		MaxSize:    LOGConfig.MaxSize,     // 每个日志文件最大 100MB
		MaxBackups: LOGConfig.MaxBackups,  // 最多保留 5 个旧文件
		MaxAge:     LOGConfig.MaxAge,      // 文件保留 30 天
		Compress:   LOGConfig.Compress,    // 压缩旧文件
	})

	// 2. 配置编码器 (JSON 格式)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	// 3. 配置日志级别
	level := zapcore.InfoLevel

	// 4. 构建核心配置
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 5. 构建 Logger
	// zap.AddCaller() 用于显示文件名和行号
	logger := zap.New(core, zap.AddCaller())

	Logger = logger.Sugar()

	// 注意：不要在这里调用 Sync()，应该在程序退出时调用
	Logger.Info("日志系统初始化成功，日志将写入 logs/app.log")
}
