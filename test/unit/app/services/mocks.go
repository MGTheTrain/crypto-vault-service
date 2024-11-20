package services

import (
	"crypto_vault_service/internal/domain/blobs"
	"fmt"

	"github.com/stretchr/testify/mock"
)

// MockBlobConnector is a mock for the BlobConnector interface
type MockBlobConnector struct {
	mock.Mock
}

func (m *MockBlobConnector) Upload(filePaths []string) ([]*blobs.BlobMeta, error) {
	args := m.Called(filePaths)
	return args.Get(0).([]*blobs.BlobMeta), args.Error(1)
}

func (m *MockBlobConnector) Delete(blobID, blobName string) error {
	args := m.Called(blobID, blobName)
	return args.Error(0)
}

func (m *MockBlobConnector) Download(blobID, blobName string) ([]byte, error) {
	args := m.Called(blobID, blobName)
	data, ok := args.Get(0).([]byte)
	if !ok {
		return nil, fmt.Errorf("expected []byte, but got %T", args.Get(0))
	}
	return data, args.Error(1)
}

// MockBlobRepository is a mock for the BlobRepository interface
type MockBlobRepository struct {
	mock.Mock
}

// Create mocks the Create method of the BlobRepository interface
func (m *MockBlobRepository) Create(blob *blobs.BlobMeta) error {
	args := m.Called(blob)
	return args.Error(0)
}

// GetById mocks the GetById method of the BlobRepository interface
func (m *MockBlobRepository) GetById(blobID string) (*blobs.BlobMeta, error) {
	args := m.Called(blobID)
	return args.Get(0).(*blobs.BlobMeta), args.Error(1)
}

// UpdateById mocks the UpdateById method of the BlobRepository interface
func (m *MockBlobRepository) UpdateById(blob *blobs.BlobMeta) error {
	args := m.Called(blob)
	return args.Error(0)
}

// DeleteById mocks the DeleteById method of the BlobRepository interface
func (m *MockBlobRepository) DeleteById(blobID string) error {
	args := m.Called(blobID)
	return args.Error(0)
}
