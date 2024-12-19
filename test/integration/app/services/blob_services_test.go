package services

import (
	"crypto_vault_service/internal/app/services"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/infrastructure/utils"
	"crypto_vault_service/test/helpers"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type BlobServicesTest struct {
	BlobUploadService      *services.BlobUploadService
	BlobDownloadService    *services.BlobDownloadService
	BlobMetadataService    *services.BlobMetadataService
	CryptoKeyUploadService *services.CryptoKeyUploadService
	DBContext              *helpers.TestDBContext
}

func NewBlobServicesTest(t *testing.T) *BlobServicesTest {

	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err, "Error creating logger")

	ctx := helpers.SetupTestDB(t)

	blobConnectorSettings := &settings.BlobConnectorSettings{
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}
	blobConnector, err := connector.NewAzureBlobConnector(blobConnectorSettings, logger)
	require.NoError(t, err, "Error creating blob connector")

	keyConnectorSettings := &settings.KeyConnectorSettings{
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}
	vaultConnector, err := connector.NewAzureVaultConnector(keyConnectorSettings, logger)
	require.NoError(t, err, "Error creating vault connector")

	blobUploadService := services.NewBlobUploadService(blobConnector, ctx.BlobRepo, vaultConnector, ctx.CryptoKeyRepo, logger)
	require.NoError(t, err, "Error creating BlobUploadService")

	blobDownloadService := services.NewBlobDownloadService(blobConnector, ctx.BlobRepo, vaultConnector, ctx.CryptoKeyRepo, logger)
	require.NoError(t, err, "Error creating BlobDownloadService")

	blobMetadataService := services.NewBlobMetadataService(ctx.BlobRepo, blobConnector, logger)
	require.NoError(t, err, "Error creating BlobMetadataService")

	cryptoKeyUploadService, err := services.NewCryptoKeyUploadService(vaultConnector, ctx.CryptoKeyRepo, logger)
	require.NoError(t, err, "Error creating CryptoKeyUploadService")

	return &BlobServicesTest{
		BlobUploadService:      blobUploadService,
		BlobDownloadService:    blobDownloadService,
		BlobMetadataService:    blobMetadataService,
		CryptoKeyUploadService: cryptoKeyUploadService,
		DBContext:              ctx,
	}
}

// Test case for successful blob upload with RSA encryption and signing
func TestBlobUploadService_Upload_With_RSA_Encryption_And_Signing_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer helpers.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := helpers.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()

	keyPairId := uuid.New().String()
	keyAlgorithm := "RSA"
	keySize := 2048

	cryptoKeyMetas, err := blobServices.CryptoKeyUploadService.Upload(userId, keyPairId, keyAlgorithm, uint(keySize))
	require.NoError(t, err)
	require.Equal(t, len(cryptoKeyMetas), 2)

	signKeyId := cryptoKeyMetas[0].ID       // private key
	encryptionKeyId := cryptoKeyMetas[1].ID // public key

	blobMetas, err := blobServices.BlobUploadService.Upload(form, userId, &encryptionKeyId, &signKeyId)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)
	require.NotEmpty(t, blobMetas[0].ID)
	require.Equal(t, userId, blobMetas[0].UserID)
}

// Test case for successful blob upload with AES encryption and ECDSA signing
func TestBlobUploadService_Upload_With_AES_Encryption_And_ECDSA_Signing_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer helpers.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := helpers.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()

	// generate signing private EC key
	signKeyPairId := uuid.New().String()
	signKeyAlgorithm := "EC"
	signKeySize := 256

	cryptoKeyMetas, err := blobServices.CryptoKeyUploadService.Upload(userId, signKeyPairId, signKeyAlgorithm, uint(signKeySize))
	require.NoError(t, err)
	require.Equal(t, len(cryptoKeyMetas), 2)

	// generate AES encryption key
	encryptionKeyPairId := uuid.New().String()
	encryptionKeyAlgorithm := "AES"
	encryptionKeySize := 256

	cryptoKeyMetas2, err := blobServices.CryptoKeyUploadService.Upload(userId, encryptionKeyPairId, encryptionKeyAlgorithm, uint(encryptionKeySize))
	require.NoError(t, err)
	require.Equal(t, len(cryptoKeyMetas2), 1)

	signKeyId := cryptoKeyMetas[0].ID        // private key
	encryptionKeyId := cryptoKeyMetas2[0].ID // symmetric key

	blobMetas, err := blobServices.BlobUploadService.Upload(form, userId, &encryptionKeyId, &signKeyId)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)
	require.NotEmpty(t, blobMetas[0].ID)
	require.Equal(t, userId, blobMetas[0].UserID)
}

// Test case for successful blob upload without encryption and signing
func TestBlobUploadService_Upload_Without_Encryption_And_Signing_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer helpers.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := helpers.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()
	var encryptionKeyId *string = nil
	var signKeyId *string = nil

	blobMetas, err := blobServices.BlobUploadService.Upload(form, userId, encryptionKeyId, signKeyId)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)
	require.NotEmpty(t, blobMetas[0].ID)
	require.Equal(t, userId, blobMetas[0].UserID)
}

// Test case for failed blob upload due to invalid encryption key
func TestBlobUploadService_Upload_Fail_InvalidEncryptionKey(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer helpers.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := helpers.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()
	invalidEncryptionKeyId := "invalid-encryption-key-id"
	signKeyId := uuid.New().String()

	blobMetas, err := blobServices.BlobUploadService.Upload(form, userId, &invalidEncryptionKeyId, &signKeyId)
	require.Error(t, err)
	require.Nil(t, blobMetas)
}

// Test case for successful blob download
func TestBlobDownloadService_Download_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer helpers.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := helpers.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()
	var encryptionKeyId *string = nil
	var signKeyId *string = nil

	blobMetas, err := blobServices.BlobUploadService.Upload(form, userId, encryptionKeyId, signKeyId)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)

	// decryptionKeyId := uuid.New().String()

	// blobData, err := blobServices.BlobDownloadService.Download(blobId, &decryptionKeyId)
	blobData, err := blobServices.BlobDownloadService.Download(blobMetas[0].ID, nil)
	require.NoError(t, err)
	require.NotNil(t, blobData)
	require.NotEmpty(t, blobData)
}

// Test case for failed blob download with invalid decryption key
func TestBlobDownloadService_Download_Fail_InvalidDecryptionKey(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer helpers.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	blobId := uuid.New().String()
	invalidDecryptionKeyId := "invalid-decryption-key-id"

	blobData, err := blobServices.BlobDownloadService.Download(blobId, &invalidDecryptionKeyId)
	require.Error(t, err)
	require.Nil(t, blobData)
}

// Test case for successful listing of blob metadata
func TestBlobMetadataService_List_Success(t *testing.T) {
	blobServices := NewBlobServicesTest(t)
	defer helpers.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"
	err := helpers.CreateTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)

	form, err := utils.CreateForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	// encryptionKeyId := uuid.New().String()
	// signKeyId := uuid.New().String()

	// blobMetas, err := blobServices.BlobUploadService.Upload(form, userId, &encryptionKeyId, &signKeyId)
	blobMetas, err := blobServices.BlobUploadService.Upload(form, userId, nil, nil)
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
	defer helpers.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := helpers.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()
	var encryptionKeyId *string = nil
	var signKeyId *string = nil

	blobMetas, err := blobServices.BlobUploadService.Upload(form, userId, encryptionKeyId, signKeyId)
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
	defer helpers.TeardownTestDB(t, blobServices.DBContext, "sqlite")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"

	form, err := helpers.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()
	var encryptionKeyId *string = nil
	var signKeyId *string = nil

	blobMetas, err := blobServices.BlobUploadService.Upload(form, userId, encryptionKeyId, signKeyId)
	require.NoError(t, err)
	require.NotNil(t, blobMetas)

	err = blobServices.BlobMetadataService.DeleteByID(blobMetas[0].ID)
	require.NoError(t, err)

	// Verify deletion
	var deletedBlobMeta blobs.BlobMeta
	err = blobServices.DBContext.DB.First(&deletedBlobMeta, "id = ?", blobMetas[0].ID).Error
	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
}
