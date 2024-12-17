package connector

import (
	"bytes"
	"context"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
)

// VaultConnector is an interface for interacting with custom key storage.
// The current implementation uses Azure Blob Storage, but this may be replaced
// with Azure Key Vault, AWS KMS, or any other cloud-based key management system in the future.
type VaultConnector interface {
	// Upload uploads single file to Blob Storage and returns their metadata.
	// In the future, this may be refactored to integrate with more advanced key storage systems like Key Vault.
	Upload(filePath, userId, keyType, keyAlgorihm string) (*keys.CryptoKeyMeta, error)

	// UploadFromForm uploads single file to Blob Storage
	// and returns the metadata for each uploaded byte stream.
	UploadFromForm(form *multipart.Form, userId, keyType, keyAlgorihm string) (*keys.CryptoKeyMeta, error)

	// Download retrieves a key's content by its ID and name and returns the data as a byte slice.
	Download(keyId, keyType string) ([]byte, error)

	// Delete deletes a key from Vault Storage by its ID and Name and returns any error encountered.
	Delete(keyId, keyType string) error
}

// AzureVaultConnector is a struct that implements the VaultConnector interface using Azure Blob Storage.
// This is a temporary implementation and may later be replaced with more specialized external key management systems
// like Azure Key Vault or AWS KMS.
type AzureVaultConnector struct {
	Client        *azblob.Client
	containerName string
	Logger        logger.Logger
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

	_, err = client.CreateContainer(context.Background(), settings.ContainerName, nil)
	if err != nil {
		log.Printf("Failed to create Azure container: %v\n", err)
	}

	return &AzureVaultConnector{
		Client:        client,
		containerName: settings.ContainerName,
		Logger:        logger,
	}, nil
}

// Upload uploads single file to Azure Blob Storage and returns their metadata.
// In the future, this may be refactored to integrate with more advanced key storage systems like Azure Key Vault.
func (vc *AzureVaultConnector) Upload(filePath, userId, keyType, keyAlgorihm string) (*keys.CryptoKeyMeta, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %w", filePath, err)
	}

	keyID := uuid.New().String()

	keyMeta := &keys.CryptoKeyMeta{
		ID:              keyID,
		Type:            keyType,
		Algorithm:       keyAlgorihm,
		DateTimeCreated: time.Now(),
		UserID:          userId,
	}

	fullKeyName := fmt.Sprintf("%s/%s", keyID, keyType)

	_, err = vc.Client.UploadBuffer(context.Background(), vc.containerName, fullKeyName, buf.Bytes(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to upload blob '%s' to storage: %w", fullKeyName, err)
	}

	vc.Logger.Info(fmt.Sprintf("uploaded blob %s", fullKeyName))
	return keyMeta, nil
}

// UploadFromForm uploads single file to Blob Storage
// and returns the metadata for each uploaded byte stream.
func (vc *AzureVaultConnector) UploadFromForm(form *multipart.Form, userId, keyType, keyAlgorihm string) (*keys.CryptoKeyMeta, error) {
	fileHeaders := form.File["files"]

	if len(fileHeaders) != 1 {
		err := fmt.Errorf("only 1 file can be uploaded")
		return nil, err
	}

	keyID := uuid.New().String()
	fullKeyName := fmt.Sprintf("%s/%s", keyID, keyType)

	file, err := fileHeaders[0].Open()
	if err != nil {
		err = fmt.Errorf("failed to open file '%s': %w", fullKeyName, err)
		return nil, err
	}
	defer file.Close()

	buffer := bytes.NewBuffer(make([]byte, 0))
	_, err = io.Copy(buffer, file)

	if err != nil {
		err = fmt.Errorf("failed create new buffer for '%s': %w", fullKeyName, err)
		return nil, err
	}

	keyMeta := &keys.CryptoKeyMeta{
		ID:              keyID,
		Type:            keyType,
		Algorithm:       keyAlgorihm,
		DateTimeCreated: time.Now(),
		UserID:          userId,
	}

	_, err = vc.Client.UploadBuffer(context.Background(), vc.containerName, fullKeyName, buffer.Bytes(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to upload blob '%s' to storage: %w", fullKeyName, err)
	}

	vc.Logger.Info(fmt.Sprintf("uploaded blob %s", fullKeyName))
	return keyMeta, nil
}

// Download retrieves a key's content by its ID and name and returns the data as a byte slice.
func (vc *AzureVaultConnector) Download(keyId, keyType string) ([]byte, error) {

	fullKeyName := fmt.Sprintf("%s/%s", keyId, keyType)

	ctx := context.Background()
	get, err := vc.Client.DownloadStream(ctx, vc.containerName, fullKeyName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob '%s': %w", fullKeyName, err)
	}

	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(get.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from blob '%s': %w", fullKeyName, err)
	}

	vc.Logger.Info(fmt.Sprintf("downloaded blob %s", fullKeyName))
	return downloadedData.Bytes(), nil
}

// Delete deletes a key from Azure Blob Storage by its ID and Name.
func (vc *AzureVaultConnector) Delete(keyId, keyType string) error {
	fullKeyName := fmt.Sprintf("%s/%s", keyId, keyType)

	ctx := context.Background()
	_, err := vc.Client.DeleteBlob(ctx, vc.containerName, fullKeyName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob '%s': %w", fullKeyName, err)
	}

	vc.Logger.Info(fmt.Sprintf("deleted blob %s", fullKeyName))
	return nil
}
