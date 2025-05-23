package v1

import (
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/utils"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BlobHandler defines the interface for handling blob-related operations
type BlobHandler interface {
	Upload(ctx *gin.Context)
	ListMetadata(ctx *gin.Context)
	GetMetadataByID(ctx *gin.Context)
	DownloadByID(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
}

// BlobHandler struct holds the services
type blobHandler struct {
	blobUploadService      blobs.BlobUploadService
	blobMetadataService    blobs.BlobMetadataService
	blobDownloadService    blobs.BlobDownloadService
	cryptoKeyUploadService keys.CryptoKeyUploadService
}

// NewBlobHandler creates a new BlobHandler
func NewBlobHandler(blobUploadService blobs.BlobUploadService, blobDownloadService blobs.BlobDownloadService, blobMetadataService blobs.BlobMetadataService, cryptoKeyUploadService keys.CryptoKeyUploadService) BlobHandler {
	return &blobHandler{
		blobUploadService:      blobUploadService,
		blobDownloadService:    blobDownloadService,
		blobMetadataService:    blobMetadataService,
		cryptoKeyUploadService: cryptoKeyUploadService,
	}
}

// Upload handles the POST request to upload a blob with optional encryption/signing
// @Summary Upload a blob with optional encryption and signing
// @Description Upload a blob to the system with optional encryption and signing using the provided keys
// @Tags Blob
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Blob File"
// @Param encryption_key_id formData string false "Encryption Key ID"
// @Param sign_key_id formData string false "Sign Key ID"
// @Success 201 {array} BlobMetaResponse
// @Failure 400 {object} ErrorResponse
// @Router /blobs [post]
func (handler *blobHandler) Upload(ctx *gin.Context) {
	var form *multipart.Form
	var encryptionKeyID *string
	var signKeyID *string
	userID := uuid.New().String() // TODO(MGTheTrain): extract user id from JWT

	form, err := ctx.MultipartForm()
	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = "invalid form data"
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if encryptionKeys := form.Value["encryption_key_id"]; len(encryptionKeys) > 0 {
		encryptionKeyID = &encryptionKeys[0]
	}

	if signKeys := form.Value["sign_key_id"]; len(signKeys) > 0 {
		signKeyID = &signKeys[0]
	}

	blobMetas, err := handler.blobUploadService.Upload(ctx, form, userID, encryptionKeyID, signKeyID)
	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("error uploading blob: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var blobMetadataResponses []BlobMetaResponse
	for _, blobMeta := range blobMetas {
		blobMetadataResponse := BlobMetaResponse{
			ID:              blobMeta.ID,
			DateTimeCreated: blobMeta.DateTimeCreated,
			UserID:          blobMeta.UserID,
			Name:            blobMeta.Name,
			Size:            blobMeta.Size,
			Type:            blobMeta.Type,
			EncryptionKeyID: nil,
			SignKeyID:       nil,
		}
		if blobMeta.EncryptionKeyID != nil {
			blobMetadataResponse.EncryptionKeyID = blobMeta.EncryptionKeyID
		}
		if blobMeta.SignKeyID != nil {
			blobMetadataResponse.SignKeyID = blobMeta.SignKeyID
		}
		blobMetadataResponses = append(blobMetadataResponses, blobMetadataResponse)
	}

	ctx.JSON(http.StatusCreated, blobMetadataResponses)
}

// ListMetadata handles the GET request to fetch metadata of blobs optionally considering query parameters
// @Summary List blob metadata based on query parameters
// @Description Fetch a list of metadata for blobs based on query filters like name, size, type, and creation date.
// @Tags Blob
// @Accept json
// @Produce json
// @Param name query string false "Blob Name"
// @Param size query int false "Blob Size"
// @Param type query string false "Blob Type"
// @Param dateTimeCreated query string false "Blob Creation Date (RFC3339)"
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset the results"
// @Success 200 {array} BlobMetaResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /blobs [get]
func (handler *blobHandler) ListMetadata(ctx *gin.Context) {
	query := blobs.NewBlobMetaQuery()

	if blobName := ctx.Query("name"); len(blobName) > 0 {
		query.Name = blobName
	}

	if blobSize := ctx.Query("size"); len(blobSize) > 0 {
		query.Size = utils.ConvertToInt64(blobSize)
	}

	if blobType := ctx.Query("type"); len(blobType) > 0 {
		query.Type = blobType
	}

	if dateTimeCreated := ctx.Query("dateTimeCreated"); len(dateTimeCreated) > 0 {
		parsedTime, err := time.Parse(time.RFC3339, dateTimeCreated)
		if err == nil {
			query.DateTimeCreated = parsedTime
		}
	}

	if limit := ctx.Query("limit"); len(limit) > 0 {
		query.Limit = utils.ConvertToInt(limit)
	}

	if offset := ctx.Query("offset"); len(offset) > 0 {
		query.Offset = utils.ConvertToInt(offset)
	}

	if sortBy := ctx.Query("sortBy"); len(sortBy) > 0 {
		query.SortBy = sortBy
	}

	if sortOrder := ctx.Query("sortOrder"); len(sortOrder) > 0 {
		query.SortOrder = sortOrder
	}

	if err := query.Validate(); err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("validation failed: %v", err.Error())
		ctx.JSON(400, errorResponse)
		return
	}

	blobMetas, err := handler.blobMetadataService.List(ctx, query)
	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("list query failed: %v", err.Error())
		ctx.JSON(http.StatusNotFound, errorResponse)
		return
	}

	var listResponse = []BlobMetaResponse{}
	for _, blobMeta := range blobMetas {
		blobMetadataResponse := BlobMetaResponse{
			ID:              blobMeta.ID,
			DateTimeCreated: blobMeta.DateTimeCreated,
			UserID:          blobMeta.UserID,
			Name:            blobMeta.Name,
			Size:            blobMeta.Size,
			Type:            blobMeta.Type,
			EncryptionKeyID: nil,
			SignKeyID:       nil,
		}
		if blobMeta.EncryptionKeyID != nil {
			blobMetadataResponse.EncryptionKeyID = blobMeta.EncryptionKeyID
		}
		if blobMeta.SignKeyID != nil {
			blobMetadataResponse.SignKeyID = blobMeta.SignKeyID
		}
		listResponse = append(listResponse, blobMetadataResponse)
	}

	ctx.JSON(http.StatusOK, listResponse)
}

// GetMetadataByID handles the GET request to fetch metadata of a blob by its ID
// @Summary Retrieve metadata of a blob by its ID
// @Description Fetch the metadata of a specific blob by its unique ID, including its name, size, type, encryption and signing key IDs, and creation date.
// @Tags Blob
// @Accept json
// @Produce json
// @Param id path string true "Blob ID"
// @Success 200 {object} BlobMetaResponse
// @Failure 404 {object} ErrorResponse
// @Router /blobs/{id} [get]
func (handler *blobHandler) GetMetadataByID(ctx *gin.Context) {
	blobID := ctx.Param("id")

	blobMeta, err := handler.blobMetadataService.GetByID(ctx, blobID)
	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("blob with id %s not found", blobID)
		ctx.JSON(http.StatusNotFound, errorResponse)
		return
	}

	blobMetadataResponse := BlobMetaResponse{
		ID:              blobMeta.ID,
		DateTimeCreated: blobMeta.DateTimeCreated,
		UserID:          blobMeta.UserID,
		Name:            blobMeta.Name,
		Size:            blobMeta.Size,
		Type:            blobMeta.Type,
		EncryptionKeyID: nil,
		SignKeyID:       nil,
	}

	if blobMeta.EncryptionKeyID != nil {
		blobMetadataResponse.EncryptionKeyID = blobMeta.EncryptionKeyID
	}
	if blobMeta.SignKeyID != nil {
		blobMetadataResponse.SignKeyID = blobMeta.SignKeyID
	}

	ctx.JSON(http.StatusOK, blobMetadataResponse)
}

// DownloadByID handles the GET request to download a blob by its ID
// @Summary Download a blob by its ID
// @Description Download the content of a specific blob by its ID, optionally decrypted with a provided decryption key ID.
// @Tags Blob
// @Accept json
// @Produce octet-stream
// @Param id path string true "Blob ID"
// @Param decryption_key_id query string false "Decryption Key ID"
// @Success 200 {file} file "Blob content"
// @Failure 404 {object} ErrorResponse
// @Router /blobs/{id}/file [get]
func (handler *blobHandler) DownloadByID(ctx *gin.Context) {
	blobID := ctx.Param("id")

	var decryptionKeyID *string
	if decryptionKeyQuery := ctx.Query("decryption_key_id"); len(decryptionKeyQuery) > 0 {
		decryptionKeyID = &decryptionKeyQuery
	}

	bytes, err := handler.blobDownloadService.DownloadByID(ctx, blobID, decryptionKeyID)
	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("could not download blob with id %s: %v", blobID, err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	blobMeta, err := handler.blobMetadataService.GetByID(ctx, blobID)
	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("blob with id %s not found", blobID)
		ctx.JSON(http.StatusNotFound, errorResponse)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Header().Set("Content-Type", "application/octet-stream; charset=utf-8")
	ctx.Writer.Header().Set("Content-Disposition", "attachment; filename="+blobMeta.Name)
	_, err = ctx.Writer.Write(bytes)

	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("could not write bytes: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}
}

// DeleteByID handles the DELETE request to delete a blob by its ID
// @Summary Delete a blob by its ID
// @Description Delete a specific blob and its associated metadata by its ID.
// @Tags Blob
// @Accept json
// @Produce json
// @Param id path string true "Blob ID"
// @Success 204 {object} InfoResponse
// @Failure 404 {object} ErrorResponse
// @Router /blobs/{id} [delete]
func (handler *blobHandler) DeleteByID(ctx *gin.Context) {
	blobID := ctx.Param("id")

	if err := handler.blobMetadataService.DeleteByID(ctx, blobID); err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("blob with id %s not found", blobID)
		ctx.JSON(http.StatusNotFound, errorResponse)
		return
	}

	var infoResponse InfoResponse
	infoResponse.Message = fmt.Sprintf("deleted blob with id %s", blobID)
	ctx.JSON(http.StatusNoContent, infoResponse)
}

// KeyHandler defines the interface for handling key-related operations
type KeyHandler interface {
	UploadKeys(ctx *gin.Context)
	ListMetadata(ctx *gin.Context)
	GetMetadataByID(ctx *gin.Context)
	DownloadByID(ctx *gin.Context)
	DeleteByID(ctx *gin.Context)
}

// KeyHandler struct holds the services
type keyHandler struct {
	cryptoKeyUploadService   keys.CryptoKeyUploadService
	cryptoKeyDownloadService keys.CryptoKeyDownloadService
	cryptoKeyMetadataService keys.CryptoKeyMetadataService
}

// NewKeyHandler creates a new KeyHandler
func NewKeyHandler(cryptoKeyUploadService keys.CryptoKeyUploadService, cryptoKeyDownloadService keys.CryptoKeyDownloadService, cryptoKeyMetadataService keys.CryptoKeyMetadataService) KeyHandler {
	return &keyHandler{
		cryptoKeyUploadService:   cryptoKeyUploadService,
		cryptoKeyDownloadService: cryptoKeyDownloadService,
		cryptoKeyMetadataService: cryptoKeyMetadataService,
	}
}

// UploadKeys handles the POST request to generate and upload cryptographic keys
// @Summary Upload cryptographic keys and metadata
// @Description Generate cryptographic keys based on provided parameters and upload them to the system.
// @Tags Key
// @Accept json
// @Produce json
// @Param requestBody body UploadKeyRequest true "Cryptographic Key Data"
// @Success 201 {array} CryptoKeyMetaResponse
// @Failure 400 {object} ErrorResponse
// @Router /keys [post]
func (handler *keyHandler) UploadKeys(ctx *gin.Context) {

	var request UploadKeyRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("invalid key data: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if err := request.Validate(); err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("validation failed: %v", err.Error())
		ctx.JSON(400, errorResponse)
		return
	}

	userID := uuid.New().String() // TODO(MGTheTrain): extract user id from JWT

	cryptoKeyMetas, err := handler.cryptoKeyUploadService.Upload(ctx, userID, request.Algorithm, request.KeySize)
	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("error uploading key: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	var listResponse = []CryptoKeyMetaResponse{}
	for _, cryptoKeyMeta := range cryptoKeyMetas {
		cryptoKeyMetadataResponse := CryptoKeyMetaResponse{
			ID:              cryptoKeyMeta.ID,
			KeyPairID:       cryptoKeyMeta.KeyPairID,
			Algorithm:       cryptoKeyMeta.Algorithm,
			KeySize:         cryptoKeyMeta.KeySize,
			Type:            cryptoKeyMeta.Type,
			DateTimeCreated: cryptoKeyMeta.DateTimeCreated,
			UserID:          cryptoKeyMeta.UserID,
		}
		listResponse = append(listResponse, cryptoKeyMetadataResponse)
	}

	ctx.JSON(http.StatusCreated, listResponse)
}

// ListMetadata handles the GET request to list cryptographic key metadata with optional query parameters
// @Summary List cryptographic key metadata based on query parameters
// @Description Fetch a list of cryptographic key metadata based on filters like algorithm, type, and creation date, with pagination and sorting options.
// @Tags Key
// @Accept json
// @Produce json
// @Param algorithm query string false "Cryptographic Algorithm"
// @Param type query string false "Key Type"
// @Param dateTimeCreated query string false "Key Creation Date (RFC3339)"
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset the results"
// @Param sortBy query string false "Sort by a specific field"
// @Param sortOrder query string false "Sort order (asc/desc)"
// @Success 200 {array} CryptoKeyMetaResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /keys [get]
func (handler *keyHandler) ListMetadata(ctx *gin.Context) {
	query := keys.NewCryptoKeyQuery()

	if keyAlgorithm := ctx.Query("algorithm"); len(keyAlgorithm) > 0 {
		query.Algorithm = keyAlgorithm
	}

	if keyType := ctx.Query("type"); len(keyType) > 0 {
		query.Type = keyType
	}

	if dateTimeCreated := ctx.Query("dateTimeCreated"); len(dateTimeCreated) > 0 {
		parsedTime, err := time.Parse(time.RFC3339, dateTimeCreated)
		if err == nil {
			query.DateTimeCreated = parsedTime
		}
	}

	if limit := ctx.Query("limit"); len(limit) > 0 {
		query.Limit = utils.ConvertToInt(limit)
	}

	if offset := ctx.Query("offset"); len(offset) > 0 {
		query.Offset = utils.ConvertToInt(offset)
	}

	if sortBy := ctx.Query("sortBy"); len(sortBy) > 0 {
		query.SortBy = sortBy
	}

	if sortOrder := ctx.Query("sortOrder"); len(sortOrder) > 0 {
		query.SortOrder = sortOrder
	}

	if err := query.Validate(); err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("validation failed: %v", err.Error())
		ctx.JSON(400, errorResponse)
		return
	}

	cryptoKeyMetas, err := handler.cryptoKeyMetadataService.List(ctx, query)
	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("list query failed: %v", err.Error())
		ctx.JSON(http.StatusNotFound, errorResponse)
		return
	}

	var listResponse = []CryptoKeyMetaResponse{}
	for _, cryptoKeyMeta := range cryptoKeyMetas {
		cryptoKeyMetadataResponse := CryptoKeyMetaResponse{
			ID:              cryptoKeyMeta.ID,
			KeyPairID:       cryptoKeyMeta.KeyPairID,
			Algorithm:       cryptoKeyMeta.Algorithm,
			KeySize:         cryptoKeyMeta.KeySize,
			Type:            cryptoKeyMeta.Type,
			DateTimeCreated: cryptoKeyMeta.DateTimeCreated,
			UserID:          cryptoKeyMeta.UserID,
		}
		listResponse = append(listResponse, cryptoKeyMetadataResponse)
	}

	ctx.JSON(http.StatusOK, listResponse)
}

// GetMetadataByID handles the GET request to retrieve metadata of a key by its ID
// @Summary Retrieve metadata of a key by its ID
// @Description Fetch the metadata of a specific cryptographic key by its unique ID, including algorithm, key size, and creation date.
// @Tags Key
// @Accept json
// @Produce json
// @Param id path string true "Key ID"
// @Success 200 {object} CryptoKeyMetaResponse
// @Failure 404 {object} ErrorResponse
// @Router /keys/{id} [get]
func (handler *keyHandler) GetMetadataByID(ctx *gin.Context) {
	keyID := ctx.Param("id")

	cryptoKeyMeta, err := handler.cryptoKeyMetadataService.GetByID(ctx, keyID)
	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("key with id %s not found", keyID)
		ctx.JSON(http.StatusNotFound, errorResponse)
		return
	}

	cryptoKeyMetadataResponse := CryptoKeyMetaResponse{
		ID:              cryptoKeyMeta.ID,
		KeyPairID:       cryptoKeyMeta.KeyPairID,
		Algorithm:       cryptoKeyMeta.Algorithm,
		KeySize:         cryptoKeyMeta.KeySize,
		Type:            cryptoKeyMeta.Type,
		DateTimeCreated: cryptoKeyMeta.DateTimeCreated,
		UserID:          cryptoKeyMeta.UserID,
	}

	ctx.JSON(http.StatusOK, cryptoKeyMetadataResponse)
}

// DownloadByID handles the GET request to download a key by its ID
// @Summary Download a cryptographic key by its ID
// @Description Download the content of a specific cryptographic key by its ID.
// @Tags Key
// @Accept json
// @Produce octet-stream
// @Param id path string true "Key ID"
// @Success 200 {file} file "Cryptographic key content"
// @Failure 404 {object} ErrorResponse
// @Router /keys/{id}/file [get]
func (handler *keyHandler) DownloadByID(ctx *gin.Context) {
	keyID := ctx.Param("id")

	bytes, err := handler.cryptoKeyDownloadService.DownloadByID(ctx, keyID)
	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("could not download key with id %s: %v", keyID, err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Header().Set("Content-Type", "application/octet-stream; charset=utf-8")
	ctx.Writer.Header().Set("Content-Disposition", "attachment; filename="+keyID)
	_, err = ctx.Writer.Write(bytes)

	if err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("could not write bytes: %v", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}
}

// DeleteByID handles the DELETE request to delete a key by its ID
// @Summary Delete a cryptographic key by its ID
// @Description Delete a specific cryptographic key and its associated metadata by its ID.
// @Tags Key
// @Accept json
// @Produce json
// @Param id path string true "Key ID"
// @Success 204 {object} InfoResponse
// @Failure 404 {object} ErrorResponse
// @Router /keys/{id} [delete]
func (handler *keyHandler) DeleteByID(ctx *gin.Context) {
	keyID := ctx.Param("id")

	if err := handler.cryptoKeyMetadataService.DeleteByID(ctx, keyID); err != nil {
		var errorResponse ErrorResponse
		errorResponse.Message = fmt.Sprintf("error deleting key with id %s", keyID)
		ctx.JSON(http.StatusNotFound, errorResponse)
		return
	}

	var infoResponse InfoResponse
	infoResponse.Message = fmt.Sprintf("deleted key with id %s", keyID)
	ctx.JSON(http.StatusNoContent, infoResponse)
}
