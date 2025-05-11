package logger

import "github.com/sirupsen/logrus"

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

// Info logs an informational message to the console using the underlying logger.
func (l *ConsoleLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Warn logs a warning message to the console using the underlying logger.
func (l *ConsoleLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Error logs an error message to the console using the underlying logger.
func (l *ConsoleLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Fatal logs a fatal message to the console using the underlying logger and then exits the program.
func (l *ConsoleLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Panic logs a panic message to the console using the underlying logger and then panics.
func (l *ConsoleLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}
