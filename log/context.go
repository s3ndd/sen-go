package log

import "context"

// contextKey is the type for context keys.
type contextKey int

const (
	// loggerContextKey is the context key for the logger.
	loggerContextKey contextKey = iota
)

// ContextLogger returns the logger stored in context or a new logger.
func ContextLogger(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerContextKey).(Logger); ok {
		return logger
	}

	return Global()
}

// ContextWithLogger returns a new context with the logger.
func ContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}
