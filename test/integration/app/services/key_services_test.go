package services

import (
	"crypto_vault_service/internal/app/services"
	"crypto_vault_service/internal/infrastructure/connector"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// Helper function to create a test file
func createTestFile(testFilePath string, testFileContent []byte) error {
	err := os.WriteFile(testFilePath, testFileContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to create test file at %s: %v", testFilePath, err)
	}
	return nil
}

// Helper function to remove the test file
func removeTestFile(testFilePath string) error {
	err := os.Remove(testFilePath)
	if err != nil {
		return fmt.Errorf("failed to remove test file at %s: %v", testFilePath, err)
	}
	return nil
}

// Test case for successful file upload and metadata creation
func TestCryptoKeyUploadService_Upload_Success(t *testing.T) {
	// Set up environment variable
	err := os.Setenv("DB_TYPE", "postgres")
	require.NoError(t, err, "Error setting environment variable")

	// Set up test context and ensure proper teardown
	ctx := setupTestDB(t)
	defer teardownTestDB(t, ctx)

	// Initialize Vault Connector
	connectionString := "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;"
	containerName := "testblobs"
	vaultConnector, err := connector.NewAzureVaultConnector(connectionString, containerName)
	require.NoError(t, err, "Error creating vault connector")

	// Set up the CryptoKeyUploadService
	cryptoKeyUploadService := &services.CryptoKeyUploadService{
		VaultConnector: vaultConnector,
		CryptoKeyRepo:  ctx.CryptoKeyRepo,
	}

	// Prepare test file
	testFilePath := "testfile.pem"
	testFileContent := []byte("This is a test file content.")
	err = createTestFile(testFilePath, testFileContent)
	require.NoError(t, err, "Error creating test file")
	defer removeTestFile(testFilePath)

	// Call the method under test
	userId := uuid.New().String()
	keyType := "private"
	keyAlgorithm := "EC"
	keyMeta, err := cryptoKeyUploadService.Upload(testFilePath, userId, keyType, keyAlgorithm)

	// Assert that no error occurred during the upload
	require.NoError(t, err, "Error uploading file")

	// Assert that keyMeta is not nil and contains expected data
	require.NotNil(t, keyMeta, "Key metadata should not be nil")
	require.NotEmpty(t, keyMeta.ID, "KeyID should not be empty")
	require.Equal(t, userId, keyMeta.UserID, "UserID does not match")
}
