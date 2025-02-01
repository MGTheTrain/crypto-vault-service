package keys

import "context"

// ICryptoKeyUploadService defines methods for uploading cryptographic keys.
type ICryptoKeyUploadService interface {
	// Upload uploads cryptographic keys
	// It returns a slice of CryptoKeyMeta and any error encountered during the upload process.
	Upload(ctx context.Context, userId, keyPairId, keyAlgorihm string, keySize uint) ([]*CryptoKeyMeta, error)
}

// ICryptoKeyMetadataService defines methods for managing cryptographic key metadata and deleting keys.
type ICryptoKeyMetadataService interface {
	// List retrieves all cryptographic keys metadata considering a query filter when set.
	// It returns a slice of CryptoKeyMeta and any error encountered during the retrieval process.
	List(query *CryptoKeyQuery) ([]*CryptoKeyMeta, error)

	// GetByID retrieves the metadata of a cryptographic key by its unique ID.
	// It returns the CryptoKeyMeta and any error encountered during the retrieval process.
	GetByID(keyID string) (*CryptoKeyMeta, error)

	// DeleteByID deletes a cryptographic key and its associated metadata by ID.
	// It returns any error encountered during the deletion process.
	DeleteByID(ctx context.Context, keyID string) error
}

// ICryptoKeyDownloadService defines methods for downloading cryptographic keys.
type ICryptoKeyDownloadService interface {
	// Download retrieves a cryptographic key by its ID
	// It returns the CryptoKeyMeta, the key data as a byte slice, and any error encountered during the download process.
	Download(ctx context.Context, keyID string) ([]byte, error)
}

// CryptoKeyRepository defines the interface for CryptoKey-related operations
type CryptoKeyRepository interface {
	Create(key *CryptoKeyMeta) error
	List(query *CryptoKeyQuery) ([]*CryptoKeyMeta, error)
	GetByID(keyId string) (*CryptoKeyMeta, error)
	UpdateByID(key *CryptoKeyMeta) error
	DeleteByID(keyId string) error
}
