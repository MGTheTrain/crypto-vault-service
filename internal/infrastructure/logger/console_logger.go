package logger

import "github.com/sirupsen/logrus"

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
