package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type logruswrapper struct {
	prefix string
	logrus *logrus.Entry
}

// WrapLogrus returns Logger based on Logrus backend.
func WrapLogrus(l *logrus.Logger) Logger {
	return &logruswrapper{
		logrus: logrus.NewEntry(l),
	}
}

func (w *logruswrapper) WithPrefix(prefix string) Logger {
	return &logruswrapper{
		prefix: w.prefix + prefix,
		logrus: w.logrus.WithContext(w.logrus.Context), // Create a new entry without modifying anything.
	}
}

func (w *logruswrapper) WithPrefixf(format string, args ...any) Logger {
	return w.WithPrefix(fmt.Sprintf(format, args...))
}

func (w *logruswrapper) WithField(key string, value any) Logger {
	return &logruswrapper{
		prefix: w.prefix,
		logrus: w.logrus.WithField(key, value),
	}
}

func (w *logruswrapper) WithError(err error) Logger {
	return &logruswrapper{
		prefix: w.prefix,
		logrus: w.logrus.WithError(err),
	}
}

func (w *logruswrapper) WithFields(fields map[string]any) Logger {
	return &logruswrapper{
		prefix: w.prefix,
		logrus: w.logrus.WithFields(fields),
	}
}

func (w *logruswrapper) Debug(args ...any) {
	w.logrus.Debug(w.prependPrefix(args)...)
}

func (w *logruswrapper) Debugf(format string, args ...any) {
	w.logrus.Debugf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Info(args ...any) {
	w.logrus.Info(w.prependPrefix(args)...)
}

func (w *logruswrapper) Infof(format string, args ...any) {
	w.logrus.Infof(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Warn(args ...any) {
	w.logrus.Warn(w.prependPrefix(args)...)
}

func (w *logruswrapper) Warnf(format string, args ...any) {
	w.logrus.Warnf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Error(args ...any) {
	w.logrus.Error(w.prependPrefix(args)...)
}

func (w *logruswrapper) Errorf(format string, args ...any) {
	w.logrus.Errorf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Print(args ...any) {
	w.logrus.Print(w.prependPrefix(args)...)
}

func (w *logruswrapper) Printf(format string, args ...any) {
	w.logrus.Printf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Println(args ...any) {
	w.logrus.Println(w.prependPrefix(args)...)
}

func (w *logruswrapper) Fatal(args ...any) {
	w.logrus.Fatal(w.prependPrefix(args)...)
}

func (w *logruswrapper) Fatalf(format string, args ...any) {
	w.logrus.Fatalf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Fatalln(args ...any) {
	w.logrus.Fatalln(w.prependPrefix(args)...)
}

func (w *logruswrapper) Panic(args ...any) {
	w.logrus.Panic(w.prependPrefix(args)...)
}

func (w *logruswrapper) Panicf(format string, args ...any) {
	w.logrus.Panicf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Panicln(args ...any) {
	w.logrus.Panicln(w.prependPrefix(args)...)
}

func (w *logruswrapper) appendPrefix(msg string) string {
	if w.prefix == "" {
		return msg
	}

	return fmt.Sprintf("%s %s", w.prefix, msg)
}

func (w *logruswrapper) prependPrefix(args []any) []any {
	if w.prefix == "" {
		return args
	}

	return append([]any{w.prefix + " "}, args...)
}
