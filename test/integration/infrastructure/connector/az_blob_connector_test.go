package connector

import (
	"os"
	"testing"

	"crypto_vault_service/internal/infrastructure/connector"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAzureBlobConnector_Upload tests the Upload method of AzureBlobConnector
func TestAzureBlobConnector_Upload(t *testing.T) {
	// Prepare environment and initialize the connector
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	abc, err := connector.NewAzureBlobConnector(connectionString, containerName)
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
	filePaths := []string{testFilePath}
	blobs, err := abc.Upload(filePaths, userId)
	require.NoError(t, err)

	// Assert that we received one blob metadata
	assert.Len(t, blobs, 1)
	blob := blobs[0]
	assert.NotEmpty(t, blob.ID)
	assert.Equal(t, testFilePath, blob.Name)
	assert.Equal(t, int64(len(testFileContent)), blob.Size)
	assert.Equal(t, ".txt", blob.Type)

	// Delete the uploaded blob
	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestAzureBlobConnector_Download tests the Download method of AzureBlobConnector
func TestAzureBlobConnector_Download(t *testing.T) {
	// Prepare environment and initialize the connector
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	abc, err := connector.NewAzureBlobConnector(connectionString, containerName)
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
	filePaths := []string{testFilePath}
	blobs, err := abc.Upload(filePaths, userId)
	require.NoError(t, err)

	// Download the uploaded file
	blob := blobs[0]
	downloadedData, err := abc.Download(blob.ID, blob.Name)
	require.NoError(t, err)

	// Assert that the downloaded content matches the original content
	assert.Equal(t, testFileContent, downloadedData)

	// Delete the uploaded blob
	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestAzureBlobConnector_Delete tests the Delete method of AzureBlobConnector
func TestAzureBlobConnector_Delete(t *testing.T) {
	// Prepare environment and initialize the connector
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	abc, err := connector.NewAzureBlobConnector(connectionString, containerName)
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
	filePaths := []string{testFilePath}
	blobs, err := abc.Upload(filePaths, userId)
	require.NoError(t, err)

	// Get the uploaded blob ID
	blob := blobs[0]

	// Now delete the uploaded blob by ID
	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)

	// Try downloading the blob to ensure it was deleted (should fail)
	_, err = abc.Download(blob.ID, blob.Name)
	assert.Error(t, err)
}
