package v1

import (
	"context"
	"crypto_vault_service/internal/app/services"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/utils"
	"fmt"

	pb "crypto_vault_service/internal/api/grpc/v1/proto"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type BlobUploadServiceServer struct {
	pb.UnimplementedBlobUploadServiceServer
	blobUploadService *services.BlobUploadService
}

type BlobDownloadServiceServer struct {
	pb.UnimplementedBlobDownloadServiceServer
	blobDownloadService *services.BlobDownloadService
}

type BlobMetadataServiceServer struct {
	pb.UnimplementedBlobMetadataServiceServer
	blobMetadataService *services.BlobMetadataService
}

type CryptoKeyUploadServiceServer struct {
	pb.UnimplementedCryptoKeyUploadServiceServer
	cryptoKeyUploadService *services.CryptoKeyUploadService
}

type CryptoKeyDownloadServiceServer struct {
	pb.UnimplementedCryptoKeyDownloadServiceServer
	cryptoKeyDownloadService *services.CryptoKeyDownloadService
}

type CryptoKeyMetadataServiceServer struct {
	pb.UnimplementedCryptoKeyMetadataServiceServer
	cryptoKeyMetadataService *services.CryptoKeyMetadataService
}

// NewBlobUploadServiceServer creates a new instance of BlobUploadServiceServer.
func NewBlobUploadServiceServer(blobUploadService *services.BlobUploadService) (*BlobUploadServiceServer, error) {
	return &BlobUploadServiceServer{
		blobUploadService: blobUploadService,
	}, nil
}

// Upload uploads a blob with optional encryption/signing
func (s BlobUploadServiceServer) Upload(req *pb.BlobUploadRequest, stream pb.BlobUploadService_UploadServer) error {
	fileContent := [][]byte{req.GetFileContent()}
	fileNames := []string{req.GetFileName()}

	var encryptionKeyId *string = nil
	var signKeyId *string = nil

	if len(req.GetEncryptionKeyId()) > 0 {
		encryptionKeyId = &req.EncryptionKeyId
	}

	if len(req.GetSignKeyId()) > 0 {
		signKeyId = &req.SignKeyId
	}

	form, err := utils.CreateMultipleFilesForm(fileContent, fileNames)
	blobMetas, err := s.blobUploadService.Upload(form, req.GetUserId(), encryptionKeyId, signKeyId)
	if err != nil {
		return fmt.Errorf("failed to upload blob: %v", err)
	}

	for _, blobMeta := range blobMetas {
		blobMetaResponse := &pb.BlobMetaResponse{
			Id:              blobMeta.ID,
			DateTimeCreated: timestamppb.New(blobMeta.DateTimeCreated),
			UserId:          blobMeta.UserID,
			Name:            blobMeta.Name,
			Size:            blobMeta.Size,
			Type:            blobMeta.Type,
			EncryptionKeyId: "",
			SignKeyId:       "",
		}

		if blobMeta.EncryptionKeyID != nil {
			blobMetaResponse.EncryptionKeyId = *blobMeta.EncryptionKeyID
		}
		if blobMeta.SignKeyID != nil {
			blobMetaResponse.SignKeyId = *blobMeta.SignKeyID
		}

		// Send the metadata response to the client
		if err := stream.Send(blobMetaResponse); err != nil {
			return fmt.Errorf("failed to send metadata response: %v", err)
		}
	}

	return nil
}

// NewBlobDownloadServiceServer creates a new instance of BlobDownloadServiceServer.
func NewBlobDownloadServiceServer(blobDownloadService *services.BlobDownloadService) (*BlobDownloadServiceServer, error) {
	return &BlobDownloadServiceServer{
		blobDownloadService: blobDownloadService,
	}, nil
}

// DownloadById downloads a blob by its ID
func (s *BlobDownloadServiceServer) DownloadById(req *pb.BlobDownloadRequest, stream pb.BlobDownloadService_DownloadByIdServer) error {
	id := req.GetId()
	var decryptionKeyId *string = nil
	if len(req.GetDecryptionKeyId()) > 0 {
		decryptionKeyId = &req.DecryptionKeyId
	}

	bytes, err := s.blobDownloadService.Download(id, decryptionKeyId)
	if err != nil {
		return fmt.Errorf("could not download blob with id %s: %v", id, err)
	}

	// If no error, stream the blob content back in chunks
	chunkSize := 1024 * 1024 // 1MB chunk size, adjust as needed
	for i := 0; i < len(bytes); i += chunkSize {
		end := i + chunkSize
		if end > len(bytes) {
			end = len(bytes)
		}

		// Create the chunk of data to send
		chunk := &pb.BlobContent{
			Content: bytes[i:end],
		}

		// Send the chunk
		if err := stream.Send(chunk); err != nil {
			return fmt.Errorf("failed to send chunk: %v", err)
		}
	}
	return nil
}

// NewBlobMetadataServiceServer creates a new instance of BlobMetadataServiceServer.
func NewBlobMetadataServiceServer(blobMetadataService *services.BlobMetadataService) (*BlobMetadataServiceServer, error) {
	return &BlobMetadataServiceServer{
		blobMetadataService: blobMetadataService,
	}, nil
}

// ListMetadata fetches metadata of blobs optionally considering query parameters
func (s *BlobMetadataServiceServer) ListMetadata(ctx context.Context, req *pb.BlobMetaQuery, stream pb.BlobMetadataService_ListMetadataServer) error {
	query := blobs.NewBlobMetaQuery()
	if len(req.GetName()) > 0 {
		query.Name = req.Name
	}
	if req.GetSize() > 0 {
		query.Size = req.Size
	}
	if len(req.GetType()) > 0 {
		query.Type = req.Type
	}
	if req.GetDateTimeCreated() != nil {
		query.DateTimeCreated = req.DateTimeCreated.AsTime()
	}
	if req.GetLimit() > -1 {
		query.Limit = int(req.GetLimit())
	}
	if req.GetOffset() > -1 {
		query.Offset = int(req.GetOffset())
	}

	blobMetas, err := s.blobMetadataService.List(query)
	if err != nil {
		return fmt.Errorf("failed to list metadata: %v", err)
	}

	for _, blobMeta := range blobMetas {
		blobMetaResponse := &pb.BlobMetaResponse{
			Id:              blobMeta.ID,
			DateTimeCreated: timestamppb.New(blobMeta.DateTimeCreated),
			UserId:          blobMeta.UserID,
			Name:            blobMeta.Name,
			Size:            blobMeta.Size,
			Type:            blobMeta.Type,
			EncryptionKeyId: "",
			SignKeyId:       "",
		}

		if blobMeta.EncryptionKeyID != nil {
			blobMetaResponse.EncryptionKeyId = *blobMeta.EncryptionKeyID
		}
		if blobMeta.SignKeyID != nil {
			blobMetaResponse.SignKeyId = *blobMeta.SignKeyID
		}

		// Send the metadata response to the client
		if err := stream.Send(blobMetaResponse); err != nil {
			return fmt.Errorf("failed to send metadata response: %v", err)
		}
	}

	return nil
}

// GetMetadataById handles the GET request to fetch metadata of a blob by its ID
func (s *BlobMetadataServiceServer) GetMetadataById(ctx context.Context, req *pb.IdRequest) (*pb.BlobMetaResponse, error) {
	blobMeta, err := s.blobMetadataService.GetByID(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata by ID: %v", err)
	}

	blobMetaResponse := &pb.BlobMetaResponse{
		Id:              blobMeta.ID,
		DateTimeCreated: timestamppb.New(blobMeta.DateTimeCreated),
		UserId:          blobMeta.UserID,
		Name:            blobMeta.Name,
		Size:            blobMeta.Size,
		Type:            blobMeta.Type,
		EncryptionKeyId: "",
		SignKeyId:       "",
	}

	if blobMeta.EncryptionKeyID != nil {
		blobMetaResponse.EncryptionKeyId = *blobMeta.EncryptionKeyID
	}
	if blobMeta.SignKeyID != nil {
		blobMetaResponse.SignKeyId = *blobMeta.SignKeyID
	}
	return blobMetaResponse, nil
}

// DeleteById handles the DELETE request to delete a blob by its ID
func (s *BlobMetadataServiceServer) DeleteById(ctx context.Context, req *pb.IdRequest) (*pb.InfoResponse, error) {
	err := s.blobMetadataService.DeleteByID(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete blob: %v", err)
	}

	return &pb.InfoResponse{
		Message: fmt.Sprintf("blob with id %s deleted successfully", req.Id),
	}, nil
}

// NewCryptoKeyUploadServiceServer creates a new instance of CryptoKeyUploadServiceServer.
func NewCryptoKeyUploadServiceServer(cryptoKeyUploadService *services.CryptoKeyUploadService) (*CryptoKeyUploadServiceServer, error) {
	return &CryptoKeyUploadServiceServer{
		cryptoKeyUploadService: cryptoKeyUploadService,
	}, nil
}

// UploadKeys generates and uploads cryptographic keys
func (s *CryptoKeyUploadServiceServer) Upload(req *pb.UploadKeyRequest, stream pb.CryptoKeyUploadService_UploadServer) error {
	cryptoKeyMetas, err := s.cryptoKeyUploadService.Upload(req.GetUserId(), req.GetAlgorithm(), uint(req.GetKeySize()))
	if err != nil {
		return fmt.Errorf("failed to generate and upload crypto keys: %v", err)
	}

	for _, cryptoKeyMeta := range cryptoKeyMetas {
		cryptoKeyMetaResponse := &pb.CryptoKeyMetaResponse{
			Id:              cryptoKeyMeta.ID,
			DateTimeCreated: timestamppb.New(cryptoKeyMeta.DateTimeCreated),
			UserId:          cryptoKeyMeta.UserID,
			Algorithm:       cryptoKeyMeta.Algorithm,
			KeySize:         uint32(cryptoKeyMeta.KeySize),
			Type:            cryptoKeyMeta.Type,
		}

		// Send the metadata response to the client
		if err := stream.Send(cryptoKeyMetaResponse); err != nil {
			return fmt.Errorf("failed to send metadata response: %v", err)
		}
	}

	return nil
}

// NewCryptoKeyDownloadServiceServer creates a new instance of CryptoKeyDownloadServiceServer.
func NewCryptoKeyDownloadServiceServer(cryptoKeyDownloadService *services.CryptoKeyDownloadService) (*CryptoKeyDownloadServiceServer, error) {
	return &CryptoKeyDownloadServiceServer{
		cryptoKeyDownloadService: cryptoKeyDownloadService,
	}, nil
}

// DownloadById downloads a key by its ID
func (s *CryptoKeyDownloadServiceServer) DownloadById(req *pb.KeyDownloadRequest, stream pb.CryptoKeyDownloadService_DownloadByIdServer) error {
	bytes, err := s.cryptoKeyDownloadService.Download(req.GetId())
	if err != nil {
		return fmt.Errorf("failed to download crypto key: %v", err)
	}

	// If no error, stream the blob content back in chunks
	chunkSize := 1024 * 1024 // 1MB chunk size, adjust as needed
	for i := 0; i < len(bytes); i += chunkSize {
		end := i + chunkSize
		if end > len(bytes) {
			end = len(bytes)
		}

		// Create the chunk of data to send
		chunk := &pb.BlobContent{
			Content: bytes[i:end],
		}

		// Send the chunk
		if err := stream.Send(chunk); err != nil {
			return fmt.Errorf("failed to send chunk: %v", err)
		}
	}

	return nil
}

// NewCryptoKeyMetadataServiceServer creates a new instance of CryptoKeyMetadataServiceServer.
func NewCryptoKeyMetadataServiceServer(cryptoKeyMetadataService *services.CryptoKeyMetadataService) (*CryptoKeyMetadataServiceServer, error) {
	return &CryptoKeyMetadataServiceServer{
		cryptoKeyMetadataService: cryptoKeyMetadataService,
	}, nil
}

// ListMetadata lists cryptographic key metadata with optional query parameters
func (s *CryptoKeyMetadataServiceServer) ListMetadata(req *pb.KeyMetadataQuery, stream pb.CryptoKeyMetadataService_ListMetadataServer) error {
	query := keys.NewCryptoKeyQuery()
	if req.Algorithm != "" {
		query.Algorithm = req.Algorithm
	}
	if req.Type != "" {
		query.Type = req.Type
	}
	if req.DateTimeCreated != nil {
		query.DateTimeCreated = req.DateTimeCreated.AsTime()
	}
	query.Limit = int(req.Limit)
	query.Offset = int(req.Offset)

	cryptoKeyMetas, err := s.cryptoKeyMetadataService.List(query)
	if err != nil {
		return fmt.Errorf("failed to list crypto key metadata: %v", err)
	}

	for _, cryptoKeyMeta := range cryptoKeyMetas {
		cryptoKeyMetaResponse := &pb.CryptoKeyMetaResponse{
			Id:              cryptoKeyMeta.ID,
			KeyPairId:       cryptoKeyMeta.KeyPairID,
			Algorithm:       cryptoKeyMeta.Algorithm,
			KeySize:         uint32(cryptoKeyMeta.KeySize),
			Type:            cryptoKeyMeta.Type,
			DateTimeCreated: timestamppb.New(cryptoKeyMeta.DateTimeCreated),
			UserId:          cryptoKeyMeta.UserID,
		}

		// Send the metadata response to the client
		if err := stream.Send(cryptoKeyMetaResponse); err != nil {
			return fmt.Errorf("failed to send metadata response: %v", err)
		}
	}

	return nil
}

// GetMetadataById handles the GET request to retrieve metadata of a key by its ID
func (s *CryptoKeyMetadataServiceServer) GetMetadataById(ctx context.Context, req *pb.IdRequest) (*pb.CryptoKeyMetaResponse, error) {
	cryptoKeyMeta, err := s.cryptoKeyMetadataService.GetByID(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get crypto key metadata by ID: %v", err)
	}

	return &pb.CryptoKeyMetaResponse{
		Id:              cryptoKeyMeta.ID,
		KeyPairId:       cryptoKeyMeta.KeyPairID,
		Algorithm:       cryptoKeyMeta.Algorithm,
		KeySize:         uint32(cryptoKeyMeta.KeySize),
		Type:            cryptoKeyMeta.Type,
		DateTimeCreated: timestamppb.New(cryptoKeyMeta.DateTimeCreated),
		UserId:          cryptoKeyMeta.UserID,
	}, nil
}

// DeleteById deletes a key by its ID
func (s *CryptoKeyMetadataServiceServer) DeleteCryptoKeyById(ctx context.Context, req *pb.IdRequest) (*pb.InfoResponse, error) {
	err := s.cryptoKeyMetadataService.DeleteByID(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete crypto key: %v", err)
	}

	return &pb.InfoResponse{
		Message: fmt.Sprintf("crypto key with id %s deleted successfully", req.Id),
	}, nil
}
