package log


import (
	"go.uber.org/zap"
	"fmt"
)

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Println(err)
	}
}

func Logger() *zap.Logger {
	return logger
}