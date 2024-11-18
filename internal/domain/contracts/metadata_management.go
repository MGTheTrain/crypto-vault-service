package contracts

import (
	"crypto_vault_service/internal/domain/model"
)

// MetadataManagement defines the methods for managing Blob and CryptographicKey metadata
type MetadataManagement interface {
	// ---CRUD operations for Blob metadata---

	// CreateBlob creates a new blob
	CreateBlob(blob *model.Blob) (*model.Blob, error)
	// GetBlob retrieves blob by ID
	GetBlob(blobID string) (*model.Blob, error)
	// UpdateBlob updates a blob's metadata
	UpdateBlob(blobID string, updates *model.Blob) (*model.Blob, error)
	// DeleteBlob deletes a blob by ID
	DeleteBlob(blobID string) error

	// ---CRUD operations for CryptographicKey metadata---

	// CreateCryptographicKey creates a new cryptographic key
	CreateCryptographicKey(key *model.CryptographicKey) (*model.CryptographicKey, error)
	// GetCryptographicKey retrieves cryptographic key by ID
	GetCryptographicKey(keyID string) (*model.CryptographicKey, error)
	// UpdateCryptographicKey updates cryptographic key metadata
	UpdateCryptographicKey(keyID string, updates *model.CryptographicKey) (*model.CryptographicKey, error)
	// DeleteCryptographicKey deletes a cryptographic key by ID
	DeleteCryptographicKey(keyID string) error
}
