package blobs

import (
	"context"
	"mime/multipart"
)

// IBlobUploadService defines methods for uploading blobs.
type IBlobUploadService interface {
	// Upload transfers blobs with the option to encrypt them using an encryption key or sign them with a signing key.
	// It returns a slice of Blob for the uploaded blobs and any error encountered during the upload process.
	Upload(ctx context.Context, form *multipart.Form, userId string, encryptionKeyId, signKeyId *string) ([]*BlobMeta, error)
}

// IBlobMetadataService defines methods for retrieving Blob and deleting a blob along with its metadata.
type IBlobMetadataService interface {
	// List retrieves all blobs' metadata considering a query filter when set.
	// It returns a slice of Blob and any error encountered during the retrieval.
	List(query *BlobMetaQuery) ([]*BlobMeta, error)

	// GetByID retrieves the metadata of a blob by its unique ID.
	// It returns the Blob and any error encountered during the retrieval process.
	GetByID(blobId string) (*BlobMeta, error)

	// DeleteByID deletes a blob and its associated metadata by ID.
	// It returns any error encountered during the deletion process.
	DeleteByID(ctx context.Context, blobId string) error
}

// IBlobDownloadService defines methods for downloading blobs.
type IBlobDownloadService interface {
	// The download function retrieves a blob's content using its ID and also enables data decryption.
	Download(ctx context.Context, blobId string, decryptionKeyId *string) ([]byte, error)
}

// BlobRepository defines the interface for Blob-related operations
type BlobRepository interface {
	Create(ctx context.Context, blob *BlobMeta) error
	List(ctx context.Context, query *BlobMetaQuery) ([]*BlobMeta, error)
	GetById(ctx context.Context, blobId string) (*BlobMeta, error)
	UpdateById(ctx context.Context, blob *BlobMeta) error
	DeleteById(ctx context.Context, blobId string) error
}
