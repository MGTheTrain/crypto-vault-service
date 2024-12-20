package repository

import (
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/logger"
	"fmt"

	"gorm.io/gorm"
)

// CryptoKeyRepository defines the interface for CryptoKey-related operations
type CryptoKeyRepository interface {
	Create(key *keys.CryptoKeyMeta) error
	List(query *keys.CryptoKeyQuery) ([]*keys.CryptoKeyMeta, error)
	GetByID(keyId string) (*keys.CryptoKeyMeta, error)
	UpdateByID(key *keys.CryptoKeyMeta) error
	DeleteByID(keyId string) error
}

// GormCryptoKeyRepository is the implementation of the CryptoKeyRepository interface
type GormCryptoKeyRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// GormCryptoKeyRepository creates a new GormCryptoKeyRepository instance
func NewGormCryptoKeyRepository(db *gorm.DB, logger logger.Logger) (*GormCryptoKeyRepository, error) {

	return &GormCryptoKeyRepository{
		db:     db,
		logger: logger,
	}, nil
}

// Create adds a new CryptoKey to the database
func (r *GormCryptoKeyRepository) Create(key *keys.CryptoKeyMeta) error {
	// Validate the CryptoKey before saving
	if err := key.Validate(); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	if err := r.db.Create(&key).Error; err != nil {
		return fmt.Errorf("failed to create cryptographic key: %w", err)
	}

	r.logger.Info(fmt.Sprintf("Created key metadata with id %s", key.ID))
	return nil
}

func (r *GormCryptoKeyRepository) List(query *keys.CryptoKeyQuery) ([]*keys.CryptoKeyMeta, error) {
	// Validate the query parameters before using them
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	// Start building the query
	var cryptoKeyMetas []*keys.CryptoKeyMeta
	dbQuery := r.db.Model(&keys.CryptoKeyMeta{})

	// Apply filters based on the query
	if query.Algorithm != "" {
		dbQuery = dbQuery.Where("algorithm = ?", query.Algorithm)
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
	if err := dbQuery.Find(&cryptoKeyMetas).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch crypto key metadata: %w", err)
	}

	// Return the list of crypto key metadata
	return cryptoKeyMetas, nil
}

// GetByID retrieves a CryptoKey by its ID from the database
func (r *GormCryptoKeyRepository) GetByID(keyId string) (*keys.CryptoKeyMeta, error) {
	var key keys.CryptoKeyMeta
	if err := r.db.Where("id = ?", keyId).First(&key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("cryptographic key with ID %s not found", keyId)
		}
		return nil, fmt.Errorf("failed to fetch cryptographic key: %w", err)
	}
	return &key, nil
}

// UpdateByID updates an existing CryptoKey in the database
func (r *GormCryptoKeyRepository) UpdateByID(key *keys.CryptoKeyMeta) error {
	// Validate the CryptoKey before updating
	if err := key.Validate(); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	if err := r.db.Save(&key).Error; err != nil {
		return fmt.Errorf("failed to update cryptographic key: %w", err)
	}

	r.logger.Info(fmt.Sprintf("Updated key metadata with id %s", key.ID))
	return nil
}

// DeleteByID removes a CryptoKey from the database by its ID
func (r *GormCryptoKeyRepository) DeleteByID(keyId string) error {
	if err := r.db.Where("id = ?", keyId).Delete(&keys.CryptoKeyMeta{}).Error; err != nil {
		return fmt.Errorf("failed to delete cryptographic key: %w", err)
	}
	r.logger.Info(fmt.Sprintf("Deleted key metadata with id %s", keyId))
	return nil
}
