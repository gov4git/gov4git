package base

import "testing"

func TestLogging(t *testing.T) {
	logger.Sugar().Infof("abc %d", 3)
}
