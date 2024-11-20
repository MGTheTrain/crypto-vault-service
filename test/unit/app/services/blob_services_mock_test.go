package services

import (
	"crypto_vault_service/internal/app/services"
	"crypto_vault_service/internal/domain/blobs"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBlobUploadService_Upload(t *testing.T) {
	// Prepare mock dependencies
	mockBlobConnector := new(MockBlobConnector)
	mockBlobRepository := new(MockBlobRepository)

	// Initialize the service with mock dependencies
	service := services.NewBlobUploadService(mockBlobConnector, mockBlobRepository)

	// Test case 1: Successful upload and metadata storage
	t.Run("successfully uploads blobs and stores metadata", func(t *testing.T) {
		// Define the test file paths
		filePaths := []string{"file1.txt", "file2.txt"}

		// Define mock return values
		mockBlobMeta := []*blobs.BlobMeta{
			{Name: "file1.txt", ID: "1"},
			{Name: "file2.txt", ID: "2"},
		}

		// Setup expectations for the mock BlobConnector
		mockBlobConnector.On("Upload", filePaths).Return(mockBlobMeta, nil)

		// Setup expectations for the mock BlobRepository
		mockBlobRepository.On("Create", mockBlobMeta[0]).Return(nil)
		mockBlobRepository.On("Create", mockBlobMeta[1]).Return(nil)

		// Call the method under test
		userId := uuid.New().String()
		uploadedBlobs, err := service.Upload(filePaths, userId)

		// Assert the results
		assert.NoError(t, err)
		assert.Equal(t, mockBlobMeta, uploadedBlobs)

		// Assert that the expectations were met
		mockBlobConnector.AssertExpectations(t)
		mockBlobRepository.AssertExpectations(t)
	})

	// Test case 2: Failed upload (BlobConnector returns error)
	t.Run("fails when BlobConnector returns error", func(t *testing.T) {
		// Define the test file paths
		filePaths := []string{"file1.txt"}

		// Setup expectations for the mock BlobConnector
		mockBlobConnector.On("Upload", filePaths).Return([]*blobs.BlobMeta(nil), fmt.Errorf("upload failed"))

		// Call the method under test
		userId := uuid.New().String()
		uploadedBlobs, err := service.Upload(filePaths, userId)

		// Assert the results
		assert.Error(t, err)
		assert.Nil(t, uploadedBlobs)

		// Assert that the expectations were met
		mockBlobConnector.AssertExpectations(t)
		mockBlobRepository.AssertExpectations(t)
	})

	// Test case 3: Failed metadata storage (BlobRepository returns error)
	t.Run("fails when BlobRepository returns error", func(t *testing.T) {
		// Define the test file paths
		filePaths := []string{"file1.txt"}

		// Define mock return values
		mockBlobMeta := []*blobs.BlobMeta{
			{Name: "file1.txt", ID: "1"},
		}

		// Setup expectations for the mock BlobConnector
		mockBlobConnector.On("Upload", filePaths).Return(mockBlobMeta, nil)

		// Setup expectations for the mock BlobRepository
		mockBlobRepository.On("Create", mockBlobMeta[0]).Return(fmt.Errorf("failed to store metadata"))

		// Call the method under test
		userId := uuid.New().String()
		uploadedBlobs, err := service.Upload(filePaths, userId)

		// Assert the results
		assert.Error(t, err)
		assert.Nil(t, uploadedBlobs)

		// Assert that the expectations were met
		mockBlobConnector.AssertExpectations(t)
		mockBlobRepository.AssertExpectations(t)
	})
}

func TestBlobMetadataService(t *testing.T) {
	// Prepare mock dependencies
	mockBlobConnector := new(MockBlobConnector)
	mockBlobRepository := new(MockBlobRepository)

	// Initialize the service with mock dependencies
	metadataService := services.NewBlobMetadataService(mockBlobRepository, mockBlobConnector)

	// Test case 1: Successfully retrieve blob metadata by ID
	t.Run("successfully retrieves blob metadata by ID", func(t *testing.T) {
		blobID := "1"
		mockBlobMeta := &blobs.BlobMeta{
			ID:   "1",
			Name: "file1.txt",
		}

		// Setup expectations for the mock BlobRepository
		mockBlobRepository.On("GetById", blobID).Return(mockBlobMeta, nil)

		// Call the method under test
		blob, err := metadataService.GetByID(blobID)

		// Assert the results
		assert.NoError(t, err)
		assert.Equal(t, mockBlobMeta, blob)

		// Assert that the expectations were met
		mockBlobRepository.AssertExpectations(t)
	})
}

func TestBlobDownloadService(t *testing.T) {
	// Prepare mock dependencies
	mockBlobConnector := new(MockBlobConnector)

	// Initialize the service with mock dependencies
	downloadService := services.NewBlobDownloadService(mockBlobConnector)

	// Test case 1: Successfully download a blob
	t.Run("successfully downloads a blob", func(t *testing.T) {
		blobID := "1"
		blobName := "file1.txt"
		mockContent := []byte("file content")

		// Setup expectations for the mock BlobConnector to return the content as []byte
		mockBlobConnector.On("Download", blobID, blobName).Return(mockContent, nil)

		// Call the method under test
		downloadedBlob, err := downloadService.Download(blobID, blobName)

		// Assert the results
		assert.NoError(t, err)
		assert.Equal(t, mockContent, downloadedBlob) // Assert equality of byte slices

		// Assert that the expectations were met
		mockBlobConnector.AssertExpectations(t)
	})
}
