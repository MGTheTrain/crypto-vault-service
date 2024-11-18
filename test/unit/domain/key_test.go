package model

import (
	"crypto_vault_service/internal/domain/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCryptographicKeyValidation tests the Validator method for CryptographicKey
func TestCryptographicKeyValidation(t *testing.T) {
	// Valid CryptographicKey
	validKey := model.CryptographicKey{
		KeyID:     uuid.New().String(), // Valid UUID
		KeyType:   "AES",               // Valid KeyType
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24), // Valid ExpiresAt
		UserID:    uuid.New().String(),            // Valid UserID
	}

	// Validate the valid CryptographicKey
	err := validKey.Validate()
	assert.Nil(t, err, "Expected no validation errors for valid CryptographicKey")

	// Invalid CryptographicKey (empty KeyID, invalid KeyType, expired)
	invalidKey := model.CryptographicKey{
		KeyID:     "",            // Invalid empty KeyID
		KeyType:   "InvalidType", // Invalid KeyType
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(-time.Hour * 24), // Invalid ExpiresAt (before CreatedAt)
		UserID:    "invalid-user-id",               // Invalid UserID
	}

	// Validate the invalid CryptographicKey
	err = invalidKey.Validate()
	assert.NotNil(t, err, "Expected validation errors for invalid CryptographicKey")
	assert.Contains(t, err.Error(), "Field: KeyID, Tag: required")
	assert.Contains(t, err.Error(), "Field: KeyType, Tag: oneof")
	assert.Contains(t, err.Error(), "Field: ExpiresAt, Tag: gtefield")
}

// TestCryptographicKeyValidations tests the validation edge cases for CryptographicKey
func TestCryptographicKeyValidations(t *testing.T) {
	// Test missing UserID (should fail)
	invalidKey := model.CryptographicKey{
		KeyID:     uuid.New().String(), // Valid UUID
		KeyType:   "AES",               // Valid KeyType
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24), // Valid ExpiresAt
		UserID:    "",                             // Invalid empty UserID
	}

	err := invalidKey.Validate()
	assert.NotNil(t, err, "Expected validation error for missing UserID")
	assert.Contains(t, err.Error(), "Field: UserID, Tag: required")
}
