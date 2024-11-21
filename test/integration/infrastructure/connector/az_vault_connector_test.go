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

// Define a struct to represent the test context for Azure Vault Connector tests.
type AzureVaultConnectorTest struct {
	Connector       *connector.AzureVaultConnector
	TestFilePath    string
	TestFileContent []byte
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
		Connector:       abc,
		TestFilePath:    "testfile.txt",
		TestFileContent: []byte("This is a test file content."),
	}
}

// TestUpload tests the Upload method of AzureVaultConnector
func TestAzureVaultConnector_Upload(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureVaultConnectorTest(t)

	// Prepare test file
	azureTest.createTestFile(t)
	// Clean up the test file and delete the uploaded key
	defer azureTest.removeTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	filePath := azureTest.TestFilePath
	keyAlgorithm := "RSA"
	keyType := "private"
	cryptoKeyMeta, err := azureTest.Connector.Upload(filePath, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	// Assert that we received key metadata
	assert.NotEmpty(t, cryptoKeyMeta.ID)
	assert.Equal(t, keyType, cryptoKeyMeta.Type)
	assert.Equal(t, userId, cryptoKeyMeta.UserID)
	assert.WithinDuration(t, time.Now(), cryptoKeyMeta.DateTimeCreated, time.Second)

	err = azureTest.Connector.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)
}

// TestDownload tests the Download method of AzureVaultConnector
func TestAzureVaultConnector_Download(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureVaultConnectorTest(t)

	// Prepare test file
	azureTest.createTestFile(t)
	// Clean up the test file and delete the uploaded key
	defer azureTest.removeTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	filePath := azureTest.TestFilePath
	keyAlgorithm := "RSA"
	keyType := "private"
	cryptoKeyMeta, err := azureTest.Connector.Upload(filePath, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	// Get the metadata of the uploaded file

	// Download the uploaded file
	downloadedData, err := azureTest.Connector.Download(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)

	// Assert that the downloaded content matches the original content
	assert.Equal(t, azureTest.TestFileContent, downloadedData)

	err = azureTest.Connector.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)
}

// TestDelete tests the Delete method of AzureVaultConnector
func TestAzureVaultConnector_Delete(t *testing.T) {
	// Initialize test struct
	azureTest := NewAzureVaultConnectorTest(t)

	// Prepare test file
	azureTest.createTestFile(t)
	// Clean up the test file and delete the uploaded key
	defer azureTest.removeTestFile(t)

	// Upload the file
	userId := uuid.New().String()
	filePath := azureTest.TestFilePath
	keyAlgorithm := "RSA"
	keyType := "private"
	cryptoKeyMeta, err := azureTest.Connector.Upload(filePath, userId, keyType, keyAlgorithm)
	require.NoError(t, err)

	// Get the metadata of the uploaded file

	// Now delete the uploaded key by ID
	err = azureTest.Connector.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	require.NoError(t, err)

	// Try downloading the key to ensure it was deleted (should fail)
	_, err = azureTest.Connector.Download(cryptoKeyMeta.ID, cryptoKeyMeta.Type)
	assert.Error(t, err)
}
