package service

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func LoggerInit() (err error) {
	config := zap.NewDevelopmentConfig()                    // 获取生产环境的默认配置
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel) // 设置日志级别为Debug
	config.Encoding = "json"                                // 设置输出格式为JSON
	config.OutputPaths = []string{"stdout", "mygo5.log"}    // 设置输出路径，同时输出到标准输出和文件
	config.ErrorOutputPaths = []string{"stderr"}            // 设置错误输出路径到标准错误输出

	Logger, err = config.Build() // 根据配置构建Logger
	return
}

func SyncLogger() {
	Logger.Sync()
}
