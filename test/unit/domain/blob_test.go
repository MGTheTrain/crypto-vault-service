package model

import (
	"crypto_vault_service/internal/domain/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestBlobValidation tests the Validator method for Blob
func TestBlobValidation(t *testing.T) {
	// Valid Blob
	validBlob := model.Blob{
		BlobID:              uuid.New().String(), // Valid UUID
		BlobStoragePath:     "/path/to/blob",
		UploadTime:          time.Now(),
		UserID:              uuid.New().String(), // Valid UUID
		BlobName:            "test_blob.txt",
		BlobSize:            12345,
		BlobType:            "text",
		EncryptionAlgorithm: "AES",
		HashAlgorithm:       "SHA256",
		IsEncrypted:         true,
		IsSigned:            false,
		CryptographicKey:    model.CryptographicKey{KeyID: "abc123", KeyType: "AES", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(24 * time.Hour), UserID: "e2b073b5-3e23-4fbd-b44b-c607b04c9c3e"},
		KeyID:               "abc123",
	}

	// Validate the valid Blob
	err := validBlob.Validate()
	assert.Nil(t, err, "Expected no validation errors for valid Blob")

	// Invalid Blob (empty BlobID, invalid BlobSize)
	invalidBlob := model.Blob{
		BlobID:              "", // Invalid empty BlobID
		BlobStoragePath:     "/path/to/blob",
		UploadTime:          time.Now(),
		UserID:              "invalid-uuid", // Invalid UserID
		BlobName:            "test_blob.txt",
		BlobSize:            -12345, // Invalid BlobSize (negative)
		BlobType:            "text",
		EncryptionAlgorithm: "AES",
		HashAlgorithm:       "SHA256",
		IsEncrypted:         true,
		IsSigned:            false,
		CryptographicKey:    model.CryptographicKey{KeyID: "abc123", KeyType: "AES", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(24 * time.Hour), UserID: "e2b073b5-3e23-4fbd-b44b-c607b04c9c3e"},
		KeyID:               "abc123",
	}

	// Validate the invalid Blob
	err = invalidBlob.Validate()
	assert.NotNil(t, err, "Expected validation errors for invalid Blob")
	assert.Contains(t, err.Error(), "Field: BlobID, Tag: required")
	assert.Contains(t, err.Error(), "Field: BlobSize, Tag: min")
	assert.Contains(t, err.Error(), "Field: UserID, Tag: uuid4")
}

// TestBlobValidationEdgeCases tests validation edge cases for Blob
func TestBlobValidationEdgeCases(t *testing.T) {
	// Test missing BlobName (should fail)
	invalidBlob := model.Blob{
		BlobID:              uuid.New().String(), // Valid UUID
		BlobStoragePath:     "/path/to/blob",
		UploadTime:          time.Now(),
		UserID:              uuid.New().String(), // Valid UUID
		BlobName:            "",                  // Invalid empty BlobName
		BlobSize:            12345,
		BlobType:            "text",
		EncryptionAlgorithm: "AES",
		HashAlgorithm:       "SHA256",
		IsEncrypted:         true,
		IsSigned:            false,
		CryptographicKey:    model.CryptographicKey{KeyID: "abc123", KeyType: "AES", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(24 * time.Hour), UserID: "e2b073b5-3e23-4fbd-b44b-c607b04c9c3e"},
		KeyID:               "abc123",
	}

	err := invalidBlob.Validate()
	assert.NotNil(t, err, "Expected validation error for missing BlobName")
	assert.Contains(t, err.Error(), "Field: BlobName, Tag: required")
}
