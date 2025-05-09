package logger

import (
	"crypto_vault_service/internal/infrastructure/settings"
	"fmt"
	"sync"

	"github.com/natefinch/lumberjack"
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

// consoleLogger is an implementation of Logger that logs to the console.
type consoleLogger struct {
	logger *logrus.Logger
}

// NewConsoleLogger creates a new console logger with the specified log level.
func NewConsoleLogger(level logrus.Level) *consoleLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(level)
	return &consoleLogger{logger: logger}
}

func (l *consoleLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *consoleLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *consoleLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *consoleLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *consoleLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

// fileLogger is an implementation of Logger that logs to a file.
type fileLogger struct {
	logger *logrus.Logger
}

// NewFileLogger creates a new file logger with the specified log level and file path.
func NewFileLogger(level logrus.Level, filePath string) *fileLogger {
	logger := logrus.New()

	logger.SetOutput(&lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // number of days to retain logs
	})

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(level)
	return &fileLogger{logger: logger}
}

func (l *fileLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *fileLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *fileLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *fileLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *fileLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
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
