package services

import (
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/persistence/repository"
	"fmt"
)

type CryptoKeyUploadService struct {
	VaultConnector connector.VaultConnector
	CryptoKeyRepo  repository.CryptoKeyRepository
}

func (s *CryptoKeyUploadService) Upload(filePaths []string, userId string) ([]*keys.CryptoKeyMeta, error) {
	// Step 1: Upload files to blob storage
	keyMetas, err := s.VaultConnector.Upload(filePaths, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to upload files: %w", err)
	}

	// Step 2: Store the metadata in the database
	for _, keyMeta := range keyMetas {
		if err := s.CryptoKeyRepo.Create(keyMeta); err != nil {
			return nil, fmt.Errorf("failed to create metadata for key of type %s: %w", keyMeta.Type, err)
		}
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
	var keyMetas []*keys.CryptoKeyMeta
	// TBD

	return keyMetas, nil
}

// GetByID retrieves the metadata of a cryptographic key by its ID.
func (s *CryptoKeyMetadataService) GetByID(keyID string) (*keys.CryptoKeyMeta, error) {
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
	blobName := ""

	if keyType == keys.AsymmetricPublic {
		blobName = "public"
	} else if keyType == keys.AsymmetricPrivate {
		blobName = "private"
	} else if keyType == keys.Symmetric {
		blobName = "symmetric"
	} else {
		return nil, fmt.Errorf("unsupported key type: %v", keyType)
	}

	blobData, err := s.VaultConnector.Download(keyID, blobName)
	if err != nil {
		return nil, fmt.Errorf("failed to download key from blob storage: %w", err)
	}

	return blobData, nil
}
