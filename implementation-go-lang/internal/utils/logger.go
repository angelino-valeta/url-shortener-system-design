package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger setup
func NewLogger(level string) *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zapLevel(level)
	logger, _ := config.Build()
	return logger
}

func zapLevel(level string) zap.AtomicLevel {
	atomicLevel := zap.NewAtomicLevel()

	switch level {
	case "debug":
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case "warn":
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case "error":
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	default:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	}

	return atomicLevel
}
