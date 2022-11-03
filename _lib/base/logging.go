package base

import (
	"fmt"
	"os"
	"runtime"

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
	_, file, line, _ := runtime.Caller(1)
	src := fmt.Sprintf("%s:%d ", file, line)
	msg := fmt.Sprintf(template, args...)
	logger.Sugar().Info(src + msg)
}

func Fatalf(template string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	src := fmt.Sprintf("%s:%d ", file, line)
	msg := fmt.Sprintf(template, args...)
	logger.Sugar().Fatal(src + msg)
	// logger.Sugar().Fatalf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	src := fmt.Sprintf("%s:%d ", file, line)
	msg := fmt.Sprintf(template, args...)
	logger.Sugar().Error(src + msg)
	// logger.Sugar().Errorf(template, args...)
}

func Sync() error {
	return logger.Sync()
}
