package services

import (
	"crypto_vault_service/internal/app/services"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/test/helpers"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Test case for successful file upload and metadata creation
func TestCryptoKeyUploadService_Upload_Success(t *testing.T) {
	// Set up test context
	ctx := helpers.SetupTestDB(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, ctx, dbType)

	// Initialize Vault Connector
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	vaultConnector, err := connector.NewAzureVaultConnector(connectionString, containerName)
	require.NoError(t, err, "Error creating vault connector")

	// Set up the CryptoKeyUploadService
	cryptoKeyUploadService := &services.CryptoKeyUploadService{
		VaultConnector: vaultConnector,
		CryptoKeyRepo:  ctx.CryptoKeyRepo,
	}

	// Prepare test file
	testFilePath := "testfile.pem"
	testFileContent := []byte("This is a test file content.")
	// Create test file
	err = os.WriteFile(testFilePath, testFileContent, 0644)
	require.NoError(t, err)

	// Clean up the test file after the test
	defer os.Remove(testFilePath)

	// Call the method under test
	userId := uuid.New().String()
	keyType := "private"
	keyAlgorithm := "EC"
	keyMeta, err := cryptoKeyUploadService.Upload(testFilePath, userId, keyType, keyAlgorithm)

	// Assert that no error occurred during the upload
	require.NoError(t, err, "Error uploading file")

	// Assert that keyMeta is not nil and contains expected data
	require.NotNil(t, keyMeta, "Key metadata should not be nil")
	require.NotEmpty(t, keyMeta.ID, "KeyID should not be empty")
	require.Equal(t, userId, keyMeta.UserID, "UserID does not match")
}

// Test case for successful retrieval of cryptographic key metadata by ID
func TestCryptoKeyMetadataService_GetByID_Success(t *testing.T) {
	// Set up test context
	ctx := helpers.SetupTestDB(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, ctx, dbType)

	// Create a key metadata instance in the database for testing
	cryptoKeyMeta := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(),
		Type:            "public",
		Algorithm:       "EC",
		DateTimeCreated: time.Now(),
		UserID:          uuid.New().String(),
	}

	// Store the cryptographic key metadata in the database
	err := ctx.CryptoKeyRepo.Create(cryptoKeyMeta)
	require.NoError(t, err, "Error creating test cryptographic key metadata")

	// Set up the CryptoKeyMetadataService
	cryptoKeyMetadataService := &services.CryptoKeyMetadataService{
		CryptoKeyRepo: ctx.CryptoKeyRepo,
	}

	// Call the method under test
	fetchedCryptoKeyMeta, err := cryptoKeyMetadataService.GetByID(cryptoKeyMeta.ID)

	// Assert that no error occurred during retrieval
	require.NoError(t, err, "Error retrieving cryptographic key metadata")

	// Assert that the fetched metadata matches the original
	require.NotNil(t, fetchedCryptoKeyMeta, "Fetched cryptographic key metadata should not be nil")
	require.Equal(t, cryptoKeyMeta.ID, fetchedCryptoKeyMeta.ID, "ID should match")
	require.Equal(t, cryptoKeyMeta.Type, fetchedCryptoKeyMeta.Type, "Type should match")
	require.Equal(t, cryptoKeyMeta.Algorithm, fetchedCryptoKeyMeta.Algorithm, "Algorithm should match")
}

// Test case for successful deletion of cryptographic key metadata by ID
func TestCryptoKeyMetadataService_DeleteByID_Success(t *testing.T) {
	// Set up test context
	ctx := helpers.SetupTestDB(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, ctx, dbType)

	// Create a key metadata instance in the database for testing
	cryptoKeyMeta := &keys.CryptoKeyMeta{
		ID:              uuid.New().String(),
		Type:            "private",
		Algorithm:       "RSA",
		DateTimeCreated: time.Now(),
		UserID:          uuid.New().String(),
	}

	// Store the cryptographic key metadata in the database
	err := ctx.CryptoKeyRepo.Create(cryptoKeyMeta)
	require.NoError(t, err, "Error creating test cryptographic key metadata")

	// Set up the CryptoKeyMetadataService
	cryptoKeyMetadataService := &services.CryptoKeyMetadataService{
		CryptoKeyRepo: ctx.CryptoKeyRepo,
	}

	// Call the method under test
	err = cryptoKeyMetadataService.DeleteByID(cryptoKeyMeta.ID)

	// Assert that no error occurred during deletion
	require.NoError(t, err, "Error deleting cryptographic key metadata")

	// Verify the metadata is deleted
	var deletedCryptoKeyMeta keys.CryptoKeyMeta
	err = ctx.DB.First(&deletedCryptoKeyMeta, "id = ?", cryptoKeyMeta.ID).Error
	require.Error(t, err, "Cryptographic key metadata should be deleted")
	require.Equal(t, gorm.ErrRecordNotFound, err, "Error should be 'record not found'")
}

// Test case for successful download of cryptographic key
func TestCryptoKeyDownloadService_Download_Success(t *testing.T) {
	// Set up test context
	ctx := helpers.SetupTestDB(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, ctx, dbType)

	// Initialize Vault Connector
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	vaultConnector, err := connector.NewAzureVaultConnector(connectionString, containerName)
	require.NoError(t, err, "Error creating vault connector")

	// Set up the CryptoKeyUploadService
	cryptoKeyUploadService := &services.CryptoKeyUploadService{
		VaultConnector: vaultConnector,
		CryptoKeyRepo:  ctx.CryptoKeyRepo,
	}

	// Prepare test file
	testFilePath := "testfile.pem"
	testFileContent := []byte("This is a test file content.")
	// Create test file
	err = os.WriteFile(testFilePath, testFileContent, 0644)
	require.NoError(t, err)

	// Clean up the test file after the test
	defer os.Remove(testFilePath)

	// Call the method under test
	userId := uuid.New().String()
	keyType := "private"
	keyAlgorithm := "EC"
	cryptoKeyMeta, err := cryptoKeyUploadService.Upload(testFilePath, userId, keyType, keyAlgorithm)
	require.NoError(t, err, "Error uploading file")

	// Set up the CryptoKeyDownloadService
	cryptoKeyDownloadService := &services.CryptoKeyDownloadService{
		VaultConnector: vaultConnector,
	}

	// Call the method under test
	blobData, err := cryptoKeyDownloadService.Download(cryptoKeyMeta.ID, keyType)

	// Assert that no error occurred during download
	require.NoError(t, err, "Error downloading cryptographic key")

	// Assert that blobData is not nil or empty
	require.NotNil(t, blobData, "Downloaded key data should not be nil")
	require.NotEmpty(t, blobData, "Downloaded key data should not be empty")
}
