package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

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

// Info logs an informational message using the underlying logger.
func (l *FileLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Warn logs a warning message using the underlying logger.
func (l *FileLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Error logs an error message using the underlying logger.
func (l *FileLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Fatal logs a fatal message using the underlying logger and then exits the program.
func (l *FileLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Panic logs a panic message using the underlying logger and then panics.
func (l *FileLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}
