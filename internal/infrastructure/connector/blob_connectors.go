package connector

import (
	"bytes"
	"context"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/google/uuid"
)

// BlobConnector is an interface for interacting with Blob storage
type BlobConnector interface {
	// Upload uploads multiple files (specified by their file paths) to Blob Storage
	// and returns the metadata for each uploaded file.
	Upload(filePaths []string, userId string) ([]*blobs.BlobMeta, error)

	// UploadFromForm uploads files from form to Blob Storage
	// and returns the metadata for each uploaded byte stream.
	UploadFromForm(form *multipart.Form, userId string) ([]*blobs.BlobMeta, error)

	// Download retrieves a blob's content by its ID and name, and returns the data as a stream.
	Download(blobId, blobName string) ([]byte, error)

	// Delete deletes a blob from Blob Storage by its ID and Name, and returns any error encountered.
	Delete(blobId, blobName string) error
}

// AzureBlobConnector is a struct that holds the Azure Blob storage client and implements the BlobConnector interfaces.
type AzureBlobConnector struct {
	Client        *azblob.Client
	containerName string
	Logger        logger.Logger
}

// NewAzureBlobConnector creates a new AzureBlobConnector instance using a connection string.
// It returns the connector and any error encountered during the initialization.
func NewAzureBlobConnector(settings *settings.BlobConnectorSettings, logger logger.Logger) (*AzureBlobConnector, error) {
	if err := settings.Validate(); err != nil {
		return nil, err
	}

	client, err := azblob.NewClientFromConnectionString(settings.ConnectionString, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Azure Blob client: %w", err)
	}

	_, err = client.CreateContainer(context.Background(), settings.ContainerName, nil)
	if err != nil {
		fmt.Printf("Failed to create Azure container: %v\n", err)
	}

	return &AzureBlobConnector{
		Client:        client,
		containerName: settings.ContainerName,
		Logger:        logger,
	}, nil
}

// Upload uploads multiple files (specified by their file paths) to Blob Storage
// and returns the metadata for each uploaded file.
func (abc *AzureBlobConnector) Upload(filePaths []string, userId string) ([]*blobs.BlobMeta, error) {
	var blobMeta []*blobs.BlobMeta
	blobID := uuid.New().String()

	for _, filePath := range filePaths {

		file, err := os.Open(filePath)
		if err != nil {
			err = fmt.Errorf("failed to open file '%s': %w", filePath, err)
			abc.rollbackUploadedBlobs(blobMeta)
			return nil, err
		}

		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			err = fmt.Errorf("failed to stat file '%s': %w", filePath, err)
			abc.rollbackUploadedBlobs(blobMeta)
			return nil, err
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(file)
		if err != nil {
			err = fmt.Errorf("failed to read file '%s': %w", filePath, err)
			abc.rollbackUploadedBlobs(blobMeta)
			return nil, err
		}

		fileExt := filepath.Ext(fileInfo.Name())

		blob := &blobs.BlobMeta{
			ID:              blobID,
			Name:            fileInfo.Name(),
			Size:            fileInfo.Size(),
			Type:            fileExt,
			DateTimeCreated: time.Now(),
			UserID:          userId,
			// Size                int64
			// EncryptionAlgorithm string
			// HashAlgorithm       string
			// CryptoKey           keys.CryptoKeyMeta
			// KeyID               string
		}

		fullBlobName := fmt.Sprintf("%s/%s", blob.ID, blob.Name)
		fullBlobName = filepath.ToSlash(fullBlobName)

		_, err = abc.Client.UploadBuffer(context.Background(), abc.containerName, fullBlobName, buf.Bytes(), nil)
		if err != nil {
			err = fmt.Errorf("failed to upload blob '%s': %w", fullBlobName, err)
			abc.rollbackUploadedBlobs(blobMeta)
			return nil, err
		}

		abc.Logger.Info(fmt.Sprintf("Blob '%s' uploaded successfully", blob.Name))

		blobMeta = append(blobMeta, blob)
	}

	return blobMeta, nil
}

// UploadFromForm uploads files from form to Blob Storage
// and returns the metadata for each uploaded byte stream.
func (abc *AzureBlobConnector) UploadFromForm(form *multipart.Form, userId string) ([]*blobs.BlobMeta, error) {
	var blobMeta []*blobs.BlobMeta

	fileHeaders := form.File["files"]

	for _, fileHeader := range fileHeaders {

		blobID := uuid.New().String()

		fileExt := filepath.Ext(fileHeader.Filename)

		blob := &blobs.BlobMeta{
			ID:              blobID,
			Name:            fileHeader.Filename,
			Size:            fileHeader.Size,
			Type:            fileExt,
			DateTimeCreated: time.Now(),
			UserID:          userId,
		}

		fullBlobName := fmt.Sprintf("%s/%s", blob.ID, blob.Name)
		fullBlobName = filepath.ToSlash(fullBlobName)

		file, err := fileHeader.Open()
		if err != nil {
			err = fmt.Errorf("failed to open file '%s': %w", fullBlobName, err)
			abc.rollbackUploadedBlobs(blobMeta)
			return nil, err
		}
		defer file.Close()

		buffer := bytes.NewBuffer(make([]byte, 0))
		_, err = io.Copy(buffer, file)

		if err != nil {
			err = fmt.Errorf("failed create new buffer for '%s': %w", fullBlobName, err)
			abc.rollbackUploadedBlobs(blobMeta)
			return nil, err
		}

		_, err = abc.Client.UploadBuffer(context.Background(), abc.containerName, fullBlobName, buffer.Bytes(), nil)
		if err != nil {
			err = fmt.Errorf("failed to upload blob '%s': %w", fullBlobName, err)
			abc.rollbackUploadedBlobs(blobMeta)
			return nil, err
		}

		abc.Logger.Info(fmt.Sprintf("Blob '%s' uploaded successfully", blob.Name))

		blobMeta = append(blobMeta, blob)
	}

	return blobMeta, nil
}

// rollbackUploadedBlobs deletes the blobs that were uploaded successfully before the error occurred
func (abc *AzureBlobConnector) rollbackUploadedBlobs(blobs []*blobs.BlobMeta) {
	for _, blob := range blobs {
		err := abc.Delete(blob.ID, blob.Name)
		if err != nil {
			abc.Logger.Info(fmt.Sprintf("Failed to delete blob '%s' during rollback: %v", blob.Name, err))
		} else {
			abc.Logger.Info(fmt.Sprintf("Blob '%s' deleted during rollback", blob.Name))
		}
	}
}

// Download retrieves a blob's content by its ID and name, and returns the data as a stream.
func (abc *AzureBlobConnector) Download(blobId, blobName string) ([]byte, error) {
	ctx := context.Background()

	fullBlobName := fmt.Sprintf("%s/%s", blobId, blobName)

	get, err := abc.Client.DownloadStream(ctx, abc.containerName, fullBlobName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob '%s': %w", fullBlobName, err)
	}

	downloadedData := bytes.Buffer{}

	retryReader := get.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from blob '%s': %w", fullBlobName, err)
	}

	err = retryReader.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close retryReader for blob '%s': %w", fullBlobName, err)
	}

	abc.Logger.Info(fmt.Sprintf("Blob '%s' downloaded successfully", fullBlobName))
	return downloadedData.Bytes(), nil
}

// Delete deletes a blob from Azure Blob Storage by its ID and Name, and returns any error encountered.
func (abc *AzureBlobConnector) Delete(blobId, blobName string) error {
	ctx := context.Background()

	fullBlobName := fmt.Sprintf("%s/%s", blobId, blobName)

	_, err := abc.Client.DeleteBlob(ctx, abc.containerName, fullBlobName, nil)
	if err != nil {
		return fmt.Errorf("failed to delete blob in %s", fullBlobName)
	}

	abc.Logger.Info(fmt.Sprintf("Blob '%s' deleted successfully", fullBlobName))
	return nil
}
