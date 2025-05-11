package v1

import (
	"context"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"fmt"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

// Mock the services
type MockBlobUploadService struct {
	mock.Mock
}

func (m *MockBlobUploadService) Upload(ctx context.Context, form *multipart.Form, userID string, encryptionKeyID, signKeyID *string) ([]*blobs.BlobMeta, error) {
	args := m.Called(ctx, form, userID, encryptionKeyID, signKeyID)

	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock Upload error: %w", err)
	}

	return args.Get(0).([]*blobs.BlobMeta), nil
}

type MockBlobMetadataService struct {
	mock.Mock
}

func (m *MockBlobMetadataService) List(ctx context.Context, query *blobs.BlobMetaQuery) ([]*blobs.BlobMeta, error) {
	args := m.Called(ctx, query)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock List error: %w", err)
	}
	return args.Get(0).([]*blobs.BlobMeta), nil
}

func (m *MockBlobMetadataService) GetByID(ctx context.Context, blobID string) (*blobs.BlobMeta, error) {
	args := m.Called(ctx, blobID)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock GetByID error: %w", err)
	}
	return args.Get(0).(*blobs.BlobMeta), nil
}

func (m *MockBlobMetadataService) DeleteByID(ctx context.Context, blobID string) error {
	args := m.Called(ctx, blobID)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("mock DeleteByID error: %w", err)
	}
	return nil
}

type MockBlobDownloadService struct {
	mock.Mock
}

func (m *MockBlobDownloadService) DownloadByID(ctx context.Context, blobID string, decryptionKeyId *string) ([]byte, error) {
	args := m.Called(ctx, blobID, decryptionKeyId)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock DownloadByID error: %w", err)
	}
	return args.Get(0).([]byte), nil
}

type MockCryptoKeyUploadService struct {
	mock.Mock
}

func (m *MockCryptoKeyUploadService) Upload(ctx context.Context, userId, keyAlgorithm string, keySize uint32) ([]*keys.CryptoKeyMeta, error) {
	args := m.Called(ctx, userId, keyAlgorithm, keySize)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock Upload error: %w", err)
	}
	return args.Get(0).([]*keys.CryptoKeyMeta), nil
}

type MockCryptoKeyMetadataService struct {
	mock.Mock
}

func (m *MockCryptoKeyMetadataService) List(ctx context.Context, query *keys.CryptoKeyQuery) ([]*keys.CryptoKeyMeta, error) {
	args := m.Called(ctx, query)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock List error: %w", err)
	}
	return args.Get(0).([]*keys.CryptoKeyMeta), nil
}

func (m *MockCryptoKeyMetadataService) GetByID(ctx context.Context, keyID string) (*keys.CryptoKeyMeta, error) {
	args := m.Called(ctx, keyID)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock GetByID error: %w", err)
	}
	return args.Get(0).(*keys.CryptoKeyMeta), nil
}

func (m *MockCryptoKeyMetadataService) DeleteByID(ctx context.Context, keyID string) error {
	args := m.Called(ctx, keyID)
	err := args.Error(0)
	if err != nil {
		return fmt.Errorf("mock DeleteByID error: %w", err)
	}
	return nil
}

type MockCryptoKeyDownloadService struct {
	mock.Mock
}

func (m *MockCryptoKeyDownloadService) DownloadByID(ctx context.Context, keyID string) ([]byte, error) {
	args := m.Called(ctx, keyID)
	err := args.Error(1)
	if err != nil {
		return nil, fmt.Errorf("mock DownloadByID error: %w", err)
	}
	return args.Get(0).([]byte), nil
}
