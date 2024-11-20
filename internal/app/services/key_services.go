package services

import (
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/persistence/repository"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CryptoKeyUploadService struct {
	VaultConnector connector.VaultConnector
	CryptoKeyRepo  repository.CryptoKeyRepository
}

func (s *CryptoKeyUploadService) Upload(filePaths []string) ([]*keys.CryptoKeyMeta, error) {
	// Step 1: Upload files to blob storage
	userId := uuid.New().String()
	blobMeta, err := s.VaultConnector.Upload(filePaths, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to upload files: %w", err)
	}

	// Step 2: Store the metadata in the database
	var keyMetas []*keys.CryptoKeyMeta
	for _, blob := range blobMeta {
		// Map Blob metadata to CryptoKey metadata
		keyMeta := &keys.CryptoKeyMeta{
			ID:              uuid.New().String(), // Generate valid UUID for ID
			Type:            "RSA",               // Example key type
			DateTimeCreated: time.Now(),          // Valid DateTimeCreated time
			UserID:          uuid.New().String(), // Generate valid UUID for UserID
		}

		// Validate CryptoKeyMeta
		if err := keyMeta.Validate(); err != nil {
			return nil, fmt.Errorf("invalid key metadata: %w", err)
		}

		// Save metadata to DB
		if err := s.CryptoKeyRepo.Create(blob); err != nil {
			return nil, fmt.Errorf("failed to create metadata for key of type %s: %w", keyMeta.Type, err)
		}

		// Append to list
		keyMetas = append(keyMetas, keyMeta)
	}

	// Return metadata
	return keyMetas, nil
}

// CryptoKeyMetadataService manages cryptographic key metadata.
type CryptoKeyMetadataService struct {
	CryptoKeyRepo repository.CryptoKeyRepository
}

// List retrieves all cryptographic key metadata based on a query.
func (s *CryptoKeyMetadataService) List(query *keys.CryptoKeyQuery) ([]*keys.CryptoKeyMeta, error) {
	// For now, let's just retrieve all metadata from the database
	var keyMetas []*keys.CryptoKeyMeta
	// TBD

	return keyMetas, nil
}

// GetByID retrieves the metadata of a cryptographic key by its ID.
func (s *CryptoKeyMetadataService) GetByID(keyID string) (*keys.CryptoKeyMeta, error) {
	// Retrieve the metadata from the database
	keyMeta, err := s.CryptoKeyRepo.GetByID(keyID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve key metadata: %w", err)
	}

	return keyMeta, nil
}

// DeleteByID deletes a cryptographic key's metadata by its ID.
func (s *CryptoKeyMetadataService) DeleteByID(keyID string) error {
	// Delete the metadata from the database
	err := s.CryptoKeyRepo.DeleteByID(keyID)
	if err != nil {
		return fmt.Errorf("failed to delete key metadata: %w", err)
	}
	return nil
}

// CryptoKeyDownloadService handles the download of cryptographic keys.
type CryptoKeyDownloadService struct {
	VaultConnector connector.VaultConnector
}

// Download retrieves a cryptographic key by its ID and type.
func (s *CryptoKeyDownloadService) Download(keyID string, keyType keys.KeyType) ([]byte, error) {
	blobName := "" // Declare the variable outside the blocks

	if keyType == keys.AsymmetricPublic {
		blobName = "asymmetric-public-key" // Assign to the already declared variable
	} else if keyType == keys.AsymmetricPrivate {
		blobName = "asymmetric-private-key" // Assign to the already declared variable
	} else if keyType == keys.Symmetric {
		blobName = "symmetric-key" // Assign to the already declared variable
	} else {
		return nil, fmt.Errorf("unsupported key type: %v", keyType)
	}

	blobData, err := s.VaultConnector.Download(keyID, blobName)
	if err != nil {
		return nil, fmt.Errorf("failed to download key from blob storage: %w", err)
	}

	// Return the metadata and the downloaded content (as a byte slice)
	return blobData, nil
}
