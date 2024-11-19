package keys

import (
	"crypto_vault_service/internal/domain/keys"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCryptoKeyValidation tests the Validator method for CryptoKey
func TestCryptoKeyValidation(t *testing.T) {
	// Valid CryptoKey
	validKey := keys.CryptoKeyMeta{
		ID:        uuid.New().String(), // Valid UUID
		Type:      "AES",               // Valid Type
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24), // Valid ExpiresAt
		UserID:    uuid.New().String(),            // Valid UserID
	}

	// Validate the valid CryptoKey
	err := validKey.Validate()
	assert.Nil(t, err, "Expected no validation errors for valid CryptoKey")

	// Invalid CryptoKey (empty ID, invalid Type, expired)
	invalidKey := keys.CryptoKeyMeta{
		ID:        "",            // Invalid empty ID
		Type:      "InvalidType", // Invalid Type
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(-time.Hour * 24), // Invalid ExpiresAt (before CreatedAt)
		UserID:    "invalid-user-id",               // Invalid UserID
	}

	// Validate the invalid CryptoKey
	err = invalidKey.Validate()
	assert.NotNil(t, err, "Expected validation errors for invalid CryptoKey")
	assert.Contains(t, err.Error(), "Field: ID, Tag: required")
	assert.Contains(t, err.Error(), "Field: Type, Tag: oneof")
	assert.Contains(t, err.Error(), "Field: ExpiresAt, Tag: gtefield")
}

// TestCryptoKeyValidations tests the validation edge cases for CryptoKey
func TestCryptoKeyValidations(t *testing.T) {
	// Test missing UserID (should fail)
	invalidKey := keys.CryptoKeyMeta{
		ID:        uuid.New().String(), // Valid UUID
		Type:      "AES",               // Valid Type
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24), // Valid ExpiresAt
		UserID:    "",                             // Invalid empty UserID
	}

	err := invalidKey.Validate()
	assert.NotNil(t, err, "Expected validation error for missing UserID")
	assert.Contains(t, err.Error(), "Field: UserID, Tag: required")
}
