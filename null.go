package logger

type null struct{}

// NewNullLogger returns null Logger.
func NewNullLogger() Logger {
	return &null{}
}

func (w *null) WithPrefix(_ string) Logger {
	return w
}

func (w *null) WithPrefixf(_ string, _ ...any) Logger {
	return w
}

func (w *null) WithField(_ string, _ any) Logger {
	return w
}

func (w *null) WithError(_ error) Logger {
	return w
}

func (w *null) WithFields(_ map[string]any) Logger {
	return w
}

func (w *null) Debug(_ ...any) {
}

func (w *null) Debugf(_ string, _ ...any) {
}

func (w *null) Info(_ ...any) {
}

func (w *null) Infof(_ string, _ ...any) {
}

func (w *null) Warn(_ ...any) {
}

func (w *null) Warnf(_ string, _ ...any) {
}

func (w *null) Error(_ ...any) {
}

func (w *null) Errorf(_ string, _ ...any) {
}

func (w *null) Print(_ ...any) {
}

func (w *null) Printf(_ string, _ ...any) {
}

func (w *null) Println(_ ...any) {
}

func (w *null) Fatal(_ ...any) {
}

func (w *null) Fatalf(_ string, _ ...any) {
}

func (w *null) Fatalln(_ ...any) {
}

func (w *null) Panic(_ ...any) {
}

func (w *null) Panicf(_ string, _ ...any) {
}

func (w *null) Panicln(_ ...any) {
}
