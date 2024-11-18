package contracts

import (
	"crypto_vault_service/internal/domain/model"
)

// BlobMetadataManagement defines the methods for managing Blob metadata
type BlobMetadataManagement interface {
	// Create creates a new blob
	Create(blob *model.Blob) (*model.Blob, error)
	// GetByID retrieves blob by ID
	GetByID(blobID string) (*model.Blob, error)
	// UpdateByID updates a blob's metadata
	UpdateByID(blobID string, updates *model.Blob) (*model.Blob, error)
	// DeleteByID deletes a blob by ID
	DeleteByID(blobID string) error
}

// CryptographicKeyMetadataManagement defines the methods for managing CryptographicKey metadata
type CryptographicKeyMetadataManagement interface {
	// Create creates a new cryptographic key
	Create(key *model.CryptographicKey) (*model.CryptographicKey, error)
	// GetByID retrieves cryptographic key by ID
	GetByID(keyID string) (*model.CryptographicKey, error)
	// UpdateByID updates cryptographic key metadata
	UpdateByID(keyID string, updates *model.CryptographicKey) (*model.CryptographicKey, error)
	// DeleteByID deletes a cryptographic key by ID
	DeleteByID(keyID string) error
}
