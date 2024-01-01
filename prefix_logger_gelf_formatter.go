package logger

import (
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/mdouchement/logger/syslog"
)

type PrefixGELF struct {
	Writer   io.Writer
	Level    slog.Level
	Hostname string
}

func (f *PrefixGELF) Enabled(l slog.Level) bool {
	return l >= f.Level
}

func (f *PrefixGELF) Format(l slog.Level, t time.Time, attrs Attrs, message string) error {
	gelf := NewBufferGELF()

	var prefix string
	attrs(func(k string, v any) bool {
		if k == KeyPrefix {
			if p, ok := v.(string); ok {
				prefix += p
			}
			return true
		}
		gelf.Add(k, v)
		return true
	})

	// Main fields.
	gelf.Host(f.Hostname)
	gelf.Timestamp(t)
	gelf.Level(f.priorities(l))

	if prefix != "" {
		message = fmt.Sprintf("%s %s", prefix, message)
	}
	gelf.Message(message)

	gelf.Add("level_name", l.String())
	// m["_file"]
	// m["_line"]

	_, err := f.Writer.Write(gelf.Complete(true))
	return err
}

func (PrefixGELF) priorities(level slog.Level) int32 {
	var p syslog.Priority

	switch level {
	case slog.LevelError:
		p = syslog.LOG_ERR
	case slog.LevelWarn:
		p = syslog.LOG_WARNING
	case slog.LevelInfo:
		p = syslog.LOG_INFO
	case slog.LevelDebug:
		p = syslog.LOG_DEBUG
	}

	return int32(p)
}
