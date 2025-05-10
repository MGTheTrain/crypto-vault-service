package v1

import (
	"bytes"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/test/testutils"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBlobHandler_Upload(t *testing.T) {
	// Set up mock services
	mockUploadService := new(MockBlobUploadService)
	mockDownloadService := new(MockBlobDownloadService)
	mockMetadataService := new(MockBlobMetadataService)
	mockCryptoKeyUploadService := new(MockCryptoKeyUploadService)

	handler := NewBlobHandler(mockUploadService, mockDownloadService, mockMetadataService, mockCryptoKeyUploadService)

	// Create mock blob metadata
	blobMeta := blobs.BlobMeta{
		ID: "123",
	}

	// Mock the Upload service call
	mockUploadService.On("Upload", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return([]*blobs.BlobMeta{&blobMeta}, nil)

	// Create test file and form data
	fileName := "testfile.txt"
	fileContent := []byte("This is a test file content")
	form, err := testutils.CreateTestFileAndForm(t, fileName, fileContent)
	require.NoError(t, err)

	// Create a test HTTP request with the file attached as multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	fileWriter, err := writer.CreateFormFile("file", fileName)
	require.NoError(t, err)

	// Write the file content to the form
	_, err = fileWriter.Write(fileContent)
	require.NoError(t, err)

	// Close the writer to finish the multipart form
	writer.Close()

	// Create the request with multipart content
	req, err := http.NewRequest("POST", "/blobs", &b)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Set up Gin context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Request.MultipartForm = form

	// Call the Upload handler
	handler.Upload(c)

	// Assert the response
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "123") // Assert that blob ID is in response

	// Ensure the mock was called
	mockUploadService.AssertExpectations(t)
}

func TestBlobHandler_ListMetadata(t *testing.T) {
	// Set up mock services
	mockUploadService := new(MockBlobUploadService)
	mockDownloadService := new(MockBlobDownloadService)
	mockMetadataService := new(MockBlobMetadataService)
	mockCryptoKeyUploadService := new(MockCryptoKeyUploadService)

	handler := NewBlobHandler(mockUploadService, mockDownloadService, mockMetadataService, mockCryptoKeyUploadService)

	// Create mock blob metadata
	blobMeta := blobs.BlobMeta{
		ID: "123",
	}

	// Mock the List service call
	mockMetadataService.On("List", mock.Anything, mock.Anything).Return([]*blobs.BlobMeta{&blobMeta}, nil)

	// Create a test HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/blobs", nil)

	// Set up Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the ListMetadata handler
	handler.ListMetadata(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "123")
	mockMetadataService.AssertExpectations(t)
}

func TestBlobHandler_GetMetadataById(t *testing.T) {
	// Set up mock services
	mockUploadService := new(MockBlobUploadService)
	mockDownloadService := new(MockBlobDownloadService)
	mockMetadataService := new(MockBlobMetadataService)
	mockCryptoKeyUploadService := new(MockCryptoKeyUploadService)

	handler := NewBlobHandler(mockUploadService, mockDownloadService, mockMetadataService, mockCryptoKeyUploadService)

	// Create mock blob metadata
	blobMeta := blobs.BlobMeta{
		ID: "123",
	}

	// Mock the GetByID service call
	mockMetadataService.On("GetByID", mock.Anything, "123").Return(&blobMeta, nil)

	// Create a test HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/blobs/123", nil)

	// Set up Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

	// Call the GetMetadataById handler
	handler.GetMetadataById(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "123")
	mockMetadataService.AssertExpectations(t)
}

func TestBlobHandler_DownloadById(t *testing.T) {
	// Set up mock services
	mockUploadService := new(MockBlobUploadService)
	mockDownloadService := new(MockBlobDownloadService)
	mockMetadataService := new(MockBlobMetadataService)
	mockCryptoKeyUploadService := new(MockCryptoKeyUploadService)

	handler := NewBlobHandler(mockUploadService, mockDownloadService, mockMetadataService, mockCryptoKeyUploadService)

	// Prepare mock data
	blobId := "123"
	blobContent := []byte("file content")
	blobMeta := &blobs.BlobMeta{
		ID:   blobId,
		Name: "testfile.txt",
	}

	// Mock the Download and GetByID service calls
	mockDownloadService.On("Download", mock.Anything, blobId, (*string)(nil)).Return(blobContent, nil)
	mockMetadataService.On("GetByID", mock.Anything, blobId).Return(blobMeta, nil)

	// Create a test HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/blobs/123/file", nil)

	// Set up Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: blobId}}

	// Call the handler
	handler.DownloadById(c)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/octet-stream; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "attachment; filename="+blobMeta.Name, w.Header().Get("Content-Disposition"))
	assert.Equal(t, string(blobContent), w.Body.String())

	mockDownloadService.AssertExpectations(t)
	mockMetadataService.AssertExpectations(t)
}

func TestBlobHandler_DeleteById(t *testing.T) {
	// Set up mock services
	mockUploadService := new(MockBlobUploadService)
	mockDownloadService := new(MockBlobDownloadService)
	mockMetadataService := new(MockBlobMetadataService)
	mockCryptoKeyUploadService := new(MockCryptoKeyUploadService)

	handler := NewBlobHandler(mockUploadService, mockDownloadService, mockMetadataService, mockCryptoKeyUploadService)

	// Mock the DeleteByID service call
	mockMetadataService.On("DeleteByID", mock.Anything, "123").Return(nil, nil)

	// Create a test HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/blobs/123", nil)

	// Set up Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

	handler.DeleteById(c)

	// Assert the response
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockMetadataService.AssertExpectations(t)
}

// Test error case in GetMetadataById
func TestBlobHandler_GetMetadataById_Error(t *testing.T) {
	// Set up mock services
	mockUploadService := new(MockBlobUploadService)
	mockDownloadService := new(MockBlobDownloadService)
	mockMetadataService := new(MockBlobMetadataService)
	mockCryptoKeyUploadService := new(MockCryptoKeyUploadService)

	handler := NewBlobHandler(mockUploadService, mockDownloadService, mockMetadataService, mockCryptoKeyUploadService)

	// Mock an error in GetByID service call
	mockMetadataService.On("GetByID", mock.Anything, "123").Return(&blobs.BlobMeta{}, errors.New("not found"))

	// Create a test HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/blobs/123", nil)

	// Set up Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}

	// Call the GetMetadataById handler
	handler.GetMetadataById(c)

	// Assert the response
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "not found")
	mockMetadataService.AssertExpectations(t)
}

func TestBlobHandler_Upload_InvalidData_Error(t *testing.T) {
	// Set up mock services
	mockBlobUploadService := new(MockBlobUploadService)
	mockBlobMetadataService := new(MockBlobMetadataService)
	mockBlobDownloadService := new(MockBlobDownloadService)
	mockCryptoKeyUploadService := new(MockCryptoKeyUploadService)

	handler := NewBlobHandler(mockBlobUploadService, mockBlobDownloadService, mockBlobMetadataService, mockCryptoKeyUploadService)

	// Mock an error in the Upload service call
	mockBlobUploadService.On("Upload", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("invalid form data"))

	// Create a test HTTP request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/upload", nil)

	// Set up Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the Upload handler
	handler.Upload(c)

	// Assert the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid form data")
	// mockBlobUploadService.AssertExpectations(t)
}

func TestKeyHandler_UploadKeys(t *testing.T) {
	mockUploadService := new(MockCryptoKeyUploadService)
	mockDownloadService := new(MockCryptoKeyDownloadService)
	mockMetadataService := new(MockCryptoKeyMetadataService)

	handler := NewKeyHandler(mockUploadService, mockDownloadService, mockMetadataService)

	keyMeta := &keys.CryptoKeyMeta{
		ID:              "abc-123",
		KeyPairID:       "pair-123",
		Algorithm:       "RSA",
		KeySize:         2048,
		Type:            "private",
		DateTimeCreated: time.Now(),
		UserID:          "user-1",
	}

	requestBody := `{"algorithm": "RSA", "key_size": 2048}`

	mockUploadService.
		On("Upload", mock.Anything, mock.AnythingOfType("string"), "RSA", uint32(2048)).
		Return([]*keys.CryptoKeyMeta{keyMeta}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/keys", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.UploadKeys(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "abc-123")
	mockUploadService.AssertExpectations(t)
}

func TestKeyHandler_ListMetadata(t *testing.T) {
	mockUploadService := new(MockCryptoKeyUploadService)
	mockDownloadService := new(MockCryptoKeyDownloadService)
	mockMetadataService := new(MockCryptoKeyMetadataService)

	handler := NewKeyHandler(mockUploadService, mockDownloadService, mockMetadataService)

	keyMeta := &keys.CryptoKeyMeta{
		ID:              "abc-123",
		KeyPairID:       "pair-123",
		Algorithm:       "RSA",
		KeySize:         2048,
		Type:            "private",
		DateTimeCreated: time.Now(),
		UserID:          "user-1",
	}

	mockMetadataService.
		On("List", mock.Anything, mock.Anything).
		Return([]*keys.CryptoKeyMeta{keyMeta}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/keys", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ListMetadata(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "abc-123")
	mockMetadataService.AssertExpectations(t)
}

func TestKeyHandler_GetMetadataById(t *testing.T) {
	mockUploadService := new(MockCryptoKeyUploadService)
	mockDownloadService := new(MockCryptoKeyDownloadService)
	mockMetadataService := new(MockCryptoKeyMetadataService)

	handler := NewKeyHandler(mockUploadService, mockDownloadService, mockMetadataService)

	keyMeta := &keys.CryptoKeyMeta{
		ID:              "abc-123",
		KeyPairID:       "pair-123",
		Algorithm:       "RSA",
		KeySize:         2048,
		Type:            "private",
		DateTimeCreated: time.Now(),
		UserID:          "user-1",
	}

	mockMetadataService.
		On("GetByID", mock.Anything, "abc-123").
		Return(keyMeta, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/keys/abc-123", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "abc-123"}}

	handler.GetMetadataById(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "abc-123")
	mockMetadataService.AssertExpectations(t)
}

func TestKeyHandler_DownloadById(t *testing.T) {
	mockUploadService := new(MockCryptoKeyUploadService)
	mockDownloadService := new(MockCryptoKeyDownloadService)
	mockMetadataService := new(MockCryptoKeyMetadataService)

	handler := NewKeyHandler(mockUploadService, mockDownloadService, mockMetadataService)

	keyID := "abc-123"
	keyContent := []byte("secret key content")

	mockDownloadService.
		On("Download", mock.Anything, keyID).
		Return(keyContent, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/keys/abc-123/file", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: keyID}}

	handler.DownloadById(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/octet-stream; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, "attachment; filename="+keyID, w.Header().Get("Content-Disposition"))
	assert.Equal(t, string(keyContent), w.Body.String())

	mockDownloadService.AssertExpectations(t)
}

func TestKeyHandler_DeleteById(t *testing.T) {
	mockUploadService := new(MockCryptoKeyUploadService)
	mockDownloadService := new(MockCryptoKeyDownloadService)
	mockMetadataService := new(MockCryptoKeyMetadataService)

	handler := NewKeyHandler(mockUploadService, mockDownloadService, mockMetadataService)

	keyID := "abc-123"

	mockMetadataService.
		On("DeleteByID", mock.Anything, keyID).
		Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/keys/abc-123", nil)

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: keyID}}

	handler.DeleteById(c)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockMetadataService.AssertExpectations(t)
}
