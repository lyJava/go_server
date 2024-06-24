package config

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// init 初始化日志库
func init() {

	// 配置文件日志的编码器（无颜色）
	fileConfig := zap.NewProductionEncoderConfig()
	fileConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	fileConfig.EncodeTime = parseTime("2006-01-02 15:04:05.000")
	fileConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// 配置控制台日志的编码器（有颜色）
	consoleConfig := zap.NewDevelopmentEncoderConfig()
	consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleConfig.EncodeTime = parseTime("2006-01-02 15:04:05.000")
	consoleConfig.EncodeCaller = zapcore.ShortCallerEncoder
	
	// 日志输出配置, 借助另外一个库 lumberjack 协助完成日志切割。
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "log/express-api.log", // 日志文件名
		MaxSize:    5,                     // 最大日志数 M为单位!!!
		MaxAge:     30,                    // 最大存在天数
		MaxBackups: 30,                    // 最大备份数量
		Compress:   false,                 // 是否压缩
	}
	fileSync := zapcore.AddSync(lumberjackLogger)
	consoleSync := zapcore.Lock(os.Stdout)

	var core zapcore.Core
	fileEncoder := zapcore.NewConsoleEncoder(fileConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
	core = zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileSync, zapcore.DebugLevel),
		zapcore.NewCore(consoleEncoder, consoleSync, zapcore.DebugLevel),
	)
	logger := zap.New(core, zap.AddCaller()) // --添加函数调用信息
	zap.ReplaceGlobals(logger)               // 替换该日志为全局日志
	defer logger.Sync()
}

// parseTime 进行时间格式处理
func parseTime(layout string) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(layout))
	}
}
