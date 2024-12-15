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
	GetByID(keyId string) (*keys.CryptoKeyMeta, error)
	UpdateByID(key *keys.CryptoKeyMeta) error
	DeleteByID(keyId string) error
}

// GormCryptoKeyRepository is the implementation of the CryptoKeyRepository interface
type GormCryptoKeyRepository struct {
	DB     *gorm.DB
	Logger logger.Logger
}

// GormCryptoKeyRepository creates a new GormCryptoKeyRepository instance
func NewGormCryptoKeyRepository(db *gorm.DB, logger logger.Logger) (*GormCryptoKeyRepository, error) {

	return &GormCryptoKeyRepository{
		DB:     db,
		Logger: logger,
	}, nil
}

// Create adds a new CryptoKey to the database
func (r *GormCryptoKeyRepository) Create(key *keys.CryptoKeyMeta) error {
	// Validate the CryptoKey before saving
	if err := key.Validate(); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	if err := r.DB.Create(&key).Error; err != nil {
		return fmt.Errorf("failed to create cryptographic key: %w", err)
	}

	r.Logger.Info(fmt.Sprintf("Created key metadata with id %s", key.ID))
	return nil
}

// GetByID retrieves a CryptoKey by its ID from the database
func (r *GormCryptoKeyRepository) GetByID(keyId string) (*keys.CryptoKeyMeta, error) {
	var key keys.CryptoKeyMeta
	if err := r.DB.Where("id = ?", keyId).First(&key).Error; err != nil {
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

	if err := r.DB.Save(&key).Error; err != nil {
		return fmt.Errorf("failed to update cryptographic key: %w", err)
	}

	r.Logger.Info(fmt.Sprintf("Updated key metadata with id %s", key.ID))
	return nil
}

// DeleteByID removes a CryptoKey from the database by its ID
func (r *GormCryptoKeyRepository) DeleteByID(keyId string) error {
	if err := r.DB.Where("id = ?", keyId).Delete(&keys.CryptoKeyMeta{}).Error; err != nil {
		return fmt.Errorf("failed to delete cryptographic key: %w", err)
	}
	r.Logger.Info(fmt.Sprintf("Deleted key metadata with id %s", keyId))
	return nil
}
