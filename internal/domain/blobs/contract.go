package blobs

// BlobUploadService defines methods for uploading blobs.
type BlobUploadService interface {
	// Upload handles the upload of blobs from the specified file paths.
	// It returns a slice of Blob metadata for the uploaded blobs and any error encountered during the upload process.
	Upload(filePaths []string) ([]*Blob, error)
}

// BlobMetadataService defines methods for retrieving Blob metadata and deleting a blob along with its metadata.
type BlobMetadataService interface {
	// List retrieves all blobs' metadata.
	// It returns a slice of Blob metadata and any error encountered during the retrieval.
	List() ([]*Blob, error)

	// GetByID retrieves the metadata of a blob by its unique ID.
	// It returns the Blob metadata and any error encountered during the retrieval process.
	GetByID(blobID string) (*Blob, error)

	// DeleteByID deletes a blob and its associated metadata by ID.
	// It returns any error encountered during the deletion process.
	DeleteByID(blobID string) error
}

// BlobDownloadService defines methods for downloading blobs.
type BlobDownloadService interface {
	// Download retrieves a blob by its ID and name.
	// It returns the Blob metadata, the file data as a byte slice, and any error encountered during the download process.
	Download(blobID, blobName string) (*Blob, []byte, error)
}
