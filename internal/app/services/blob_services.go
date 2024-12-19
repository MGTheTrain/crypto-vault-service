package services

import (
	"bytes"
	"crypto/x509"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/utils"
	"crypto_vault_service/internal/persistence/repository"
	"encoding/pem"
	"fmt"
	"io"
	"mime/multipart"
)

// BlobUploadService implements the BlobUploadService interface for handling blob uploads
type BlobUploadService struct {
	BlobConnector  connector.BlobConnector
	BlobRepository repository.BlobRepository
	VaultConnector connector.VaultConnector
	CryptoKeyRepo  repository.CryptoKeyRepository
	Logger         logger.Logger
}

// NewBlobUploadService creates a new instance of BlobUploadService
func NewBlobUploadService(blobConnector connector.BlobConnector, blobRepository repository.BlobRepository, vaultConnector connector.VaultConnector, cryptoKeyRepo repository.CryptoKeyRepository, logger logger.Logger) *BlobUploadService {
	return &BlobUploadService{
		BlobConnector:  blobConnector,
		BlobRepository: blobRepository,
		CryptoKeyRepo:  cryptoKeyRepo,
		VaultConnector: vaultConnector,
		Logger:         logger,
	}
}

// Upload transfers blobs with the option to encrypt them using an encryption key or sign them with a signing key.
// It returns a slice of Blob for the uploaded blobs and any error encountered during the upload process.
func (s *BlobUploadService) Upload(form *multipart.Form, userId string, encryptionKeyId, signKeyId *string) ([]*blobs.BlobMeta, error) {
	var newForm *multipart.Form

	// Process encryptionKeyId if provided
	if encryptionKeyId != nil {
		keyBytes, cryptoKeyMeta, err := s.getCryptoKeyAndData(*encryptionKeyId)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		cryptoOperation := "encryption"
		contents, fileNames, err := s.applyCryptographicOperation(form, cryptoKeyMeta.Algorithm, cryptoOperation, keyBytes)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		newForm, err = utils.CreateMultipleFilesForm(contents, fileNames)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	// Process signKeyId if provided
	if signKeyId != nil {
		keyBytes, cryptoKeyMeta, err := s.getCryptoKeyAndData(*signKeyId)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		cryptoOperation := "signing"
		contents, fileNames, err := s.applyCryptographicOperation(form, cryptoKeyMeta.Algorithm, cryptoOperation, keyBytes)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		newForm, err = utils.CreateMultipleFilesForm(contents, fileNames)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	if signKeyId != nil || encryptionKeyId != nil {
		//
		blobMetas, err := s.BlobConnector.Upload(newForm, userId)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		for _, blobMeta := range blobMetas {
			err := s.BlobRepository.Create(blobMeta)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}
		}
		return blobMetas, nil
	}

	//
	blobMetas, err := s.BlobConnector.Upload(form, userId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	for _, blobMeta := range blobMetas {
		err := s.BlobRepository.Create(blobMeta)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	return blobMetas, nil
}

// getCryptoKeyAndData retrieves the encryption or signing key along with its metadata by ID.
// It downloads the key from the vault and returns the key bytes and associated metadata.
func (s *BlobUploadService) getCryptoKeyAndData(cryptoKeyId string) ([]byte, *keys.CryptoKeyMeta, error) {
	// Get meta info
	cryptoKeyMeta, err := s.CryptoKeyRepo.GetByID(cryptoKeyId)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	// Download key
	keyBytes, err := s.VaultConnector.Download(cryptoKeyMeta.ID, cryptoKeyMeta.KeyPairID, cryptoKeyMeta.Type)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	return keyBytes, cryptoKeyMeta, nil
}

// applyCryptographicOperation performs cryptographic operations (encryption or signing)
// on files within a multipart form using the specified algorithm and key.
func (s *BlobUploadService) applyCryptographicOperation(form *multipart.Form, algorithm, operation string, keyBytes []byte) ([][]byte, []string, error) {
	var contents [][]byte
	var fileNames []string

	fileHeaders := form.File["files"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, nil, fmt.Errorf("%w", err)
		}
		defer file.Close()

		buffer := bytes.NewBuffer(make([]byte, 0))
		_, err = io.Copy(buffer, file)
		if err != nil {
			return nil, nil, fmt.Errorf("%w", err)
		}
		data := buffer.Bytes()

		var processedBytes []byte

		switch algorithm {
		case "AES":
			if operation == "encryption" {
				aes, err := cryptography.NewAES(s.Logger)
				if err != nil {
					return nil, nil, fmt.Errorf("%w", err)
				}
				processedBytes, err = aes.Encrypt(data, keyBytes)
				if err != nil {
					return nil, nil, fmt.Errorf("%w", err)
				}
			}
		case "RSA":
			rsa, err := cryptography.NewRSA(s.Logger)
			if err != nil {
				return nil, nil, fmt.Errorf("%w", err)
			}
			if operation == "encryption" {
				block, _ := pem.Decode(keyBytes)
				if block == nil || block.Type != "RSA PUBLIC KEY" {
					return nil, nil, fmt.Errorf("invalid public key PEM block")
				}
				publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
				if err != nil {
					return nil, nil, fmt.Errorf("error parsing public key: %w", err)
				}
				processedBytes, err = rsa.Encrypt(data, publicKey)
				if err != nil {
					return nil, nil, fmt.Errorf("%w", err)
				}
			} else if operation == "signing" {
				block, _ := pem.Decode(keyBytes)
				if block == nil || block.Type != "RSA PRIVATE KEY" {
					return nil, nil, fmt.Errorf("invalid private key PEM block")
				}
				privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
				if err != nil {
					return nil, nil, fmt.Errorf("error parsing private key: %w", err)
				}
				processedBytes, err = rsa.Sign(data, privateKey)
				if err != nil {
					return nil, nil, fmt.Errorf("%w", err)
				}
			}
		case "EC":
			if operation == "signing" {
				ec, err := cryptography.NewEC(s.Logger)
				if err != nil {
					return nil, nil, fmt.Errorf("%w", err)
				}
				block, _ := pem.Decode(keyBytes)
				if block == nil || block.Type != "EC PRIVATE KEY" {
					return nil, nil, fmt.Errorf("invalid private key PEM block")
				}
				privateKey, err := x509.ParseECPrivateKey(block.Bytes)
				if err != nil {
					return nil, nil, fmt.Errorf("error parsing private key: %w", err)
				}
				processedBytes, err = ec.Sign(data, privateKey)
				if err != nil {
					return nil, nil, fmt.Errorf("%w", err)
				}
			}
		default:
			return nil, nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
		}

		contents = append(contents, processedBytes)
		fileNames = append(fileNames, fileHeader.Filename)
	}

	return contents, fileNames, nil
}

// BlobMetadataService implements the BlobMetadataService interface for retrieving and deleting blob metadata
type BlobMetadataService struct {
	BlobConnector  connector.BlobConnector
	BlobRepository repository.BlobRepository
	Logger         logger.Logger
}

// NewBlobMetadataService creates a new instance of BlobMetadataService
func NewBlobMetadataService(blobRepository repository.BlobRepository, blobConnector connector.BlobConnector, logger logger.Logger) *BlobMetadataService {
	return &BlobMetadataService{
		BlobConnector:  blobConnector,
		BlobRepository: blobRepository,
		Logger:         logger,
	}
}

// List retrieves all blobs' metadata considering a query filter
func (s *BlobMetadataService) List(query *blobs.BlobMetaQuery) ([]*blobs.BlobMeta, error) {
	blobMetas, err := s.BlobRepository.List(query)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return blobMetas, nil
}

// GetByID retrieves a blob's metadata by its unique ID
func (s *BlobMetadataService) GetByID(blobId string) (*blobs.BlobMeta, error) {
	blobMeta, err := s.BlobRepository.GetById(blobId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return blobMeta, nil
}

// DeleteByID deletes a blob and its associated metadata by ID
func (s *BlobMetadataService) DeleteByID(blobId string) error {

	blobMeta, err := s.BlobRepository.GetById(blobId)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = s.BlobRepository.DeleteById(blobId)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = s.BlobConnector.Delete(blobMeta.ID, blobMeta.Name)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// BlobDownloadService implements the BlobDownloadService interface for downloading blobs
type BlobDownloadService struct {
	BlobConnector  connector.BlobConnector
	BlobRepository repository.BlobRepository
	VaultConnector connector.VaultConnector
	CryptoKeyRepo  repository.CryptoKeyRepository
	Logger         logger.Logger
}

// NewBlobDownloadService creates a new instance of BlobDownloadService
func NewBlobDownloadService(blobConnector connector.BlobConnector, blobRepository repository.BlobRepository, vaultConnector connector.VaultConnector, cryptoKeyRepo repository.CryptoKeyRepository, logger logger.Logger) *BlobDownloadService {
	return &BlobDownloadService{
		BlobConnector:  blobConnector,
		BlobRepository: blobRepository,
		CryptoKeyRepo:  cryptoKeyRepo,
		VaultConnector: vaultConnector,
		Logger:         logger,
	}
}

// The download function retrieves a blob's content using its ID and also enables data decryption.
// NOTE: Signing should be performed locally by first downloading the associated key, followed by verification.
// Optionally, a verify endpoint will be available soon for optional use.
func (s *BlobDownloadService) Download(blobId string, decryptionKeyId *string) ([]byte, error) {

	blobMeta, err := s.BlobRepository.GetById(blobId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	blobBytes, err := s.BlobConnector.Download(blobId, blobMeta.Name)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if decryptionKeyId != nil {
		var processedBytes []byte
		keyBytes, cryptoKeyMeta, err := s.getCryptoKeyAndData(*decryptionKeyId)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		switch cryptoKeyMeta.Algorithm {
		case "AES":
			aes, err := cryptography.NewAES(s.Logger)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}
			processedBytes, err = aes.Decrypt(blobBytes, keyBytes)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}
		case "RSA":
			rsa, err := cryptography.NewRSA(s.Logger)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}
			block, _ := pem.Decode(keyBytes)
			if block == nil || block.Type != "RSA PRIVATE KEY" {
				return nil, fmt.Errorf("invalid private key PEM block")
			}
			privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("error parsing private key: %w", err)
			}
			processedBytes, err = rsa.Decrypt(blobBytes, privateKey)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}
		default:
			return nil, fmt.Errorf("unsupported algorithm: %s", cryptoKeyMeta.Algorithm)
		}
		return processedBytes, nil
	}
	return blobBytes, nil
}

// getCryptoKeyAndData retrieves the encryption or signing key along with its metadata by ID.
// It downloads the key from the vault and returns the key bytes and associated metadata.
func (s *BlobDownloadService) getCryptoKeyAndData(cryptoKeyId string) ([]byte, *keys.CryptoKeyMeta, error) {
	// Get meta info
	cryptoKeyMeta, err := s.CryptoKeyRepo.GetByID(cryptoKeyId)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	// Download key
	keyBytes, err := s.VaultConnector.Download(cryptoKeyMeta.ID, cryptoKeyMeta.KeyPairID, cryptoKeyMeta.Type)
	if err != nil {
		return nil, nil, fmt.Errorf("%w", err)
	}

	return keyBytes, cryptoKeyMeta, nil
}
