package connector

import (
	"crypto_vault_service/internal/infrastructure/connector"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define a struct for the test context to reuse across multiple tests
type AzureVaultConnectorTest struct {
	Connector        *connector.AzureVaultConnector
	TestFilePath     string
	TestFileContent  []byte
	ContainerName    string
	ConnectionString string
}

// Helper function to create a test file
func (abt *AzureVaultConnectorTest) createTestFile(t *testing.T) {
	err := os.WriteFile(abt.TestFilePath, abt.TestFileContent, 0644)
	require.NoError(t, err)
}

// Helper function to remove the test file
func (abt *AzureVaultConnectorTest) removeTestFile(t *testing.T) {
	err := os.Remove(abt.TestFilePath)
	require.NoError(t, err)
}

// Helper function to create a new AzureVaultConnectorTest instance
func NewAzureVaultConnectorTest(t *testing.T) *AzureVaultConnectorTest {
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	// Create the Azure Vault Connector
	abc, err := connector.NewAzureVaultConnector(connectionString, containerName)
	require.NoError(t, err)

	return &AzureVaultConnectorTest{
		Connector:        abc,
		TestFilePath:     "testfile.txt",
		TestFileContent:  []byte("This is a test file content."),
		ContainerName:    containerName,
		ConnectionString: connectionString,
	}
}

// TestUpload tests the Upload method of AzureVaultConnector
func TestAzureVaultConnector_Upload(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureVaultConnectorTest(t)

	// Prepare test file
	azureTest.createTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	uploadedKeys, err := azureTest.Connector.Upload([]string{azureTest.TestFilePath}, userId)
	require.NoError(t, err)

	// Assert that we received key metadata
	require.Len(t, uploadedKeys, 1)
	keyMeta := uploadedKeys[0]
	assert.NotEmpty(t, keyMeta.ID)
	assert.Equal(t, "testfile.txt", keyMeta.Type) // Type based on file name
	assert.Equal(t, userId, keyMeta.UserID)
	assert.WithinDuration(t, time.Now(), keyMeta.DateTimeCreated, time.Second)

	// Clean up the test file and delete the uploaded key
	azureTest.removeTestFile(t)
	err = azureTest.Connector.Delete(keyMeta.ID, keyMeta.Type)
	require.NoError(t, err)
}

// TestDownload tests the Download method of AzureVaultConnector
func TestAzureVaultConnector_Download(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureVaultConnectorTest(t)

	// Prepare test file
	azureTest.createTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	uploadedKeys, err := azureTest.Connector.Upload([]string{azureTest.TestFilePath}, userId)
	require.NoError(t, err)

	// Get the metadata of the uploaded file
	keyMeta := uploadedKeys[0]

	// Download the uploaded file
	downloadedData, err := azureTest.Connector.Download(keyMeta.ID, keyMeta.Type)
	require.NoError(t, err)

	// Assert that the downloaded content matches the original content
	assert.Equal(t, azureTest.TestFileContent, downloadedData)

	// Clean up the test file and delete the uploaded key
	azureTest.removeTestFile(t)
	err = azureTest.Connector.Delete(keyMeta.ID, keyMeta.Type)
	require.NoError(t, err)
}

// TestDelete tests the Delete method of AzureVaultConnector
func TestAzureVaultConnector_Delete(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureVaultConnectorTest(t)

	// Prepare test file
	azureTest.createTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	uploadedKeys, err := azureTest.Connector.Upload([]string{azureTest.TestFilePath}, userId)
	require.NoError(t, err)

	// Get the metadata of the uploaded file
	keyMeta := uploadedKeys[0]

	// Now delete the uploaded key by ID
	err = azureTest.Connector.Delete(keyMeta.ID, keyMeta.Type)
	require.NoError(t, err)

	// Try downloading the key to ensure it was deleted (should fail)
	_, err = azureTest.Connector.Download(keyMeta.ID, keyMeta.Type)
	assert.Error(t, err)

	// Clean up the test file (key is already deleted)
	azureTest.removeTestFile(t)
}
