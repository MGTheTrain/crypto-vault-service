//go:build integration
// +build integration

package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/persistence/repository"
)

type KeyServicesTest struct {
	CryptoKeyUploadService   *CryptoKeyUploadService
	CryptoKeyMetadataService *CryptoKeyMetadataService
	CryptoKeyDownloadService *CryptoKeyDownloadService
	DBContext                *repository.TestDBContext
}

func NewKeyServicesTest(t *testing.T) *KeyServicesTest {
	ctx := context.Background()
	// Set up logger
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err, "Error creating logger")

	// Set up DB context (sqlite)
	dbContext := repository.SetupTestDB(t)

	// Set up connector
	keyConnectorSettings := &settings.KeyConnectorSettings{
		CloudProvider:    "azure",
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}
	vaultConnector, err := connector.NewAzureVaultConnector(ctx, keyConnectorSettings, logger)
	require.NoError(t, err, "Error creating vault connector")

	// Initialize services
	cryptoKeyUploadService, err := NewCryptoKeyUploadService(vaultConnector, dbContext.CryptoKeyRepo, logger)
	require.NoError(t, err, "Error creating CryptoKeyUploadService")

	cryptoKeyMetadataService, err := NewCryptoKeyMetadataService(vaultConnector, dbContext.CryptoKeyRepo, logger)
	require.NoError(t, err, "Error creating CryptoKeyMetadataService")

	cryptoKeyDownloadService, err := NewCryptoKeyDownloadService(vaultConnector, dbContext.CryptoKeyRepo, logger)
	require.NoError(t, err, "Error creating CryptoKeyDownloadService")

	// Return struct with services and context
	return &KeyServicesTest{
		CryptoKeyUploadService:   cryptoKeyUploadService,
		CryptoKeyMetadataService: cryptoKeyMetadataService,
		CryptoKeyDownloadService: cryptoKeyDownloadService,
		DBContext:                dbContext,
	}
}

// Test case for successful file upload and metadata creation
func TestCryptoKeyUploadService_Upload_Success(t *testing.T) {
	keyServices := NewKeyServicesTest(t)
	dbType := "sqlite"
	defer repository.TeardownTestDB(t, keyServices.DBContext, dbType)

	userId := uuid.New().String()
	keyAlgorithm := "EC"
	keySize := 256
	ctx := context.Background()

	cryptoKeyMetas, err := keyServices.CryptoKeyUploadService.Upload(ctx, userId, keyAlgorithm, uint(keySize))
	require.NoError(t, err)
	require.Equal(t, len(cryptoKeyMetas), 2)
	require.NotNil(t, cryptoKeyMetas)
	require.NotEmpty(t, cryptoKeyMetas[0].ID)
	require.Equal(t, userId, cryptoKeyMetas[0].UserID)
	require.NotEmpty(t, cryptoKeyMetas[1].ID)
	require.Equal(t, userId, cryptoKeyMetas[1].UserID)
}

// Test case for successful retrieval of cryptographic key metadata by ID
func TestCryptoKeyMetadataService_GetByID_Success(t *testing.T) {

	keyServices := NewKeyServicesTest(t)
	dbType := "sqlite"
	defer repository.TeardownTestDB(t, keyServices.DBContext, dbType)

	userId := uuid.New().String()
	keyAlgorithm := "EC"
	keySize := 256
	ctx := context.Background()

	cryptoKeyMetas, err := keyServices.CryptoKeyUploadService.Upload(ctx, userId, keyAlgorithm, uint(keySize))
	require.NoError(t, err)

	fetchedCryptoKeyMeta, err := keyServices.CryptoKeyMetadataService.GetByID(cryptoKeyMetas[0].ID)
	require.NoError(t, err)
	require.NotNil(t, fetchedCryptoKeyMeta)
	require.Equal(t, cryptoKeyMetas[0].ID, fetchedCryptoKeyMeta.ID)
}

// Test case for successful deletion of cryptographic key metadata by ID
func TestCryptoKeyMetadataService_DeleteByID_Success(t *testing.T) {
	keyServices := NewKeyServicesTest(t)
	dbType := "sqlite"
	defer repository.TeardownTestDB(t, keyServices.DBContext, dbType)

	userId := uuid.New().String()
	keyAlgorithm := "EC"
	keySize := 521
	ctx := context.Background()

	cryptoKeyMetas, err := keyServices.CryptoKeyUploadService.Upload(ctx, userId, keyAlgorithm, uint(keySize))
	require.NoError(t, err)

	err = keyServices.CryptoKeyMetadataService.DeleteByID(ctx, cryptoKeyMetas[0].ID)
	require.NoError(t, err)

	var deletedCryptoKeyMeta keys.CryptoKeyMeta
	err = keyServices.DBContext.DB.First(&deletedCryptoKeyMeta, "id = ?", cryptoKeyMetas[0].ID).Error
	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
}

// Test case for successful download of cryptographic key
func TestCryptoKeyDownloadService_Download_Success(t *testing.T) {
	keyServices := NewKeyServicesTest(t)
	dbType := "sqlite"
	defer repository.TeardownTestDB(t, keyServices.DBContext, dbType)

	userId := uuid.New().String()
	keyAlgorithm := "EC"
	keySize := 256
	ctx := context.Background()

	cryptoKeyMetas, err := keyServices.CryptoKeyUploadService.Upload(ctx, userId, keyAlgorithm, uint(keySize))
	require.NoError(t, err)

	blobData, err := keyServices.CryptoKeyDownloadService.Download(ctx, cryptoKeyMetas[0].ID)
	require.NoError(t, err)
	require.NotNil(t, blobData)
	require.NotEmpty(t, blobData)
}
