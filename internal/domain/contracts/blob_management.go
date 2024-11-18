package contracts

import (
	"crypto_vault_service/internal/domain/model"
	"mime/multipart"
)

// BlobManagement defines methods for managing blob operations.
type BlobManagement interface {
	// Upload handles the upload of a blob from a multipart form.
	// Returns the created Blob metadata and any error encountered.
	Upload(form *multipart.Form) (*model.Blob, error)

	// DownloadByID retrieves a blob by its ID, returning the metadata and file data.
	// Returns the Blob metadata, file data as a byte slice, and any error.
	DownloadByID(blobId string) (*model.Blob, []byte, error)

	// DeleteByID removes a blob by its ID.
	// Returns any error encountered.
	DeleteByID(blobId string) error
}
