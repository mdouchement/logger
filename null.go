package logger

type null struct{}

// NewNullLogger returns null Logger.
func NewNullLogger() Logger {
	return &null{}
}

func (w *null) WithPrefix(prefix string) Logger {
	return w
}

func (w *null) WithPrefixf(prefix string, args ...any) Logger {
	return w
}

func (w *null) WithField(key string, value any) Logger {
	return w
}

func (w *null) WithError(err error) Logger {
	return w
}

func (w *null) WithFields(fields map[string]any) Logger {
	return w
}

func (w *null) Debug(args ...any) {
}

func (w *null) Debugf(format string, args ...any) {
}

func (w *null) Info(args ...any) {
}

func (w *null) Infof(format string, args ...any) {
}

func (w *null) Warn(args ...any) {
}

func (w *null) Warnf(format string, args ...any) {
}

func (w *null) Error(args ...any) {
}

func (w *null) Errorf(format string, args ...any) {
}

func (w *null) Print(args ...any) {
}

func (w *null) Printf(format string, args ...any) {
}

func (w *null) Println(args ...any) {
}

func (w *null) Fatal(args ...any) {
}

func (w *null) Fatalf(format string, args ...any) {
}

func (w *null) Fatalln(args ...any) {
}

func (w *null) Panic(args ...any) {
}

func (w *null) Panicf(format string, args ...any) {
}

func (w *null) Panicln(args ...any) {
}
