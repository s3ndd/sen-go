package log

import (
	"go.uber.org/zap/zapcore"
)

// Level is the log level
type Level string

// Fields is a map of fields
type Fields map[string]interface{}

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelFatal Level = "FATAL"
	LevelPanic Level = "PANIC"
)

func (l Level) String() string {
	return string(l)
}

// LoggerConfig is the configuration for the logger.
type LoggerConfig interface {
	Level() string
	Format() string
	SourceProgram() string
	IsDev() bool
}

// Logger is the interface for the logger.
type Logger interface {
	WithLevel(level Level) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithError(err error) Logger
	Debug(message string, args ...zapcore.Field)
	Info(message string, args ...zapcore.Field)
	Warn(message string, args ...zapcore.Field)
	Error(message string, args ...zapcore.Field)
	Fatal(message string, args ...zapcore.Field)
	Panic(message string, args ...zapcore.Field)
}

// NewLogger returns a new logger.
func NewLogger(config *Config) Logger {
	return NewZapLogger(config)
}
