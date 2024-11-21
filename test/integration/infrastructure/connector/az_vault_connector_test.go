package connector

import (
	"os"
	"testing"
	"time"

	"crypto_vault_service/internal/infrastructure/connector"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUpload tests the Upload method of AzureVaultConnector
func TestAzureVaultConnector_Upload(t *testing.T) {
	// Prepare environment and initialize the connector
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	abc, err := connector.NewAzureVaultConnector(connectionString, containerName)
	require.NoError(t, err)

	// Prepare test file content
	testFilePath := "testfile.txt"
	testFileContent := []byte("This is a test file content.")

	// Create test file
	err = os.WriteFile(testFilePath, testFileContent, 0644)
	require.NoError(t, err)

	// Clean up the test file after the test
	defer os.Remove(testFilePath)

	// Upload the file
	userId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	cryptoKeyMeta, err := abc.Upload(testFilePath, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	// Assert that we received key metadata
	assert.NotEmpty(t, cryptoKeyMeta.ID)
	assert.Equal(t, keyType, cryptoKeyMeta.Type)
	assert.Equal(t, userId, cryptoKeyMeta.UserID)
	assert.WithinDuration(t, time.Now(), cryptoKeyMeta.DateTimeCreated, time.Second)

	// Delete the uploaded key
	err = abc.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)
}

// TestDownload tests the Download method of AzureVaultConnector
func TestAzureVaultConnector_Download(t *testing.T) {
	// Prepare environment and initialize the connector
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	abc, err := connector.NewAzureVaultConnector(connectionString, containerName)
	require.NoError(t, err)

	// Prepare test file content
	testFilePath := "testfile.pem"
	testFileContent := []byte("This is a test file content.")

	// Create test file
	err = os.WriteFile(testFilePath, testFileContent, 0644)
	require.NoError(t, err)

	// Clean up the test file after the test
	defer os.Remove(testFilePath)

	// Upload the file
	userId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	cryptoKeyMeta, err := abc.Upload(testFilePath, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	// Download the uploaded file
	downloadedData, err := abc.Download(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)

	// Assert that the downloaded content matches the original content
	assert.Equal(t, testFileContent, downloadedData)

	// Delete the uploaded key
	err = abc.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)
}

// TestDelete tests the Delete method of AzureVaultConnector
func TestAzureVaultConnector_Delete(t *testing.T) {
	// Prepare environment and initialize the connector
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	abc, err := connector.NewAzureVaultConnector(connectionString, containerName)
	require.NoError(t, err)

	// Prepare test file content
	testFilePath := "testfile.pem"
	testFileContent := []byte("This is a test file content.")

	// Create test file
	err = os.WriteFile(testFilePath, testFileContent, 0644)
	require.NoError(t, err)

	// Clean up the test file after the test
	defer os.Remove(testFilePath)

	// Upload the file
	userId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	cryptoKeyMeta, err := abc.Upload(testFilePath, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	// Now delete the uploaded key by ID
	err = abc.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)

	// Try downloading the key to ensure it was deleted (should fail)
	_, err = abc.Download(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	assert.Error(t, err)
}
