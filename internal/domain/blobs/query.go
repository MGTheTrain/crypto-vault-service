package blobs

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// BlobMetaQuery represents metadata on the actual blob being stored.
type BlobMetaQuery struct {
	UploadTime time.Time `validate:"required"`               // UploadTime is required
	Name       string    `validate:"required,min=1,max=255"` // Name is required and its length must be between 1 and 255 characters
	Size       int64     `validate:"required,min=1"`         // Size must be greater than 0
	Type       string    `validate:"required,min=1,max=50"`  // Type is required, and its length must be between 1 and 50 characters

	// Pagination properties
	Limit  int `validate:"omitempty,min=1"` // Limit is optional but if provided, should be at least 1
	Offset int `validate:"omitempty,min=0"` // Offset is optional but should be 0 or greater for pagination

	// Sorting properties
	SortBy    string `validate:"omitempty,oneof=ID Type CreatedAt ExpiresAt"` // SortBy is optional but can be one of the fields to sort by
	SortOrder string `validate:"omitempty,oneof=asc desc"`                    // SortOrder is optional, default is ascending ('asc'), can also be 'desc'
}

// NewBlobMetaQuery creates a BlobMetaQuery with default values.
func NewBlobMetaQuery() *BlobMetaQuery {
	return &BlobMetaQuery{
		Limit:     10,          // Default limit to 10 results per page
		Offset:    0,           // Default offset to 0 for pagination
		SortBy:    "CreatedAt", // Default sort by CreatedAt
		SortOrder: "asc",       // Default sort order ascending
	}
}

// Validate validates the BlobMetaQuery struct based on the defined rules.
func (b *BlobMetaQuery) Validate() error {
	// Initialize the validator
	validate := validator.New()

	// Validate the struct fields
	err := validate.Struct(b)
	if err != nil {
		// Collect all validation errors
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field: %s, Tag: %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("Validation failed: %v", validationErrors)
	}

	// Custom validation logic: Check that ExpiresAt (if exists) is after CreatedAt
	if !b.UploadTime.IsZero() && b.UploadTime.After(time.Now()) {
		return fmt.Errorf("UploadTime cannot be in the future")
	}

	// Return nil if no validation errors are found
	return nil
}
