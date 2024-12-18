package services

import (
	"crypto/elliptic"
	"crypto/x509"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/persistence/repository"
	"fmt"
)

type CryptoKeyUploadService struct {
	VaultConnector connector.VaultConnector
	CryptoKeyRepo  repository.CryptoKeyRepository
	Logger         logger.Logger
}

// NewCryptoKeyUploadService creates a new CryptoKeyUploadService instance
func NewCryptoKeyUploadService(vaultConnector connector.VaultConnector, cryptoKeyRepo repository.CryptoKeyRepository, logger logger.Logger) (*CryptoKeyUploadService, error) {
	return &CryptoKeyUploadService{
		VaultConnector: vaultConnector,
		CryptoKeyRepo:  cryptoKeyRepo,
		Logger:         logger,
	}, nil
}

// Upload uploads cryptographic keys
// It returns a slice of CryptoKeyMeta and any error encountered during the upload process.
func (s *CryptoKeyUploadService) Upload(userId, keyAlgorihm string, keySize uint) ([]*keys.CryptoKeyMeta, error) {
	var cryptKeyMetas []*keys.CryptoKeyMeta

	if keyAlgorihm == "AES" {
		aes, err := cryptography.NewAES(s.Logger)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		symmetricKeyBytes, err := aes.GenerateKey(int(keySize))
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		keyType := "symmetric"
		cryptoKeyMeta, err := s.VaultConnector.Upload(symmetricKeyBytes, userId, keyType, keyAlgorihm, keySize)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		if err := s.CryptoKeyRepo.Create(cryptoKeyMeta); err != nil {
			return nil, fmt.Errorf("failed to create metadata for key of type %s: %w", cryptoKeyMeta.Type, err)
		}

		cryptKeyMetas = append(cryptKeyMetas, cryptoKeyMeta)
	} else if keyAlgorihm == "EC" {
		var curve elliptic.Curve
		if keySize == 224 {
			curve = elliptic.P224()
		} else if keySize == 256 {
			curve = elliptic.P256()
		} else if keySize == 384 {
			curve = elliptic.P384()
		} else if keySize == 521 {
			curve = elliptic.P521()
		} else {
			return nil, fmt.Errorf("key size %v not supported", keySize)
		}
		ec, err := cryptography.NewEC(s.Logger)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		privateKey, publicKey, err := ec.GenerateKeys(curve)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		// PRIV
		privateKeyBytes := append(privateKey.D.Bytes(), privateKey.PublicKey.X.Bytes()...)
		privateKeyBytes = append(privateKeyBytes, privateKey.PublicKey.Y.Bytes()...)
		keyType := "private"

		cryptoKeyMeta, err := s.VaultConnector.Upload(privateKeyBytes, userId, keyType, keyAlgorihm, keySize)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		if err := s.CryptoKeyRepo.Create(cryptoKeyMeta); err != nil {
			return nil, fmt.Errorf("failed to create metadata for key of type %s: %w", cryptoKeyMeta.Type, err)
		}

		cryptKeyMetas = append(cryptKeyMetas, cryptoKeyMeta)

		// PUB
		publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
		keyType = "public"

		cryptoKeyMeta, err = s.VaultConnector.Upload(publicKeyBytes, userId, keyType, keyAlgorihm, keySize)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		if err := s.CryptoKeyRepo.Create(cryptoKeyMeta); err != nil {
			return nil, fmt.Errorf("failed to create metadata for key of type %s: %w", cryptoKeyMeta.Type, err)
		}

		cryptKeyMetas = append(cryptKeyMetas, cryptoKeyMeta)
	} else if keyAlgorihm == "RSA" {
		rsa, err := cryptography.NewRSA(s.Logger)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		privateKey, publicKey, err := rsa.GenerateKeys(int(keySize))
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		// PRIV
		privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		keyType := "private"

		cryptoKeyMeta, err := s.VaultConnector.Upload(privateKeyBytes, userId, keyType, keyAlgorihm, keySize)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		if err := s.CryptoKeyRepo.Create(cryptoKeyMeta); err != nil {
			return nil, fmt.Errorf("failed to create metadata for key of type %s: %w", cryptoKeyMeta.Type, err)
		}

		cryptKeyMetas = append(cryptKeyMetas, cryptoKeyMeta)

		// PUB
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
		keyType = "public"

		if err != nil {
			return nil, fmt.Errorf("failed to marshal public key: %v", err)
		}

		cryptoKeyMeta, err = s.VaultConnector.Upload(publicKeyBytes, userId, keyType, keyAlgorihm, keySize)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		if err := s.CryptoKeyRepo.Create(cryptoKeyMeta); err != nil {
			return nil, fmt.Errorf("failed to create metadata for key of type %s: %w", cryptoKeyMeta.Type, err)
		}

		cryptKeyMetas = append(cryptKeyMetas, cryptoKeyMeta)
	}

	return cryptKeyMetas, nil
}

// CryptoKeyMetadataService manages cryptographic key metadata.
type CryptoKeyMetadataService struct {
	VaultConnector connector.VaultConnector
	CryptoKeyRepo  repository.CryptoKeyRepository
	Logger         logger.Logger
}

// NewCryptoKeyMetadataService creates a new CryptoKeyMetadataService instance
func NewCryptoKeyMetadataService(vaultConnector connector.VaultConnector, cryptoKeyRepo repository.CryptoKeyRepository, logger logger.Logger) (*CryptoKeyMetadataService, error) {
	return &CryptoKeyMetadataService{
		VaultConnector: vaultConnector,
		CryptoKeyRepo:  cryptoKeyRepo,
		Logger:         logger,
	}, nil
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
		return nil, fmt.Errorf("%w", err)
	}

	return keyMeta, nil
}

// DeleteByID deletes a cryptographic key's metadata by its ID.
func (s *CryptoKeyMetadataService) DeleteByID(keyId string) error {
	keyMeta, err := s.GetByID(keyId)
	if err != nil {
		return fmt.Errorf("failed to%w", err)
	}

	err = s.VaultConnector.Delete(keyId, keyMeta.Type)
	if err != nil {
		return fmt.Errorf("failed to%w", err)
	}

	err = s.CryptoKeyRepo.DeleteByID(keyId)
	if err != nil {
		return fmt.Errorf("failed to%w", err)
	}
	return nil
}

// CryptoKeyDownloadService handles the download of cryptographic keys.
type CryptoKeyDownloadService struct {
	VaultConnector connector.VaultConnector
	logger         logger.Logger
}

// NewCryptoKeyDownloadService creates a new CryptoKeyDownloadService instance
func NewCryptoKeyDownloadService(vaultConnector connector.VaultConnector, logger logger.Logger) (*CryptoKeyDownloadService, error) {
	return &CryptoKeyDownloadService{
		VaultConnector: vaultConnector,
		logger:         logger,
	}, nil
}

// Download retrieves a cryptographic key by its ID and type.
func (s *CryptoKeyDownloadService) Download(keyId, keyType string) ([]byte, error) {
	blobData, err := s.VaultConnector.Download(keyId, keyType)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return blobData, nil
}
