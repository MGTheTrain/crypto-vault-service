package keys

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// CryptoKeyQuery represents the parameters used to query encryption keys.
type CryptoKeyQuery struct {
	Algorithm       string    `validate:"omitempty,oneof=AES RSA EC"`               // Type is optional but if provided, must be one of the listed types (AES, RSA, EC)
	Type            string    `validate:"omitempty,oneof=private public symmetric"` // Type is optional but if provided, must be one of the listed types (private-key, public-key, symmetric-key)
	DateTimeCreated time.Time `validate:"omitempty,gtefield=DateTimeCreated"`       // DateTimeCreated is optional, but can be used for filtering

	// Pagination properties
	Limit  int `validate:"omitempty,min=1"` // Limit is optional but if provided, should be at least 1
	Offset int `validate:"omitempty,min=0"` // Offset is optional but should be 0 or greater for pagination

	// Sorting properties
	SortBy    string `validate:"omitempty,oneof=ID Type DateTimeCreated"` // SortBy is optional but can be one of the fields to sort by
	SortOrder string `validate:"omitempty,oneof=asc desc"`                // SortOrder is optional, default is ascending ('asc'), can also be 'desc'
}

// New function to create a CryptoKeyQuery with default values
func NewCryptoKeyQuery() *CryptoKeyQuery {
	return &CryptoKeyQuery{
		Limit:     10,                // Default limit to 10 results per page
		Offset:    0,                 // Default offset to 0 for pagination
		SortBy:    "DateTimeCreated", // Default sort by DateTimeCreated
		SortOrder: "asc",             // Default sort order ascending
	}
}

// Validate for validating CryptoKeyQuery struct
func (k *CryptoKeyQuery) Validate() error {
	// Initialize the validator
	validate := validator.New()

	// Validate the struct
	err := validate.Struct(k)
	if err != nil {
		// If validation fails, return a formatted error
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, fmt.Sprintf("Field: %s, Tag: %s", err.Field(), err.Tag()))
		}
		return fmt.Errorf("Validation failed: %v", validationErrors)
	}
	return nil // Return nil if validation passes
}
