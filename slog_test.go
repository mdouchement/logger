package logger_test

import (
	"log/slog"
	"testing"

	"github.com/mdouchement/logger"
)

func TestSlogCompliance(t *testing.T) {
	l := slog.New(nil)
	var _ logger.Logger = logger.WrapSlog(l)
}
