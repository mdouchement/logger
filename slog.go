package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type slogwrapper struct {
	slog *slog.Logger
}

// WrapSlog returns Logger based on log/slog backend.
func WrapSlog(l *slog.Logger) Logger {
	return &slogwrapper{
		slog: l,
	}
}

func (w *slogwrapper) WithPrefix(prefix string) Logger {
	return &slogwrapper{
		slog: w.slog.With(KeyPrefix, prefix),
	}
}

func (w *slogwrapper) WithPrefixf(format string, args ...any) Logger {
	return w.WithPrefix(fmt.Sprintf(format, args...))
}

func (w *slogwrapper) WithField(key string, value any) Logger {
	return &slogwrapper{
		slog: w.slog.With(key, value),
	}
}

func (w *slogwrapper) WithError(err error) Logger {
	return &slogwrapper{
		slog: w.slog.With("error", err),
	}
}

func (w *slogwrapper) WithFields(fields map[string]any) Logger {
	l := w.slog
	for k, v := range fields {
		l = l.With(k, v)
	}

	return &slogwrapper{
		slog: l,
	}
}

func (w *slogwrapper) Debug(args ...any) {
	w.slog.Debug(fmt.Sprint(args...))
}

func (w *slogwrapper) Debugf(format string, args ...any) {
	w.slog.Debug(fmt.Sprintf(format, args...))
}

func (w *slogwrapper) Info(args ...any) {
	w.slog.Info(fmt.Sprint(args...))
}

func (w *slogwrapper) Infof(format string, args ...any) {
	w.slog.Info(fmt.Sprintf(format, args...))
}

func (w *slogwrapper) Warn(args ...any) {
	w.slog.Warn(fmt.Sprint(args...))
}

func (w *slogwrapper) Warnf(format string, args ...any) {
	w.slog.Warn(fmt.Sprintf(format, args...))
}

func (w *slogwrapper) Error(args ...any) {
	w.slog.Error(fmt.Sprint(args...))
}

func (w *slogwrapper) Errorf(format string, args ...any) {
	w.slog.Error(fmt.Sprintf(format, args...))
}

func (w *slogwrapper) Print(args ...any) {
	w.slog.Info(fmt.Sprint(args...))
}

func (w *slogwrapper) Printf(format string, args ...any) {
	w.Infof(format, args...)
}

func (w *slogwrapper) Println(args ...any) {
	w.slog.Info(w.argsln(args))
}

func (w *slogwrapper) Fatal(args ...any) {
	w.Error(args...)
	os.Exit(1)
}

func (w *slogwrapper) Fatalf(format string, args ...any) {
	w.Errorf(format, args...)
	os.Exit(1)
}

func (w *slogwrapper) Fatalln(args ...any) {
	w.slog.Error(w.argsln(args))
	os.Exit(1)
}

func (w *slogwrapper) Panic(args ...any) {
	w.Error(args...)
	panic(w)
}

func (w *slogwrapper) Panicf(format string, args ...any) {
	w.Errorf(format, args...)
	panic(w)
}

func (w *slogwrapper) Panicln(args ...any) {
	w.slog.Error(w.argsln(args))
	panic(w)
}

//
//
//
//

func (w *slogwrapper) argsln(args []any) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}
