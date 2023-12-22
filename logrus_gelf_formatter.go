package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mdouchement/logger/syslog"
	"github.com/sirupsen/logrus"
)

// A LogrusGELFFormatter is GELF formatter for Logrus.
type LogrusGELFFormatter struct {
	sync.Once
	Hostname string
}

func (f *LogrusGELFFormatter) init() {
	var err error

	if f.Hostname == "" {
		f.Hostname, err = os.Hostname()
		if err != nil {
			f.Hostname = "localhost"
		}
	}
}

// Format implements logrus.Formatter.
func (f *LogrusGELFFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	f.Do(f.init)

	m := make(map[string]any)
	for k, v := range entry.Data {
		// skip id if present
		if k == "id" || k == "_id" {
			continue
		}

		if k == KeyPrefix {
			entry.Message = fmt.Sprintf("%s %s", v, entry.Message)
			continue
		}

		// otherwise convert if necessary
		switch value := v.(type) {
		case time.Time:
			m["_"+k] = value.Format(time.RFC3339)
		case uint, uint8, uint16, uint32,
			int, int8, int16, int32, int64:
			m["_"+k] = value
		case uint64, float32, float64:
			// NOTE: uint64 is not supported by graylog due to java limitation
			//       so we're sending them as double for the time being
			m["_"+k] = value
		case bool:
			m["_"+k] = fmt.Sprintf("%t", value)
		case string:
			m["_"+k] = value
		default:
			m["_"+k] = fmt.Sprint(value)
		}
	}

	m["version"] = "1.1"
	m["host"] = f.Hostname
	m["timestamp"] = float64(entry.Time.UnixNano()) / 1e9 // Unix epoch timestamp in seconds with decimals for nanoseconds.
	m["level"] = f.priorities(entry.Level)

	m["short_message"] = entry.Message
	// If there are newlines in the message, use the first line
	// for the short_message and set the full_message to the
	// original input. If the input has no newlines, stick the
	// whole thing in short_message.
	if i := strings.IndexRune(entry.Message, '\n'); i > 0 {
		m["short_message"] = entry.Message[:i]
		m["full_message"] = entry.Message
	}

	m["_level_name"] = entry.Level.String()
	if entry.Caller != nil {
		m["_file"] = entry.Caller.File
		m["_line"] = entry.Caller.Line
	}

	payload, err := json.Marshal(m)
	return append(payload, '\n'), err
}

func (f *LogrusGELFFormatter) priorities(level logrus.Level) int32 {
	var p syslog.Priority

	switch level {
	case logrus.PanicLevel:
		p = syslog.LOG_ALERT
	case logrus.FatalLevel:
		p = syslog.LOG_CRIT
	case logrus.ErrorLevel:
		p = syslog.LOG_ERR
	case logrus.WarnLevel:
		p = syslog.LOG_WARNING
	case logrus.InfoLevel:
		p = syslog.LOG_INFO
	case logrus.DebugLevel:
		p = syslog.LOG_DEBUG
	}

	return int32(p)
}
