package logger

import (
	"route256/loms/internal/config"

	"go.uber.org/zap"
)

var globalLogger *zap.Logger

func Init(dev bool) {
	globalLogger = NewLogger(dev)
}

func GetLogger() *zap.Logger {
	if globalLogger == nil {
		Init(config.ConfigData.Dev)
	}

	return globalLogger
}

func NewLogger(dev bool) *zap.Logger {
	var logger *zap.Logger
	var err error
	if dev {
		cfg := zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		logger, err = cfg.Build()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

		logger, err = cfg.Build()
	}

	if err != nil {
		panic(err)
	}

	return logger
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}
