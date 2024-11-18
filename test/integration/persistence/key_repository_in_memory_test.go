package repository

import (
	"crypto_vault_service/internal/domain/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestCryptographicKeyRepository_Create tests the Create method of CryptographicKeyRepositoryImpl
func TestCryptographicKeyRepository_Create(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptographicKey object
	cryptographicKey := &model.CryptographicKey{
		ID:        uuid.New().String(),            // Generate valid UUID for ID
		Type:      "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err := ctx.CryptoKeyRepo.Create(cryptographicKey)
	assert.NoError(t, err, "Create should not return an error")

	// Verify the cryptographic key is created and exists in DB
	var createdKey model.CryptographicKey
	err = ctx.DB.First(&createdKey, "id = ?", cryptographicKey.ID).Error
	assert.NoError(t, err, "Failed to find created cryptographic key")
	assert.Equal(t, cryptographicKey.ID, createdKey.ID, "ID should match")
	assert.Equal(t, cryptographicKey.Type, createdKey.Type, "Type should match")
}

// TestCryptographicKeyRepository_GetByID tests the GetByID method of CryptographicKeyRepositoryImpl
func TestCryptographicKeyRepository_GetByID(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptographicKey object
	cryptographicKey := &model.CryptographicKey{
		ID:        uuid.New().String(),            // Generate valid UUID for ID
		Type:      "RSA",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err := ctx.CryptoKeyRepo.Create(cryptographicKey)
	assert.NoError(t, err, "Create should not return an error")

	// Get the cryptographic key by ID
	fetchedKey, err := ctx.CryptoKeyRepo.GetByID(cryptographicKey.ID)
	assert.NoError(t, err, "GetByID should not return an error")
	assert.NotNil(t, fetchedKey, "Fetched cryptographic key should not be nil")
	assert.Equal(t, cryptographicKey.ID, fetchedKey.ID, "ID should match")
}

// TestCryptographicKeyRepository_UpdateByID tests the UpdateByID method of CryptographicKeyRepositoryImpl
func TestCryptographicKeyRepository_UpdateByID(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptographicKey object
	cryptographicKey := &model.CryptographicKey{
		ID:        uuid.New().String(),            // Generate valid UUID for ID
		Type:      "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err := ctx.CryptoKeyRepo.Create(cryptographicKey)
	assert.NoError(t, err, "Create should not return an error")

	// Update the cryptographic key's type
	cryptographicKey.Type = "ECDSA"
	err = ctx.CryptoKeyRepo.UpdateByID(cryptographicKey)
	assert.NoError(t, err, "UpdateByID should not return an error")

	// Verify the cryptographic key is updated
	var updatedKey model.CryptographicKey
	err = ctx.DB.First(&updatedKey, "id = ?", cryptographicKey.ID).Error
	assert.NoError(t, err, "Failed to find updated cryptographic key")
	assert.Equal(t, "ECDSA", updatedKey.Type, "Type should be updated")
}

// TestCryptographicKeyRepository_DeleteByID tests the DeleteByID method of CryptographicKeyRepositoryImpl
func TestCryptographicKeyRepository_DeleteByID(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptographicKey object
	cryptographicKey := &model.CryptographicKey{
		ID:        uuid.New().String(),            // Generate valid UUID for ID
		Type:      "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err := ctx.CryptoKeyRepo.Create(cryptographicKey)
	assert.NoError(t, err, "Create should not return an error")

	// Delete the cryptographic key by ID
	err = ctx.CryptoKeyRepo.DeleteByID(cryptographicKey.ID)
	assert.NoError(t, err, "DeleteByID should not return an error")

	// Verify the cryptographic key is deleted
	var deletedKey model.CryptographicKey
	err = ctx.DB.First(&deletedKey, "id = ?", cryptographicKey.ID).Error
	assert.Error(t, err, "Cryptographic key should be deleted")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Error should be 'record not found'")
}
