package repository

import (
	"crypto_vault_service/internal/domain/keys"
	"fmt"

	"gorm.io/gorm"
)

// CryptographicKeyRepository defines the interface for CryptographicKey-related operations
type CryptographicKeyRepository interface {
	Create(key *keys.CryptographicKey) error
	GetByID(keyID string) (*keys.CryptographicKey, error)
	UpdateByID(key *keys.CryptographicKey) error
	DeleteByID(keyID string) error
}

// CryptographicKeyRepositoryImpl is the implementation of the CryptographicKeyRepository interface
type CryptographicKeyRepositoryImpl struct {
	DB *gorm.DB
}

// Create adds a new CryptographicKey to the database
func (r *CryptographicKeyRepositoryImpl) Create(key *keys.CryptographicKey) error {
	// Validate the CryptographicKey before saving
	if err := key.Validate(); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	// Save the key to the database
	if err := r.DB.Create(&key).Error; err != nil {
		return fmt.Errorf("failed to create cryptographic key: %w", err)
	}
	return nil
}

// GetByID retrieves a CryptographicKey by its ID from the database
func (r *CryptographicKeyRepositoryImpl) GetByID(keyID string) (*keys.CryptographicKey, error) {
	var key keys.CryptographicKey
	if err := r.DB.Where("id = ?", keyID).First(&key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("cryptographic key with ID %s not found", keyID)
		}
		return nil, fmt.Errorf("failed to fetch cryptographic key: %w", err)
	}
	return &key, nil
}

// UpdateByID updates an existing CryptographicKey in the database
func (r *CryptographicKeyRepositoryImpl) UpdateByID(key *keys.CryptographicKey) error {
	// Validate the CryptographicKey before updating
	if err := key.Validate(); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	// Update the key in the database
	if err := r.DB.Save(&key).Error; err != nil {
		return fmt.Errorf("failed to update cryptographic key: %w", err)
	}
	return nil
}

// DeleteByID removes a CryptographicKey from the database by its ID
func (r *CryptographicKeyRepositoryImpl) DeleteByID(keyID string) error {
	if err := r.DB.Where("id = ?", keyID).Delete(&keys.CryptographicKey{}).Error; err != nil {
		return fmt.Errorf("failed to delete cryptographic key: %w", err)
	}
	return nil
}
