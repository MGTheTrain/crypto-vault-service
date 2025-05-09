package settings

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type DatabaseSettings struct {
	Type string `mapstructure:"type" validate:"required"`
	DSN  string `mapstructure:"dsn" validate:"required"`
	Name string `mapstructure:"name" validate:"required"`
}

// Validate checks that all fields in DatabaseSettings are valid (non-empty in this case)
func (settings *DatabaseSettings) Validate() error {
	validate := validator.New()

	err := validate.Struct(settings)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
