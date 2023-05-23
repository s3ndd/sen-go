package log

import (
	"fmt"
	"github.com/s3ndd/sen-go/config"
	"sync/atomic"
)

// sharedLogger holds the global Logger
var sharedLogger atomic.Value

// Global returns the global logger
func defaultConfig() *Config {
	return &Config{
		LogLevel:  config.String("LOG_LEVEL", "INFO"),
		LogFormat: config.String("LOG_FORMAT", "json"),
		Program:   config.String("SOURCE_PROGRAM", "unknown"),
		Env:       config.String("ENV", "dev"),
	}
}

// SetLogger sets the default global logger
func SetLogger(logger Logger) error {
	if sharedLogger.Load() != nil {
		return fmt.Errorf("Shared logger exists, cannot be modified")
	}
	sharedLogger.Store(logger)
	return nil
}

// SetGlobalFields sets fields on the log entry that should be global to all requests
func SetGlobalFields(fields Fields) {
	_ = SetLogger(Global().WithFields(fields))
}

// Global returns the global log entry.
// Use SetGlobalFields to configure this entry with global fields.
// Use ForRequest if you want a log entry pre-configured with relevant request metadata
func Global() Logger {
	if sharedLogger.Load() == nil {
		_ = SetLogger(NewLogger(defaultConfig()))
	}

	return sharedLogger.Load().(Logger)
}
