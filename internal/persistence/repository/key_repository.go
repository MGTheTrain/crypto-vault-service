package repository

import (
	"crypto_vault_service/internal/domain/keys"
	"fmt"

	"gorm.io/gorm"
)

// CryptoKeyRepository defines the interface for CryptoKey-related operations
type CryptoKeyRepository interface {
	Create(key *keys.CryptoKeyMeta) error
	GetByID(keyID string) (*keys.CryptoKeyMeta, error)
	UpdateByID(key *keys.CryptoKeyMeta) error
	DeleteByID(keyID string) error
}

// GormCryptoKeyRepository is the implementation of the CryptoKeyRepository interface
type GormCryptoKeyRepository struct {
	DB *gorm.DB
}

// Create adds a new CryptoKey to the database
func (r *GormCryptoKeyRepository) Create(key *keys.CryptoKeyMeta) error {
	// Validate the CryptoKey before saving
	if err := key.Validate(); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	// Save the key to the database
	if err := r.DB.Create(&key).Error; err != nil {
		return fmt.Errorf("failed to create cryptographic key: %w", err)
	}
	return nil
}

// GetByID retrieves a CryptoKey by its ID from the database
func (r *GormCryptoKeyRepository) GetByID(keyID string) (*keys.CryptoKeyMeta, error) {
	var key keys.CryptoKeyMeta
	if err := r.DB.Where("id = ?", keyID).First(&key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("cryptographic key with ID %s not found", keyID)
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

	// Update the key in the database
	if err := r.DB.Save(&key).Error; err != nil {
		return fmt.Errorf("failed to update cryptographic key: %w", err)
	}
	return nil
}

// DeleteByID removes a CryptoKey from the database by its ID
func (r *GormCryptoKeyRepository) DeleteByID(keyID string) error {
	if err := r.DB.Where("id = ?", keyID).Delete(&keys.CryptoKeyMeta{}).Error; err != nil {
		return fmt.Errorf("failed to delete cryptographic key: %w", err)
	}
	return nil
}
