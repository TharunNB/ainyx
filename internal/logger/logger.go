package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
)

func Init(env string) {
	once.Do(func() {
		var (
			log *zap.Logger
			err error
		)

		if env == "development" {
			log, err = buildDevelopmentLogger()
		} else {
			log, err = buildProductionLogger()
		}

		if err != nil {
			panic(fmt.Sprintf("logger: failed to initialise zap: %v", err))
		}

		instance = log
	})
}
func Get() *zap.Logger {
	if instance == nil {
		panic("logger: Get() called before Init()")
	}
	return instance
}

func Sync() {
	if instance != nil {
		_ = instance.Sync()
	}
}

func buildDevelopmentLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return cfg.Build()
}

func buildProductionLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return cfg.Build()
}
