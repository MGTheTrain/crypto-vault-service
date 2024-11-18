package repository

import (
	"fmt"

	"crypto_vault_service/internal/domain/model"

	"gorm.io/gorm"
)

// BlobRepository defines the interface for Blob-related operations
type BlobRepository interface {
	Create(blob *model.Blob) error
	GetById(blobID string) (*model.Blob, error)
	UpdateById(blob *model.Blob) error
	DeleteById(blobID string) error
}

// BlobRepositoryImpl is the implementation of the BlobRepository interface
type BlobRepositoryImpl struct {
	DB *gorm.DB
}

// Create adds a new Blob to the database
func (r *BlobRepositoryImpl) Create(blob *model.Blob) error {
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
func (r *BlobRepositoryImpl) GetById(blobID string) (*model.Blob, error) {
	var blob model.Blob
	if err := r.DB.Where("blob_id = ?", blobID).First(&blob).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("blob with ID %s not found", blobID)
		}
		return nil, fmt.Errorf("failed to fetch blob: %w", err)
	}
	return &blob, nil
}

// UpdateById updates an existing Blob in the database
func (r *BlobRepositoryImpl) UpdateById(blob *model.Blob) error {
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
func (r *BlobRepositoryImpl) DeleteById(blobID string) error {
	if err := r.DB.Where("blob_id = ?", blobID).Delete(&model.Blob{}).Error; err != nil {
		return fmt.Errorf("failed to delete blob: %w", err)
	}
	return nil
}
