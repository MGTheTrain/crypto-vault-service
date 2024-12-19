package connector

import (
	"testing"
	"time"

	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AzureVaultConnectorTest encapsulates common logic for tests
type AzureVaultConnectorTest struct {
	Connector *connector.AzureVaultConnector
}

// NewAzureVaultConnectorTest initializes and returns a new AzureVaultConnectorTest
func NewAzureVaultConnectorTest(t *testing.T, connectionString, containerName string) *AzureVaultConnectorTest {

	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}
	logger, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err)

	keyConnectorSettings := &settings.KeyConnectorSettings{
		ConnectionString: connectionString,
		ContainerName:    containerName,
	}
	abc, err := connector.NewAzureVaultConnector(keyConnectorSettings, logger)
	require.NoError(t, err)

	return &AzureVaultConnectorTest{
		Connector: abc,
	}
}

func TestAzureVaultConnector_Upload(t *testing.T) {
	helper := NewAzureVaultConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFileContent := []byte("This is a test file content.")

	userId := uuid.New().String()
	keyPairId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	keySize := 2048

	cryptoKeyMeta, err := helper.Connector.Upload(testFileContent, userId, keyPairId, keyType, keyAlgorithm, uint(keySize))
	require.NoError(t, err)

	assert.NotEmpty(t, cryptoKeyMeta.ID)
	assert.Equal(t, keyType, cryptoKeyMeta.Type)
	assert.Equal(t, userId, cryptoKeyMeta.UserID)
	assert.WithinDuration(t, time.Now(), cryptoKeyMeta.DateTimeCreated, time.Second)

	err = helper.Connector.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.KeyPairID, cryptoKeyMeta.Type)
	require.NoError(t, err)
}

func TestAzureVaultConnector_Download(t *testing.T) {
	helper := NewAzureVaultConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFileContent := []byte("This is a test file content.")

	userId := uuid.New().String()
	keyPairId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	keySize := 2048

	cryptoKeyMeta, err := helper.Connector.Upload(testFileContent, userId, keyPairId, keyType, keyAlgorithm, uint(keySize))
	require.NoError(t, err)

	downloadedData, err := helper.Connector.Download(cryptoKeyMeta.ID, cryptoKeyMeta.KeyPairID, cryptoKeyMeta.Type)
	require.NoError(t, err)

	assert.Equal(t, testFileContent, downloadedData)

	err = helper.Connector.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.KeyPairID, cryptoKeyMeta.Type)
	require.NoError(t, err)
}

func TestAzureVaultConnector_Delete(t *testing.T) {
	helper := NewAzureVaultConnectorTest(t, "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;", "testblobs")

	testFileContent := []byte("This is a test file content.")

	userId := uuid.New().String()
	keyPairId := uuid.New().String()
	keyAlgorithm := "RSA"
	keyType := "private"
	keySize := 2048

	cryptoKeyMeta, err := helper.Connector.Upload(testFileContent, userId, keyPairId, keyType, keyAlgorithm, uint(keySize))
	require.NoError(t, err)

	err = helper.Connector.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.KeyPairID, cryptoKeyMeta.Type)
	require.NoError(t, err)

	_, err = helper.Connector.Download(cryptoKeyMeta.ID, cryptoKeyMeta.KeyPairID, cryptoKeyMeta.Type)
	assert.Error(t, err)
}
