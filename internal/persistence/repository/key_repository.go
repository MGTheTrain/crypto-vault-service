package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/logger"
)

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
func (r *GormCryptoKeyRepository) Create(ctx context.Context, key *keys.CryptoKeyMeta) error {
	// Validate the CryptoKey before saving
	if err := key.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(&key).Error; err != nil {
		return fmt.Errorf("failed to create cryptographic key: %w", err)
	}

	r.logger.Info(fmt.Sprintf("Created key metadata with id %s", key.ID))
	return nil
}

func (r *GormCryptoKeyRepository) List(ctx context.Context, query *keys.CryptoKeyQuery) ([]*keys.CryptoKeyMeta, error) {
	// Validate the query parameters before using them
	if err := query.Validate(); err != nil {
		return nil, fmt.Errorf("invalid query parameters: %w", err)
	}

	// Start building the query
	var cryptoKeyMetas []*keys.CryptoKeyMeta
	dbQuery := r.db.WithContext(ctx).Model(&keys.CryptoKeyMeta{})

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
func (r *GormCryptoKeyRepository) GetByID(ctx context.Context, keyId string) (*keys.CryptoKeyMeta, error) {
	var key keys.CryptoKeyMeta
	if err := r.db.WithContext(ctx).Where("id = ?", keyId).First(&key).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("cryptographic key with ID %s not found", keyId)
		}

		return nil, fmt.Errorf("failed to fetch cryptographic key: %w", err)
	}
	return &key, nil
}

// UpdateByID updates an existing CryptoKey in the database
func (r *GormCryptoKeyRepository) UpdateByID(ctx context.Context, key *keys.CryptoKeyMeta) error {
	// Validate the CryptoKey before updating
	if err := key.Validate(); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if err := r.db.WithContext(ctx).Save(&key).Error; err != nil {
		return fmt.Errorf("failed to update cryptographic key: %w", err)
	}

	r.logger.Info(fmt.Sprintf("Updated key metadata with id %s", key.ID))
	return nil
}

// DeleteByID removes a CryptoKey from the database by its ID
func (r *GormCryptoKeyRepository) DeleteByID(ctx context.Context, keyId string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", keyId).Delete(&keys.CryptoKeyMeta{}).Error; err != nil {
		return fmt.Errorf("failed to delete cryptographic key: %w", err)
	}
	r.logger.Info(fmt.Sprintf("Deleted key metadata with id %s", keyId))
	return nil
}
