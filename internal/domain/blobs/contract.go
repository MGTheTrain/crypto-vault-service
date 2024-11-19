package blobs

// BlobManagement defines methods for managing blob operations.
type BlobManagement interface {
	// Upload handles the upload of blobs from file paths.
	// Returns the created Blobs metadata and any error encountered.
	Upload(filePath []string) ([]*Blob, error)

	// Download retrieves a blob by its ID and name, returning the metadata and file data.
	// Returns the Blob metadata, file data as a byte slice, and any error.
	Download(blobId, blobName string) (*Blob, []byte, error)

	// DeleteByID removes a blob by its ID.
	// Returns any error encountered.
	DeleteByID(blobId string) error
}

// BlobMetadataManagement defines the methods for managing Blob metadata
type BlobMetadataManagement interface {
	// Create creates a new blob
	Create(blob *Blob) (*Blob, error)
	// GetByID retrieves blob by ID
	GetByID(blobID string) (*Blob, error)
	// UpdateByID updates a blob's metadata
	UpdateByID(blobID string, updates *Blob) (*Blob, error)
	// DeleteByID deletes a blob by ID
	DeleteByID(blobID string) error
}
