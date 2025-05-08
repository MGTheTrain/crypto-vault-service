//go:build unit
// +build unit

package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerSettingsValidation(t *testing.T) {
	tests := []struct {
		name          string
		settings      *LoggerSettings
		expectedError bool
	}{
		{
			name: "Valid Settings",
			settings: &LoggerSettings{
				LogLevel: "info",
				LogType:  "console",
				FilePath: "",
			},
			expectedError: false,
		},
		{
			name: "Valid Settings with File",
			settings: &LoggerSettings{
				LogLevel: "info",
				LogType:  "file",
				FilePath: "/path/to/log/file",
			},
			expectedError: false,
		},
		{
			name: "Missing LogLevel",
			settings: &LoggerSettings{
				LogType:  "console",
				FilePath: "",
			},
			expectedError: true,
		},
		{
			name: "Missing LogType",
			settings: &LoggerSettings{
				LogLevel: "info",
				FilePath: "",
			},
			expectedError: true,
		},
		{
			name: "Invalid LogType",
			settings: &LoggerSettings{
				LogLevel: "info",
				LogType:  "invalid", // Invalid log type
				FilePath: "",
			},
			expectedError: true,
		},
		{
			name: "Missing FilePath when LogType is file",
			settings: &LoggerSettings{
				LogLevel: "info",
				LogType:  "file",
				FilePath: "", // Missing file path when LogType is "file"
			},
			expectedError: true,
		},
		{
			name: "FilePath provided when LogType is console",
			settings: &LoggerSettings{
				LogLevel: "info",
				LogType:  "console",
				FilePath: "/path/to/log/file", // FilePath should not be used for console log type
			},
			expectedError: false,
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.settings.Validate()

			if tt.expectedError {
				assert.Errorf(t, err, "expected an error, got nil")
			} else {
				assert.NoError(t, err, "expected no error, got: %v", err)
			}
		})
	}
}
