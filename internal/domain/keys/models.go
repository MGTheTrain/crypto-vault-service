package keys

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// CryptoKeyMeta represents the encryption key entity
type CryptoKeyMeta struct {
	ID              string    `gorm:"primaryKey" validate:"required,uuid4"` // ID is required and must be a valid UUID
	Type            string    `validate:"required,oneof=AES RSA ECDSA"`     // Type is required and must be one of the listed types
	DateTimeCreated time.Time `validate:"required"`                         // DateTimeCreated is required
	UserID          string    `gorm:"index" validate:"required,uuid4"`      // UserID is required and must be a valid UUID
}

// Validate for validating CryptoKeyMeta struct
func (k *CryptoKeyMeta) Validate() error {
	// Initialize the validator
	validate := validator.New()

	// Validate the struct
	err := validate.Struct(k)
	if err != nil {
		// If validation fails, return a formatted error
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field: %s, Tag: %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("Validation failed: %v", validationErrors)
	}
	return nil // Return nil if validation passes
}
