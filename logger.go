package logger

// KeyPrefix is used to tell that en given attr/field is a prefix.
const KeyPrefix = "__prefix"

// A Logger is the interface used in this package for logging,
// so that any backend can be plugged in.
type Logger interface {
	WithPrefix(prefix string) Logger
	WithPrefixf(format string, args ...any) Logger
	WithField(key string, value any) Logger
	WithError(error error) Logger
	WithFields(fields map[string]any) Logger
	//
	Debug(args ...any)
	Debugf(format string, args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	//
	Print(args ...any)
	Printf(format string, args ...any)
	Println(args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Fatalln(args ...any)
	Panic(args ...any)
	Panicf(format string, args ...any)
	Panicln(args ...any)
}

// M defines a generic map of type `map[string]interface{}`.
type M map[string]any
