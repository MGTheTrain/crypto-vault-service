package model

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// CryptographicKey represents the encryption key entity
type CryptographicKey struct {
	KeyID     string    `gorm:"primaryKey" validate:"required,uuid4"` // KeyID is required and must be a valid UUID
	KeyType   string    `validate:"required,oneof=AES RSA ECDSA"`     // KeyType is required and must be one of the listed types
	CreatedAt time.Time `validate:"required"`                         // CreatedAt is required
	ExpiresAt time.Time `validate:"required,gtefield=CreatedAt"`      // ExpiresAt is required and must be after CreatedAt
	UserID    string    `gorm:"index" validate:"required,uuid4"`      // UserID is required and must be a valid UUID
}

// Validate for validating CryptographicKey struct
func (k *CryptographicKey) Validate() error {
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
