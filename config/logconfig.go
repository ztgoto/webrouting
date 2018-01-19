package config

import (
	"log"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// DefaultLogger 默认日志输出
	DefaultLogger *zap.Logger
)

var one sync.Once

// InitDefaultLog 初始化日志
func InitDefaultLog() {
	one.Do(initDefaultLogger)
}

// GetLogger 创建日志
func GetLogger(lv string, output ...string) (*zap.Logger, error) {
	logConf := &zap.Config{
		Level:       zap.NewAtomicLevelAt(unmarshalText(lv)),
		Encoding:    "json",
		OutputPaths: output,
		// ErrorOutputPaths:output,
	}
	logConf.EncoderConfig = zap.NewProductionEncoderConfig()
	logConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return logConf.Build()
}

func initDefaultLogger() {
	var logger *zap.Logger
	var e error
	var filePath = strings.TrimSpace(GlobalConfig.HTTP.LogPath)
	if len(filePath) == 0 {
		logger, e = GetLogger(GlobalConfig.HTTP.LogLevel, "stdout")
	} else {
		logger, e = GetLogger(GlobalConfig.HTTP.LogLevel, "stdout", filePath)
	}

	if e != nil {
		panic("log init fail!")
	}
	DefaultLogger = logger

	log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
}

func unmarshalText(text string) zapcore.Level {
	level := zapcore.InfoLevel
	switch text {
	case "debug", "DEBUG":
		level = zapcore.DebugLevel
	case "info", "INFO", "": // make the zero value useful
		level = zapcore.InfoLevel
	case "warn", "WARN":
		level = zapcore.WarnLevel
	case "error", "ERROR":
		level = zapcore.ErrorLevel
	case "dpanic", "DPANIC":
		level = zapcore.DPanicLevel
	case "panic", "PANIC":
		level = zapcore.PanicLevel
	case "fatal", "FATAL":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}
	return level
}
