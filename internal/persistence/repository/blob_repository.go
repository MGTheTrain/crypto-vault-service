package repository

import (
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/infrastructure/logger"
	"fmt"

	"gorm.io/gorm"
)

// BlobRepository defines the interface for Blob-related operations
type BlobRepository interface {
	Create(blob *blobs.BlobMeta) error
	List(query *blobs.BlobMetaQuery) ([]*blobs.BlobMeta, error)
	GetById(blobId string) (*blobs.BlobMeta, error)
	UpdateById(blob *blobs.BlobMeta) error
	DeleteById(blobId string) error
}

// GormBlobRepository is the implementation of the BlobRepository interface
type GormBlobRepository struct {
	DB     *gorm.DB
	Logger logger.Logger
}

// NewGormBlobRepository creates a new GormBlobRepository instance
func NewGormBlobRepository(db *gorm.DB, logger logger.Logger) (*GormBlobRepository, error) {

	return &GormBlobRepository{
		DB:     db,
		Logger: logger,
	}, nil
}

// Create adds a new Blob to the database
func (r *GormBlobRepository) Create(blob *blobs.BlobMeta) error {
	// Validate the Blob before saving
	if err := blob.Validate(); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	if err := r.DB.Create(&blob).Error; err != nil {
		return fmt.Errorf("failed to create blob: %w", err)
	}
	r.Logger.Info(fmt.Sprintf("Created blob metadata with id %s", blob.ID))
	return nil
}

func (r *GormBlobRepository) List(query *blobs.BlobMetaQuery) ([]*blobs.BlobMeta, error) {
	// Validate the query parameters before using them
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	// Start building the query
	var blobMetas []*blobs.BlobMeta
	dbQuery := r.DB.Model(&blobs.BlobMeta{})

	// Apply filters based on the query
	if query.Name != "" {
		dbQuery = dbQuery.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Size > 0 {
		dbQuery = dbQuery.Where("size = ?", query.Size)
	}
	if query.Type != "" {
		dbQuery = dbQuery.Where("type = ?", query.Type)
	}
	if !query.DateTimeCreated.IsZero() {
		dbQuery = dbQuery.Where("date_time_created >= ?", query.DateTimeCreated)
	}

	// Sorting
	if query.SortBy != "" {
		order := query.SortOrder // Default to ascending if not specified
		if query.SortOrder == "" {
			order = "asc"
		}
		dbQuery = dbQuery.Order(fmt.Sprintf("%s %s", query.SortBy, order))
	}

	// Pagination
	if query.Limit > 0 {
		dbQuery = dbQuery.Limit(query.Limit)
	}
	if query.Offset > 0 {
		dbQuery = dbQuery.Offset(query.Offset)
	}

	// Execute the query
	if err := dbQuery.Find(&blobMetas).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch blobs: %w", err)
	}

	// Return the list of blob metadata
	return blobMetas, nil
}

// GetById retrieves a Blob by its ID from the database
func (r *GormBlobRepository) GetById(blobId string) (*blobs.BlobMeta, error) {
	var blob blobs.BlobMeta
	if err := r.DB.Where("id = ?", blobId).First(&blob).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("blob with ID %s not found", blobId)
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

	if err := r.DB.Save(&blob).Error; err != nil {
		return fmt.Errorf("failed to update blob: %w", err)
	}
	r.Logger.Info(fmt.Sprintf("Updated blob metadata with id %s", blob.ID))
	return nil
}

// DeleteById removes a Blob from the database by its ID
func (r *GormBlobRepository) DeleteById(blobId string) error {
	if err := r.DB.Where("id = ?", blobId).Delete(&blobs.BlobMeta{}).Error; err != nil {
		return fmt.Errorf("failed to delete blob: %w", err)
	}
	r.Logger.Info(fmt.Sprintf("Deleted blob metadata with id %s", blobId))
	return nil
}
