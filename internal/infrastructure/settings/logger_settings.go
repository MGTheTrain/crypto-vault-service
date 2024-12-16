package settings

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type LoggerSettings struct {
	LogLevel string `validate:"required,oneof=info debug error warning critical"`
	LogType  string `validate:"required,oneof=console file"`
	FilePath string `validate:"required_if=LogType file"` // File path is required only if LogType is "file"
}

// Validate checks that all fields in LoggerSettings are valid (non-empty in this case)
func (settings *LoggerSettings) Validate() error {
	validate := validator.New()

	err := validate.Struct(settings)
	if err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}
	return nil
}
