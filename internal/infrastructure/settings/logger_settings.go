package settings

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type LoggerSettings struct {
	LogLevel string `validate:"required"`
	LogType  string `validate:"required,oneof=console file"` // Log type must be either "console" or "file"
	FilePath string `validate:"required_if=LogType file"`    // File path is required only if LogType is "file"
}

// Validate checks that all fields in PKCS11Settings are valid (non-empty in this case)
func (settings *LoggerSettings) Validate() error {
	validate := validator.New()

	err := validate.Struct(settings)
	if err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}
	return nil
}
