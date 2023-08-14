package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type slogwrapper struct {
	prefix string
	slog   *slog.Logger
}

// WrapSlog returns Logger based on log/slog backend.
func WrapSlog(l *slog.Logger) Logger {
	return &slogwrapper{
		slog: cloneslog(l),
	}
}

func (w *slogwrapper) WithPrefix(prefix string) Logger {
	return &slogwrapper{
		prefix: w.prefix + prefix,
		slog:   cloneslog(w.slog),
	}
}

func (w *slogwrapper) WithPrefixf(format string, args ...any) Logger {
	return w.WithPrefix(fmt.Sprintf(format, args...))
}

func (w *slogwrapper) WithField(key string, value any) Logger {
	return &slogwrapper{
		prefix: w.prefix,
		slog:   w.slog.With(key, value),
	}
}

func (w *slogwrapper) WithError(err error) Logger {
	return &slogwrapper{
		prefix: w.prefix,
		slog:   w.slog.With("error", err),
	}
}

func (w *slogwrapper) WithFields(fields map[string]any) Logger {
	l := &slogwrapper{
		prefix: w.prefix,
		slog:   cloneslog(w.slog),
	}

	for k, v := range fields {
		l.slog = l.slog.With(k, v)
	}

	return l
}

func (w *slogwrapper) Debug(args ...any) {
	w.slog.Debug(w.args(args))
}

func (w *slogwrapper) Debugf(format string, args ...any) {
	w.slog.Debug(w.format(format, args))
}

func (w *slogwrapper) Info(args ...any) {
	w.slog.Info(w.args(args))
}

func (w *slogwrapper) Infof(format string, args ...any) {
	w.slog.Info(w.format(format, args))
}

func (w *slogwrapper) Warn(args ...any) {
	w.slog.Warn(w.args(args))
}

func (w *slogwrapper) Warnf(format string, args ...any) {
	w.slog.Warn(w.format(format, args))
}

func (w *slogwrapper) Error(args ...any) {
	w.slog.Error(w.args(args))
}

func (w *slogwrapper) Errorf(format string, args ...any) {
	w.slog.Error(w.format(format, args))
}

func (w *slogwrapper) Print(args ...any) {
	w.slog.Info(w.args(args))
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

func (w *slogwrapper) format(format string, args []any) string {
	return fmt.Sprintf(w.appendPrefix(format), args...)
}

func (w *slogwrapper) args(args []any) string {
	args = w.prependPrefix(args)
	return fmt.Sprint(args...)
}

func (w *slogwrapper) argsln(args []any) string {
	args = w.prependPrefix(args)
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}

func (w *slogwrapper) appendPrefix(msg string) string {
	if w.prefix == "" {
		return msg
	}

	return fmt.Sprintf("%s %s", w.prefix, msg)
}

func (w *slogwrapper) prependPrefix(args []any) []any {
	if w.prefix == "" {
		return args
	}

	return append([]any{w.prefix + " "}, args...)
}

// aka l.clone()
func cloneslog(l *slog.Logger) *slog.Logger {
	c := *l
	return &c
}
