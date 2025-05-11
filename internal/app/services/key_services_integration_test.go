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
	cryptoKeyUploadService   keys.CryptoKeyUploadService
	cryptoKeyMetadataService keys.CryptoKeyMetadataService
	cryptoKeyDownloadService keys.CryptoKeyDownloadService
	dbContext                *repository.TestDBContext
}

func NewKeyServicesTest(t *testing.T, dbType string) *KeyServicesTest {
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
	dbContext := repository.SetupTestDB(t, dbType)

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
		cryptoKeyUploadService:   cryptoKeyUploadService,
		cryptoKeyMetadataService: cryptoKeyMetadataService,
		cryptoKeyDownloadService: cryptoKeyDownloadService,
		dbContext:                dbContext,
	}
}

// Test case for successful file upload and metadata creation
func TestCryptoKeyUploadService_Upload_Success(t *testing.T) {
	dbType := "sqlite"
	keyServices := NewKeyServicesTest(t, dbType)
	defer repository.TeardownTestDB(t, keyServices.dbContext, dbType)

	userID := uuid.New().String()
	keyAlgorithm := "EC"
	var keySize uint32 = 256
	ctx := context.Background()

	cryptoKeyMetas, err := keyServices.cryptoKeyUploadService.Upload(ctx, userID, keyAlgorithm, keySize)
	require.NoError(t, err)
	require.Equal(t, len(cryptoKeyMetas), 2)
	require.NotNil(t, cryptoKeyMetas)
	require.NotEmpty(t, cryptoKeyMetas[0].ID)
	require.Equal(t, userID, cryptoKeyMetas[0].UserID)
	require.NotEmpty(t, cryptoKeyMetas[1].ID)
	require.Equal(t, userID, cryptoKeyMetas[1].UserID)
}

// Test case for successful retrieval of cryptographic key metadata by ID
func TestCryptoKeyMetadataService_GetByID_Success(t *testing.T) {

	dbType := "sqlite"
	keyServices := NewKeyServicesTest(t, dbType)
	defer repository.TeardownTestDB(t, keyServices.dbContext, dbType)

	userID := uuid.New().String()
	keyAlgorithm := "EC"
	var keySize uint32 = 256
	ctx := context.Background()

	cryptoKeyMetas, err := keyServices.cryptoKeyUploadService.Upload(ctx, userID, keyAlgorithm, keySize)
	require.NoError(t, err)

	fetchedCryptoKeyMeta, err := keyServices.cryptoKeyMetadataService.GetByID(ctx, cryptoKeyMetas[0].ID)
	require.NoError(t, err)
	require.NotNil(t, fetchedCryptoKeyMeta)
	require.Equal(t, cryptoKeyMetas[0].ID, fetchedCryptoKeyMeta.ID)
}

// Test case for successful deletion of cryptographic key metadata by ID
func TestCryptoKeyMetadataService_DeleteByID_Success(t *testing.T) {
	dbType := "sqlite"
	keyServices := NewKeyServicesTest(t, dbType)
	defer repository.TeardownTestDB(t, keyServices.dbContext, dbType)

	userID := uuid.New().String()
	keyAlgorithm := "EC"
	var keySize uint32 = 521
	ctx := context.Background()

	cryptoKeyMetas, err := keyServices.cryptoKeyUploadService.Upload(ctx, userID, keyAlgorithm, keySize)
	require.NoError(t, err)

	err = keyServices.cryptoKeyMetadataService.DeleteByID(ctx, cryptoKeyMetas[0].ID)
	require.NoError(t, err)

	var deletedCryptoKeyMeta keys.CryptoKeyMeta
	err = keyServices.dbContext.DB.First(&deletedCryptoKeyMeta, "id = ?", cryptoKeyMetas[0].ID).Error
	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
}

// Test case for successful download of cryptographic key
func TestCryptoKeyDownloadService_Download_Success(t *testing.T) {
	dbType := "sqlite"
	keyServices := NewKeyServicesTest(t, dbType)
	defer repository.TeardownTestDB(t, keyServices.dbContext, dbType)

	userID := uuid.New().String()
	keyAlgorithm := "EC"
	var keySize uint32 = 256
	ctx := context.Background()

	cryptoKeyMetas, err := keyServices.cryptoKeyUploadService.Upload(ctx, userID, keyAlgorithm, keySize)
	require.NoError(t, err)

	blobData, err := keyServices.cryptoKeyDownloadService.DownloadByID(ctx, cryptoKeyMetas[0].ID)
	require.NoError(t, err)
	require.NotNil(t, blobData)
	require.NotEmpty(t, blobData)
}
