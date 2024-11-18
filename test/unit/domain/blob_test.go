package model

import (
	"crypto_vault_service/internal/domain/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// BlobValidationTests struct encapsulates the test data and methods for blob validation
type BlobValidationTests struct {
	// TestData can be used for holding the Blob and CryptographicKey data
	validBlob    model.Blob
	invalidBlob  model.Blob
	invalidBlob2 model.Blob
}

// NewBlobValidationTests is a constructor to create a new instance of BlobValidationTests
func NewBlobValidationTests() *BlobValidationTests {
	// Create valid and invalid test data
	validBlob := model.Blob{
		BlobID:              uuid.New().String(),
		BlobStoragePath:     "/path/to/blob",
		UploadTime:          time.Now(),
		UserID:              uuid.New().String(),
		BlobName:            "test_blob.txt",
		BlobSize:            12345,
		BlobType:            "text",
		EncryptionAlgorithm: "AES",
		HashAlgorithm:       "SHA256",
		IsEncrypted:         true,
		IsSigned:            false,
		CryptographicKey:    model.CryptographicKey{KeyID: uuid.New().String(), KeyType: "AES", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(24 * time.Hour), UserID: uuid.New().String()},
		KeyID:               uuid.New().String(),
	}

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
		CryptographicKey:    model.CryptographicKey{KeyID: uuid.New().String(), KeyType: "AES", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(24 * time.Hour), UserID: uuid.New().String()},
		KeyID:               uuid.New().String(),
	}

	invalidBlob2 := model.Blob{
		BlobID:              uuid.New().String(),
		BlobStoragePath:     "/path/to/blob",
		UploadTime:          time.Now(),
		UserID:              uuid.New().String(),
		BlobName:            "", // Invalid empty BlobName
		BlobSize:            12345,
		BlobType:            "text",
		EncryptionAlgorithm: "AES",
		HashAlgorithm:       "SHA256",
		IsEncrypted:         true,
		IsSigned:            false,
		CryptographicKey:    model.CryptographicKey{KeyID: uuid.New().String(), KeyType: "AES", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(24 * time.Hour), UserID: uuid.New().String()},
		KeyID:               uuid.New().String(),
	}

	return &BlobValidationTests{
		validBlob:    validBlob,
		invalidBlob:  invalidBlob,
		invalidBlob2: invalidBlob2,
	}
}

// TestBlobValidation tests the Validator method for Blob
func (bt *BlobValidationTests) TestBlobValidation(t *testing.T) {
	// Validate the valid Blob
	err := bt.validBlob.Validate()
	assert.Nil(t, err, "Expected no validation errors for valid Blob")

	// Validate the invalid Blob
	err = bt.invalidBlob.Validate()
	assert.NotNil(t, err, "Expected validation errors for invalid Blob")
	assert.Contains(t, err.Error(), "Field: BlobID, Tag: required")
	assert.Contains(t, err.Error(), "Field: BlobSize, Tag: min")
	assert.Contains(t, err.Error(), "Field: UserID, Tag: uuid4")
}

// TestBlobValidationEdgeCases tests validation edge cases for Blob
func (bt *BlobValidationTests) TestBlobValidationEdgeCases(t *testing.T) {
	// Validate the invalid Blob with empty BlobName
	err := bt.invalidBlob2.Validate()
	assert.NotNil(t, err, "Expected validation error for missing BlobName")
	assert.Contains(t, err.Error(), "Field: BlobName, Tag: required")
}

// TestBlobValidation is the entry point to run the Blob validation tests
func TestBlobValidation(t *testing.T) {
	// Create a new BlobValidationTests instance
	bt := NewBlobValidationTests()

	// Run each test method
	t.Run("TestBlobValidation", bt.TestBlobValidation)
	t.Run("TestBlobValidationEdgeCases", bt.TestBlobValidationEdgeCases)
}
