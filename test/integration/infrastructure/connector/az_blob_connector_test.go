package connector

import (
	"log"
	"os"
	"testing"

	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAzureBlobConnector_Upload tests the Upload method of AzureBlobConnector
func TestAzureBlobConnector_Upload(t *testing.T) {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	blobConnectorSettings := &settings.BlobConnectorSettings{
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}

	abc, err := connector.NewAzureBlobConnector(blobConnectorSettings, logger)
	require.NoError(t, err)

	testFilePath := "testfile.txt"
	testFileContent := []byte("This is a test file content.")

	err = os.WriteFile(testFilePath, testFileContent, 0644)
	require.NoError(t, err)

	defer os.Remove(testFilePath)

	userId := uuid.New().String()
	filePaths := []string{testFilePath}
	blobs, err := abc.Upload(filePaths, userId)
	require.NoError(t, err)

	assert.Len(t, blobs, 1)
	blob := blobs[0]
	assert.NotEmpty(t, blob.ID)
	assert.Equal(t, testFilePath, blob.Name)
	assert.Equal(t, int64(len(testFileContent)), blob.Size)
	assert.Equal(t, ".txt", blob.Type)

	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestAzureBlobConnector_Download tests the Download method of AzureBlobConnector
func TestAzureBlobConnector_Download(t *testing.T) {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	blobConnectorSettings := &settings.BlobConnectorSettings{
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}
	abc, err := connector.NewAzureBlobConnector(blobConnectorSettings, logger)
	require.NoError(t, err)

	testFilePath := "testfile.pem"
	testFileContent := []byte("This is a test file content.")

	err = os.WriteFile(testFilePath, testFileContent, 0644)
	require.NoError(t, err)

	defer os.Remove(testFilePath)

	userId := uuid.New().String()
	filePaths := []string{testFilePath}
	blobs, err := abc.Upload(filePaths, userId)
	require.NoError(t, err)

	blob := blobs[0]
	downloadedData, err := abc.Download(blob.ID, blob.Name)
	require.NoError(t, err)

	assert.Equal(t, testFileContent, downloadedData)

	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestAzureBlobConnector_Delete tests the Delete method of AzureBlobConnector
func TestAzureBlobConnector_Delete(t *testing.T) {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	blobConnectorSettings := &settings.BlobConnectorSettings{
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}
	abc, err := connector.NewAzureBlobConnector(blobConnectorSettings, logger)
	require.NoError(t, err)

	testFilePath := "testfile.pem"
	testFileContent := []byte("This is a test file content.")

	err = os.WriteFile(testFilePath, testFileContent, 0644)
	require.NoError(t, err)

	defer os.Remove(testFilePath)

	userId := uuid.New().String()
	filePaths := []string{testFilePath}
	blobs, err := abc.Upload(filePaths, userId)
	require.NoError(t, err)

	blob := blobs[0]

	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)

	_, err = abc.Download(blob.ID, blob.Name)
	assert.Error(t, err)
}
