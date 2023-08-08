package logger

import (
	"context"
)

// A ContextKey is used to add data into a context.Context.
type ContextKey string

// LoggerKey is the storing key used for storing and retrieve the logger from a context.
const LoggerKey ContextKey = "_logger"

// WithLogger returns a new context that embeds the given logger.
func WithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, l)
}

// LogWith returns the logger extracted from the context.
// It panics if no logger inside the context.
func LogWith(ctx context.Context) Logger {
	l, ok := ctx.Value(LoggerKey).(Logger)
	if !ok {
		panic("no logger defined in the given context")
	}
	return l
}
