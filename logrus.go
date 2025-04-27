package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type logruswrapper struct {
	logrus *logrus.Entry
}

// WrapLogrus returns Logger based on Logrus backend.
//
// Deprecated: Use WrapSlog instead.
func WrapLogrus(l *logrus.Logger) Logger {
	return &logruswrapper{
		logrus: logrus.NewEntry(l),
	}
}

func (w *logruswrapper) WithPrefix(prefix string) Logger {
	p, _ := w.logrus.Data[KeyPrefix].(string) // If the map value is not a string, p is an empty string (no panic).
	return &logruswrapper{
		logrus: w.logrus.WithField(KeyPrefix, p+prefix),
	}
}

func (w *logruswrapper) WithPrefixf(format string, args ...any) Logger {
	return w.WithPrefix(fmt.Sprintf(format, args...))
}

func (w *logruswrapper) WithField(key string, value any) Logger {
	return &logruswrapper{
		logrus: w.logrus.WithField(key, value),
	}
}

func (w *logruswrapper) WithError(err error) Logger {
	return &logruswrapper{
		logrus: w.logrus.WithError(err),
	}
}

func (w *logruswrapper) WithFields(fields map[string]any) Logger {
	return &logruswrapper{
		logrus: w.logrus.WithFields(fields),
	}
}

func (w *logruswrapper) Debug(args ...any) {
	w.logrus.Debug(args...)
}

func (w *logruswrapper) Debugf(format string, args ...any) {
	w.logrus.Debugf(format, args...)
}

func (w *logruswrapper) Info(args ...any) {
	w.logrus.Info(args...)
}

func (w *logruswrapper) Infof(format string, args ...any) {
	w.logrus.Infof(format, args...)
}

func (w *logruswrapper) Warn(args ...any) {
	w.logrus.Warn(args...)
}

func (w *logruswrapper) Warnf(format string, args ...any) {
	w.logrus.Warnf(format, args...)
}

func (w *logruswrapper) Error(args ...any) {
	w.logrus.Error(args...)
}

func (w *logruswrapper) Errorf(format string, args ...any) {
	w.logrus.Errorf(format, args...)
}

func (w *logruswrapper) Print(args ...any) {
	w.logrus.Print(args...)
}

func (w *logruswrapper) Printf(format string, args ...any) {
	w.logrus.Printf(format, args...)
}

func (w *logruswrapper) Println(args ...any) {
	w.logrus.Println(args...)
}

func (w *logruswrapper) Fatal(args ...any) {
	w.logrus.Fatal(args...)
}

func (w *logruswrapper) Fatalf(format string, args ...any) {
	w.logrus.Fatalf(format, args...)
}

func (w *logruswrapper) Fatalln(args ...any) {
	w.logrus.Fatalln(args...)
}

func (w *logruswrapper) Panic(args ...any) {
	w.logrus.Panic(args...)
}

func (w *logruswrapper) Panicf(format string, args ...any) {
	w.logrus.Panicf(format, args...)
}

func (w *logruswrapper) Panicln(args ...any) {
	w.logrus.Panicln(args...)
}
