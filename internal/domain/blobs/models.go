package blobs

import (
	"crypto_vault_service/internal/domain/keys"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// BlobMeta represents metadata on the actual blob metadata being stored
type BlobMeta struct {
	ID                  string             `gorm:"primaryKey" validate:"required,uuid4"`  // ID is required and must be a valid UUID
	DateTimeCreated     time.Time          `validate:"required"`                          // DateTimeCreated is required
	UserID              string             `validate:"required,uuid4"`                    // UserID is required and must be a valid UUID
	Name                string             `validate:"required,min=1,max=255"`            // Name is required, and its length must be between 1 and 255 characters
	Size                int64              `validate:"required,min=1"`                    // Size must be greater than 0
	Type                string             `validate:"required,min=1,max=50"`             // Type is required, and its length must be between 1 and 50 characters
	EncryptionAlgorithm string             `validate:"omitempty,oneof=AES RSA ECDSA"`     // EncryptionAlgorithm is optional and if set must be one of the listed algorithms
	HashAlgorithm       string             `validate:"omitempty,oneof=SHA256 SHA512 MD5"` // HashAlgorithm is optional and if set must be one of the listed algorithms
	CryptoKey           keys.CryptoKeyMeta `gorm:"foreignKey:KeyID" validate:"omitempty"` // CryptoKey is optional
	KeyID               string             `validate:"omitempty,uuid4"`                   // KeyID is optional and if set must be a valid UUID
}

// Validate for validating BlobMeta struct
func (b *BlobMeta) Validate() error {

	validate := validator.New()

	err := validate.Struct(b)
	if err != nil {

		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field: %s, Tag: %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("Validation failed: %v", validationErrors)
	}
	return nil
}
