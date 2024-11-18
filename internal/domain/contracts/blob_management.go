package contracts

import (
	"crypto_vault_service/internal/domain/model"
)

// BlobManagement defines methods for managing blob operations.
type BlobManagement interface {
	// Upload handles the upload of blobs from file paths.
	// Returns the created Blobs metadata and any error encountered.
	Upload(filePath []string) ([]*model.Blob, error)

	// Download retrieves a blob by its ID and name, returning the metadata and file data.
	// Returns the Blob metadata, file data as a byte slice, and any error.
	Download(blobId, blobName string) (*model.Blob, []byte, error)

	// DeleteByID removes a blob by its ID.
	// Returns any error encountered.
	DeleteByID(blobId string) error
}
