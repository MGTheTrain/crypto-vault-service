package connector

import (
	"testing"

	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AzureBlobConnectorTest encapsulates common logic for tests
type AzureBlobConnectorTest struct {
	BlobConnector *connector.AzureBlobConnector
}

func NewAzureBlobConnectorTest(t *testing.T, connectionString string, containerName string) *AzureBlobConnectorTest {

	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}
	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err)
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

func TestAzureBlobConnector_Upload(t *testing.T) {

	abct := NewAzureBlobConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.txt"
	form, err := helpers.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()

	blobs, err := abct.BlobConnector.Upload(form, userId)
	require.NoError(t, err)

	require.Len(t, blobs, 1)
	blob := blobs[0]
	assert.NotEmpty(t, blob.ID)
	assert.Equal(t, testFileName, blob.Name)
	assert.Equal(t, int64(len(testFileContent)), blob.Size)
	assert.Equal(t, ".txt", blob.Type)

	err = abct.BlobConnector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

func TestAzureBlobConnector_Download(t *testing.T) {

	abct := NewAzureBlobConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.pem"
	form, err := helpers.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()

	blobs, err := abct.BlobConnector.Upload(form, userId)
	require.NoError(t, err)

	blob := blobs[0]

	downloadedData, err := abct.BlobConnector.Download(blob.ID, blob.Name)
	require.NoError(t, err)

	assert.Equal(t, testFileContent, downloadedData)

	err = abct.BlobConnector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)
}

func TestAzureBlobConnector_Delete(t *testing.T) {

	abct := NewAzureBlobConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFileContent := []byte("This is test file content")
	testFileName := "testfile.pem"
	form, err := helpers.CreateTestFileAndForm(t, testFileName, testFileContent)
	require.NoError(t, err)

	userId := uuid.New().String()

	blobs, err := abct.BlobConnector.Upload(form, userId)
	require.NoError(t, err)

	blob := blobs[0]

	err = abct.BlobConnector.Delete(blob.ID, blob.Name)
	require.NoError(t, err)

	_, err = abct.BlobConnector.Download(blob.ID, blob.Name)
	assert.Error(t, err)
}
