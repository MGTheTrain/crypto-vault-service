package settings

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type BlobConnectorSettings struct {
	ConnectionString string `mapstructure:"connectionString" validate:"required"`
	ContainerName    string `mapstructure:"containerName" validate:"required"`
}

// Validate checks that all fields in BlobConnectorSettings are valid (non-empty in this case)
func (settings *BlobConnectorSettings) Validate() error {
	validate := validator.New()

	err := validate.Struct(settings)
	if err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}
	return nil
}
