package blobs

import "mime/multipart"

// IBlobUploadService defines methods for uploading blobs.
type IBlobUploadService interface {
	// Upload transfers blobs with the option to encrypt them using an encryption key or sign them with a signing key.
	// It returns a slice of Blob for the uploaded blobs and any error encountered during the upload process.
	Upload(form *multipart.Form, userId string, encryptionKeyId, signKeyId *string) ([]*BlobMeta, error)
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
	DeleteByID(blobId string) error
}

// IBlobDownloadService defines methods for downloading blobs.
type IBlobDownloadService interface {
	// Download retrieves a blob by its ID and name.
	// It returns the file data as a byte slice, and any error encountered during the download process.
	Download(blobId string, decryptionKeyId, verificationKeyId *string) ([]byte, error)
}
