package connector

import (
	"bytes"
	"context"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
)

// VaultConnector is an interface for interacting with custom key storage.
// The current implementation uses Azure Blob Storage, but this may be replaced
// with Azure Key Vault, AWS KMS, or any other cloud-based key management system in the future.
type VaultConnector interface {
	// Upload uploads bytes of a single file to Blob Storage
	// and returns the metadata for each uploaded byte stream.
	Upload(bytes []byte, userId, keyPairId, keyType, keyAlgorihm string, keySize uint) (*keys.CryptoKeyMeta, error)

	// Download retrieves a key's content by its IDs and type and returns the data as a byte slice.
	Download(keyId, keyPairId, keyType string) ([]byte, error)

	// Delete deletes a key from Vault Storage by its IDs and type and returns any error encountered.
	Delete(keyId, keyPairId, keyType string) error
}

// AzureVaultConnector is a struct that implements the VaultConnector interface using Azure Blob Storage.
// This is a temporary implementation and may later be replaced with more specialized external key management systems
// like Azure Key Vault or AWS KMS.
type AzureVaultConnector struct {
	client        *azblob.Client
	containerName string
	logger        logger.Logger
}

// NewAzureVaultConnector creates a new instance of AzureVaultConnector, which connects to Azure Blob Storage.
// This method can be updated in the future to support a more sophisticated key management system like Azure Key Vault.
func NewAzureVaultConnector(settings *settings.KeyConnectorSettings, logger logger.Logger) (*AzureVaultConnector, error) {
	if err := settings.Validate(); err != nil {
		return nil, err
	}

	client, err := azblob.NewClientFromConnectionString(settings.ConnectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Blob client: %w", err)
	}

	_, _ = client.CreateContainer(context.Background(), settings.ContainerName, nil)
	// if err != nil {
	// 	log.Printf("Failed to create Azure container: %v\n", err)
	// }

	return &AzureVaultConnector{
		client:        client,
		containerName: settings.ContainerName,
		logger:        logger,
	}, nil
}

// Upload uploads bytes of a single file to Blob Storage
// and returns the metadata for each uploaded byte stream.
func (vc *AzureVaultConnector) Upload(bytes []byte, userId, keyPairId, keyType, keyAlgorihm string, keySize uint) (*keys.CryptoKeyMeta, error) {
	keyId := uuid.New().String()
	fullKeyName := fmt.Sprintf("%s/%s-%s", keyPairId, keyId, keyType)

	cryptoKeyMeta := &keys.CryptoKeyMeta{
		ID:              keyId,
		KeyPairID:       keyPairId,
		Type:            keyType,
		Algorithm:       keyAlgorihm,
		KeySize:         keySize,
		DateTimeCreated: time.Now(),
		UserID:          userId,
	}

	_, err := vc.client.UploadBuffer(context.Background(), vc.containerName, fullKeyName, bytes, nil)
	if err != nil {
		vc.rollbackUploadedBlobs(cryptoKeyMeta)
		return nil, fmt.Errorf("failed to upload blob '%s' to storage: %w", fullKeyName, err)
	}

	vc.logger.Info(fmt.Sprintf("uploaded blob %s", fullKeyName))
	return cryptoKeyMeta, nil
}

// rollbackUploadedBlobs deletes the blobs that were uploaded successfully before the error occurred
func (vc *AzureVaultConnector) rollbackUploadedBlobs(cryptoKeyMeta *keys.CryptoKeyMeta) {
	err := vc.Delete(cryptoKeyMeta.ID, cryptoKeyMeta.KeyPairID, cryptoKeyMeta.Type)
	if err != nil {
		vc.logger.Info(fmt.Sprintf("Failed to delete key '%s' during rollback: %v", cryptoKeyMeta.ID, err))
	} else {
		vc.logger.Info(fmt.Sprintf("Key '%s' deleted during rollback", cryptoKeyMeta.ID))
	}
}

// Download retrieves a key's content by its IDs and Type and returns the data as a byte slice.
func (vc *AzureVaultConnector) Download(keyId, keyPairId, keyType string) ([]byte, error) {

	fullKeyName := fmt.Sprintf("%s/%s-%s", keyPairId, keyId, keyType)

	ctx := context.Background()
	get, err := vc.client.DownloadStream(ctx, vc.containerName, fullKeyName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob '%s': %w", fullKeyName, err)
	}

	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(get.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from blob '%s': %w", fullKeyName, err)
	}

	vc.logger.Info(fmt.Sprintf("downloaded blob %s", fullKeyName))
	return downloadedData.Bytes(), nil
}

// Delete deletes a key from Azure Blob Storage by its IDs and Type.
func (vc *AzureVaultConnector) Delete(keyId, keyPairId, keyType string) error {
	fullKeyName := fmt.Sprintf("%s/%s-%s", keyPairId, keyId, keyType)

	ctx := context.Background()
	_, err := vc.client.DeleteBlob(ctx, vc.containerName, fullKeyName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob '%s': %w", fullKeyName, err)
	}

	vc.logger.Info(fmt.Sprintf("deleted blob %s", fullKeyName))
	return nil
}
