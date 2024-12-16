package connector

import (
	"log"
	"os"
	"testing"
	"time"

	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUpload tests the Upload method of AzureVaultConnector
func TestAzureVaultConnector_Upload(t *testing.T) {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	abc, err := connector.NewAzureVaultConnector(connectionString, containerName, logger)
	require.NoError(t, err)

	testFilePath := "testfile.txt"
	testFileContent := []byte("This is a test file content.")

	err = os.WriteFile(testFilePath, testFileContent, 0644)
	require.NoError(t, err)

	defer os.Remove(testFilePath)

	userId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	cryptoKeyMeta, err := abc.Upload(testFilePath, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	assert.NotEmpty(t, cryptoKeyMeta.ID)
	assert.Equal(t, keyType, cryptoKeyMeta.Type)
	assert.Equal(t, userId, cryptoKeyMeta.UserID)
	assert.WithinDuration(t, time.Now(), cryptoKeyMeta.DateTimeCreated, time.Second)

	err = abc.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)
}

// TestDownload tests the Download method of AzureVaultConnector
func TestAzureVaultConnector_Download(t *testing.T) {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	abc, err := connector.NewAzureVaultConnector(connectionString, containerName, logger)
	require.NoError(t, err)

	testFilePath := "testfile.pem"
	testFileContent := []byte("This is a test file content.")

	err = os.WriteFile(testFilePath, testFileContent, 0644)
	require.NoError(t, err)

	defer os.Remove(testFilePath)

	userId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	cryptoKeyMeta, err := abc.Upload(testFilePath, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	downloadedData, err := abc.Download(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)

	assert.Equal(t, testFileContent, downloadedData)

	err = abc.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)
}

// TestDelete tests the Delete method of AzureVaultConnector
func TestAzureVaultConnector_Delete(t *testing.T) {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	abc, err := connector.NewAzureVaultConnector(connectionString, containerName, logger)
	require.NoError(t, err)

	testFilePath := "testfile.pem"
	testFileContent := []byte("This is a test file content.")

	err = os.WriteFile(testFilePath, testFileContent, 0644)
	require.NoError(t, err)

	defer os.Remove(testFilePath)

	userId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	cryptoKeyMeta, err := abc.Upload(testFilePath, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	err = abc.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)

	_, err = abc.Download(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	assert.Error(t, err)
}
