package services

import (
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/persistence/repository"
	"fmt"
	"mime/multipart"
)

type CryptoKeyUploadService struct {
	VaultConnector connector.VaultConnector
	CryptoKeyRepo  repository.CryptoKeyRepository
}

// Upload uploads cryptographic keys
// It returns a slice of CryptoKeyMeta and any error encountered during the upload process.
func (s *CryptoKeyUploadService) Upload(form *multipart.Form, userId, keyType, keyAlgorihm string) (*keys.CryptoKeyMeta, error) {

	keyMeta, err := s.VaultConnector.UploadFromForm(form, userId, keyType, keyAlgorihm)
	if err != nil {
		return nil, fmt.Errorf("failed to upload files: %w", err)
	}

	if err := s.CryptoKeyRepo.Create(keyMeta); err != nil {
		return nil, fmt.Errorf("failed to create metadata for key of type %s: %w", keyMeta.Type, err)
	}

	return keyMeta, nil
}

// CryptoKeyMetadataService manages cryptographic key metadata.
type CryptoKeyMetadataService struct {
	VaultConnector connector.VaultConnector
	CryptoKeyRepo  repository.CryptoKeyRepository
}

// List retrieves all cryptographic key metadata based on a query.
func (s *CryptoKeyMetadataService) List(query *keys.CryptoKeyQuery) ([]*keys.CryptoKeyMeta, error) {
	var keyMetas []*keys.CryptoKeyMeta
	// TBD

	return keyMetas, nil
}

// GetByID retrieves the metadata of a cryptographic key by its ID.
func (s *CryptoKeyMetadataService) GetByID(keyId string) (*keys.CryptoKeyMeta, error) {
	keyMeta, err := s.CryptoKeyRepo.GetByID(keyId)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve key metadata: %w", err)
	}

	return keyMeta, nil
}

// DeleteByID deletes a cryptographic key's metadata by its ID.
func (s *CryptoKeyMetadataService) DeleteByID(keyId string) error {
	keyMeta, err := s.GetByID(keyId)
	if err != nil {
		return fmt.Errorf("failed to retrieve key metadata: %w", err)
	}

	err = s.VaultConnector.Delete(keyId, keyMeta.Type)
	if err != nil {
		return fmt.Errorf("failed to delete key blob: %w", err)
	}

	err = s.CryptoKeyRepo.DeleteByID(keyId)
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
func (s *CryptoKeyDownloadService) Download(keyId, keyType string) ([]byte, error) {
	blobData, err := s.VaultConnector.Download(keyId, keyType)
	if err != nil {
		return nil, fmt.Errorf("failed to download key from blob storage: %w", err)
	}

	return blobData, nil
}
