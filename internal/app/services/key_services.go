package services

import (
	"context"
	"crypto/elliptic"
	"crypto/x509"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/logger"
	"fmt"

	"github.com/google/uuid"
)

// cryptoKeyUploadService implements the CryptoKeyUploadService interface for handling blob uploads
type cryptoKeyUploadService struct {
	vaultConnector connector.VaultConnector
	cryptoKeyRepo  keys.CryptoKeyRepository
	logger         logger.Logger
}

// NewCryptoKeyUploadService creates a new cryptoKeyUploadService instance
func NewCryptoKeyUploadService(vaultConnector connector.VaultConnector, cryptoKeyRepo keys.CryptoKeyRepository, logger logger.Logger) (*cryptoKeyUploadService, error) {
	return &cryptoKeyUploadService{
		vaultConnector: vaultConnector,
		cryptoKeyRepo:  cryptoKeyRepo,
		logger:         logger,
	}, nil
}

// Upload uploads cryptographic keys
// It returns a slice of CryptoKeyMeta and any error encountered during the upload process.
func (s *cryptoKeyUploadService) Upload(ctx context.Context, userId, keyAlgorithm string, keySize uint32) ([]*keys.CryptoKeyMeta, error) {
	var cryptKeyMetas []*keys.CryptoKeyMeta

	keyPairId := uuid.New().String()
	var err error
	switch keyAlgorithm {
	case "AES":
		cryptKeyMetas, err = s.uploadAESKey(ctx, userId, keyPairId, keyAlgorithm, keySize)
	case "EC":
		cryptKeyMetas, err = s.uploadECKey(ctx, userId, keyPairId, keyAlgorithm, keySize)
	case "RSA":
		cryptKeyMetas, err = s.uploadRSAKey(ctx, userId, keyPairId, keyAlgorithm, keySize)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", keyAlgorithm)
	}

	if err != nil {
		return nil, err
	}

	return cryptKeyMetas, nil
}

// Helper function for uploading AES key
func (s *cryptoKeyUploadService) uploadAESKey(ctx context.Context, userId, keyPairId, keyAlgorithm string, keySize uint32) ([]*keys.CryptoKeyMeta, error) {
	var keyMetas []*keys.CryptoKeyMeta

	aesProcessor, err := cryptography.NewAESProcessor(s.logger)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	var keySizeInBytes int
	switch keySize {
	case 128:
		keySizeInBytes = 16
	case 192:
		keySizeInBytes = 24
	case 256:
		keySizeInBytes = 32
	default:
		return nil, fmt.Errorf("key size %v not supported for AES", keySize)
	}

	symmetricKeyBytes, err := aesProcessor.GenerateKey(keySizeInBytes)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	keyType := "symmetric"
	cryptoKeyMeta, err := s.vaultConnector.Upload(ctx, symmetricKeyBytes, userId, keyPairId, keyType, keyAlgorithm, keySize)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := s.cryptoKeyRepo.Create(ctx, cryptoKeyMeta); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	keyMetas = append(keyMetas, cryptoKeyMeta)
	return keyMetas, nil
}

// Helper function for uploading EC key pair (private and public)
func (s *cryptoKeyUploadService) uploadECKey(ctx context.Context, userId, keyPairId, keyAlgorithm string, keySize uint32) ([]*keys.CryptoKeyMeta, error) {
	var keyMetas []*keys.CryptoKeyMeta

	var curve elliptic.Curve
	switch keySize {
	case 224:
		curve = elliptic.P224()
	case 256:
		curve = elliptic.P256()
	case 384:
		curve = elliptic.P384()
	case 521:
		curve = elliptic.P521()
	default:
		return nil, fmt.Errorf("key size %v not supported for EC", keySize)
	}

	ecProcessor, err := cryptography.NewECProcessor(s.logger)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	privateKey, publicKey, err := ecProcessor.GenerateKeys(curve)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Upload Private Key
	privateKeyBytes := append(privateKey.D.Bytes(), privateKey.X.Bytes()...)
	privateKeyBytes = append(privateKeyBytes, privateKey.Y.Bytes()...)
	keyType := "private"
	cryptoKeyMeta, err := s.vaultConnector.Upload(ctx, privateKeyBytes, userId, keyPairId, keyType, keyAlgorithm, keySize)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := s.cryptoKeyRepo.Create(ctx, cryptoKeyMeta); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	keyMetas = append(keyMetas, cryptoKeyMeta)

	// Upload Public Key
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	keyType = "public"
	cryptoKeyMeta, err = s.vaultConnector.Upload(ctx, publicKeyBytes, userId, keyPairId, keyType, keyAlgorithm, keySize)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := s.cryptoKeyRepo.Create(ctx, cryptoKeyMeta); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	keyMetas = append(keyMetas, cryptoKeyMeta)
	return keyMetas, nil
}

// Helper function for uploading RSA key pair (private and public)
func (s *cryptoKeyUploadService) uploadRSAKey(ctx context.Context, userId, keyPairId, keyAlgorithm string, keySize uint32) ([]*keys.CryptoKeyMeta, error) {
	var keyMetas []*keys.CryptoKeyMeta

	rsaProcessor, err := cryptography.NewRSAProcessor(s.logger)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	privateKey, publicKey, err := rsaProcessor.GenerateKeys(int(keySize))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Upload Private Key
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	keyType := "private"
	cryptoKeyMeta, err := s.vaultConnector.Upload(ctx, privateKeyBytes, userId, keyPairId, keyType, keyAlgorithm, keySize)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := s.cryptoKeyRepo.Create(ctx, cryptoKeyMeta); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	keyMetas = append(keyMetas, cryptoKeyMeta)

	// Upload Public Key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %w", err)
	}
	keyType = "public"
	cryptoKeyMeta, err = s.vaultConnector.Upload(ctx, publicKeyBytes, userId, keyPairId, keyType, keyAlgorithm, keySize)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err := s.cryptoKeyRepo.Create(ctx, cryptoKeyMeta); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	keyMetas = append(keyMetas, cryptoKeyMeta)
	return keyMetas, nil
}

// cryptoKeyMetadataService implements the CryptoKeyMetadataService interface to manages cryptographic key metadata.
type cryptoKeyMetadataService struct {
	vaultConnector connector.VaultConnector
	cryptoKeyRepo  keys.CryptoKeyRepository
	logger         logger.Logger
}

// NewCryptoKeyMetadataService creates a new cryptoKeyMetadataService instance
func NewCryptoKeyMetadataService(vaultConnector connector.VaultConnector, cryptoKeyRepo keys.CryptoKeyRepository, logger logger.Logger) (*cryptoKeyMetadataService, error) {
	return &cryptoKeyMetadataService{
		vaultConnector: vaultConnector,
		cryptoKeyRepo:  cryptoKeyRepo,
		logger:         logger,
	}, nil
}

// List retrieves all cryptographic key metadata based on a query.
func (s *cryptoKeyMetadataService) List(ctx context.Context, query *keys.CryptoKeyQuery) ([]*keys.CryptoKeyMeta, error) {
	crypoKeyMetas, err := s.cryptoKeyRepo.List(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return crypoKeyMetas, nil
}

// GetByID retrieves the metadata of a cryptographic key by its ID.
func (s *cryptoKeyMetadataService) GetByID(ctx context.Context, keyID string) (*keys.CryptoKeyMeta, error) {
	keyMeta, err := s.cryptoKeyRepo.GetByID(ctx, keyID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return keyMeta, nil
}

// DeleteByID deletes a cryptographic key's metadata by its ID.
func (s *cryptoKeyMetadataService) DeleteByID(ctx context.Context, keyID string) error {
	keyMeta, err := s.GetByID(ctx, keyID)
	if err != nil {
		return fmt.Errorf("failed to%w", err)
	}

	err = s.vaultConnector.Delete(ctx, keyID, keyMeta.KeyPairID, keyMeta.Type)
	if err != nil {
		return fmt.Errorf("failed to%w", err)
	}

	err = s.cryptoKeyRepo.DeleteByID(ctx, keyID)
	if err != nil {
		return fmt.Errorf("failed to%w", err)
	}
	return nil
}

// cryptoKeyDownloadService implements the CryptoKeyDownloadService interface to handle the download of cryptographic keys.
type cryptoKeyDownloadService struct {
	vaultConnector connector.VaultConnector
	cryptoKeyRepo  keys.CryptoKeyRepository
	logger         logger.Logger
}

// NewCryptoKeyDownloadService creates a new cryptoKeyDownloadService instance
func NewCryptoKeyDownloadService(vaultConnector connector.VaultConnector, cryptoKeyRepo keys.CryptoKeyRepository, logger logger.Logger) (*cryptoKeyDownloadService, error) {
	return &cryptoKeyDownloadService{
		vaultConnector: vaultConnector,
		cryptoKeyRepo:  cryptoKeyRepo,
		logger:         logger,
	}, nil
}

// Download retrieves a cryptographic key by its ID.
func (s *cryptoKeyDownloadService) DownloadByID(ctx context.Context, keyID string) ([]byte, error) {
	keyMeta, err := s.cryptoKeyRepo.GetByID(ctx, keyID)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	blobData, err := s.vaultConnector.Download(ctx, keyMeta.ID, keyMeta.KeyPairID, keyMeta.Type)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return blobData, nil
}
