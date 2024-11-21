package connector

import (
	"bytes"
	"context"
	"crypto_vault_service/internal/domain/keys"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
)

// VaultConnector is an interface for interacting with custom key storage.
// The current implementation uses Azure Blob Storage, but this may be replaced
// with Azure Key Vault, AWS KMS, or any other cloud-based key management system in the future.
type VaultConnector interface {
	// Upload uploads multiple files to the Vault Storage and returns their metadata.
	Upload(filePaths []string, userId string) ([]*keys.CryptoKeyMeta, error)

	// Download retrieves a key's content by its ID and name, and returns the data as a byte slice.
	Download(blobId, blobName string) ([]byte, error)

	// Delete deletes a key from Vault Storage by its ID and Name, and returns any error encountered.
	Delete(blobId, blobName string) error
}

// AzureVaultConnector is a struct that implements the VaultConnector interface using Azure Blob Storage.
// This is a temporary implementation and may later be replaced with more specialized external key management systems
// like Azure Key Vault or AWS KMS.
type AzureVaultConnector struct {
	Client        *azblob.Client
	ContainerName string
}

// NewAzureVaultConnector creates a new instance of AzureVaultConnector, which connects to Azure Blob Storage.
// This method can be updated in the future to support a more sophisticated key management system like Azure Key Vault.
func NewAzureVaultConnector(connectionString string, containerName string) (*AzureVaultConnector, error) {
	// Create a new Azure Blob client using the provided connection string
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Blob client: %w", err)
	}

	_, err = client.CreateContainer(context.Background(), containerName, nil)
	if err != nil {
		log.Printf("Failed to create Azure container: %v\n", err) // The container may already exist, so we should not return an error in this case.
	}

	return &AzureVaultConnector{
		Client:        client,
		ContainerName: containerName,
	}, nil
}

// Upload uploads multiple files to Azure Blob Storage and returns their metadata.
// In the future, this may be refactored to integrate with more advanced key storage systems like Azure Key Vault.
func (vc *AzureVaultConnector) Upload(filePaths []string, userId string) ([]*keys.CryptoKeyMeta, error) {
	var keyMetas []*keys.CryptoKeyMeta

	// Iterate through all file paths and upload each file
	for _, filePath := range filePaths {
		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file '%s': %w", filePath, err)
		}
		defer file.Close()

		// Get file information (name, size, etc.)
		fileInfo, err := file.Stat()
		if err != nil {
			return nil, fmt.Errorf("failed to stat file '%s': %w", filePath, err)
		}

		// Read the file content into a buffer
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file '%s': %w", filePath, err)
		}

		// Generate a unique ID for the key
		keyID := uuid.New().String()

		// Create metadata for the uploaded key
		keyMeta := &keys.CryptoKeyMeta{
			ID:              keyID,
			Type:            fileInfo.Name(), // one of asymmetric-public-key, asymmetric-private-key, symmetric-key
			DateTimeCreated: time.Now(),
			UserID:          userId,
		}

		// Construct the full blob name (ID and Name)
		fullBlobName := fmt.Sprintf("%s/%s", keyID, fileInfo.Name())

		// Upload the blob (file) to Azure Blob Storage
		_, err = vc.Client.UploadBuffer(context.Background(), vc.ContainerName, fullBlobName, buf.Bytes(), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to upload blob '%s' to storage: %w", fullBlobName, err)
		}

		// Add the metadata to the list
		keyMetas = append(keyMetas, keyMeta)
	}

	// Return the metadata of the uploaded keys
	return keyMetas, nil
}

// Download retrieves a key's content by its ID and name, and returns the data as a byte slice.
func (vc *AzureVaultConnector) Download(blobId, blobName string) ([]byte, error) {
	// Construct the full blob path by combining blob ID and name
	fullBlobName := fmt.Sprintf("%s/%s", blobId, blobName)

	// Download the blob from Azure Blob Storage
	ctx := context.Background()
	get, err := vc.Client.DownloadStream(ctx, vc.ContainerName, fullBlobName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob '%s': %w", fullBlobName, err)
	}

	// Read the content into a buffer
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(get.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from blob '%s': %w", fullBlobName, err)
	}

	// Return the downloaded data
	return downloadedData.Bytes(), nil
}

// Delete deletes a key from Azure Blob Storage by its ID and Name.
func (vc *AzureVaultConnector) Delete(blobId, blobName string) error {
	// Construct the full blob path by combining blob ID and name
	fullBlobName := fmt.Sprintf("%s/%s", blobId, blobName)

	// Delete the blob from Azure Blob Storage
	ctx := context.Background()
	_, err := vc.Client.DeleteBlob(ctx, vc.ContainerName, fullBlobName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob '%s': %w", fullBlobName, err)
	}

	// Log the successful deletion
	log.Printf("Deleted blob '%s' from storage.\n", fullBlobName)
	return nil
}
