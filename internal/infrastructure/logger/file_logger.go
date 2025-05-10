package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

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
