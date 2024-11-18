package connector

import (
	"bytes"
	"context"
	"crypto_vault_service/internal/domain/model"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
)

// AzureBlobConnector is an interface for interacting with Azure Blob storage
type AzureBlobConnector interface {
	// Upload uploads multiple files to Azure Blob Storage and returns their metadata.
	Upload(filePaths []string) ([]*model.Blob, error)
	// Download retrieves a blob's content by its ID and name, and returns the data as a stream.
	Download(blobId, blobName string) (*bytes.Buffer, error)
	// DeleteById deletes a blob from Azure Blob Storage by its ID and returns any error encountered.
	DeleteById(blobId string) error
}

// AzureBlobConnectorImpl is a struct that holds the Azure Blob storage client.
type AzureBlobConnectorImpl struct {
	Client        *azblob.Client
	ContainerName string
}

// NewAzureBlobConnector creates a new AzureBlobConnectorImpl instance using a connection string.
// It returns the connector and any error encountered during the initialization.
func NewAzureBlobConnector(connectionString string, containerName string) (*AzureBlobConnectorImpl, error) {
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Blob client: %w", err)
	}

	// Create the container if it doesn't already exist
	ctx := context.Background()
	_, err = client.CreateContainer(ctx, containerName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Blob container: %w", err)
	}
	fmt.Printf("Created container named %s\n", containerName)

	return &AzureBlobConnectorImpl{
		Client:        client,
		ContainerName: containerName,
	}, nil
}

// Upload uploads multiple files to Azure Blob Storage and returns their metadata.
func (abc *AzureBlobConnectorImpl) Upload(filePaths []string) ([]*model.Blob, error) {
	var blobs []*model.Blob
	blobId := uuid.New().String()

	// Iterate through all file paths and upload each file
	for _, filePath := range filePaths {
		// Open the file from the given filePath
		file, err := os.Open(filePath)
		if err != nil {
			abc.DeleteById(blobId)
			return nil, fmt.Errorf("failed to open file '%s': %w", filePath, err)
		}
		defer file.Close()

		// Get file info (name, size, etc.)
		fileInfo, err := file.Stat()
		if err != nil {
			abc.DeleteById(blobId)
			return nil, fmt.Errorf("failed to stat file '%s': %w", filePath, err)
		}

		// Read the file into a byte slice
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(file)
		if err != nil {
			abc.DeleteById(blobId)
			return nil, fmt.Errorf("failed to read file '%s': %w", filePath, err)
		}

		// Extract the file extension (type)
		fileExt := filepath.Ext(fileInfo.Name()) // Gets the file extension (e.g. ".txt", ".jpg")

		// Create a Blob object for metadata
		blob := &model.Blob{
			ID:   blobId,
			Name: fileInfo.Name(),
			Size: fileInfo.Size(),
			Type: fileExt,
		}

		fullBlobName := fmt.Sprintf("%s/%s", blob.ID, blob.Name) // Combine ID and name to form a full path
		fullBlobName = filepath.ToSlash(fullBlobName)            // Ensure consistent slash usage across platforms

		// Upload the blob to Azure
		_, err = abc.Client.UploadBuffer(context.Background(), abc.ContainerName, fullBlobName, buf.Bytes(), nil)
		if err != nil {
			abc.DeleteById(blobId)
			return nil, fmt.Errorf("failed to upload blob '%s': %w", fullBlobName, err)
		}

		log.Printf("Blob '%s' uploaded successfully.\n", blob.Name)

		// Add the successfully uploaded blob to the list
		blobs = append(blobs, blob)
	}

	// Return the list of blobs after successful upload.
	return blobs, nil
}

// Download retrieves a blob's content by its ID and name, and returns the data as a stream.
func (abc *AzureBlobConnectorImpl) Download(blobId, blobName string) (*bytes.Buffer, error) {
	ctx := context.Background()

	// Construct the full blob path by combining blob ID and name
	fullBlobName := fmt.Sprintf("%s/%s", blobId, blobName) // Combine ID and name to form a full path
	fullBlobName = filepath.ToSlash(fullBlobName)          // Ensure consistent slash usage across platforms

	// Download the blob as a stream
	get, err := abc.Client.DownloadStream(ctx, abc.ContainerName, fullBlobName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob '%s': %w", fullBlobName, err)
	}

	// Prepare the buffer to hold the downloaded data
	downloadedData := bytes.Buffer{}

	// Create a retryable reader in case of network or temporary failures
	retryReader := get.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from blob '%s': %w", fullBlobName, err)
	}

	// Close the retryReader stream after reading
	err = retryReader.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close retryReader for blob '%s': %w", fullBlobName, err)
	}

	// Return the buffer containing the downloaded data
	return &downloadedData, nil
}

// DeleteById deletes a blob from Azure Blob Storage by its ID and returns any error encountered.
func (abc *AzureBlobConnectorImpl) DeleteById(blobId string) error {
	ctx := context.Background()

	// Delete the blob
	_, err := abc.Client.DeleteBlob(ctx, abc.ContainerName, blobId, nil)
	if err != nil {
		return fmt.Errorf("failed to delete all blobs in %s", blobId)
	}
	fmt.Printf("Deleted all blobs in %s folder", blobId)
	return nil
}
