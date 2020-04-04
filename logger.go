package logger

// A Logger is the interface used in this package for logging,
// so that any backend can be plugged in.
type Logger interface {
	WithPrefix(prefix string) Logger
	WithField(key string, value interface{}) Logger
	WithError(error error) Logger
	WithFields(fields map[string]interface{}) Logger
	//
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	//
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
}

// M defines a generic map of type `map[string]interface{}`.
type M map[string]interface{}
