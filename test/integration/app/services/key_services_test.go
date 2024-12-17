package services_test

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"crypto_vault_service/internal/app/services"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/test/helpers"
)

type KeyServicesTest struct {
	CryptoKeyUploadService   *services.CryptoKeyUploadService
	CryptoKeyMetadataService *services.CryptoKeyMetadataService
	CryptoKeyDownloadService *services.CryptoKeyDownloadService
	DBContext                *helpers.TestDBContext
}

func NewKeyServicesTest(t *testing.T) *KeyServicesTest {
	// Set up logger
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err, "Error creating logger")

	// Set up DB context (sqlite)
	ctx := helpers.SetupTestDB(t)

	// Set up connector
	keyConnectorSettings := &settings.KeyConnectorSettings{
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}
	vaultConnector, err := connector.NewAzureVaultConnector(keyConnectorSettings, logger)
	require.NoError(t, err, "Error creating vault connector")

	// Initialize services
	cryptoKeyUploadService, err := services.NewCryptoKeyUploadService(vaultConnector, ctx.CryptoKeyRepo)
	require.NoError(t, err, "Error creating CryptoKeyUploadService")

	cryptoKeyMetadataService, err := services.NewCryptoKeyMetadataService(vaultConnector, ctx.CryptoKeyRepo)
	require.NoError(t, err, "Error creating CryptoKeyMetadataService")

	cryptoKeyDownloadService, err := services.NewCryptoKeyDownloadService(vaultConnector)
	require.NoError(t, err, "Error creating CryptoKeyDownloadService")

	// Return struct with services and context
	return &KeyServicesTest{
		CryptoKeyUploadService:   cryptoKeyUploadService,
		CryptoKeyMetadataService: cryptoKeyMetadataService,
		CryptoKeyDownloadService: cryptoKeyDownloadService,
		DBContext:                ctx,
	}
}

// Test case for successful file upload and metadata creation
func TestCryptoKeyUploadService_Upload_Success(t *testing.T) {
	keyServices := NewKeyServicesTest(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, keyServices.DBContext, dbType)

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"
	err := helpers.CreateTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)
	form, err := helpers.CreateForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	keyType := "private"
	keyAlgorithm := "EC"

	keyMeta, err := keyServices.CryptoKeyUploadService.Upload(form, userId, keyType, keyAlgorithm)
	require.NoError(t, err)
	require.NotNil(t, keyMeta)
	require.NotEmpty(t, keyMeta.ID)
	require.Equal(t, userId, keyMeta.UserID)
}

// Test case for successful retrieval of cryptographic key metadata by ID
func TestCryptoKeyMetadataService_GetByID_Success(t *testing.T) {

	keyServices := NewKeyServicesTest(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, keyServices.DBContext, dbType)

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"
	err := helpers.CreateTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)
	form, err := helpers.CreateForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	keyType := "private"
	keyAlgorithm := "EC"

	cryptoKeyMeta, err := keyServices.CryptoKeyUploadService.Upload(form, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	fetchedCryptoKeyMeta, err := keyServices.CryptoKeyMetadataService.GetByID(cryptoKeyMeta.ID)
	require.NoError(t, err)
	require.NotNil(t, fetchedCryptoKeyMeta)
	require.Equal(t, cryptoKeyMeta.ID, fetchedCryptoKeyMeta.ID)
}

// Test case for successful deletion of cryptographic key metadata by ID
func TestCryptoKeyMetadataService_DeleteByID_Success(t *testing.T) {
	keyServices := NewKeyServicesTest(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, keyServices.DBContext, dbType)

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"
	err := helpers.CreateTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)
	form, err := helpers.CreateForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	keyType := "private"
	keyAlgorithm := "EC"

	cryptoKeyMeta, err := keyServices.CryptoKeyUploadService.Upload(form, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	err = keyServices.CryptoKeyMetadataService.DeleteByID(cryptoKeyMeta.ID)
	require.NoError(t, err)

	var deletedCryptoKeyMeta keys.CryptoKeyMeta
	err = keyServices.DBContext.DB.First(&deletedCryptoKeyMeta, "id = ?", cryptoKeyMeta.ID).Error
	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
}

// Test case for successful download of cryptographic key
func TestCryptoKeyDownloadService_Download_Success(t *testing.T) {
	keyServices := NewKeyServicesTest(t)
	dbType := "sqlite"
	defer helpers.TeardownTestDB(t, keyServices.DBContext, dbType)

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"
	err := helpers.CreateTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)

	form, err := helpers.CreateForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	keyType := "private"
	keyAlgorithm := "EC"

	cryptoKeyMeta, err := keyServices.CryptoKeyUploadService.Upload(form, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	blobData, err := keyServices.CryptoKeyDownloadService.Download(cryptoKeyMeta.ID, keyType)
	require.NoError(t, err)
	require.NotNil(t, blobData)
	require.NotEmpty(t, blobData)
}
