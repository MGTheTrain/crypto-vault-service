package repository

import (
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/test/helpers"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCryptoKeySqliteRepository_Create(t *testing.T) {
	ctx := helpers.SetupTestDB(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, ctx, dbType)

	cryptoKeyMeta := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(),
		KeyPairID:       uuid.New().String(),
		Type:            "public",
		KeySize:         521,
		Algorithm:       "EC",
		DateTimeCreated: time.Now(),
		UserID:          uuid.New().String(),
	}

	err := ctx.CryptoKeyRepo.Create(cryptoKeyMeta)
	assert.NoError(t, err, "Create should not return an error")

	var createdCryptoKeyMeta keys.CryptoKeyMeta
	err = ctx.DB.First(&createdCryptoKeyMeta, "id = ?", cryptoKeyMeta.ID).Error
	assert.NoError(t, err, "Failed to find created cryptographic key")
	assert.Equal(t, cryptoKeyMeta.ID, createdCryptoKeyMeta.ID, "ID should match")
	assert.Equal(t, cryptoKeyMeta.Type, createdCryptoKeyMeta.Type, "Type should match")
}

func TestCryptoKeySqliteRepository_GetByID(t *testing.T) {
	ctx := helpers.SetupTestDB(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, ctx, dbType)

	cryptoKeyMeta := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(),
		KeyPairID:       uuid.New().String(),
		Type:            "private",
		KeySize:         2048,
		Algorithm:       "RSA",
		DateTimeCreated: time.Now(),
		UserID:          uuid.New().String(),
	}

	err := ctx.CryptoKeyRepo.Create(cryptoKeyMeta)
	assert.NoError(t, err, "Create should not return an error")

	fetchedCryptoKeyMeta, err := ctx.CryptoKeyRepo.GetByID(cryptoKeyMeta.ID)
	assert.NoError(t, err, "GetByID should not return an error")
	assert.NotNil(t, fetchedCryptoKeyMeta, "Fetched cryptographic key should not be nil")
	assert.Equal(t, cryptoKeyMeta.ID, fetchedCryptoKeyMeta.ID, "ID should match")
}

func TestCryptoKeySqliteRepository_List(t *testing.T) {
	ctx := helpers.SetupTestDB(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, ctx, dbType)

	cryptoKeyMeta1 := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(),
		KeyPairID:       uuid.New().String(),
		Type:            "private",
		KeySize:         2048,
		Algorithm:       "RSA",
		DateTimeCreated: time.Now(),
		UserID:          uuid.New().String(),
	}

	cryptoKeyMeta2 := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(),
		KeyPairID:       uuid.New().String(),
		Type:            "public",
		KeySize:         521,
		Algorithm:       "EC",
		DateTimeCreated: time.Now(),
		UserID:          cryptoKeyMeta1.UserID, // Same UserID for listing purpose
	}

	// Create crypto keys
	err := ctx.CryptoKeyRepo.Create(cryptoKeyMeta1)
	assert.NoError(t, err, "Create should not return an error")
	err = ctx.CryptoKeyRepo.Create(cryptoKeyMeta2)
	assert.NoError(t, err, "Create should not return an error")

	// List crypto keys
	query := &keys.CryptoKeyQuery{}
	cryptoKeys, err := ctx.CryptoKeyRepo.List(query)
	assert.NoError(t, err, "List should not return an error")
	assert.Len(t, cryptoKeys, 2, "There should be two cryptographic keys in the list")
}

func TestCryptoKeySqliteRepository_UpdateByID(t *testing.T) {
	ctx := helpers.SetupTestDB(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, ctx, dbType)

	cryptoKeyMeta := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(),
		KeyPairID:       uuid.New().String(),
		Type:            "public",
		KeySize:         521,
		Algorithm:       "EC",
		DateTimeCreated: time.Now(),
		UserID:          uuid.New().String(),
	}

	err := ctx.CryptoKeyRepo.Create(cryptoKeyMeta)
	assert.NoError(t, err, "Create should not return an error")

	// Update the key's type
	cryptoKeyMeta.Type = "private"
	err = ctx.CryptoKeyRepo.UpdateByID(cryptoKeyMeta)
	assert.NoError(t, err, "UpdateByID should not return an error")

	var updatedCryptoKeyMeta keys.CryptoKeyMeta
	err = ctx.DB.First(&updatedCryptoKeyMeta, "id = ?", cryptoKeyMeta.ID).Error
	assert.NoError(t, err, "Failed to find updated cryptographic key")
	assert.Equal(t, "private", updatedCryptoKeyMeta.Type, "Type should be updated")
}

func TestCryptoKeySqliteRepository_DeleteByID(t *testing.T) {
	ctx := helpers.SetupTestDB(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, ctx, dbType)

	cryptoKeyMeta := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(),
		KeyPairID:       uuid.New().String(),
		Type:            "public",
		KeySize:         256,
		Algorithm:       "EC",
		DateTimeCreated: time.Now(),
		UserID:          uuid.New().String(),
	}

	err := ctx.CryptoKeyRepo.Create(cryptoKeyMeta)
	assert.NoError(t, err, "Create should not return an error")

	err = ctx.CryptoKeyRepo.DeleteByID(cryptoKeyMeta.ID)
	assert.NoError(t, err, "DeleteByID should not return an error")

	var deletedCryptoKeyMeta keys.CryptoKeyMeta
	err = ctx.DB.First(&deletedCryptoKeyMeta, "id = ?", cryptoKeyMeta.ID).Error
	assert.Error(t, err, "Cryptographic key should be deleted")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Error should be 'record not found'")
}
