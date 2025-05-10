//go:build unit
// +build unit

package logger

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNewConsoleLoggerAndLogging(t *testing.T) {
	var buffer bytes.Buffer

	// Create console logger
	l := NewConsoleLogger(logrus.InfoLevel)
	require.NotNil(t, l, "expected console logger to be initialized")

	// Redirect logger output to buffer for inspection
	l.logger.SetOutput(&buffer)

	// These should not panic
	require.NotPanics(t, func() { l.Info("console info message") })
	require.NotPanics(t, func() { l.Warn("console warn message") })
	require.NotPanics(t, func() { l.Error("console error message") })

	// Assert log contents
	logOutput := buffer.String()
	require.Contains(t, logOutput, "console info message")
	require.Contains(t, logOutput, "console warn message")
	require.Contains(t, logOutput, "console error message")
}
