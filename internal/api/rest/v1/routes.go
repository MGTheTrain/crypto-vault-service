package v1

import (
	"crypto_vault_service/internal/app/services"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all the API routes for version 1.
func SetupRoutes(r *gin.Engine,
	blobUploadService *services.BlobUploadService,
	blobDownloadService *services.BlobDownloadService,
	blobMetadataService *services.BlobMetadataService,
	cryptoKeyUploadService *services.CryptoKeyUploadService,
	cryptoKeyDownloadService *services.CryptoKeyDownloadService,
	cryptoKeyMetadataService *services.CryptoKeyMetadataService) {

	v1 := r.Group(BasePath) // lookup in version file

	// Blobs Routes
	blobHandler := NewBlobHandler(blobUploadService, blobDownloadService, blobMetadataService, cryptoKeyUploadService)
	v1.POST("/blobs", blobHandler.Upload)
	v1.GET("/blobs", blobHandler.ListMetadata)
	v1.GET("/blobs/:id", blobHandler.GetMetadataById)
	v1.GET("/blobs/:id/file", blobHandler.DownloadById)
	v1.DELETE("/blobs/:id", blobHandler.DeleteById)

	// Keys Routes
	keyHandler := NewKeyHandler(cryptoKeyUploadService, cryptoKeyDownloadService, cryptoKeyMetadataService)
	v1.POST("/keys", keyHandler.UploadKeys)
	v1.GET("/keys", keyHandler.ListMetadata)
	v1.GET("/keys/:id", keyHandler.GetMetadataById)
	v1.GET("/keys/:id/file", keyHandler.DownloadById)
	v1.DELETE("/keys/:id", keyHandler.DeleteById)
}
