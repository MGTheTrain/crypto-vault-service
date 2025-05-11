package connector

import (
	"context"
	"crypto_vault_service/internal/domain/keys"
)

// VaultConnector is an interface for interacting with custom key storage.
// The current implementation uses Azure Blob Storage, but this may be replaced
// with Azure Key Vault, AWS KMS, or any other cloud-based key management system in the future.
type VaultConnector interface {
	// Upload uploads bytes of a single file to Blob Storage
	// and returns the metadata for each uploaded byte stream.
	Upload(ctx context.Context, bytes []byte, userID, keyPairID, keyType, keyAlgorihm string, keySize uint32) (*keys.CryptoKeyMeta, error)

	// Download retrieves a key's content by its IDs and type and returns the data as a byte slice.
	Download(ctx context.Context, keyID, keyPairID, keyType string) ([]byte, error)

	// Delete deletes a key from Vault Storage by its IDs and type and returns any error encountered.
	Delete(ctx context.Context, keyID, keyPairID, keyType string) error
}
