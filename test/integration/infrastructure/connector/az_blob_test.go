package connector

import (
	"os"
	"testing"

	"crypto_vault_service/internal/infrastructure/connector"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var connectionString = "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"

var containerName = "blobs"

// Helper function to create a test file
func createTestFile(t *testing.T, filePath string, content []byte) {
	err := os.WriteFile(filePath, content, 0644)
	require.NoError(t, err)
}

// TestUpload tests the Upload method of AzureBlobConnectorImpl
func TestUpload(t *testing.T) {
	// Create a connector instance using a local Azure Blob emulator connection string
	abc, err := connector.NewAzureBlobConnector(connectionString, containerName)
	require.NoError(t, err)

	// Prepare test files
	testFilePath := "testfile.txt"
	testContent := []byte("This is a test file content.")
	createTestFile(t, testFilePath, testContent)

	// Upload the file
	blobs, err := abc.Upload([]string{testFilePath})
	require.NoError(t, err)

	// Assert that we received one blob metadata
	assert.Len(t, blobs, 1)
	blob := blobs[0]
	assert.NotEmpty(t, blob.ID)
	assert.Equal(t, "testfile.txt", blob.Name)
	assert.Equal(t, int64(len(testContent)), blob.Size)
	assert.Equal(t, ".txt", blob.Type)

	// Clean up the test file
	err = os.Remove(testFilePath)
	require.NoError(t, err)

	// Clean up the blob in the Azure Blob storage (delete by ID)
	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestDownload tests the Download method of AzureBlobConnectorImpl
func TestDownload(t *testing.T) {
	// Create a connector instance using a local Azure Blob emulator connection string
	abc, err := connector.NewAzureBlobConnector(connectionString, containerName)
	require.NoError(t, err)

	// Upload a test file
	testFilePath := "testfile.txt"
	testContent := []byte("This is a test file content.")
	createTestFile(t, testFilePath, testContent)

	blobs, err := abc.Upload([]string{testFilePath})
	require.NoError(t, err)

	// Download the uploaded file
	blob := blobs[0]
	downloadedData, err := abc.Download(blob.ID, blob.Name)
	require.NoError(t, err)

	// Assert that the downloaded content is the same as the original file content
	assert.Equal(t, string(testContent), downloadedData.String())

	// Clean up the test file
	err = os.Remove(testFilePath)
	require.NoError(t, err)

	// Clean up the blob in the Azure Blob storage (delete by ID)
	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestDelete tests the Delete method of AzureBlobConnectorImpl
func TestDelete(t *testing.T) {
	// Create a connector instance using a local Azure Blob emulator connection string
	abc, err := connector.NewAzureBlobConnector(connectionString, containerName)
	require.NoError(t, err)

	// Upload a test file
	testFilePath := "testfile.txt"
	testContent := []byte("This is a test file content.")
	createTestFile(t, testFilePath, testContent)

	blobs, err := abc.Upload([]string{testFilePath})
	require.NoError(t, err)

	// Get the uploaded blob ID
	blob := blobs[0]

	// Now delete the uploaded blob by ID
	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)

	// Try downloading the blob to ensure it was deleted (should fail)
	_, err = abc.Download(blob.ID, blob.Name)
	assert.Error(t, err)

	// Clean up the test file
	err = os.Remove(testFilePath)
	require.NoError(t, err)
}
