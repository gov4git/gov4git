package base

import (
	"os"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	LogQuietly()
}

func newQuietConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	return cfg
}

func LogQuietly() {
	l, err := newQuietConfig().Build()
	if err != nil {
		println("cannot create logger:", err)
		os.Exit(1)
	}
	logger = l
}

func LogVerbosely() {
	l, err := zap.NewDevelopmentConfig().Build()
	if err != nil {
		println("cannot create logger:", err)
		os.Exit(1)
	}
	logger = l
}

func AssertNoErr(err error) {
	if err == nil {
		return
	}
	Fatalf("encountered %v", err)
}

func Infof(template string, args ...interface{}) {
	logger.Sugar().Infof(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Sugar().Fatalf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Sugar().Errorf(template, args...)
}

func Sync() error {
	return logger.Sync()
}
