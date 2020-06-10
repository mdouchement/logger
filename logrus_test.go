package logger_test

import (
	"testing"

	"github.com/mdouchement/logger"
	"github.com/sirupsen/logrus"
)

func TestLogrusCompliance(t *testing.T) {
	l := logrus.New()
	var _ logger.Logger = logger.WrapLogrus(l)
}
