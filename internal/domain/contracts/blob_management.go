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

	// Download retrieves a blob by its ID, returning the metadata and file data.
	// Returns the Blob metadata, file data as a byte slice, and any error.
	Download(blobId string) (*model.Blob, []byte, error)

	// Delete removes a blob by its ID.
	// Returns the deleted Blob metadata and any error encountered.
	Delete(blobId string) (*model.Blob, error)
}
