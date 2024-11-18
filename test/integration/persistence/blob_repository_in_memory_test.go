package repository

import (
	"crypto_vault_service/internal/domain/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestBlobRepository_Create tests the Create method of BlobRepositoryImpl
func TestBlobRepository_Create(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptographicKey object
	cryptographicKey := model.CryptographicKey{
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &model.Blob{
		BlobID:           uuid.New().String(), // Generate valid UUID
		UploadTime:       time.Now(),
		UserID:           uuid.New().String(), // Generate valid UUID for UserID
		BlobName:         "test-blob",
		BlobSize:         1024,
		BlobType:         "text",
		IsEncrypted:      false,
		IsSigned:         true,
		CryptographicKey: cryptographicKey,       // Set the CryptographicKey
		KeyID:            cryptographicKey.KeyID, // Ensure KeyID is set
	}

	// Call the Create method
	err := ctx.BlobRepo.Create(blob)
	assert.NoError(t, err, "Create should not return an error")

	// Verify the blob is created and exists in DB
	var createdBlob model.Blob
	err = ctx.DB.First(&createdBlob, "blob_id = ?", blob.BlobID).Error
	assert.NoError(t, err, "Failed to find created blob")
	assert.Equal(t, blob.BlobID, createdBlob.BlobID, "BlobID should match")
	assert.Equal(t, blob.BlobName, createdBlob.BlobName, "BlobName should match")
}

// TestBlobRepository_GetById tests the GetById method of BlobRepositoryImpl
func TestBlobRepository_GetById(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptographicKey object
	cryptographicKey := model.CryptographicKey{
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &model.Blob{
		BlobID:           uuid.New().String(), // Generate valid UUID
		UploadTime:       time.Now(),
		UserID:           cryptographicKey.UserID, // Link to valid UserID from CryptographicKey
		BlobName:         "test-blob",
		BlobSize:         1024,
		BlobType:         "text",
		IsEncrypted:      false,
		IsSigned:         true,
		CryptographicKey: cryptographicKey,       // Set the CryptographicKey
		KeyID:            cryptographicKey.KeyID, // Ensure KeyID is set
	}

	// Create the blob in DB
	err := ctx.BlobRepo.Create(blob)
	assert.NoError(t, err, "Create should not return an error")

	// Get the blob by ID
	fetchedBlob, err := ctx.BlobRepo.GetById(blob.BlobID)
	assert.NoError(t, err, "GetById should not return an error")
	assert.NotNil(t, fetchedBlob, "Fetched blob should not be nil")
	assert.Equal(t, blob.BlobID, fetchedBlob.BlobID, "BlobID should match")
}

// TestBlobRepository_UpdateById tests the UpdateById method of BlobRepositoryImpl
func TestBlobRepository_UpdateById(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptographicKey object
	cryptographicKey := model.CryptographicKey{
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &model.Blob{
		BlobID:           uuid.New().String(), // Generate valid UUID
		UploadTime:       time.Now(),
		UserID:           uuid.New().String(), // Generate valid UUID for UserID
		BlobName:         "test-blob",
		BlobSize:         1024,
		BlobType:         "text",
		IsEncrypted:      false,
		IsSigned:         true,
		CryptographicKey: cryptographicKey,       // Set the CryptographicKey
		KeyID:            cryptographicKey.KeyID, // Ensure KeyID is set
	}

	// Create the blob in DB
	err := ctx.BlobRepo.Create(blob)
	assert.NoError(t, err)

	// Update the blob's name
	blob.BlobName = "updated-blob-name"
	err = ctx.BlobRepo.UpdateById(blob)
	assert.NoError(t, err)

	// Verify the blob is updated
	var updatedBlob model.Blob
	err = ctx.DB.First(&updatedBlob, "blob_id = ?", blob.BlobID).Error
	assert.NoError(t, err)
	assert.Equal(t, "updated-blob-name", updatedBlob.BlobName, "BlobName should be updated")
}

// TestBlobRepository_DeleteById tests the DeleteById method of BlobRepositoryImpl
func TestBlobRepository_DeleteById(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptographicKey object
	cryptographicKey := model.CryptographicKey{
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &model.Blob{
		BlobID:           uuid.New().String(), // Generate valid UUID
		UploadTime:       time.Now(),
		UserID:           uuid.New().String(), // Generate valid UUID for UserID
		BlobName:         "test-blob",
		BlobSize:         1024,
		BlobType:         "text",
		IsEncrypted:      false,
		IsSigned:         true,
		CryptographicKey: cryptographicKey,       // Set the CryptographicKey
		KeyID:            cryptographicKey.KeyID, // Ensure KeyID is set
	}

	// Create the blob in DB
	err := ctx.BlobRepo.Create(blob)
	assert.NoError(t, err)

	// Delete the blob
	err = ctx.BlobRepo.DeleteById(blob.BlobID)
	assert.NoError(t, err)

	// Verify the blob is deleted
	var deletedBlob model.Blob
	err = ctx.DB.First(&deletedBlob, "blob_id = ?", blob.BlobID).Error
	assert.Error(t, err, "Blob should be deleted")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Error should be 'record not found'")
}
