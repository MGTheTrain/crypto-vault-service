package keys

import (
	"crypto_vault_service/internal/domain/validators"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// CryptoKeyMeta represents the encryption key entity
type CryptoKeyMeta struct {
	ID              string    `gorm:"primaryKey" validate:"required,uuid4"`
	KeyPairID       string    `gorm:"index" validate:"required,uuid4"`
	Algorithm       string    `validate:"omitempty,oneof=AES RSA EC"`
	KeySize         uint      `json:"key_size" validate:"omitempty,keySizeValidation"`
	Type            string    `validate:"omitempty,oneof=private public symmetric"`
	DateTimeCreated time.Time `validate:"required"`
	UserID          string    `gorm:"index" validate:"required,uuid4"`
}

// Validate method for CryptoKeyMeta struct
func (k *CryptoKeyMeta) Validate() error {
	validate := validator.New()

	err := validate.RegisterValidation("keySizeValidation", validators.KeySizeValidation)
	if err != nil {
		return fmt.Errorf("failed to register custom validator: %w", err)
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
