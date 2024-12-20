package v1

import (
	"crypto_vault_service/internal/app/services"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BlobHandler struct holds the services
type BlobHandler struct {
	blobUploadService      *services.BlobUploadService
	blobDownloadService    *services.BlobDownloadService
	blobMetadataService    *services.BlobMetadataService
	cryptoKeyUploadService *services.CryptoKeyUploadService
}

// NewBlobHandler creates a new BlobHandler
func NewBlobHandler(blobUploadService *services.BlobUploadService, blobDownloadService *services.BlobDownloadService, blobMetadataService *services.BlobMetadataService, cryptoKeyUploadService *services.CryptoKeyUploadService) *BlobHandler {
	return &BlobHandler{
		blobUploadService:      blobUploadService,
		blobDownloadService:    blobDownloadService,
		blobMetadataService:    blobMetadataService,
		cryptoKeyUploadService: cryptoKeyUploadService,
	}
}

// Upload handles the POST request to upload a blob with optional encryption/signing
func (handler *BlobHandler) Upload(c *gin.Context) {
	var form *multipart.Form
	var encryptionKeyId, signKeyId *string
	userId := uuid.New().String() // TBD: extract user id from JWT

	form, err := c.MultipartForm()
	if err != nil {
		var errorResponseDto ErrorResponseDto
		errorResponseDto.Message = "invalid form data"
		c.JSON(http.StatusBadRequest, errorResponseDto)
		return
	}

	if encryptionKeys := form.Value["encryption_key_id"]; len(encryptionKeys) > 0 {
		encryptionKeyId = &encryptionKeys[0]
	}

	if signKeys := form.Value["sign_key_id"]; len(signKeys) > 0 {
		signKeyId = &signKeys[0]
	}

	blobMetas, err := handler.blobUploadService.Upload(form, userId, encryptionKeyId, signKeyId)
	if err != nil {
		var errorResponseDto ErrorResponseDto
		errorResponseDto.Message = "Error uploading blob"
		c.JSON(http.StatusBadRequest, errorResponseDto)
		return
	}

	var blobMetadataResponses []BlobMetaResponseDto
	for _, blobMeta := range blobMetas {
		blobMetadataResponse := BlobMetaResponseDto{
			ID:              blobMeta.ID,
			DateTimeCreated: blobMeta.DateTimeCreated,
			UserID:          blobMeta.UserID,
			Name:            blobMeta.Name,
			Size:            blobMeta.Size,
			Type:            blobMeta.Type,
			EncryptionKeyID: blobMeta.EncryptionKeyID,
			SignKeyID:       blobMeta.SignKeyID,
		}
		blobMetadataResponses = append(blobMetadataResponses, blobMetadataResponse)
	}

	c.JSON(http.StatusCreated, blobMetadataResponses)
}

// ListMetadata handles the GET request to fetch metadata of blobs by query
func (handler *BlobHandler) ListMetadata(c *gin.Context) {
	var query *blobs.BlobMetaQuery = nil

	// TBD: extract query parameters with Gin

	blobMetas, err := handler.blobMetadataService.List(query)
	if err != nil {
		var errorResponseDto ErrorResponseDto
		errorResponseDto.Message = "List query failed"
		c.JSON(http.StatusNotFound, errorResponseDto)
		return
	}

	var blobMetadataResponses []BlobMetaResponseDto
	for _, blobMeta := range blobMetas {
		blobMetadataResponse := BlobMetaResponseDto{
			ID:              blobMeta.ID,
			DateTimeCreated: blobMeta.DateTimeCreated,
			UserID:          blobMeta.UserID,
			Name:            blobMeta.Name,
			Size:            blobMeta.Size,
			Type:            blobMeta.Type,
			EncryptionKeyID: blobMeta.EncryptionKeyID,
			SignKeyID:       blobMeta.SignKeyID,
		}
		blobMetadataResponses = append(blobMetadataResponses, blobMetadataResponse)
	}

	c.JSON(http.StatusOK, blobMetadataResponses)
}

// GetMetadataById handles the GET request to fetch metadata of a blob by its ID
func (handler *BlobHandler) GetMetadataById(c *gin.Context) {
	blobId := c.Param("id")

	blobMeta, err := handler.blobMetadataService.GetByID(blobId)
	if err != nil {
		var errorResponseDto ErrorResponseDto
		errorResponseDto.Message = fmt.Sprintf("Blob with id %s not found", blobId)
		c.JSON(http.StatusNotFound, errorResponseDto)
		return
	}

	blobMetadataResponse := BlobMetaResponseDto{
		ID:              blobMeta.ID,
		DateTimeCreated: blobMeta.DateTimeCreated,
		UserID:          blobMeta.UserID,
		Name:            blobMeta.Name,
		Size:            blobMeta.Size,
		Type:            blobMeta.Type,
		EncryptionKeyID: blobMeta.EncryptionKeyID,
		SignKeyID:       blobMeta.SignKeyID,
	}

	c.JSON(http.StatusOK, blobMetadataResponse)
}

// DeleteById handles the DELETE request to delete a blob by its ID
func (handler *BlobHandler) DeleteById(c *gin.Context) {
	blobId := c.Param("id")

	if err := handler.blobMetadataService.DeleteByID(blobId); err != nil {
		var errorResponseDto ErrorResponseDto
		errorResponseDto.Message = fmt.Sprintf("Error deleting blob with id %s", blobId)
		c.JSON(http.StatusNotFound, errorResponseDto)
		return
	}

	var infoResponseDto InfoResponseDto
	infoResponseDto.Message = fmt.Sprintf("Deleted blob with id %s", blobId)
	c.JSON(http.StatusNoContent, infoResponseDto)
}

// KeyHandler struct holds the services
type KeyHandler struct {
	cryptoKeyUploadService   *services.CryptoKeyUploadService
	cryptoKeyDownloadService *services.CryptoKeyDownloadService
	cryptoKeyMetadataService *services.CryptoKeyMetadataService
}

// NewKeyHandler creates a new KeyHandler with a registered custom validator
func NewKeyHandler(cryptoKeyUploadService *services.CryptoKeyUploadService, cryptoKeyDownloadService *services.CryptoKeyDownloadService, cryptoKeyMetadataService *services.CryptoKeyMetadataService) *KeyHandler {

	return &KeyHandler{
		cryptoKeyUploadService:   cryptoKeyUploadService,
		cryptoKeyDownloadService: cryptoKeyDownloadService,
		cryptoKeyMetadataService: cryptoKeyMetadataService,
	}
}

// UploadKeys handles the POST request to internally generate cryptographic keys and upload those
func (handler *KeyHandler) UploadKeys(c *gin.Context) {

	var requestDto UploadKeyRequestDto

	if err := c.ShouldBindJSON(&requestDto); err != nil {
		var errorResponseDto ErrorResponseDto
		errorResponseDto.Message = "Invalid key data"
		c.JSON(http.StatusBadRequest, errorResponseDto)
		return
	}

	if err := requestDto.Validate(); err != nil {
		var errorResponseDto ErrorResponseDto
		errorResponseDto.Message = "Validation failed"
		c.JSON(400, errorResponseDto)
		return
	}

	userId := uuid.New().String() // TBD: extract user id from JWT

	cryptoKeyMetas, err := handler.cryptoKeyUploadService.Upload(userId, requestDto.Algorithm, requestDto.KeySize)
	if err != nil {
		var errorResponseDto ErrorResponseDto
		errorResponseDto.Message = "Error uploading key"
		c.JSON(http.StatusBadRequest, errorResponseDto)
		return
	}

	var cryptoKeyMetadataResponses []CryptoKeyMetaResponseDto
	for _, cryptoKeyMeta := range cryptoKeyMetas {
		cryptoKeyMetadataResponse := CryptoKeyMetaResponseDto{
			ID:              cryptoKeyMeta.ID,
			KeyPairID:       cryptoKeyMeta.KeyPairID,
			Algorithm:       cryptoKeyMeta.Algorithm,
			KeySize:         cryptoKeyMeta.KeySize,
			Type:            cryptoKeyMeta.Type,
			DateTimeCreated: cryptoKeyMeta.DateTimeCreated,
			UserID:          cryptoKeyMeta.UserID,
		}
		cryptoKeyMetadataResponses = append(cryptoKeyMetadataResponses, cryptoKeyMetadataResponse)
	}

	c.JSON(http.StatusCreated, cryptoKeyMetadataResponses)
}

// GetMetadataById handles the GET request to retrieve metadata of a key by its ID
func (handler *KeyHandler) GetMetadataById(c *gin.Context) {
	keyId := c.Param("id")

	cryptoKeyMeta, err := handler.cryptoKeyMetadataService.GetByID(keyId)
	if err != nil {
		var errorResponseDto ErrorResponseDto
		errorResponseDto.Message = fmt.Sprintf("key with id %s not found", keyId)
		c.JSON(http.StatusNotFound, errorResponseDto)
		return
	}

	cryptoKeyMetadataResponse := CryptoKeyMetaResponseDto{
		ID:              cryptoKeyMeta.ID,
		KeyPairID:       cryptoKeyMeta.KeyPairID,
		Algorithm:       cryptoKeyMeta.Algorithm,
		KeySize:         cryptoKeyMeta.KeySize,
		Type:            cryptoKeyMeta.Type,
		DateTimeCreated: cryptoKeyMeta.DateTimeCreated,
		UserID:          cryptoKeyMeta.UserID,
	}

	c.JSON(http.StatusOK, cryptoKeyMetadataResponse)
}

// ListMetadata handles the GET request to list cryptographic keys metadata
func (handler *KeyHandler) ListMetadata(c *gin.Context) {
	var query *keys.CryptoKeyQuery = nil

	cryptoKeyMetas, err := handler.cryptoKeyMetadataService.List(query)
	if err != nil {
		var errorResponseDto ErrorResponseDto
		errorResponseDto.Message = "List query failed"
		c.JSON(http.StatusNotFound, errorResponseDto)
		return
	}

	var cryptoKeyMetadataResponses []CryptoKeyMetaResponseDto
	for _, cryptoKeyMeta := range cryptoKeyMetas {
		cryptoKeyMetadataResponse := CryptoKeyMetaResponseDto{
			ID:              cryptoKeyMeta.ID,
			KeyPairID:       cryptoKeyMeta.KeyPairID,
			Algorithm:       cryptoKeyMeta.Algorithm,
			KeySize:         cryptoKeyMeta.KeySize,
			Type:            cryptoKeyMeta.Type,
			DateTimeCreated: cryptoKeyMeta.DateTimeCreated,
			UserID:          cryptoKeyMeta.UserID,
		}
		cryptoKeyMetadataResponses = append(cryptoKeyMetadataResponses, cryptoKeyMetadataResponse)
	}

	c.JSON(http.StatusOK, cryptoKeyMetadataResponses)
}

// DeleteById handles the DELETE request to delete a key by its ID
func (handler *KeyHandler) DeleteById(c *gin.Context) {
	keyId := c.Param("id")

	if err := handler.cryptoKeyMetadataService.DeleteByID(keyId); err != nil {
		var errorResponseDto ErrorResponseDto
		errorResponseDto.Message = fmt.Sprintf("Error deleting key with id %s", keyId)
		c.JSON(http.StatusNotFound, errorResponseDto)
		return
	}

	var infoResponseDto InfoResponseDto
	infoResponseDto.Message = fmt.Sprintf("Deleted key with id %s", keyId)
	c.JSON(http.StatusNoContent, infoResponseDto)
}
