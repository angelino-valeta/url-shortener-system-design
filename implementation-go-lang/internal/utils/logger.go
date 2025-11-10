package utils

import "go.uber.org/zap"

// NewLogger setup
func NewLogger(level string) *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapLevel(level))
	logger, _ := config.Build()
	return logger
}

func zapLevel(level string) zap.AtomicLevel {
	switch level {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}
