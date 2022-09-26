package base

import (
	"os"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	l, err := zap.NewDevelopment()
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
