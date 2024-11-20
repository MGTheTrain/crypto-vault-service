package connector

import (
	"bytes"
	"context"
	"crypto_vault_service/internal/domain/blobs"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
)

// BlobConnector is an interface for interacting with Blob storage
type BlobConnector interface {
	// Upload uploads multiple files to Blob Storage and returns their metadata.
	Upload(filePaths []string) ([]*blobs.BlobMeta, error)
	// Download retrieves a blob's content by its ID and name, and returns the data as a stream.
	Download(blobId, blobName string) ([]byte, error)
	// Delete deletes a blob from Blob Storage by its ID and Name, and returns any error encountered.
	Delete(blobId, blobName string) error
}

// AzureBlobConnector is a struct that holds the Azure Blob storage client and implements the BlobConnector interfaces.
type AzureBlobConnector struct {
	Client        *azblob.Client
	ContainerName string
}

// NewAzureBlobConnector creates a new AzureBlobConnector instance using a connection string.
// It returns the connector and any error encountered during the initialization.
func NewAzureBlobConnector(connectionString string, containerName string) (*AzureBlobConnector, error) {
	client, err := azblob.NewClientFromConnectionString(connectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Blob client: %w", err)
	}

	_, err = client.CreateContainer(context.Background(), containerName, nil)
	if err != nil {
		fmt.Printf("Failed to create Azure container: %v\n", err) // The container may already exist, so we should not return an error in this case.
	}

	return &AzureBlobConnector{
		Client:        client,
		ContainerName: containerName,
	}, nil
}

// Upload uploads multiple files to Azure Blob Storage and returns their metadata.
func (abc *AzureBlobConnector) Upload(filePaths []string) ([]*blobs.BlobMeta, error) {
	var blobMeta []*blobs.BlobMeta
	blobID := uuid.New().String()

	// Iterate through all file paths and upload each file
	for _, filePath := range filePaths {
		// Open the file from the given filePath
		file, err := os.Open(filePath)
		if err != nil {
			err = fmt.Errorf("failed to open file '%s': %w", filePath, err)
			abc.rollbackUploadedBlobs(blobMeta) // Rollback previously uploaded blobs
			return nil, err
		}
		// Ensure file is closed after processing
		defer file.Close()

		// Get file info (name, size, etc.)
		fileInfo, err := file.Stat()
		if err != nil {
			err = fmt.Errorf("failed to stat file '%s': %w", filePath, err)
			abc.rollbackUploadedBlobs(blobMeta)
			return nil, err
		}

		// Read the file into a byte slice
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(file)
		if err != nil {
			err = fmt.Errorf("failed to read file '%s': %w", filePath, err)
			abc.rollbackUploadedBlobs(blobMeta)
			return nil, err
		}

		// Extract the file extension (type)
		fileExt := filepath.Ext(fileInfo.Name()) // Gets the file extension (e.g. ".txt", ".jpg")

		// Create a Blob object for metadata (Fill in missing fields)
		blob := &blobs.BlobMeta{
			ID:              blobID,
			Name:            fileInfo.Name(),
			Size:            fileInfo.Size(),
			Type:            fileExt,
			DateTimeCreated: time.Now(), // Set the current time
		}

		// Construct the full blob name (ID and Name)
		fullBlobName := fmt.Sprintf("%s/%s", blob.ID, blob.Name) // Combine ID and name to form a full path
		fullBlobName = filepath.ToSlash(fullBlobName)            // Ensure consistent slash usage across platforms

		// Upload the blob to Azure
		_, err = abc.Client.UploadBuffer(context.Background(), abc.ContainerName, fullBlobName, buf.Bytes(), nil)
		if err != nil {
			err = fmt.Errorf("failed to upload blob '%s': %w", fullBlobName, err)
			abc.rollbackUploadedBlobs(blobMeta)
			return nil, err
		}

		log.Printf("Blob '%s' uploaded successfully.\n", blob.Name)

		// Add the successfully uploaded blob to the list
		blobMeta = append(blobMeta, blob)
	}

	// Return the list of blobs after successful upload.
	return blobMeta, nil
}

// rollbackUploadedBlobs deletes the blobs that were uploaded successfully before the error occurred
func (abc *AzureBlobConnector) rollbackUploadedBlobs(blobs []*blobs.BlobMeta) {
	for _, blob := range blobs {
		err := abc.Delete(blob.ID, blob.Name)
		if err != nil {
			log.Printf("Failed to delete blob '%s' during rollback: %v", blob.Name, err)
		} else {
			log.Printf("Blob '%s' deleted during rollback.\n", blob.Name)
		}
	}
}

// Download retrieves a blob's content by its ID and name, and returns the data as a stream.
func (abc *AzureBlobConnector) Download(blobId, blobName string) ([]byte, error) {
	ctx := context.Background()

	// Construct the full blob path by combining blob ID and name
	fullBlobName := fmt.Sprintf("%s/%s", blobId, blobName) // Combine ID and name to form a full path

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
	return downloadedData.Bytes(), nil
}

// Delete deletes a blob from Azure Blob Storage by its ID and Name, and returns any error encountered.
func (abc *AzureBlobConnector) Delete(blobId, blobName string) error {
	ctx := context.Background()

	// Construct the full blob path by combining blob ID and name
	fullBlobName := fmt.Sprintf("%s/%s", blobId, blobName) // Combine ID and name to form a full path

	// Delete the blob
	_, err := abc.Client.DeleteBlob(ctx, abc.ContainerName, fullBlobName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete all blobs in %s", blobId)
	}
	fmt.Printf("Deleted all blobs in %s folder", blobId)
	return nil
}
