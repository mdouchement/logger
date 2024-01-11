package logger

import (
	"fmt"
	"log/slog"
	"strings"
)

// ParseSlogLevel takes a string level and returns the slog.Level constant.
func ParseSlogLevel(lvl string) (slog.Level, error) {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	}

	return slog.LevelError, fmt.Errorf("not a valid slog level: %s", lvl)
}
