//go:build unit
// +build unit

package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TODO(MGTheTrain): Modify tests
func TestSetupRoutes(t *testing.T) {
	mockBlobUploadService := new(MockBlobUploadService)
	mockBlobDownloadService := new(MockBlobDownloadService)
	mockBlobMetadataService := new(MockBlobMetadataService)
	mockCryptoKeyUploadService := new(MockCryptoKeyUploadService)
	mockCryptoKeyDownloadService := new(MockCryptoKeyDownloadService)
	mockCryptoKeyMetadataService := new(MockCryptoKeyMetadataService)

	// Create Gin engine
	r := gin.Default()

	mockBlobUploadService.On("Upload", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)
	mockBlobMetadataService.On("List", mock.Anything, mock.Anything).Return(nil, nil)
	mockBlobMetadataService.On("GetByID", mock.Anything, mock.Anything).Return(nil, nil)
	mockBlobMetadataService.On("DeleteByID", mock.Anything, mock.Anything).Return(nil, nil)
	mockBlobDownloadService.On("DownloadById", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	mockCryptoKeyUploadService.
		On("Upload", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil)
	mockCryptoKeyMetadataService.
		On("List", mock.Anything, mock.Anything).
		Return(nil, nil)
	mockCryptoKeyMetadataService.
		On("GetByID", mock.Anything, mock.Anything).
		Return(nil, nil)
	mockCryptoKeyDownloadService.
		On("DownloadById", mock.Anything, mock.Anything).
		Return(nil, nil)
	mockCryptoKeyMetadataService.
		On("DeleteByID", mock.Anything, mock.Anything).
		Return(nil)

	// Call SetupRoutes to register routes
	SetupRoutes(r, mockBlobUploadService, mockBlobDownloadService, mockBlobMetadataService, mockCryptoKeyUploadService, mockCryptoKeyDownloadService, mockCryptoKeyMetadataService)

	// Define test cases for different routes
	tests := []struct {
		method         string
		url            string
		expectedStatus int
	}{
		{"POST", "/api/v1/cvs/blobs", http.StatusBadRequest},
		// {"GET", "/api/v1/cvs/blobs", http.StatusBadRequest},
		// {"GET", "/api/v1/cvs/blobs/123", http.StatusBadRequest},
		// {"GET", "/api/v1/cvs/blobs/123/file", http.StatusBadRequest},
		// {"DELETE", "/api/v1/cvs/blobs/123", http.StatusNoContent},
		{"POST", "/api/v1/cvs/keys", http.StatusBadRequest},
		// {"GET", "/api/v1/cvs/keys", http.StatusOK},
		// {"GET", "/api/v1/cvs/keys/123", http.StatusOK},
		// {"GET", "/api/v1/cvs/keys/123/file", http.StatusOK},
		// {"DELETE", "/api/v1/cvs/keys/123", http.StatusNoContent},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.url, func(t *testing.T) {
			// Create a request
			req, _ := http.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			// Simulate the request
			r.ServeHTTP(w, req)

			// Assert the status code
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
