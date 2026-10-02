// Package logger 提供基于 zap 的日志初始化（构造注入，不依赖 global）。
package logger

import (
	"log"
	"os"

	"server/config"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init 根据配置初始化并返回一个新的 zap.Logger 实例。
func Init(cfg *config.Zap) *zap.Logger {
	// 创建用于日志输出的 writeSyncer
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)

	// 如果配置了控制台输出，则添加控制台输出
	if cfg.IsConsolePrint {
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
	}

	// 创建日志格式化的编码器
	encoder := getEncoder()

	// 根据配置确定日志级别
	var logLevel zapcore.Level
	if err := logLevel.UnmarshalText([]byte(cfg.Level)); err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}

	// 创建核心和日志实例
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	return zap.New(core, zap.AddCaller())
}

// getLogWriter 返回一个 zapcore.WriteSyncer，利用 lumberjack 实现日志滚动记录。
func getLogWriter(filename string, maxSize, maxBackups, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,   // 日志文件的位置
		MaxSize:    maxSize,    // 切割前日志文件最大大小（MB）
		MaxBackups: maxBackups, // 保留旧文件最大个数
		MaxAge:     maxAge,     // 保留旧文件最大天数
	}
	return zapcore.AddSync(lumberJackLogger)
}

// getEncoder 返回一个为生产日志配置的 JSON 编码器。
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
