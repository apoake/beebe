package log


import (
	"go.uber.org/zap"
	"beebe/config"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger
var Mlog *zap.Logger

func init() {
	Log = configLogger()
	Mlog = configMacaronLogger()
}

func configLogger() *zap.Logger {
	conf := config.GetConfig().Log
	logConfig := zapConfig(conf.IsDebug, []string{conf.LogPath, "stderr"})
	logger, _ := logConfig.Build()
	return logger
}

func configMacaronLogger() *zap.Logger {
	conf := config.GetConfig().Log
	logConfig := zapConfig(conf.IsDebug, []string{conf.MacaronPath, "stderr"})
	logger, _ := logConfig.Build()
	return logger
}

func zapConfig(isDevelop bool, output []string) zap.Config {
	outputPath := []string{"stdout"}
	if output == nil {
		outputPath = output
	}
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: isDevelop,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      outputPath,
		ErrorOutputPaths: []string{"stderr"},
	}
}