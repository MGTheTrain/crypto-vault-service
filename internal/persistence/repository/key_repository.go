package repository

import (
	"fmt"

	"crypto_vault_service/internal/domain/model"

	"gorm.io/gorm"
)

// CryptographicKeyRepository defines the interface for CryptographicKey-related operations
type CryptographicKeyRepository interface {
	CreateKey(key *model.CryptographicKey) error
	GetKeyByID(keyID string) (*model.CryptographicKey, error)
	UpdateKey(key *model.CryptographicKey) error
	DeleteKey(keyID string) error
}

// CryptographicKeyRepositoryImpl is the implementation of the CryptographicKeyRepository interface
type CryptographicKeyRepositoryImpl struct {
	DB *gorm.DB
}

// CreateKey adds a new CryptographicKey to the database
func (r *CryptographicKeyRepositoryImpl) CreateKey(key *model.CryptographicKey) error {
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

// GetKeyByID retrieves a CryptographicKey by its ID from the database
func (r *CryptographicKeyRepositoryImpl) GetKeyByID(keyID string) (*model.CryptographicKey, error) {
	var key model.CryptographicKey
	if err := r.DB.Where("key_id = ?", keyID).First(&key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("cryptographic key with ID %s not found", keyID)
		}
		return nil, fmt.Errorf("failed to fetch cryptographic key: %w", err)
	}
	return &key, nil
}

// UpdateKey updates an existing CryptographicKey in the database
func (r *CryptographicKeyRepositoryImpl) UpdateKey(key *model.CryptographicKey) error {
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

// DeleteKey removes a CryptographicKey from the database by its ID
func (r *CryptographicKeyRepositoryImpl) DeleteKey(keyID string) error {
	if err := r.DB.Where("key_id = ?", keyID).Delete(&model.CryptographicKey{}).Error; err != nil {
		return fmt.Errorf("failed to delete cryptographic key: %w", err)
	}
	return nil
}
