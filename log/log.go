package log


import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
}

func Logger() *zap.Logger {
	return logger
}