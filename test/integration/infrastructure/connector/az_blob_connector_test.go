package connector

import (
	"os"
	"testing"

	"crypto_vault_service/internal/infrastructure/connector"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define a struct for the test context to reuse across multiple tests
type AzureBlobTest struct {
	Connector       *connector.AzureBlobConnector
	TestFilePath    string
	TestFileContent []byte
	ContainerName   string
}

// Helper function to create a test file
func (abt *AzureBlobTest) createTestFile(t *testing.T) {
	err := os.WriteFile(abt.TestFilePath, abt.TestFileContent, 0644)
	require.NoError(t, err)
}

// Helper function to remove the test file
func (abt *AzureBlobTest) removeTestFile(t *testing.T) {
	err := os.Remove(abt.TestFilePath)
	require.NoError(t, err)
}

// Helper function to create a new AzureBlobTest instance
func NewAzureBlobTest(t *testing.T) *AzureBlobTest {
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	abc, err := connector.NewAzureBlobConnector(connectionString, containerName)
	require.NoError(t, err)

	return &AzureBlobTest{
		Connector:       abc,
		TestFilePath:    "testfile.txt",
		TestFileContent: []byte("This is a test file content."),
		ContainerName:   containerName,
	}
}

// TestUpload tests the Upload method of AzureBlobConnector
func TestUpload(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureBlobTest(t)

	// Prepare test file
	azureTest.createTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	blobs, err := azureTest.Connector.Upload([]string{azureTest.TestFilePath}, userId)
	require.NoError(t, err)

	// Assert that we received one blob metadata
	assert.Len(t, blobs, 1)
	blob := blobs[0]
	assert.NotEmpty(t, blob.ID)
	assert.Equal(t, "testfile.txt", blob.Name)
	assert.Equal(t, int64(len(azureTest.TestFileContent)), blob.Size)
	assert.Equal(t, ".txt", blob.Type)

	// Clean up the test file and blob
	azureTest.removeTestFile(t)
	err = azureTest.Connector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestDownload tests the Download method of AzureBlobConnector
func TestDownload(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureBlobTest(t)

	// Prepare test file
	azureTest.createTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	blobs, err := azureTest.Connector.Upload([]string{azureTest.TestFilePath}, userId)
	require.NoError(t, err)

	// Download the uploaded file
	blob := blobs[0]
	downloadedData, err := azureTest.Connector.Download(blob.ID, blob.Name)
	require.NoError(t, err)

	// Assert that the downloaded content matches the original content
	assert.Equal(t, azureTest.TestFileContent, downloadedData)

	// Clean up the test file and blob
	azureTest.removeTestFile(t)
	err = azureTest.Connector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestDelete tests the Delete method of AzureBlobConnector
func TestDelete(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureBlobTest(t)

	// Prepare test file
	azureTest.createTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	blobs, err := azureTest.Connector.Upload([]string{azureTest.TestFilePath}, userId)
	require.NoError(t, err)

	// Get the uploaded blob ID
	blob := blobs[0]

	// Now delete the uploaded blob by ID
	err = azureTest.Connector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)

	// Try downloading the blob to ensure it was deleted (should fail)
	_, err = azureTest.Connector.Download(blob.ID, blob.Name)
	assert.Error(t, err)

	// Clean up the test file
	azureTest.removeTestFile(t)
}
