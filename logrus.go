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

func (w *logruswrapper) WithField(key string, value interface{}) Logger {
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

func (w *logruswrapper) WithFields(fields map[string]interface{}) Logger {
	return &logruswrapper{
		prefix: w.prefix,
		logrus: w.logrus.WithFields(fields),
	}
}

func (w *logruswrapper) Debug(args ...interface{}) {
	w.logrus.Debug(w.prependPrefix(args)...)
}

func (w *logruswrapper) Debugf(format string, args ...interface{}) {
	w.logrus.Debugf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Info(args ...interface{}) {
	w.logrus.Info(w.prependPrefix(args)...)
}

func (w *logruswrapper) Infof(format string, args ...interface{}) {
	w.logrus.Infof(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Warn(args ...interface{}) {
	w.logrus.Warn(w.prependPrefix(args)...)
}

func (w *logruswrapper) Warnf(format string, args ...interface{}) {
	w.logrus.Warnf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Error(args ...interface{}) {
	w.logrus.Error(w.prependPrefix(args)...)
}

func (w *logruswrapper) Errorf(format string, args ...interface{}) {
	w.logrus.Errorf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Print(args ...interface{}) {
	w.logrus.Print(w.prependPrefix(args)...)
}

func (w *logruswrapper) Printf(format string, args ...interface{}) {
	w.logrus.Printf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Println(args ...interface{}) {
	w.logrus.Println(w.prependPrefix(args)...)
}

func (w *logruswrapper) Fatal(args ...interface{}) {
	w.logrus.Fatal(w.prependPrefix(args)...)
}

func (w *logruswrapper) Fatalf(format string, args ...interface{}) {
	w.logrus.Fatalf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Fatalln(args ...interface{}) {
	w.logrus.Fatalln(w.prependPrefix(args)...)
}

func (w *logruswrapper) Panic(args ...interface{}) {
	w.logrus.Panic(w.prependPrefix(args)...)
}

func (w *logruswrapper) Panicf(format string, args ...interface{}) {
	w.logrus.Panicf(w.appendPrefix(format), args...)
}

func (w *logruswrapper) Panicln(args ...interface{}) {
	w.logrus.Panicln(w.prependPrefix(args)...)
}

func (w *logruswrapper) appendPrefix(msg string) string {
	if w.prefix == "" {
		return msg
	}

	return fmt.Sprintf("%s %s", w.prefix, msg)
}

func (w *logruswrapper) prependPrefix(args []interface{}) []interface{} {
	if w.prefix == "" {
		return args
	}

	return append([]interface{}{w.prefix + " "}, args...)
}
