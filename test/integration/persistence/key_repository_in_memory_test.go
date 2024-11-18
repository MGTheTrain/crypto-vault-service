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
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err := ctx.CryptoKeyRepo.Create(cryptographicKey)
	assert.NoError(t, err, "Create should not return an error")

	// Verify the cryptographic key is created and exists in DB
	var createdKey model.CryptographicKey
	err = ctx.DB.First(&createdKey, "key_id = ?", cryptographicKey.KeyID).Error
	assert.NoError(t, err, "Failed to find created cryptographic key")
	assert.Equal(t, cryptographicKey.KeyID, createdKey.KeyID, "KeyID should match")
	assert.Equal(t, cryptographicKey.KeyType, createdKey.KeyType, "KeyType should match")
}

// TestCryptographicKeyRepository_GetByID tests the GetByID method of CryptographicKeyRepositoryImpl
func TestCryptographicKeyRepository_GetByID(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptographicKey object
	cryptographicKey := &model.CryptographicKey{
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "RSA",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err := ctx.CryptoKeyRepo.Create(cryptographicKey)
	assert.NoError(t, err, "Create should not return an error")

	// Get the cryptographic key by ID
	fetchedKey, err := ctx.CryptoKeyRepo.GetByID(cryptographicKey.KeyID)
	assert.NoError(t, err, "GetByID should not return an error")
	assert.NotNil(t, fetchedKey, "Fetched cryptographic key should not be nil")
	assert.Equal(t, cryptographicKey.KeyID, fetchedKey.KeyID, "KeyID should match")
}

// TestCryptographicKeyRepository_UpdateByID tests the UpdateByID method of CryptographicKeyRepositoryImpl
func TestCryptographicKeyRepository_UpdateByID(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptographicKey object
	cryptographicKey := &model.CryptographicKey{
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err := ctx.CryptoKeyRepo.Create(cryptographicKey)
	assert.NoError(t, err, "Create should not return an error")

	// Update the cryptographic key's type
	cryptographicKey.KeyType = "ECDSA"
	err = ctx.CryptoKeyRepo.UpdateByID(cryptographicKey)
	assert.NoError(t, err, "UpdateByID should not return an error")

	// Verify the cryptographic key is updated
	var updatedKey model.CryptographicKey
	err = ctx.DB.First(&updatedKey, "key_id = ?", cryptographicKey.KeyID).Error
	assert.NoError(t, err, "Failed to find updated cryptographic key")
	assert.Equal(t, "ECDSA", updatedKey.KeyType, "KeyType should be updated")
}

// TestCryptographicKeyRepository_DeleteByID tests the DeleteByID method of CryptographicKeyRepositoryImpl
func TestCryptographicKeyRepository_DeleteByID(t *testing.T) {
	// Set up test context
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Create a valid CryptographicKey object
	cryptographicKey := &model.CryptographicKey{
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err := ctx.CryptoKeyRepo.Create(cryptographicKey)
	assert.NoError(t, err, "Create should not return an error")

	// Delete the cryptographic key by ID
	err = ctx.CryptoKeyRepo.DeleteByID(cryptographicKey.KeyID)
	assert.NoError(t, err, "DeleteByID should not return an error")

	// Verify the cryptographic key is deleted
	var deletedKey model.CryptographicKey
	err = ctx.DB.First(&deletedKey, "key_id = ?", cryptographicKey.KeyID).Error
	assert.Error(t, err, "Cryptographic key should be deleted")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Error should be 'record not found'")
}
