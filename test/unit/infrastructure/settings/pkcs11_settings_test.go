package settings

import (
	"crypto_vault_service/internal/infrastructure/settings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPKCS11SettingsValidation(t *testing.T) {
	tests := []struct {
		name          string
		settings      *settings.PKCS11Settings
		expectedError bool
	}{
		{
			name: "Valid Settings",
			settings: &settings.PKCS11Settings{
				ModulePath: "/path/to/module",
				SOPin:      "1234",
				UserPin:    "5678",
				SlotId:     "1",
			},
			expectedError: false,
		},
		{
			name: "Missing ModulePath",
			settings: &settings.PKCS11Settings{
				SOPin:   "1234",
				UserPin: "5678",
				SlotId:  "1",
			},
			expectedError: true,
		},
		{
			name: "Missing SOPin",
			settings: &settings.PKCS11Settings{
				ModulePath: "/path/to/module",
				UserPin:    "5678",
				SlotId:     "1",
			},
			expectedError: true,
		},
		{
			name: "Missing UserPin",
			settings: &settings.PKCS11Settings{
				ModulePath: "/path/to/module",
				SOPin:      "1234",
				SlotId:     "1",
			},
			expectedError: true,
		},
		{
			name: "Missing SlotId",
			settings: &settings.PKCS11Settings{
				ModulePath: "/path/to/module",
				SOPin:      "1234",
				UserPin:    "5678",
			},
			expectedError: true,
		},
		{
			name:          "All Fields Missing",
			settings:      &settings.PKCS11Settings{},
			expectedError: true,
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
