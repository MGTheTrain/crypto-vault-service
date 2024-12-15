package logger

import (
	"crypto_vault_service/internal/infrastructure/settings"
	"fmt"

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

// ConsoleLogger is an implementation of Logger that logs to the console.
type ConsoleLogger struct {
	logger *logrus.Logger
}

// NewConsoleLogger creates a new console logger with the specified log level.
func NewConsoleLogger(level logrus.Level) *ConsoleLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(level)
	return &ConsoleLogger{logger: logger}
}

func (l *ConsoleLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *ConsoleLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *ConsoleLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *ConsoleLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *ConsoleLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

// FileLogger is an implementation of Logger that logs to a file.
type FileLogger struct {
	logger *logrus.Logger
}

// NewFileLogger creates a new file logger with the specified log level and file path.
func NewFileLogger(level logrus.Level, filePath string) *FileLogger {
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
	return &FileLogger{logger: logger}
}

func (l *FileLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *FileLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *FileLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *FileLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *FileLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

// LoggerFactory is a factory that creates different types of loggers.
type LoggerFactory struct{}

// NewLogger creates a logger based on the given configuration.
func (f *LoggerFactory) NewLogger(config *settings.LoggerSettings) (Logger, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, err
	}

	if config.LogType == "console" {
		return NewConsoleLogger(level), nil
	} else if config.LogType == "file" {
		return NewFileLogger(level, config.FilePath), nil
	}

	return nil, fmt.Errorf("unsupported log type: %s", config.LogType)
}
