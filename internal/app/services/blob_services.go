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

	// Upload the new form
	blobMetas, err := s.BlobConnector.Upload(newForm, userId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	// Handle empty uploads
	if len(blobMetas) == 0 {
		return nil, fmt.Errorf("no blobs uploaded")
	}

	// Save blob meta data
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
				publicKey, err := x509.ParsePKCS1PublicKey(keyBytes)
				if err != nil {
					return nil, nil, fmt.Errorf("error parsing public key: %w", err)
				}
				processedBytes, err = rsa.Encrypt(data, publicKey)
				if err != nil {
					return nil, nil, fmt.Errorf("%w", err)
				}
			} else if operation == "signing" {
				privateKey, err := x509.ParsePKCS1PrivateKey(keyBytes)
				if err != nil {
					return nil, nil, fmt.Errorf("error parsing public key: %w", err)
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
				privateKey, err := x509.ParseECPrivateKey(keyBytes)
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
	var blobMetas []*blobs.BlobMeta

	// TBD

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
	Logger         logger.Logger
}

// NewBlobDownloadService creates a new instance of BlobDownloadService
func NewBlobDownloadService(blobConnector connector.BlobConnector, blobRepository repository.BlobRepository, logger logger.Logger) *BlobDownloadService {
	return &BlobDownloadService{
		BlobConnector:  blobConnector,
		BlobRepository: blobRepository,
		Logger:         logger,
	}
}

// Download retrieves a blob's content by its ID and name
func (s *BlobDownloadService) Download(blobId string, decryptionKeyId, verificationKeyId *string) ([]byte, error) {
	blobMeta, err := s.BlobRepository.GetById(blobId)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	blob, err := s.BlobConnector.Download(blobId, blobMeta.Name)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return blob, nil
}
