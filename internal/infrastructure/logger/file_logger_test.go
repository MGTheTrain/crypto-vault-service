//go:build unit
// +build unit

package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNewFileLoggerAndLogging(t *testing.T) {
	// Create a temporary file path (won't actually open the file now)
	tmpDir := t.TempDir()
	logFilePath := filepath.Join(tmpDir, "test.log")

	// Initialize file logger
	l := NewFileLogger(logrus.InfoLevel, logFilePath)
	require.NotNil(t, l, "expected logger to be initialized")

	// These should not panic
	require.NotPanics(t, func() { l.Info("This is an info message") })
	require.NotPanics(t, func() { l.Warn("This is a warning message") })
	require.NotPanics(t, func() { l.Error("This is an error message") })

	// Optional: Check if the file was created
	_, err := os.Stat(logFilePath)
	require.NoError(t, err, "expected log file to be created")
}
