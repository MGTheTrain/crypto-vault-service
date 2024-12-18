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

	blobMetas, err := s.BlobConnector.Upload(filePaths, userId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if len(blobMetas) == 0 {
		return nil, fmt.Errorf("no blobs uploaded")
	}

	for _, blobMeta := range blobMetas {
		err := s.BlobRepository.Create(blobMeta)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	return blobMetas, nil
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
	var blobMetas []*blobs.BlobMeta

	// TBD

	return blobMetas, nil
}

// GetByID retrieves a blob's metadata by its unique ID
func (s *BlobMetadataService) GetByID(blobID string) (*blobs.BlobMeta, error) {
	// Retrieve the blob metadata using the BlobRepository
	blobMeta, err := s.BlobRepository.GetById(blobID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return blobMeta, nil
}

// DeleteByID deletes a blob and its associated metadata by ID
func (s *BlobMetadataService) DeleteByID(blobID string) error {
	// Retrieve the blob metadata to ensure it exists
	blobMeta, err := s.BlobRepository.GetById(blobID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Delete the blob from Blob Storage using the BlobConnector
	err = s.BlobRepository.DeleteById(blobID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// Now, delete the actual blob from the Blob Storage
	err = s.BlobConnector.Delete(blobMeta.ID, blobMeta.Name)
	if err != nil {
		return fmt.Errorf("%w", err)
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
		return nil, fmt.Errorf("%w", err)
	}

	return blob, nil
}
