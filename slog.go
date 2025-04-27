package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

type slogwrapper struct {
	handler slog.Handler
}

// WrapSlog returns Logger based on log/slog backend.
func WrapSlog(l *slog.Logger) Logger {
	return &slogwrapper{
		handler: l.Handler(),
	}
}

// UnwrapSlog returns the embedded slog.Logger if WrapSlog/WrapSlogHandler was used to wrap.
func UnwrapSlog(l Logger) (*slog.Logger, bool) {
	w, ok := l.(*slogwrapper)
	if !ok {
		return nil, false
	}

	return slog.New(w.handler), true
}

// WrapSlogHandler returns Logger based on log/slog's handler backend.
func WrapSlogHandler(h slog.Handler) Logger {
	return &slogwrapper{
		handler: h,
	}
}

// UnwrapSlogHandler returns the embedded slog.Handler if WrapSlog/WrapSlogHandler was used to wrap.
func UnwrapSlogHandler(l Logger) (slog.Handler, bool) {
	w, ok := l.(*slogwrapper)
	if !ok {
		return nil, false
	}

	return w.handler, true
}

func (w *slogwrapper) WithPrefix(prefix string) Logger {
	return &slogwrapper{
		handler: w.handler.WithAttrs([]slog.Attr{slog.Any(KeyPrefix, prefix)}),
	}
}

func (w *slogwrapper) WithPrefixf(format string, args ...any) Logger {
	return w.WithPrefix(fmt.Sprintf(format, args...))
}

func (w *slogwrapper) WithField(key string, value any) Logger {
	return &slogwrapper{
		handler: w.handler.WithAttrs([]slog.Attr{slog.Any(key, value)}),
	}
}

func (w *slogwrapper) WithError(err error) Logger {
	return &slogwrapper{
		handler: w.handler.WithAttrs([]slog.Attr{slog.Any("error", err)}),
	}
}

func (w *slogwrapper) WithFields(fields map[string]any) Logger {
	attrs := make([]slog.Attr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}

	return &slogwrapper{
		handler: w.handler.WithAttrs(attrs),
	}
}

func (w *slogwrapper) Debug(args ...any) {
	w.logln(slog.LevelDebug, args)
}

func (w *slogwrapper) Debugf(format string, args ...any) {
	w.logf(slog.LevelDebug, format, args)
}

func (w *slogwrapper) Info(args ...any) {
	w.logln(slog.LevelInfo, args)
}

func (w *slogwrapper) Infof(format string, args ...any) {
	w.logf(slog.LevelInfo, format, args)
}

func (w *slogwrapper) Warn(args ...any) {
	w.logln(slog.LevelWarn, args)
}

func (w *slogwrapper) Warnf(format string, args ...any) {
	w.logf(slog.LevelWarn, format, args)
}

func (w *slogwrapper) Error(args ...any) {
	w.logln(slog.LevelError, args)
}

func (w *slogwrapper) Errorf(format string, args ...any) {
	w.logf(slog.LevelError, format, args)
}

func (w *slogwrapper) Print(args ...any) {
	w.log(slog.LevelInfo, fmt.Sprint(args...))
}

func (w *slogwrapper) Printf(format string, args ...any) {
	w.Infof(format, args...)
}

func (w *slogwrapper) Println(args ...any) {
	w.Info(args...)
}

func (w *slogwrapper) Fatal(args ...any) {
	w.log(slog.LevelError, fmt.Sprint(args...))
	os.Exit(1)
}

func (w *slogwrapper) Fatalf(format string, args ...any) {
	w.Errorf(format, args...)
	os.Exit(1)
}

func (w *slogwrapper) Fatalln(args ...any) {
	w.Error(args...)
	os.Exit(1)
}

func (w *slogwrapper) Panic(args ...any) {
	w.log(slog.LevelError, fmt.Sprint(args...))
	panic(w)
}

func (w *slogwrapper) Panicf(format string, args ...any) {
	w.Errorf(format, args...)
	panic(w)
}

func (w *slogwrapper) Panicln(args ...any) {
	w.Error(args...)
	panic(w)
}

//
//
//
//

// join args with spaces. The \n at the end of string is trimed.
func (w *slogwrapper) logln(level slog.Level, args []any) {
	msg := fmt.Sprintln(args...)
	w.log(level, msg[:len(msg)-1])
}

func (w *slogwrapper) logf(level slog.Level, msg string, args []any) {
	w.log(level, fmt.Sprintf(msg, args...))
}

func (w *slogwrapper) log(level slog.Level, msg string) {
	if !w.handler.Enabled(void, level) {
		return
	}

	var pc uintptr
	r := slog.NewRecord(time.Now(), level, msg, pc)
	w.handler.Handle(void, r)
}
