package contracts

import (
	"crypto_vault_service/internal/domain/model"
	"mime/multipart"
)

// KeyManagement defines methods for managing cryptographic key operations.
type KeyManagement interface {
	// Upload handles the upload of a cryptographic key from a multipart form.
	// Returns the created key metadata and any error encountered.
	Upload(form *multipart.Form) (*model.CryptographicKey, error)

	// DownloadByID retrieves a cryptographic key by its ID, returning the metadata and key data.
	// Returns the key metadata, key data as a byte slice, and any error.
	DownloadByID(keyId string) (*model.CryptographicKey, []byte, error)

	// DeleteByID removes a cryptographic key by its ID.
	// Returns any error encountered.
	DeleteByID(keyId string) error
}
