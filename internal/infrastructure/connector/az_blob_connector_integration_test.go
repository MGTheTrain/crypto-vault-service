//go:build integration
// +build integration

package connector

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"testing"

	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/infrastructure/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create test files
func CreateTestFile(fileName string, content []byte) error {
	err := os.WriteFile(fileName, content, 0600)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	return nil
}

// Helper function to create a test file and form
func CreateTestFileAndForm(t *testing.T, fileName string, fileContent []byte) (*multipart.Form, error) {
	err := CreateTestFile(fileName, fileContent)
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := os.Remove(fileName); err != nil {
			t.Logf("failed to remove temporary file %s: %v", fileName, err)
		}
	})

	form, err := utils.CreateForm(fileContent, fileName)
	require.NoError(t, err)

	return form, nil
}

// AzureBlobConnectorTest encapsulates common logic for tests
type AzureBlobConnectorTest struct {
	blobConnector BlobConnector
}

func NewAzureBlobConnectorTest(t *testing.T, cloudProvider, connectionString string, containerName string) *AzureBlobConnectorTest {

	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}
	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err)
	blobConnectorSettings := &settings.BlobConnectorSettings{
		CloudProvider:    cloudProvider,
		ConnectionString: connectionString,
		ContainerName:    containerName,
	}

	ctx := context.Background()
	blobConnector, err := NewAzureBlobConnector(ctx, blobConnectorSettings, logger)
	require.NoError(t, err)

	return &AzureBlobConnectorTest{
		blobConnector: blobConnector,
	}
}

func TestAzureBlobConnector_Upload(t *testing.T) {

	abct := NewAzureBlobConnectorTest(t, "azure", "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"
	form, err := CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()

	var encryptionKeyId *string = nil
	var signKeyId *string = nil
	ctx := context.Background()

	blobs, err := abct.blobConnector.Upload(ctx, form, userId, encryptionKeyId, signKeyId)
	require.NoError(t, err)

	require.Len(t, blobs, 1)
	blob := blobs[0]
	assert.NotEmpty(t, blob.ID)
	assert.Equal(t, testFileName, blob.Name)
	assert.Equal(t, int64(len(testFileContent)), blob.Size)
	assert.Equal(t, ".txt", blob.Type)

	err = abct.blobConnector.Delete(ctx, blob.ID, blob.Name)
	require.NoError(t, err)
}

func TestAzureBlobConnector_Download(t *testing.T) {

	abct := NewAzureBlobConnectorTest(t, "azure", "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.pem"
	form, err := CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()

	var encryptionKeyId *string = nil
	var signKeyId *string = nil
	ctx := context.Background()
	blobs, err := abct.blobConnector.Upload(ctx, form, userId, encryptionKeyId, signKeyId)
	require.NoError(t, err)

	blob := blobs[0]

	downloadedData, err := abct.blobConnector.Download(ctx, blob.ID, blob.Name)
	require.NoError(t, err)

	assert.Equal(t, testFileContent, downloadedData)

	err = abct.blobConnector.Delete(ctx, blob.ID, blob.Name)
	require.NoError(t, err)
}

func TestAzureBlobConnector_Delete(t *testing.T) {

	abct := NewAzureBlobConnectorTest(t, "azure", "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.pem"
	form, err := CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()

	var encryptionKeyId *string = nil
	var signKeyId *string = nil
	ctx := context.Background()

	blobs, err := abct.blobConnector.Upload(ctx, form, userId, encryptionKeyId, signKeyId)
	require.NoError(t, err)

	blob := blobs[0]

	err = abct.blobConnector.Delete(ctx, blob.ID, blob.Name)
	require.NoError(t, err)

	_, err = abct.blobConnector.Download(ctx, blob.ID, blob.Name)
	assert.Error(t, err)
}
