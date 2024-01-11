package logger_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/mdouchement/logger"
)

func TestSlogCompliance(t *testing.T) {
	l := slog.New(slog.NewJSONHandler(new(bytes.Buffer), nil))
	var _ logger.Logger = logger.WrapSlog(l)
}
