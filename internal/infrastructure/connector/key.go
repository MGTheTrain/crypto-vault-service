package connector

import (
	"crypto_vault_service/internal/domain/keys"
)

// VaultConnector is an interface for interacting with key storages
type VaultConnector interface {
	// Upload uploads multiple files to Vault Storage and returns their metadata.
	Upload(filePaths []string) ([]*keys.CryptoKeyMeta, error)

	// Download retrieves a blob's content by its ID and name, and returns the data as a stream.
	Download(blobId, blobName string) ([]byte, error)

	//
	// Rotate()

	// Delete deletes a blob from Vault Storage by its ID and Name, and returns any error encountered.
	Delete(blobId, blobName string) error
}
