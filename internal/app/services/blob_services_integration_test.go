//go:build integration
// +build integration

package services

import (
	"context"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/infrastructure/utils"
	"crypto_vault_service/internal/persistence/repository"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type BlobServicesTest struct {
	BlobUploadService      *BlobUploadService
	BlobDownloadService    *BlobDownloadService
	BlobMetadataService    *BlobMetadataService
	CryptoKeyUploadService *CryptoKeyUploadService
	DBContext              *repository.TestDBContext
}

func NewBlobServicesTest(t *testing.T) *BlobServicesTest {
	ctx := context.Background()

	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err, "Error creating logger")

	dbContext := repository.SetupTestDB(t)

	blobConnectorSettings := &settings.BlobConnectorSettings{
		CloudProvider:    "azure",
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}
	blobConnector, err := connector.NewAzureBlobConnector(ctx, blobConnectorSettings, logger)
	require.NoError(t, err, "Error creating blob connector")

	keyConnectorSettings := &settings.KeyConnectorSettings{
		CloudProvider:    "azure",
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}
	vaultConnector, err := connector.NewAzureVaultConnector(ctx, keyConnectorSettings, logger)
	require.NoError(t, err, "Error creating vault connector")

	blobUploadService, err := NewBlobUploadService(blobConnector, dbContext.BlobRepo, vaultConnector, dbContext.CryptoKeyRepo, logger)
	require.NoError(t, err, "Error creating BlobUploadService")

	blobDownloadService, err := NewBlobDownloadService(blobConnector, dbContext.BlobRepo, vaultConnector, dbContext.CryptoKeyRepo, logger)
	require.NoError(t, err, "Error creating BlobDownloadService")

	blobMetadataService, err := NewBlobMetadataService(dbContext.BlobRepo, blobConnector, logger)
	require.NoError(t, err, "Error creating BlobMetadataService")

	cryptoKeyUploadService, err := NewCryptoKeyUploadService(vaultConnector, dbContext.CryptoKeyRepo, logger)
	require.NoError(t, err, "Error creating CryptoKeyUploadService")

	return &BlobServicesTest{
		BlobUploadService:      blobUploadService,
		BlobDownloadService:    blobDownloadService,
		BlobMetadataService:    blobMetadataService,
		CryptoKeyUploadService: cryptoKeyUploadService,
		DBContext:              dbContext,
	}
}

// Test case for successful blob upload with RSA encryption and signing
func TestBlobUploadService_Upload_With_RSA_Encryption_And_Signing_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer repository.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := utils.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()

	keyAlgorithm := "RSA"
	keySize := 2048
	ctx := context.Background()

	cryptoKeyMetas, err := blobServices.CryptoKeyUploadService.Upload(ctx, userId, keyAlgorithm, uint(keySize))
	require.NoError(t, err)
	require.Equal(t, len(cryptoKeyMetas), 2)

	signKeyId := cryptoKeyMetas[0].ID       // private key
	encryptionKeyId := cryptoKeyMetas[1].ID // public key

	blobMetas, err := blobServices.BlobUploadService.Upload(ctx, form, userId, &encryptionKeyId, &signKeyId)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)
	require.NotEmpty(t, blobMetas[0].ID)
	require.Equal(t, userId, blobMetas[0].UserID)
}

// Test case for successful blob upload with AES encryption and ECDSA signing
func TestBlobUploadService_Upload_With_AES_Encryption_And_ECDSA_Signing_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer repository.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := utils.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()

	// generate signing private EC key
	signKeyAlgorithm := "EC"
	signKeySize := 256
	ctx := context.Background()

	cryptoKeyMetas, err := blobServices.CryptoKeyUploadService.Upload(ctx, userId, signKeyAlgorithm, uint(signKeySize))
	require.NoError(t, err)
	require.Equal(t, len(cryptoKeyMetas), 2)

	// generate AES encryption key
	encryptionKeyAlgorithm := "AES"
	encryptionKeySize := 256

	cryptoKeyMetas2, err := blobServices.CryptoKeyUploadService.Upload(ctx, userId, encryptionKeyAlgorithm, uint(encryptionKeySize))
	require.NoError(t, err)
	require.Equal(t, len(cryptoKeyMetas2), 1)

	signKeyId := cryptoKeyMetas[0].ID        // private key
	encryptionKeyId := cryptoKeyMetas2[0].ID // symmetric key

	blobMetas, err := blobServices.BlobUploadService.Upload(ctx, form, userId, &encryptionKeyId, &signKeyId)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)
	require.NotEmpty(t, blobMetas[0].ID)
	require.Equal(t, userId, blobMetas[0].UserID)
}

// Test case for successful blob upload without encryption and signing
func TestBlobUploadService_Upload_Without_Encryption_And_Signing_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer repository.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := utils.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()
	var encryptionKeyId *string = nil
	var signKeyId *string = nil
	ctx := context.Background()

	blobMetas, err := blobServices.BlobUploadService.Upload(ctx, form, userId, encryptionKeyId, signKeyId)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)
	require.NotEmpty(t, blobMetas[0].ID)
	require.Equal(t, userId, blobMetas[0].UserID)
}

// Test case for failed blob upload due to invalid encryption key
func TestBlobUploadService_Upload_Fail_InvalidEncryptionKey(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer repository.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := utils.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()
	invalidEncryptionKeyId := "invalid-encryption-key-id"
	signKeyId := uuid.New().String()
	ctx := context.Background()

	blobMetas, err := blobServices.BlobUploadService.Upload(ctx, form, userId, &invalidEncryptionKeyId, &signKeyId)
	require.Error(t, err)
	require.Nil(t, blobMetas)
}

// Test case for successful blob download
func TestBlobDownloadService_Download_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer repository.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := utils.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()
	var encryptionKeyId *string = nil
	var signKeyId *string = nil
	ctx := context.Background()

	blobMetas, err := blobServices.BlobUploadService.Upload(ctx, form, userId, encryptionKeyId, signKeyId)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)

	// decryptionKeyId := uuid.New().String()

	// blobData, err := blobServices.BlobDownloadService.Download(blobId, &decryptionKeyId)
	blobData, err := blobServices.BlobDownloadService.Download(ctx, blobMetas[0].ID, nil)
	require.NoError(t, err)
	require.NotNil(t, blobData)
	require.NotEmpty(t, blobData)
}

// Test case for failed blob download with invalid decryption key
func TestBlobDownloadService_Download_Fail_InvalidDecryptionKey(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer repository.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	blobId := uuid.New().String()
	invalidDecryptionKeyId := "invalid-decryption-key-id"
	ctx := context.Background()

	blobData, err := blobServices.BlobDownloadService.Download(ctx, blobId, &invalidDecryptionKeyId)
	require.Error(t, err)
	require.Nil(t, blobData)
}

// Test case for successful listing of blob metadata
func TestBlobMetadataService_List_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer repository.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"
	err := utils.CreateTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)

	form, err := utils.CreateForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	// encryptionKeyId := uuid.New().String()
	// signKeyId := uuid.New().String()
	ctx := context.Background()

	// blobMetas, err := blobServices.BlobUploadService.Upload(form, userId, &encryptionKeyId, &signKeyId)
	blobMetas, err := blobServices.BlobUploadService.Upload(ctx, form, userId, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)

	query := &blobs.BlobMetaQuery{}
	blobMetas, err = blobServices.BlobMetadataService.List(query)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)
	require.Greater(t, len(blobMetas), 0)
}

// Test case for successful retrieval of blob metadata by ID
func TestBlobMetadataService_GetByID_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer repository.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := utils.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()
	var encryptionKeyId *string = nil
	var signKeyId *string = nil
	ctx := context.Background()

	blobMetas, err := blobServices.BlobUploadService.Upload(ctx, form, userId, encryptionKeyId, signKeyId)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)

	blobMeta, err := blobServices.BlobMetadataService.GetByID(blobMetas[0].ID)
	require.NoError(t, err)
	require.NotNil(t, blobMeta)
	require.Equal(t, blobMetas[0].ID, blobMeta.ID)
}

// Test case for successful deletion of blob metadata by ID
func TestBlobMetadataService_DeleteByID_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer repository.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := utils.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()
	var encryptionKeyId *string = nil
	var signKeyId *string = nil
	ctx := context.Background()

	blobMetas, err := blobServices.BlobUploadService.Upload(ctx, form, userId, encryptionKeyId, signKeyId)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)

	err = blobServices.BlobMetadataService.DeleteByID(ctx, blobMetas[0].ID)
	require.NoError(t, err)

	// Verify deletion
	var deletedBlobMeta blobs.BlobMeta
	err = blobServices.DBContext.DB.First(&deletedBlobMeta, "id = ?", blobMetas[0].ID).Error
	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
}
