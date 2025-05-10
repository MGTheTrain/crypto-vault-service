package v1

import (
	"context"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"mime/multipart"

	"github.com/stretchr/testify/mock"
)

// Mock the services
type MockBlobUploadService struct {
	mock.Mock
}

func (m *MockBlobUploadService) Upload(ctx context.Context, form *multipart.Form, userID string, encryptionKeyID, signKeyID *string) ([]*blobs.BlobMeta, error) {
	args := m.Called(ctx, form, userID, encryptionKeyID, signKeyID)
	return args.Get(0).([]*blobs.BlobMeta), args.Error(1)
}

type MockBlobMetadataService struct {
	mock.Mock
}

func (m *MockBlobMetadataService) List(ctx context.Context, query *blobs.BlobMetaQuery) ([]*blobs.BlobMeta, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*blobs.BlobMeta), args.Error(1)
}

func (m *MockBlobMetadataService) GetByID(ctx context.Context, blobId string) (*blobs.BlobMeta, error) {
	args := m.Called(ctx, blobId)
	return args.Get(0).(*blobs.BlobMeta), args.Error(1)
}

func (m *MockBlobMetadataService) DeleteByID(ctx context.Context, blobId string) error {
	args := m.Called(ctx, blobId)
	return args.Error(0)
}

type MockBlobDownloadService struct {
	mock.Mock
}

func (m *MockBlobDownloadService) DownloadById(ctx context.Context, blobId string, decryptionKeyId *string) ([]byte, error) {
	args := m.Called(ctx, blobId, decryptionKeyId)
	return args.Get(0).([]byte), args.Error(1)
}

type MockCryptoKeyUploadService struct {
	mock.Mock
}

func (m *MockCryptoKeyUploadService) Upload(ctx context.Context, userId, keyAlgorithm string, keySize uint32) ([]*keys.CryptoKeyMeta, error) {
	args := m.Called(ctx, userId, keyAlgorithm, keySize)
	return args.Get(0).([]*keys.CryptoKeyMeta), args.Error(1)
}

type MockCryptoKeyMetadataService struct {
	mock.Mock
}

func (m *MockCryptoKeyMetadataService) List(ctx context.Context, query *keys.CryptoKeyQuery) ([]*keys.CryptoKeyMeta, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*keys.CryptoKeyMeta), args.Error(1)
}

func (m *MockCryptoKeyMetadataService) GetByID(ctx context.Context, keyID string) (*keys.CryptoKeyMeta, error) {
	args := m.Called(ctx, keyID)
	return args.Get(0).(*keys.CryptoKeyMeta), args.Error(1)
}

func (m *MockCryptoKeyMetadataService) DeleteByID(ctx context.Context, keyID string) error {
	args := m.Called(ctx, keyID)
	return args.Error(0)
}

type MockCryptoKeyDownloadService struct {
	mock.Mock
}

func (m *MockCryptoKeyDownloadService) DownloadById(ctx context.Context, keyID string) ([]byte, error) {
	args := m.Called(ctx, keyID)
	return args.Get(0).([]byte), args.Error(1)
}
