package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

type (
	Formatter interface {
		Enabled(l slog.Level) bool
		Format(l slog.Level, t time.Time, attrs Attrs, message string) error
	}

	Attr     func(k string, v any) bool
	Attrs    func(Attr)
	NextAttr func() ([]slog.Attr, bool)
)

type PrefixLogger struct {
	parent    *PrefixLogger // TODO: is it a good idea for the GC?
	formatter Formatter
	attrs     []slog.Attr
}

func NewPrefixLogger(f Formatter) *PrefixLogger {
	return &PrefixLogger{
		formatter: f,
	}
}

func (l *PrefixLogger) WithPrefix(prefix string) Logger {
	nl := l.Clone()
	nl.attrs = append(nl.attrs, slog.Any(KeyPrefix, prefix))
	return nl
}

func (l *PrefixLogger) WithPrefixf(format string, args ...any) Logger {
	return l.WithPrefix(fmt.Sprintf(format, args...))
}

func (l *PrefixLogger) WithField(key string, value any) Logger {
	nl := l.Clone()
	nl.attrs = append(nl.attrs, slog.Any(key, value))
	return nl
}

func (l *PrefixLogger) WithError(err error) Logger {
	nl := l.Clone()
	nl.attrs = append(nl.attrs, slog.Any("error", err))
	return nl
}

func (l *PrefixLogger) WithFields(fields map[string]any) Logger {
	nl := l.Clone()
	nl.attrs = make([]slog.Attr, 0, len(fields))
	for k, v := range fields {
		nl.attrs = append(nl.attrs, slog.Any(k, v))
	}

	return nl
}

func (l *PrefixLogger) Debug(args ...any) {
	l.logln(slog.LevelDebug, args)
}

func (l *PrefixLogger) Debugf(format string, args ...any) {
	l.logf(slog.LevelDebug, format, args)
}

func (l *PrefixLogger) Info(args ...any) {
	l.logln(slog.LevelInfo, args)
}

func (l *PrefixLogger) Infof(format string, args ...any) {
	l.logf(slog.LevelInfo, format, args)
}

func (l *PrefixLogger) Warn(args ...any) {
	l.logln(slog.LevelWarn, args)
}

func (l *PrefixLogger) Warnf(format string, args ...any) {
	l.logf(slog.LevelWarn, format, args)
}

func (l *PrefixLogger) Error(args ...any) {
	l.logln(slog.LevelError, args)
}

func (l *PrefixLogger) Errorf(format string, args ...any) {
	l.logf(slog.LevelError, format, args)
}

func (l *PrefixLogger) Print(args ...any) {
	l.log(slog.LevelInfo, fmt.Sprint(args...))
}

func (l *PrefixLogger) Printf(format string, args ...any) {
	l.Infof(format, args...)
}

func (l *PrefixLogger) Println(args ...any) {
	l.Info(args...)
}

func (l *PrefixLogger) Fatal(args ...any) {
	l.log(slog.LevelError, fmt.Sprint(args...))
	os.Exit(1)
}

func (l *PrefixLogger) Fatalf(format string, args ...any) {
	l.Errorf(format, args...)
	os.Exit(1)
}

func (l *PrefixLogger) Fatalln(args ...any) {
	l.Error(args...)
	os.Exit(1)
}

func (l *PrefixLogger) Panic(args ...any) {
	l.log(slog.LevelError, fmt.Sprint(args...))
	panic(l)
}

func (l *PrefixLogger) Panicf(format string, args ...any) {
	l.Errorf(format, args...)
	panic(l)
}

func (l *PrefixLogger) Panicln(args ...any) {
	l.Error(args...)
	panic(l)
}

// Clone clones the entry, it creates a new instance and linking the parent to it.
func (l *PrefixLogger) Clone() *PrefixLogger {
	return &PrefixLogger{
		parent:    l,
		formatter: l.formatter,
	}
}

func (l *PrefixLogger) Attrs(attr Attr) {
	// Get all parents in a list.
	ilineage := make([]*PrefixLogger, 0, 100)
	p := l
	for p != nil {
		ilineage = append(ilineage, p)
		p = p.parent
	}

	// Process attrs.
	for i := len(ilineage) - 1; i >= 0; i-- {
		for _, a := range ilineage[i].attrs {
			ok := attr(a.Key, a.Value.Any())
			if !ok {
				i = 0
				break
			}
		}
	}
}

func (l *PrefixLogger) NextAttr() NextAttr {
	// Get all parents in a list.
	ilineage := make([]*PrefixLogger, 0, 100)
	p := l
	for p != nil {
		ilineage = append(ilineage, p)
		p = p.parent
	}

	i := len(ilineage)
	return func() ([]slog.Attr, bool) {
		i--
		return ilineage[i].attrs, i != 0
	}
}

//
//
//
//

// join args with spaces. The \n at the end of string is trimed.
func (l *PrefixLogger) logln(level slog.Level, args []any) {
	msg := fmt.Sprintln(args...)
	l.log(level, msg[:len(msg)-1])
}

func (l *PrefixLogger) logf(level slog.Level, msg string, args []any) {
	l.log(level, fmt.Sprintf(msg, args...))
}

func (l *PrefixLogger) log(level slog.Level, msg string) {
	if !l.formatter.Enabled(level) {
		return
	}

	l.formatter.Format(level, time.Now(), l.Attrs, msg)
}
