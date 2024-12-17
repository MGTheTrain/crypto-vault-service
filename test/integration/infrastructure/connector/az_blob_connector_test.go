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

// AzureBlobConnectorTest encapsulates common logic for tests
type AzureBlobConnectorTest struct {
	BlobConnector *connector.AzureBlobConnector
}

func NewAzureBlobConnectorTest(t *testing.T, connectionString string, containerName string) *AzureBlobConnectorTest {
	// Create logger
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}
	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err)

	// Create blob connector
	blobConnectorSettings := &settings.BlobConnectorSettings{
		ConnectionString: connectionString,
		ContainerName:    containerName,
	}

	blobConnector, err := connector.NewAzureBlobConnector(blobConnectorSettings, logger)
	require.NoError(t, err)

	return &AzureBlobConnectorTest{
		BlobConnector: blobConnector,
	}
}

// Helper function to create test files
func createTestFile(fileName string, content []byte) error {
	err := os.WriteFile(fileName, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to create test file: %w", err)
	}
	return nil
}

// Helper function to create multipart form
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

// Tests

func TestAzureBlobConnector_UploadWithFileHeaders(t *testing.T) {
	// Setup the test helper
	abct := NewAzureBlobConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	// Test file content
	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"
	err := createTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)

	// Create multipart form with file header
	form, err := createForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	// Upload the file using UploadFromForm method
	blobs, err := abct.BlobConnector.UploadFromForm(form, userId)
	require.NoError(t, err)

	// Check the results
	require.Len(t, blobs, 1)
	blob := blobs[0]
	assert.NotEmpty(t, blob.ID)
	assert.Equal(t, testFileName, blob.Name)
	assert.Equal(t, int64(len(testFileContent)), blob.Size)
	assert.Equal(t, ".txt", blob.Type)

	// Clean up (Delete the uploaded blob)
	err = abct.BlobConnector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

func TestAzureBlobConnector_DownloadWithFileHeaders(t *testing.T) {
	// Setup the test helper
	abct := NewAzureBlobConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	// Test file content
	testFileContent := []byte("This is test file content")
	testFileName := "testfile.pem"
	err := createTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)

	// Create multipart form with file header
	form, err := createForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	// Upload the file using UploadFromForm method
	blobs, err := abct.BlobConnector.UploadFromForm(form, userId)
	require.NoError(t, err)

	blob := blobs[0]

	// Download the uploaded file
	downloadedData, err := abct.BlobConnector.Download(blob.ID, blob.Name)
	require.NoError(t, err)

	// Check if downloaded content matches the original content
	assert.Equal(t, testFileContent, downloadedData)

	// Clean up (Delete the uploaded blob)
	err = abct.BlobConnector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

func TestAzureBlobConnector_DeleteWithFileHeaders(t *testing.T) {
	// Setup the test helper
	abct := NewAzureBlobConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	// Test file content
	testFileContent := []byte("This is test file content")
	testFileName := "testfile.pem"
	err := createTestFile(testFileName, testFileContent)
	require.NoError(t, err)
	defer os.Remove(testFileName)

	// Create multipart form with file header
	form, err := createForm(testFileContent, testFileName)
	require.NoError(t, err)

	userId := uuid.New().String()
	// Upload the file using UploadFromForm method
	blobs, err := abct.BlobConnector.UploadFromForm(form, userId)
	require.NoError(t, err)

	blob := blobs[0]

	// Delete the uploaded blob
	err = abct.BlobConnector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)

	// Attempt to download the deleted blob (should result in an error)
	_, err = abct.BlobConnector.Download(blob.ID, blob.Name)
	assert.Error(t, err)
}
