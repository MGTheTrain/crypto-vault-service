package logger

import (
	"crypto_vault_service/internal/infrastructure/settings"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

// Logger defines the logging interface
type Logger interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

var (
	// Singleton logger instance, shared across the application
	loggerInstance Logger
	loggerOnce     sync.Once  // Guarantees that the logger is created only once
	loggerMutex    sync.Mutex // Ensures thread-safe access to the logger instance
)

// GetLogger returns a singleton logger instance, shared across the application.
func GetLogger(settings *settings.LoggerSettings) (Logger, error) {
	// Ensure that the logger is created only once
	loggerOnce.Do(func() {
		// Create the logger based on the config
		logger, err := newLogger(settings)
		if err == nil {
			loggerInstance = logger
		}
	})

	// Lock access to loggerInstance to ensure thread safety when returning it
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	if loggerInstance != nil {
		return loggerInstance, nil
	}

	return nil, fmt.Errorf("failed to create logger")
}

// newLogger creates a logger based on the given configuration.
func newLogger(config *settings.LoggerSettings) (Logger, error) {
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level '%s': %w", config.LogLevel, err)
	}

	switch config.LogType {
	case "console":
		return NewConsoleLogger(level), nil
	case "file":
		return NewFileLogger(level, config.FilePath), nil
	default:
		return nil, fmt.Errorf("unsupported log type: %s", config.LogType)
	}
}
