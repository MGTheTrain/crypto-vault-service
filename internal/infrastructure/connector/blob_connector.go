package connector

import (
	"context"
	"crypto_vault_service/internal/domain/blobs"
	"mime/multipart"
)

// BlobConnector is an interface for interacting with Blob storage
type BlobConnector interface {
	// UploadFromForm uploads files to a Blob Storage
	// and returns the metadata for each uploaded byte stream.
	Upload(ctx context.Context, form *multipart.Form, userId string, encryptionKeyId, signKeyId *string) ([]*blobs.BlobMeta, error)

	// Download retrieves a blob's content by its ID and name, and returns the data as a stream.
	Download(ctx context.Context, blobID, blobName string) ([]byte, error)

	// Delete deletes a blob from Blob Storage by its ID and Name, and returns any error encountered.
	Delete(ctx context.Context, blobID, blobName string) error
}
