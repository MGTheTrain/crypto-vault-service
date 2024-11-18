package contracts

import (
	"crypto_vault_service/internal/domain/model"
)

// BlobMetadataManagement defines the methods for managing Blob metadata
type BlobMetadataManagement interface {
	// Create creates a new blob
	Create(blob *model.Blob) (*model.Blob, error)
	// Get retrieves blob by ID
	Get(blobID string) (*model.Blob, error)
	// Update updates a blob's metadata
	Update(blobID string, updates *model.Blob) (*model.Blob, error)
	// Delete deletes a blob by ID
	Delete(blobID string) error
}

// CryptographicKeyMetadataManagement defines the methods for managing CryptographicKey metadata
type CryptographicKeyMetadataManagement interface {
	// Create creates a new cryptographic key
	Create(key *model.CryptographicKey) (*model.CryptographicKey, error)
	// Get retrieves cryptographic key by ID
	Get(keyID string) (*model.CryptographicKey, error)
	// Update updates cryptographic key metadata
	Update(keyID string, updates *model.CryptographicKey) (*model.CryptographicKey, error)
	// Delete deletes a cryptographic key by ID
	Delete(keyID string) error
}
