package connector

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"os"
	"testing"

	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createForm(content []byte, fileName string) (*multipart.Form, error) {

	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)

	fileWriter, err := writer.CreateFormFile("files", fileName)
	if err != nil {
		return nil, err
	}

	_, err = fileWriter.Write(content)
	if err != nil {
		return nil, err
	}

	writer.Close()

	mr := multipart.NewReader(&buf, writer.Boundary())

	form, err := mr.ReadForm(10 << 20)
	if err != nil {
		return nil, err
	}

	return form, nil
}

// Helper function to create test files
func createTestFile(fileName string, content []byte) error {
	err := os.WriteFile(fileName, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	return nil
}

// TestAzureBlobConnector_UploadWithFileHeaders tests the UploadWithFileHeaders method of AzureBlobConnector
func TestAzureBlobConnector_UploadWithFileHeaders(t *testing.T) {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err)

	blobConnectorSettings := &settings.BlobConnectorSettings{
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}

	abc, err := connector.NewAzureBlobConnector(blobConnectorSettings, logger)
	require.NoError(t, err)

	// Test file content
	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"
	err = createTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)

	// Create multipart form with file header
	form, err := createForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	// Upload the file using UploadFromForm method
	blobs, err := abc.UploadFromForm(form, userId)
	require.NoError(t, err)

	// Check the results
	require.Len(t, blobs, 1)
	blob := blobs[0]
	assert.NotEmpty(t, blob.ID)
	assert.Equal(t, testFileName, blob.Name)
	assert.Equal(t, int64(len(testFileContent)), blob.Size)
	assert.Equal(t, ".txt", blob.Type)

	// Clean up (Delete the uploaded blob)
	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestAzureBlobConnector_DownloadWithFileHeaders tests the Download method of AzureBlobConnector
func TestAzureBlobConnector_DownloadWithFileHeaders(t *testing.T) {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err)

	blobConnectorSettings := &settings.BlobConnectorSettings{
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}

	abc, err := connector.NewAzureBlobConnector(blobConnectorSettings, logger)
	require.NoError(t, err)

	// Test file content
	testFileContent := []byte("This is test file content")
	testFileName := "testfile.pem"
	err = createTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)

	// Create multipart form with file header
	form, err := createForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	// Upload the file using UploadFromForm method
	blobs, err := abc.UploadFromForm(form, userId)
	require.NoError(t, err)

	blob := blobs[0]

	// Download the uploaded file
	downloadedData, err := abc.Download(blob.ID, blob.Name)
	require.NoError(t, err)

	// Check if downloaded content matches the original content
	assert.Equal(t, testFileContent, downloadedData)

	// Clean up (Delete the uploaded blob)
	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

// TestAzureBlobConnector_DeleteWithFileHeaders tests the Delete method of AzureBlobConnector
func TestAzureBlobConnector_DeleteWithFileHeaders(t *testing.T) {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err)

	blobConnectorSettings := &settings.BlobConnectorSettings{
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}

	abc, err := connector.NewAzureBlobConnector(blobConnectorSettings, logger)
	require.NoError(t, err)

	// Test file content
	testFileContent := []byte("This is test file content")
	testFileName := "testfile.pem"
	err = createTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)

	// Create multipart form with file header
	form, err := createForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	// Upload the file using UploadFromForm method
	blobs, err := abc.UploadFromForm(form, userId)
	require.NoError(t, err)

	blob := blobs[0]

	// Delete the uploaded blob
	err = abc.Delete(blob.ID, blob.Name)
	require.NoError(t, err)

	// Attempt to download the deleted blob (should result in an error)
	_, err = abc.Download(blob.ID, blob.Name)
	assert.Error(t, err)
}
