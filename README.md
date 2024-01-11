# Logger

A simple prefixer logger that can wrap several backends. The formatters/handlers from this repository must be used in order to make prefix work.


## Interface

```go
type Logger interface {
	WithPrefix(prefix string) Logger
	WithPrefixf(format string, args ...interface{}) Logger
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
```


## Backends

- [Logrus](https://github.com/sirupsen/logrus) with `logger.WrapLogrus(l *logrus.Logger)` function
- [slog](https://pkg.go.dev/log/slog) with `logger.WrapSlog(l *slog.Logger)` function
  - [slog_gelf_handler.go](https://github.com/mdouchement/logger/blob/master/slog_gelf_handler.go) can be used independently
  - [slog_text_handler.go](https://github.com/mdouchement/logger/blob/master/slog_text_handler.go) (logrus format) can be used independently


## Helpers

- `type M map[string]interface{}` used as a shorthand for a map of interfaces used by `WithFields` method
- `func WithLogger(ctx context.Context, l Logger) context.Context` to embeds the logger inside a context
- `func LogWith(ctx context.Context) Logger` to extract a logger from a context

## License

**MIT**


## Contributing

All PRs are welcome.

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
5. Push to the branch (git push origin my-new-feature)
6. Create new Pull Request
