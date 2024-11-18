package contracts

import (
	"crypto_vault_service/internal/domain/model"
)

// Define KeyType as a custom type (based on int)
type KeyType int

// Enum-like values using iota
const (
	AsymmetricPublic KeyType = iota
	AsymmetricPrivate
	Symmetric
)

// KeyManagement defines methods for managing cryptographic key operations.
type KeyManagement interface {
	// Upload handles the upload of blobs from file paths.
	// Returns the created Blobs metadata and any error encountered.
	Upload(filePath []string) ([]*model.CryptographicKey, error)

	// Download retrieves a cryptographic key by its ID and key type, returning the metadata and key data.
	// Returns the key metadata, key data as a byte slice, and any error.
	Download(keyId string, keyType KeyType) (*model.CryptographicKey, []byte, error)

	// DeleteByID removes a cryptographic key by its ID.
	// Returns any error encountered.
	DeleteByID(keyId string) error
}
