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

// 没有文件输出的日志
// package bootstrap

// import (
// 	"go.uber.org/zap"
// 	"go.uber.org/zap/zapcore"
// )

// var Logger *zap.SugaredLogger

// // InitLogger 初始化日志库
// // 在实际项目中，你可以从 Viper 读取配置来决定日志级别和输出文件
// func InitLogger() {
// 	// 这里演示一个生产环境的配置示例
// 	config := zap.NewProductionConfig()

// 	// 自定义配置
// 	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel) // 设置级别
// 	config.Encoding = "json"                               // 输出 JSON 格式
// 	config.EncoderConfig.TimeKey = "time"
// 	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 时间格式
// 	config.EncoderConfig.CallerKey = "file"
// 	config.EncoderConfig.MessageKey = "msg"

// 	// 构建 Logger
// 	logger, err := config.Build(zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
// 	if err != nil {
// 		panic(err)
// 	}

// 	// 转换为 SugaredLogger 以便使用
// 	Logger = logger.Sugar()

// 	Sync()
// 	Logger.Info("服务启动成功", "port", 8080)
// }

// // Sync 确保所有缓冲的日志都被写入
// func Sync() {
// 	Logger.Sync()
// }
