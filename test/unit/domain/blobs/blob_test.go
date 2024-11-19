package blobs

import (
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// BlobValidationTests struct encapsulates the test data and methods for blob validation
type BlobValidationTests struct {
	// TestData can be used for holding the Blob and CryptographicKey data
	validBlob    blobs.Blob
	invalidBlob  blobs.Blob
	invalidBlob2 blobs.Blob
}

// NewBlobValidationTests is a constructor to create a new instance of BlobValidationTests
func NewBlobValidationTests() *BlobValidationTests {
	// Create valid and invalid test data
	validBlob := blobs.Blob{
		ID:                  uuid.New().String(),
		UploadTime:          time.Now(),
		UserID:              uuid.New().String(),
		Name:                "test_blobs.txt",
		Size:                12345,
		Type:                "text",
		EncryptionAlgorithm: "AES",
		HashAlgorithm:       "SHA256",
		IsEncrypted:         true,
		IsSigned:            false,
		CryptographicKey:    keys.CryptographicKey{ID: uuid.New().String(), Type: "AES", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(24 * time.Hour), UserID: uuid.New().String()},
		KeyID:               uuid.New().String(),
	}

	invalidBlob := blobs.Blob{
		ID:                  "", // Invalid empty ID
		UploadTime:          time.Now(),
		UserID:              "invalid-uuid", // Invalid UserID
		Name:                "test_blobs.txt",
		Size:                -12345, // Invalid Size (negative)
		Type:                "text",
		EncryptionAlgorithm: "AES",
		HashAlgorithm:       "SHA256",
		IsEncrypted:         true,
		IsSigned:            false,
		CryptographicKey:    keys.CryptographicKey{ID: uuid.New().String(), Type: "AES", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(24 * time.Hour), UserID: uuid.New().String()},
		KeyID:               uuid.New().String(),
	}

	invalidBlob2 := blobs.Blob{
		ID:                  uuid.New().String(),
		UploadTime:          time.Now(),
		UserID:              uuid.New().String(),
		Name:                "", // Invalid empty Name
		Size:                12345,
		Type:                "text",
		EncryptionAlgorithm: "AES",
		HashAlgorithm:       "SHA256",
		IsEncrypted:         true,
		IsSigned:            false,
		CryptographicKey:    keys.CryptographicKey{ID: uuid.New().String(), Type: "AES", CreatedAt: time.Now(), ExpiresAt: time.Now().Add(24 * time.Hour), UserID: uuid.New().String()},
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
	assert.Contains(t, err.Error(), "Field: ID, Tag: required")
	assert.Contains(t, err.Error(), "Field: Size, Tag: min")
	assert.Contains(t, err.Error(), "Field: UserID, Tag: uuid4")
}

// TestBlobValidationEdgeCases tests validation edge cases for Blob
func (bt *BlobValidationTests) TestBlobValidationEdgeCases(t *testing.T) {
	// Validate the invalid Blob with empty Name
	err := bt.invalidBlob2.Validate()
	assert.NotNil(t, err, "Expected validation error for missing Name")
	assert.Contains(t, err.Error(), "Field: Name, Tag: required")
}

// TestBlobValidation is the entry point to run the Blob validation tests
func TestBlobValidation(t *testing.T) {
	// Create a new BlobValidationTests instance
	bt := NewBlobValidationTests()

	// Run each test method
	t.Run("TestBlobValidation", bt.TestBlobValidation)
	t.Run("TestBlobValidationEdgeCases", bt.TestBlobValidationEdgeCases)
}
