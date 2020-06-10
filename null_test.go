package logger_test

import (
	"testing"

	"github.com/mdouchement/logger"
)

func TestNullCompliance(t *testing.T) {
	var _ logger.Logger = logger.NewNullLogger()
}
