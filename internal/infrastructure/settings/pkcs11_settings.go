package settings

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// PKCS11Settings holds the configuration settings required to interact with a PKCS#11 module
type PKCS11Settings struct {
	ModulePath string `validate:"required"`
	SOPin      string `validate:"required"`
	UserPin    string `validate:"required"`
	SlotId     string `validate:"required"`
}

// Validate checks that all fields in PKCS11Settings are valid (non-empty in this case)
func (settings *PKCS11Settings) Validate() error {
	validate := validator.New()

	err := validate.Struct(settings)
	if err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}
	return nil
}
