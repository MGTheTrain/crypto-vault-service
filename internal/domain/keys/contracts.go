package keys

import "context"

// CryptoKeyUploadService defines methods for uploading cryptographic keys.
type CryptoKeyUploadService interface {
	// Upload uploads cryptographic keys
	// It returns a slice of CryptoKeyMeta and any error encountered during the upload process.
	Upload(ctx context.Context, userId, keyAlgorihm string, keySize uint32) ([]*CryptoKeyMeta, error)
}

// CryptoKeyMetadataService defines methods for managing cryptographic key metadata and deleting keys.
type CryptoKeyMetadataService interface {
	// List retrieves all cryptographic keys metadata considering a query filter when set.
	// It returns a slice of CryptoKeyMeta and any error encountered during the retrieval process.
	List(ctx context.Context, query *CryptoKeyQuery) ([]*CryptoKeyMeta, error)

	// GetByID retrieves the metadata of a cryptographic key by its unique ID.
	// It returns the CryptoKeyMeta and any error encountered during the retrieval process.
	GetByID(ctx context.Context, keyID string) (*CryptoKeyMeta, error)

	// DeleteByID deletes a cryptographic key and its associated metadata by ID.
	// It returns any error encountered during the deletion process.
	DeleteByID(ctx context.Context, keyID string) error
}

// CryptoKeyDownloadService defines methods for downloading cryptographic keys.
type CryptoKeyDownloadService interface {
	// Download retrieves a cryptographic key by its ID
	// It returns the CryptoKeyMeta, the key data as a byte slice, and any error encountered during the download process.
	DownloadById(ctx context.Context, keyID string) ([]byte, error)
}

// CryptoKeyRepository defines the interface for CryptoKey-related operations
type CryptoKeyRepository interface {
	Create(ctx context.Context, key *CryptoKeyMeta) error
	List(ctx context.Context, query *CryptoKeyQuery) ([]*CryptoKeyMeta, error)
	GetByID(ctx context.Context, keyID string) (*CryptoKeyMeta, error)
	UpdateByID(ctx context.Context, key *CryptoKeyMeta) error
	DeleteByID(ctx context.Context, keyID string) error
}
