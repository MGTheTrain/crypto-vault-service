package v1

import (
	"crypto_vault_service/internal/domain/validators"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type UploadKeyRequestDto struct {
	Algorithm string `json:"algorithm" validate:"omitempty,oneof=AES RSA EC"`
	KeySize   uint   `json:"key_size" validate:"omitempty,keySizeValidation"`
}

// Validate method for UploadKeyRequestDto struct
func (k *UploadKeyRequestDto) Validate() error {
	validate := validator.New()

	err := validate.RegisterValidation("keySizeValidation", validators.KeySizeValidation)
	if err != nil {
		return fmt.Errorf("failed to register custom validator: %v", err)
	}
	err = validate.Struct(k)
	if err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field: %s, Tag: %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("validation failed: %v", validationErrors)
	}
	return nil
}

type ErrorResponseDto struct {
	Message string `json:"message"`
}

type InfoResponseDto struct {
	Message string `json:"message"`
}

type BlobMetaResponseDto struct {
	ID              string    `json:"id"`
	DateTimeCreated time.Time `json:"dateTimeCreated"`
	UserID          string    `json:"userID"`
	Name            string    `json:"name"`
	Size            int64     `json:"size"`
	Type            string    `json:"type"`
	EncryptionKeyID *string   `json:"encryptionKeyID"`
	SignKeyID       *string   `json:"signKeyID"`
}

type CryptoKeyMetaResponseDto struct {
	ID              string    `json:"id"`
	KeyPairID       string    `json:"keyPairID"`
	Algorithm       string    `json:"algorithm"`
	KeySize         uint      `json:"keySize"`
	Type            string    `json:"type"`
	DateTimeCreated time.Time `json:"dateTimeCreated"`
	UserID          string    `json:"userID"`
}
