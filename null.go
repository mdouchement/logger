package logger

type null struct{}

// NewNullLogger returns null Logger.
func NewNullLogger() Logger {
	return &null{}
}

func (w *null) WithPrefix(prefix string) Logger {
	return w
}

func (w *null) WithPrefixf(prefix string, args ...interface{}) Logger {
	return w
}

func (w *null) WithField(key string, value interface{}) Logger {
	return w
}

func (w *null) WithError(err error) Logger {
	return w
}

func (w *null) WithFields(fields map[string]interface{}) Logger {
	return w
}

func (w *null) Debug(args ...interface{}) {
}

func (w *null) Debugf(format string, args ...interface{}) {
}

func (w *null) Info(args ...interface{}) {
}

func (w *null) Infof(format string, args ...interface{}) {
}

func (w *null) Warn(args ...interface{}) {
}

func (w *null) Warnf(format string, args ...interface{}) {
}

func (w *null) Error(args ...interface{}) {
}

func (w *null) Errorf(format string, args ...interface{}) {
}

func (w *null) Print(args ...interface{}) {
}

func (w *null) Printf(format string, args ...interface{}) {
}

func (w *null) Println(args ...interface{}) {
}

func (w *null) Fatal(args ...interface{}) {
}

func (w *null) Fatalf(format string, args ...interface{}) {
}

func (w *null) Fatalln(args ...interface{}) {
}

func (w *null) Panic(args ...interface{}) {
}

func (w *null) Panicf(format string, args ...interface{}) {
}

func (w *null) Panicln(args ...interface{}) {
}
