package log


import (
	"go.uber.org/zap"
	"beebe/config"
	"time"
	"strings"
)

var Log *zap.Logger
var Mlog *zap.Logger

func init() {
	Log = configLogger()
	Mlog = configMacaronLogger()
}

func configLogger() *zap.Logger {
	logConfig := zap.NewProductionConfig()
	today := time.Now().Format("2006-01-02")
	conf := config.GetConfig().Log
	index := strings.LastIndex(conf.LogPath, ".")
	logConfig.OutputPaths = []string{conf.LogPath[0: index] + today + conf.LogPath[index:], "stderr"}
	//index = strings.LastIndex(conf.ErrorPath, ".")
	//logConfig.ErrorOutputPaths = []string{conf.ErrorPath[0: index] + today + conf.ErrorPath[index:], "stderr"}
	logger, _ := logConfig.Build()
	return logger
}

func configMacaronLogger() *zap.Logger {
	logConfig := zap.NewProductionConfig()
	today := time.Now().Format("2006-01-02")
	conf := config.GetConfig().Log
	index := strings.LastIndex(conf.MacaronPath, ".")
	logConfig.OutputPaths = []string{conf.MacaronPath[0: index] + today + conf.MacaronPath[index:], "stderr"}
	logger, _ := logConfig.Build()
	return logger
}

func zapConfig(isDevelop bool) zap.Config {
	if isDevelop {

	}
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: isDevelop,
		Sampling: &SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}