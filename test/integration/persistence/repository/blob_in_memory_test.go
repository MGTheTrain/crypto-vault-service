package repository

import (
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
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

	// Create a valid CryptoKey object
	cryptographicKey := keys.CryptoKeyMeta{
		ID:        uuid.New().String(),            // Generate valid UUID for ID
		Type:      "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &blobs.BlobMeta{
		ID:          uuid.New().String(), // Generate valid UUID
		UploadTime:  time.Now(),
		UserID:      uuid.New().String(), // Generate valid UUID for UserID
		Name:        "test-blob",
		Size:        1024,
		Type:        "text",
		IsEncrypted: false,
		IsSigned:    true,
		CryptoKey:   cryptographicKey,    // Set the CryptoKey
		KeyID:       cryptographicKey.ID, // Ensure ID is set
	}

	// Call the Create method
	err := ctx.BlobRepo.Create(blob)
	assert.NoError(t, err, "Create should not return an error")

	// Verify the blob is created and exists in DB
	var createdBlob blobs.BlobMeta
	err = ctx.DB.First(&createdBlob, "id = ?", blob.ID).Error
	assert.NoError(t, err, "Failed to find created blob")
	assert.Equal(t, blob.ID, createdBlob.ID, "ID should match")
	assert.Equal(t, blob.Name, createdBlob.Name, "Name should match")
}

// TestBlobRepository_GetById tests the GetById method of BlobRepositoryImpl
func TestBlobRepository_GetById(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptoKey object
	cryptographicKey := keys.CryptoKeyMeta{
		ID:        uuid.New().String(),            // Generate valid UUID for ID
		Type:      "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &blobs.BlobMeta{
		ID:          uuid.New().String(), // Generate valid UUID
		UploadTime:  time.Now(),
		UserID:      cryptographicKey.UserID, // Link to valid UserID from CryptoKey
		Name:        "test-blob",
		Size:        1024,
		Type:        "text",
		IsEncrypted: false,
		IsSigned:    true,
		CryptoKey:   cryptographicKey,    // Set the CryptoKey
		KeyID:       cryptographicKey.ID, // Ensure ID is set
	}

	// Create the blob in DB
	err := ctx.BlobRepo.Create(blob)
	assert.NoError(t, err, "Create should not return an error")

	// Get the blob by ID
	fetchedBlob, err := ctx.BlobRepo.GetById(blob.ID)
	assert.NoError(t, err, "GetById should not return an error")
	assert.NotNil(t, fetchedBlob, "Fetched blob should not be nil")
	assert.Equal(t, blob.ID, fetchedBlob.ID, "ID should match")
}

// TestBlobRepository_UpdateById tests the UpdateById method of BlobRepositoryImpl
func TestBlobRepository_UpdateById(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptoKey object
	cryptographicKey := keys.CryptoKeyMeta{
		ID:        uuid.New().String(),            // Generate valid UUID for ID
		Type:      "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &blobs.BlobMeta{
		ID:          uuid.New().String(), // Generate valid UUID
		UploadTime:  time.Now(),
		UserID:      uuid.New().String(), // Generate valid UUID for UserID
		Name:        "test-blob",
		Size:        1024,
		Type:        "text",
		IsEncrypted: false,
		IsSigned:    true,
		CryptoKey:   cryptographicKey,    // Set the CryptoKey
		KeyID:       cryptographicKey.ID, // Ensure ID is set
	}

	// Create the blob in DB
	err := ctx.BlobRepo.Create(blob)
	assert.NoError(t, err)

	// Update the blob's name
	blob.Name = "updated-blob-name"
	err = ctx.BlobRepo.UpdateById(blob)
	assert.NoError(t, err)

	// Verify the blob is updated
	var updatedBlob blobs.BlobMeta
	err = ctx.DB.First(&updatedBlob, "id = ?", blob.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, "updated-blob-name", updatedBlob.Name, "Name should be updated")
}

// TestBlobRepository_DeleteById tests the DeleteById method of BlobRepositoryImpl
func TestBlobRepository_DeleteById(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptoKey object
	cryptographicKey := keys.CryptoKeyMeta{
		ID:        uuid.New().String(),            // Generate valid UUID for ID
		Type:      "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &blobs.BlobMeta{
		ID:          uuid.New().String(), // Generate valid UUID
		UploadTime:  time.Now(),
		UserID:      uuid.New().String(), // Generate valid UUID for UserID
		Name:        "test-blob",
		Size:        1024,
		Type:        "text",
		IsEncrypted: false,
		IsSigned:    true,
		CryptoKey:   cryptographicKey,    // Set the CryptoKey
		KeyID:       cryptographicKey.ID, // Ensure ID is set
	}

	// Create the blob in DB
	err := ctx.BlobRepo.Create(blob)
	assert.NoError(t, err)

	// Delete the blob
	err = ctx.BlobRepo.DeleteById(blob.ID)
	assert.NoError(t, err)

	// Verify the blob is deleted
	var deletedBlob blobs.BlobMeta
	err = ctx.DB.First(&deletedBlob, "id = ?", blob.ID).Error
	assert.Error(t, err, "Blob should be deleted")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Error should be 'record not found'")
}
