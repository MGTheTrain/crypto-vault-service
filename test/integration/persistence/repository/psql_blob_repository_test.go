// repository/blobrepository_test.go
package repository

import (
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/test/helpers"

	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBlobPsqlRepository_Create(t *testing.T) {
	//
	ctx := helpers.SetupTestDB(t)
	dbType := "postgres"
	defer helpers.TeardownTestDB(t, ctx, dbType)

	cryptographicKey := keys.CryptoKeyMeta{
		ID:              uuid.New().String(),
		Type:            "public",
		Algorithm:       "EC",
		DateTimeCreated: time.Now(),
		UserID:          uuid.New().String(),
	}

	blob := &blobs.BlobMeta{
		ID:              uuid.New().String(),
		DateTimeCreated: time.Now(),
		UserID:          cryptographicKey.UserID,
		Name:            "test-blob",
		Size:            1024,
		Type:            "text",
		CryptoKey:       cryptographicKey,
		KeyID:           cryptographicKey.ID,
	}

	err := ctx.BlobRepo.Create(blob)
	assert.NoError(t, err, "Create should not return an error")

	var createdBlob blobs.BlobMeta
	err = ctx.DB.First(&createdBlob, "id = ?", blob.ID).Error
	assert.NoError(t, err, "Failed to find created blob")
	assert.Equal(t, blob.ID, createdBlob.ID, "ID should match")
	assert.Equal(t, blob.Name, createdBlob.Name, "Name should match")
}

func TestBlobPsqlRepository_GetById(t *testing.T) {

	ctx := helpers.SetupTestDB(t)
	dbType := "postgres"
	defer helpers.TeardownTestDB(t, ctx, dbType)
	cryptographicKey := keys.CryptoKeyMeta{
		ID:              uuid.New().String(),
		Type:            "public",
		Algorithm:       "EC",
		DateTimeCreated: time.Now(),
		UserID:          uuid.New().String(),
	}
	blob := &blobs.BlobMeta{
		ID:              uuid.New().String(),
		DateTimeCreated: time.Now(),
		UserID:          cryptographicKey.UserID,
		Name:            "test-blob",
		Size:            1024,
		Type:            "text",
		CryptoKey:       cryptographicKey,
		KeyID:           cryptographicKey.ID,
	}
	err := ctx.BlobRepo.Create(blob)
	assert.NoError(t, err, "Create should not return an error")
	fetchedBlob, err := ctx.BlobRepo.GetById(blob.ID)
	assert.NoError(t, err, "GetById should not return an error")
	assert.NotNil(t, fetchedBlob, "Fetched blob should not be nil")
	assert.Equal(t, blob.ID, fetchedBlob.ID, "ID should match")
}
