package connector

import (
	"os"
	"testing"

	"crypto_vault_service/internal/infrastructure/connector"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define a struct to represent the test context for Azure Blob Connector tests.
type AzureBlobConnectorTest struct {
	Connector       *connector.AzureBlobConnector
	TestFilePath    string
	TestFileContent []byte
}

// Helper function to create a test file
func (abt *AzureBlobConnectorTest) createTestFile(t *testing.T) {
	err := os.WriteFile(abt.TestFilePath, abt.TestFileContent, 0644)
	require.NoError(t, err)
}

// Helper function to remove the test file
func (abt *AzureBlobConnectorTest) removeTestFile(t *testing.T) {
	err := os.Remove(abt.TestFilePath)
	require.NoError(t, err)
}

// Helper function to create a new AzureBlobConnectorTest instance
func NewAzureBlobConnectorTest(t *testing.T) *AzureBlobConnectorTest {
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	abc, err := connector.NewAzureBlobConnector(connectionString, containerName)
	require.NoError(t, err)

	return &AzureBlobConnectorTest{
		Connector:       abc,
		TestFilePath:    "testfile.txt",
		TestFileContent: []byte("This is a test file content."),
	}
}

// TestUpload tests the Upload method of AzureBlobConnector
func TestUpload(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureBlobConnectorTest(t)

	// Prepare test file
	azureTest.createTestFile(t)
	// Clean up the test file and blob
	defer azureTest.removeTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	filePaths := []string{azureTest.TestFilePath}
	blobs, err := azureTest.Connector.Upload(filePaths, userId)
	require.NoError(t, err)

	// Assert that we received one blob metadata
	assert.Len(t, blobs, 1)
	blob := blobs[0]
	assert.NotEmpty(t, blob.ID)
	assert.Equal(t, azureTest.TestFilePath, blob.Name)
	assert.Equal(t, int64(len(azureTest.TestFileContent)), blob.Size)
	assert.Equal(t, ".txt", blob.Type)

	err = azureTest.Connector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestDownload tests the Download method of AzureBlobConnector
func TestDownload(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureBlobConnectorTest(t)

	// Prepare test file
	azureTest.createTestFile(t)
	// Clean up the test file and blob
	defer azureTest.removeTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	filePaths := []string{azureTest.TestFilePath}
	blobs, err := azureTest.Connector.Upload(filePaths, userId)
	require.NoError(t, err)

	// Download the uploaded file
	blob := blobs[0]
	downloadedData, err := azureTest.Connector.Download(blob.ID, blob.Name)
	require.NoError(t, err)

	// Assert that the downloaded content matches the original content
	assert.Equal(t, azureTest.TestFileContent, downloadedData)

	err = azureTest.Connector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestDelete tests the Delete method of AzureBlobConnector
func TestDelete(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureBlobConnectorTest(t)

	// Prepare test file
	azureTest.createTestFile(t)
	// Clean up the test file and blob
	defer azureTest.removeTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	filePaths := []string{azureTest.TestFilePath}
	blobs, err := azureTest.Connector.Upload(filePaths, userId)
	require.NoError(t, err)

	// Get the uploaded blob ID
	blob := blobs[0]

	// Now delete the uploaded blob by ID
	err = azureTest.Connector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)

	// Try downloading the blob to ensure it was deleted (should fail)
	_, err = azureTest.Connector.Download(blob.ID, blob.Name)
	assert.Error(t, err)
}