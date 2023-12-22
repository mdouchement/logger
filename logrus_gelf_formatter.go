package logger

import (
	"fmt"
	"os"
	"sync"

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
	gelf := NewBufferGELF()

	for k, v := range entry.Data {
		if k == KeyPrefix {
			entry.Message = fmt.Sprintf("%s %s", v, entry.Message)
			continue
		}
		gelf.Add(k, v)
	}

	gelf.Host(f.Hostname)
	gelf.Timestamp(entry.Time)
	gelf.Level(f.priorities(entry.Level))
	gelf.Message(entry.Message)
	gelf.Add("level_name", entry.Level.String())

	if entry.Caller != nil {
		gelf.Add("file", entry.Caller.File)
	}

	return gelf.Complete(true), nil
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
