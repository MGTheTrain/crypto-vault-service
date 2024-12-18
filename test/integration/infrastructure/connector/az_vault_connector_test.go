package connector

import (
	"os"
	"testing"
	"time"

	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AzureVaultConnectorTest encapsulates common logic for tests
type AzureVaultConnectorTest struct {
	Connector *connector.AzureVaultConnector
}

// NewAzureVaultConnectorTest initializes and returns a new AzureVaultConnectorTest
func NewAzureVaultConnectorTest(t *testing.T, connectionString, containerName string) *AzureVaultConnectorTest {
	// Create logger
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}
	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err)

	// Create connector
	keyConnectorSettings := &settings.KeyConnectorSettings{
		ConnectionString: connectionString,
		ContainerName:    containerName,
	}
	abc, err := connector.NewAzureVaultConnector(keyConnectorSettings, logger)
	require.NoError(t, err)

	return &AzureVaultConnectorTest{
		Connector: abc,
	}
}

func TestAzureVaultConnector_Upload(t *testing.T) {
	helper := NewAzureVaultConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFilePath := "testfile.txt"
	testFileContent := []byte("This is a test file content.")
	err := helpers.CreateTestFile(testFilePath, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFilePath)

	userId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	keySize := 2048

	// Upload the file
	cryptoKeyMeta, err := helper.Connector.Upload(testFilePath, userId, keyType, keyAlgorithm, uint(keySize))
	require.NoError(t, err)

	assert.NotEmpty(t, cryptoKeyMeta.ID)
	assert.Equal(t, keyType, cryptoKeyMeta.Type)
	assert.Equal(t, userId, cryptoKeyMeta.UserID)
	assert.WithinDuration(t, time.Now(), cryptoKeyMeta.DateTimeCreated, time.Second)

	// Clean up
	err = helper.Connector.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)
}

func TestAzureVaultConnector_UploadFromForm(t *testing.T) {
	helper := NewAzureVaultConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFileContent := []byte("This is a test file content.")

	userId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	keySize := 2048

	// Upload the file using UploadFromForm method
	cryptoKeyMeta, err := helper.Connector.UploadBytes(testFileContent, userId, keyType, keyAlgorithm, uint(keySize))
	require.NoError(t, err)

	assert.NotEmpty(t, cryptoKeyMeta.ID)
	assert.Equal(t, keyType, cryptoKeyMeta.Type)
	assert.Equal(t, userId, cryptoKeyMeta.UserID)
	assert.WithinDuration(t, time.Now(), cryptoKeyMeta.DateTimeCreated, time.Second)

	// Clean up
	err = helper.Connector.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)
}

func TestAzureVaultConnector_Download(t *testing.T) {
	helper := NewAzureVaultConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFilePath := "testfile.pem"
	testFileContent := []byte("This is a test file content.")
	err := helpers.CreateTestFile(testFilePath, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFilePath)

	userId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	keySize := 2048

	// Upload the file
	cryptoKeyMeta, err := helper.Connector.Upload(testFilePath, userId, keyType, keyAlgorithm, uint(keySize))
	require.NoError(t, err)

	// Download the file
	downloadedData, err := helper.Connector.Download(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)

	assert.Equal(t, testFileContent, downloadedData)

	// Clean up
	err = helper.Connector.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)
}

func TestAzureVaultConnector_Delete(t *testing.T) {
	helper := NewAzureVaultConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFilePath := "testfile.pem"
	testFileContent := []byte("This is a test file content.")
	err := helpers.CreateTestFile(testFilePath, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFilePath)

	userId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	keySize := 2048

	// Upload the file
	cryptoKeyMeta, err := helper.Connector.Upload(testFilePath, userId, keyType, keyAlgorithm, uint(keySize))
	require.NoError(t, err)

	// Delete the file
	err = helper.Connector.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)

	// Attempt to download the deleted file (should result in an error)
	_, err = helper.Connector.Download(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	assert.Error(t, err)
}
