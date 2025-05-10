package v1

import (
	"crypto_vault_service/internal/domain/validators"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type UploadKeyRequest struct {
	Algorithm string `json:"algorithm" validate:"omitempty,oneof=AES RSA EC"`
	KeySize   uint32 `json:"key_size" validate:"omitempty,keySizeValidation"`
}

// Validate method for UploadKeyRequest struct
func (k *UploadKeyRequest) Validate() error {
	validate := validator.New()

	if err := validate.RegisterValidation("keySizeValidation", validators.KeySizeValidation); err != nil {
		return fmt.Errorf("failed to register custom validator: %w", err)
	}

	err := validate.Struct(k)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var messages []string
			for _, fieldErr := range validationErrors {
				messages = append(messages, fmt.Sprintf("Field: %s, Tag: %s", fieldErr.Field(), fieldErr.Tag()))
			}
			return fmt.Errorf("validation failed: %v", messages)
		}
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type InfoResponse struct {
	Message string `json:"message"`
}

type BlobMetaResponse struct {
	ID              string    `json:"id"`
	DateTimeCreated time.Time `json:"dateTimeCreated"`
	UserID          string    `json:"userID"`
	Name            string    `json:"name"`
	Size            int64     `json:"size"`
	Type            string    `json:"type"`
	EncryptionKeyID *string   `json:"encryptionKeyID"`
	SignKeyID       *string   `json:"signKeyID"`
}

type CryptoKeyMetaResponse struct {
	ID              string    `json:"id"`
	KeyPairID       string    `json:"keyPairID"`
	Algorithm       string    `json:"algorithm"`
	KeySize         uint32    `json:"keySize"`
	Type            string    `json:"type"`
	DateTimeCreated time.Time `json:"dateTimeCreated"`
	UserID          string    `json:"userID"`
}
