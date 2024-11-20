package repository

import (
	"crypto_vault_service/internal/domain/blobs"
	"fmt"

	"gorm.io/gorm"
)

// BlobRepository defines the interface for Blob-related operations
type BlobRepository interface {
	Create(blob *blobs.BlobMeta) error
	GetById(blobID string) (*blobs.BlobMeta, error)
	UpdateById(blob *blobs.BlobMeta) error
	DeleteById(blobID string) error
}

// GormBlobRepository is the implementation of the BlobRepository interface
type GormBlobRepository struct {
	DB *gorm.DB
}

// Create adds a new Blob to the database
func (r *GormBlobRepository) Create(blob *blobs.BlobMeta) error {
	// Validate the Blob before saving
	if err := blob.Validate(); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	// Save the blob to the database
	if err := r.DB.Create(&blob).Error; err != nil {
		return fmt.Errorf("failed to create blob: %w", err)
	}
	return nil
}

// GetById retrieves a Blob by its ID from the database
func (r *GormBlobRepository) GetById(blobID string) (*blobs.BlobMeta, error) {
	var blob blobs.BlobMeta
	if err := r.DB.Where("id = ?", blobID).First(&blob).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("blob with ID %s not found", blobID)
		}
		return nil, fmt.Errorf("failed to fetch blob: %w", err)
	}
	return &blob, nil
}

// UpdateById updates an existing Blob in the database
func (r *GormBlobRepository) UpdateById(blob *blobs.BlobMeta) error {
	// Validate the Blob before updating
	if err := blob.Validate(); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	// Update the blob in the database
	if err := r.DB.Save(&blob).Error; err != nil {
		return fmt.Errorf("failed to update blob: %w", err)
	}
	return nil
}

// DeleteById removes a Blob from the database by its ID
func (r *GormBlobRepository) DeleteById(blobID string) error {
	if err := r.DB.Where("id = ?", blobID).Delete(&blobs.BlobMeta{}).Error; err != nil {
		return fmt.Errorf("failed to delete blob: %w", err)
	}
	return nil
}
