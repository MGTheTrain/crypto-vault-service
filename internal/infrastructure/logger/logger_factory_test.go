//go:build unit
// +build unit

package logger

import (
	"crypto_vault_service/internal/infrastructure/settings"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// ResetLogger is a helper to reset the singleton logger instance between tests
func resetLoggerSingleton() {
	loggerInstance = nil
	loggerOnce = sync.Once{}
}

func TestGetLogger_ConsoleLogger(t *testing.T) {
	resetLoggerSingleton()

	cfg := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
	}

	logger, err := GetLogger(cfg)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Confirm singleton behavior
	logger2, err2 := GetLogger(cfg)
	require.NoError(t, err2)
	require.Same(t, logger, logger2, "expected singleton logger instance")
}

func TestGetLogger_FileLogger(t *testing.T) {
	resetLoggerSingleton()

	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "app.log")

	cfg := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "file",
		FilePath: logPath,
	}

	logger, err := GetLogger(cfg)
	require.NoError(t, err)
	require.NotNil(t, logger)

	// Log a message to force file creation
	logger.Info("File logger test")

	_, statErr := os.Stat(logPath)
	require.NoError(t, statErr, "log file should be created")
}

func TestGetLogger_InvalidLogLevel(t *testing.T) {
	resetLoggerSingleton()

	cfg := &settings.LoggerSettings{
		LogLevel: "invalid-level",
		LogType:  "console",
	}

	logger, err := GetLogger(cfg)
	require.Error(t, err)
	require.Nil(t, logger)
}

func TestGetLogger_UnsupportedLogType(t *testing.T) {
	resetLoggerSingleton()

	cfg := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "unknown",
	}

	logger, err := GetLogger(cfg)
	require.Error(t, err)
	require.Nil(t, logger)
}
