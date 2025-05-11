package v1

import (
	"context"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"fmt"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

// MockBlobUploadService is a mock implementation of the BlobUploadService used for testing.
// It simulates the Upload method for uploading blobs.
type MockBlobUploadService struct {
	mock.Mock
}

// Upload simulates uploading a blob and returns mocked blob metadata or an error.
func (m *MockBlobUploadService) Upload(ctx context.Context, form *multipart.Form, userID string, encryptionKeyID, signKeyID *string) ([]*blobs.BlobMeta, error) {
	args := m.Called(ctx, form, userID, encryptionKeyID, signKeyID)

	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock Upload error: %w", err)
	}

	return args.Get(0).([]*blobs.BlobMeta), nil
}

// MockBlobMetadataService is a mock implementation of the BlobMetadataService used for testing.
// It simulates operations related to listing, fetching by ID, and deleting blob metadata.
type MockBlobMetadataService struct {
	mock.Mock
}

// List simulates listing blob metadata based on a query.
func (m *MockBlobMetadataService) List(ctx context.Context, query *blobs.BlobMetaQuery) ([]*blobs.BlobMeta, error) {
	args := m.Called(ctx, query)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock List error: %w", err)
	}
	return args.Get(0).([]*blobs.BlobMeta), nil
}

// GetByID simulates fetching a blob's metadata by its ID.
func (m *MockBlobMetadataService) GetByID(ctx context.Context, blobID string) (*blobs.BlobMeta, error) {
	args := m.Called(ctx, blobID)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock GetByID error: %w", err)
	}
	return args.Get(0).(*blobs.BlobMeta), nil
}

// DeleteByID simulates deleting a blob's metadata by its ID.
func (m *MockBlobMetadataService) DeleteByID(ctx context.Context, blobID string) error {
	args := m.Called(ctx, blobID)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("mock DeleteByID error: %w", err)
	}
	return nil
}

// MockBlobDownloadService is a mock implementation of the BlobDownloadService used for testing.
// It simulates downloading a blob by ID.
type MockBlobDownloadService struct {
	mock.Mock
}

// DownloadByID simulates downloading a blob by its ID, possibly using a decryption key.
func (m *MockBlobDownloadService) DownloadByID(ctx context.Context, blobID string, decryptionKeyID *string) ([]byte, error) {
	args := m.Called(ctx, blobID, decryptionKeyID)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock DownloadByID error: %w", err)
	}
	return args.Get(0).([]byte), nil
}

// MockCryptoKeyUploadService is a mock implementation of the CryptoKeyUploadService used for testing.
// It simulates the upload of cryptographic keys.
type MockCryptoKeyUploadService struct {
	mock.Mock
}

// Upload simulates uploading a cryptographic key and returns mocked key metadata or an error.
func (m *MockCryptoKeyUploadService) Upload(ctx context.Context, userID, keyAlgorithm string, keySize uint32) ([]*keys.CryptoKeyMeta, error) {
	args := m.Called(ctx, userID, keyAlgorithm, keySize)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock Upload error: %w", err)
	}
	return args.Get(0).([]*keys.CryptoKeyMeta), nil
}

// MockCryptoKeyMetadataService is a mock implementation of the CryptoKeyMetadataService used for testing.
// It simulates operations related to listing, fetching by ID, and deleting cryptographic key metadata.
type MockCryptoKeyMetadataService struct {
	mock.Mock
}

// List simulates listing cryptographic key metadata based on a query.
func (m *MockCryptoKeyMetadataService) List(ctx context.Context, query *keys.CryptoKeyQuery) ([]*keys.CryptoKeyMeta, error) {
	args := m.Called(ctx, query)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock List error: %w", err)
	}
	return args.Get(0).([]*keys.CryptoKeyMeta), nil
}

// GetByID simulates fetching a cryptographic key's metadata by its ID.
func (m *MockCryptoKeyMetadataService) GetByID(ctx context.Context, keyID string) (*keys.CryptoKeyMeta, error) {
	args := m.Called(ctx, keyID)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock GetByID error: %w", err)
	}
	return args.Get(0).(*keys.CryptoKeyMeta), nil
}

// DeleteByID simulates deleting a cryptographic key's metadata by its ID.
func (m *MockCryptoKeyMetadataService) DeleteByID(ctx context.Context, keyID string) error {
	args := m.Called(ctx, keyID)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("mock DeleteByID error: %w", err)
	}
	return nil
}

// MockCryptoKeyDownloadService is a mock implementation of the CryptoKeyDownloadService used for testing.
// It simulates downloading a cryptographic key by ID.
type MockCryptoKeyDownloadService struct {
	mock.Mock
}

// DownloadByID simulates downloading a cryptographic key by its ID.
func (m *MockCryptoKeyDownloadService) DownloadByID(ctx context.Context, keyID string) ([]byte, error) {
	args := m.Called(ctx, keyID)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock DownloadByID error: %w", err)
	}
	return args.Get(0).([]byte), nil
}
