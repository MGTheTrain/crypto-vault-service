package repository

import (
	"crypto_vault_service/internal/domain/keys"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// TestCryptoKeySqliteRepository_Create tests the Create method of GormCryptoKeyRepository
func TestCryptoKeySqliteRepository_Create(t *testing.T) {
	err := os.Setenv("DB_TYPE", "sqlite")
	if err != nil {
		t.Fatalf("Error setting environment variable: %v", err)
	}

	// Set up test context
	ctx := SetupTestDB(t)
	defer TeardownTestDB(t, ctx)

	// Create a valid CryptoKey object
	cryptoKeyMeta := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(), // Generate valid UUID for ID
		Type:            "public",            // Example key type
		Algorithm:       "EC",                // Example algorithm
		DateTimeCreated: time.Now(),          // Valid DateTimeCreated time
		UserID:          uuid.New().String(), // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err = ctx.CryptoKeyRepo.Create(cryptoKeyMeta)
	assert.NoError(t, err, "Create should not return an error")

	// Verify the cryptographic key is created and exists in DB
	var createdCryptKeyMeta keys.CryptoKeyMeta
	err = ctx.DB.First(&createdCryptKeyMeta, "id = ?", cryptoKeyMeta.ID).Error
	assert.NoError(t, err, "Failed to find created cryptographic key")
	assert.Equal(t, cryptoKeyMeta.ID, createdCryptKeyMeta.ID, "ID should match")
	assert.Equal(t, cryptoKeyMeta.Type, createdCryptKeyMeta.Type, "Type should match")
}

// TestCryptoKeySqliteRepository_GetByID tests the GetByID method of GormCryptoKeyRepository
func TestCryptoKeySqliteRepository_GetByID(t *testing.T) {
	err := os.Setenv("DB_TYPE", "sqlite")
	if err != nil {
		t.Fatalf("Error setting environment variable: %v", err)
	}

	// Set up test context
	ctx := SetupTestDB(t)
	defer TeardownTestDB(t, ctx)

	// Create a valid CryptoKey object
	cryptoKeyMeta := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(), // Generate valid UUID for ID
		Type:            "private",           // Example key type
		Algorithm:       "RSA",               // Example algorithm
		DateTimeCreated: time.Now(),          // Valid DateTimeCreated time
		UserID:          uuid.New().String(), // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err = ctx.CryptoKeyRepo.Create(cryptoKeyMeta)
	assert.NoError(t, err, "Create should not return an error")

	// Get the cryptographic key by ID
	fetchedCryptoKeyMeta, err := ctx.CryptoKeyRepo.GetByID(cryptoKeyMeta.ID)
	assert.NoError(t, err, "GetByID should not return an error")
	assert.NotNil(t, fetchedCryptoKeyMeta, "Fetched cryptographic key should not be nil")
	assert.Equal(t, cryptoKeyMeta.ID, fetchedCryptoKeyMeta.ID, "ID should match")
}

// TestCryptoKeySqliteRepository_UpdateByID tests the UpdateByID method of GormCryptoKeyRepository
func TestCryptoKeySqliteRepository_UpdateByID(t *testing.T) {
	err := os.Setenv("DB_TYPE", "sqlite")
	if err != nil {
		t.Fatalf("Error setting environment variable: %v", err)
	}

	// Set up test context
	ctx := SetupTestDB(t)
	defer TeardownTestDB(t, ctx)

	// Create a valid CryptoKey object
	cryptoKeyMeta := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(), // Generate valid UUID for ID
		Type:            "public",            // Example key type
		Algorithm:       "EC",                // Example algorithm
		DateTimeCreated: time.Now(),          // Valid DateTimeCreated time
		UserID:          uuid.New().String(), // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err = ctx.CryptoKeyRepo.Create(cryptoKeyMeta)
	assert.NoError(t, err, "Create should not return an error")

	// Update the cryptographic key's type
	cryptoKeyMeta.Type = "public"
	err = ctx.CryptoKeyRepo.UpdateByID(cryptoKeyMeta)
	assert.NoError(t, err, "UpdateByID should not return an error")

	// Verify the cryptographic key is updated
	var updatedCryptoKeyMeta keys.CryptoKeyMeta
	err = ctx.DB.First(&updatedCryptoKeyMeta, "id = ?", cryptoKeyMeta.ID).Error
	assert.NoError(t, err, "Failed to find updated cryptographic key")
	assert.Equal(t, cryptoKeyMeta.Type, updatedCryptoKeyMeta.Type, "Type should be updated")
}

// TestCryptoKeySqliteRepository_DeleteByID tests the DeleteByID method of GormCryptoKeyRepository
func TestCryptoKeySqliteRepository_DeleteByID(t *testing.T) {
	// Set up test context
	ctx := SetupTestDB(t)
	defer TeardownTestDB(t, ctx)

	// Create a valid CryptoKey object
	cryptoKeyMeta := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(), // Generate valid UUID for ID
		Type:            "public",            // Example key type
		Algorithm:       "EC",                // Example algorithm
		DateTimeCreated: time.Now(),          // Valid DateTimeCreated time
		UserID:          uuid.New().String(), // Generate valid UUID for UserID
	}

	// Create the cryptographic key in DB
	err := ctx.CryptoKeyRepo.Create(cryptoKeyMeta)
	assert.NoError(t, err, "Create should not return an error")

	// Delete the cryptographic key by ID
	err = ctx.CryptoKeyRepo.DeleteByID(cryptoKeyMeta.ID)
	assert.NoError(t, err, "DeleteByID should not return an error")

	// Verify the cryptographic key is deleted
	var deletedCryptoKeyMeta keys.CryptoKeyMeta
	err = ctx.DB.First(&deletedCryptoKeyMeta, "id = ?", cryptoKeyMeta.ID).Error
	assert.Error(t, err, "Cryptographic key should be deleted")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Error should be 'record not found'")
}
