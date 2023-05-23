package log

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RequestIDHeaderName is the header name for the request id
const RequestIDHeaderName = "X-Request-Id"

// ZapLogger is a wrapper around the zap logger
type ZapLogger struct {
	config *zap.Config
	*zap.Logger
}

// NewZapLogger creates a new zap logger
func NewZapLogger(config *Config) Logger {
	zapConfig := newZapConfig(config)
	zapCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
		zapcore.Lock(os.Stdout),
		zapConfig.Level,
	)

	logger := zap.New(zapCore, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// attach the source name to the logger
	logger = logger.With(zap.String("source_program", config.SourceProgram()))

	zapLogger := &ZapLogger{
		zapConfig,
		logger,
	}

	return zapLogger
}

// NewContext updates the logger stored in the context by adding new zap.Fields to it
func NewContext(ctx context.Context, fields ...zap.Field) context.Context {
	return ContextWithLogger(ctx, WithContext(ctx).With(fields...))
}

// WithContext returns the logger in the context or the global logger
func WithContext(ctx context.Context) *ZapLogger {
	logger := ContextLogger(ctx)
	if logger == nil {
		logger = Global()
	}
	return logger.(*ZapLogger)
}

// ForRequest returns a Logger for the request context.
func ForRequest(ctx context.Context) Logger {
	logger := ContextLogger(ctx)
	logger = logger.WithField("request_id", ctx.Value(RequestIDHeaderName))
	return logger
}

// newZapConfig creates a new zap config
func newZapConfig(config *Config) *zap.Config {
	logLevel := zapcore.Level(0)
	if err := logLevel.UnmarshalText([]byte(config.Level())); err != nil {
		panic(fmt.Errorf("Failed to set the zap log level from config. %s", config.Level()))
	}

	return &zap.Config{
		Development:      config.IsDev(),
		Encoding:         config.Format(),
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},

		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			CallerKey:    "caller",
			TimeKey:      "time",
			NameKey:      "name",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
		},
	}
}

// WithLevel returns the logger at the supplied level.
func (l *ZapLogger) WithLevel(level Level) Logger {
	logLevel := zapcore.Level(0)
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		panic(fmt.Errorf("Failed to set the zap log level from config. %s", level))
	}

	newLogger := l.Logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewCore(
			zapcore.NewJSONEncoder(l.config.EncoderConfig),
			zapcore.Lock(os.Stdout),
			logLevel,
		)
	}))

	return &ZapLogger{l.config, newLogger}
}

// WithField returns the logger at the supplied field.
func (l *ZapLogger) WithField(key string, value interface{}) Logger {
	newLogger := l.Logger.WithOptions(zap.Fields(zap.Any(key, value)))
	return &ZapLogger{l.config, newLogger}
}

// WithFields returns the logger at the supplied fields.
func (l *ZapLogger) WithFields(fields Fields) Logger {
	zapFields := []zap.Field{}
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	newLogger := l.Logger.WithOptions(zap.Fields(zapFields...))
	return &ZapLogger{l.config, newLogger}
}

// With returns the logger at the supplied fields.
func (l *ZapLogger) With(fields ...zap.Field) Logger {
	newLogger := l.Logger.With(fields...)
	return &ZapLogger{l.config, newLogger}
}

// WithError returns the logger with the supplied error.
func (l *ZapLogger) WithError(err error) Logger {
	newLogger := l.Logger.WithOptions(zap.Fields(zap.Error(err)))
	return &ZapLogger{l.config, newLogger}
}
