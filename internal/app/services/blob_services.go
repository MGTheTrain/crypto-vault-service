package services

import (
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/persistence/repository"
	"fmt"
)

// BlobUploadService implements the BlobUploadService interface for handling blob uploads
type BlobUploadService struct {
	BlobConnector  connector.BlobConnector
	BlobRepository repository.BlobRepository
}

// NewBlobUploadService creates a new instance of BlobUploadService
func NewBlobUploadService(blobConnector connector.BlobConnector, blobRepository repository.BlobRepository) *BlobUploadService {
	return &BlobUploadService{
		BlobConnector:  blobConnector,
		BlobRepository: blobRepository,
	}
}

// Upload handles the upload of blobs and stores their metadata in the database.
func (s *BlobUploadService) Upload(filePaths []string, userId string) ([]*blobs.BlobMeta, error) {

	// Use the BlobConnector to upload the files to Azure Blob Storage
	blobMeta, err := s.BlobConnector.Upload(filePaths, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to upload blobs: %w", err)
	}

	// If no blobs are uploaded, return early
	if len(blobMeta) == 0 {
		return nil, fmt.Errorf("no blobs uploaded")
	}

	// Store the metadata in the database using the BlobRepository
	for _, blob := range blobMeta {
		err := s.BlobRepository.Create(blob)
		if err != nil {
			// Rollback any previously uploaded blobs if the metadata fails to store
			// (you can call delete method to handle this as needed)
			return nil, fmt.Errorf("failed to store metadata for blob '%s': %w", blob.Name, err)
		}
	}

	// Return the metadata of uploaded blobs
	return blobMeta, nil
}

// BlobMetadataService implements the BlobMetadataService interface for retrieving and deleting blob metadata
type BlobMetadataService struct {
	BlobConnector  connector.BlobConnector
	BlobRepository repository.BlobRepository
}

// NewBlobMetadataService creates a new instance of BlobMetadataService
func NewBlobMetadataService(blobRepository repository.BlobRepository, blobConnector connector.BlobConnector) *BlobMetadataService {
	return &BlobMetadataService{
		BlobConnector:  blobConnector,
		BlobRepository: blobRepository,
	}
}

// List retrieves all blobs' metadata considering a query filter
func (s *BlobMetadataService) List(query *blobs.BlobMetaQuery) ([]*blobs.BlobMeta, error) {
	// Assuming BlobRepository has a method to query metadata, you can adapt to GORM queries.
	var blobsList []*blobs.BlobMeta

	// TBD

	return blobsList, nil
}

// GetByID retrieves a blob's metadata by its unique ID
func (s *BlobMetadataService) GetByID(blobID string) (*blobs.BlobMeta, error) {
	// Retrieve the blob metadata using the BlobRepository
	blob, err := s.BlobRepository.GetById(blobID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve blob metadata by ID '%s': %w", blobID, err)
	}
	return blob, nil
}

// DeleteByID deletes a blob and its associated metadata by ID
func (s *BlobMetadataService) DeleteByID(blobID string) error {
	// Retrieve the blob metadata to ensure it exists
	blob, err := s.BlobRepository.GetById(blobID)
	if err != nil {
		return fmt.Errorf("failed to retrieve blob metadata by ID '%s' for deletion: %w", blobID, err)
	}

	// Delete the blob from Blob Storage using the BlobConnector
	err = s.BlobRepository.DeleteById(blobID)
	if err != nil {
		return fmt.Errorf("failed to delete blob metadata by ID '%s': %w", blobID, err)
	}

	// Now, delete the actual blob from the Blob Storage
	err = s.BlobConnector.Delete(blob.ID, blob.Name)
	if err != nil {
		return fmt.Errorf("failed to delete blob '%s' from Blob Storage: %w", blob.Name, err)
	}

	return nil
}

// BlobDownloadService implements the BlobDownloadService interface for downloading blobs
type BlobDownloadService struct {
	BlobConnector connector.BlobConnector
}

// NewBlobDownloadService creates a new instance of BlobDownloadService
func NewBlobDownloadService(blobConnector connector.BlobConnector) *BlobDownloadService {
	return &BlobDownloadService{
		BlobConnector: blobConnector,
	}
}

// Download retrieves a blob's content by its ID and name
func (s *BlobDownloadService) Download(blobID, blobName string) ([]byte, error) {
	// Retrieve the blob metadata from the BlobRepository to ensure it exists
	// Here you might want to consider validating the blob's existence.
	blob, err := s.BlobConnector.Download(blobID, blobName)
	if err != nil {
		return nil, fmt.Errorf("failed to download blob '%s': %w", blobName, err)
	}

	// Return the metadata and content of the downloaded blob
	return blob, nil
}
