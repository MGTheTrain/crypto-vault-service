package keys

// ICryptoKeyUploadService defines methods for uploading cryptographic keys.
type ICryptoKeyUploadService interface {
	// Upload uploads cryptographic keys from specified file paths.
	// It returns a slice of CryptoKeyMeta and any error encountered during the upload process.
	Upload(filePath, userId, keyType, keyAlgorihm string) (*CryptoKeyMeta, error)
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
	DeleteByID(keyID string) error
}

// ICryptoKeyDownloadService defines methods for downloading cryptographic keys.
type ICryptoKeyDownloadService interface {
	// Download retrieves a cryptographic key by its ID and type.
	// It returns the CryptoKeyMeta, the key data as a byte slice, and any error encountered during the download process.
	Download(keyID, keyType string) ([]byte, error)
}
