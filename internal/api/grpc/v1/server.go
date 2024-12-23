package v1

import (
	"context"
	"crypto_vault_service/internal/app/services"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/utils"
	"fmt"
	"log"

	pb "proto"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BlobUploadServer struct {
	pb.UnimplementedBlobUploadServer
	blobUploadService *services.BlobUploadService
}

type BlobDownloadServer struct {
	pb.UnimplementedBlobDownloadServer
	blobDownloadService *services.BlobDownloadService
}

type BlobMetadataServer struct {
	pb.UnimplementedBlobMetadataServer
	blobMetadataService *services.BlobMetadataService
}

type CryptoKeyUploadServer struct {
	pb.UnimplementedCryptoKeyUploadServer
	cryptoKeyUploadService *services.CryptoKeyUploadService
}

type CryptoKeyDownloadServer struct {
	pb.UnimplementedCryptoKeyDownloadServer
	cryptoKeyDownloadService *services.CryptoKeyDownloadService
}

type CryptoKeyMetadataServer struct {
	pb.UnimplementedCryptoKeyMetadataServer
	cryptoKeyMetadataService *services.CryptoKeyMetadataService
}

// NewBlobUploadServer creates a new instance of BlobUploadServer.
func NewBlobUploadServer(blobUploadService *services.BlobUploadService) (*BlobUploadServer, error) {
	return &BlobUploadServer{
		blobUploadService: blobUploadService,
	}, nil
}

// Upload uploads blobs with optional encryption/signing
func (s BlobUploadServer) Upload(req *pb.BlobUploadRequest, stream pb.BlobUpload_UploadServer) error {
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

	userId := uuid.New().String() // TBD: extract user id from JWT
	form, err := utils.CreateMultipleFilesForm(fileContent, fileNames)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	blobMetas, err := s.blobUploadService.Upload(form, userId, encryptionKeyId, signKeyId)
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

// NewBlobDownloadServer creates a new instance of BlobDownloadServer.
func NewBlobDownloadServer(blobDownloadService *services.BlobDownloadService) (*BlobDownloadServer, error) {
	return &BlobDownloadServer{
		blobDownloadService: blobDownloadService,
	}, nil
}

// DownloadById downloads a blob by its ID
func (s *BlobDownloadServer) DownloadById(req *pb.BlobDownloadRequest, stream pb.BlobDownload_DownloadByIdServer) error {
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

// NewBlobMetadataServer creates a new instance of BlobMetadataServer.
func NewBlobMetadataServer(blobMetadataService *services.BlobMetadataService) (*BlobMetadataServer, error) {
	return &BlobMetadataServer{
		blobMetadataService: blobMetadataService,
	}, nil
}

// ListMetadata fetches metadata of blobs optionally considering query parameters
func (s *BlobMetadataServer) ListMetadata(req *pb.BlobMetaQuery, stream pb.BlobMetadata_ListMetadataServer) error {
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
	if req.GetLimit() > 0 {
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
func (s *BlobMetadataServer) GetMetadataById(ctx context.Context, req *pb.IdRequest) (*pb.BlobMetaResponse, error) {
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
func (s *BlobMetadataServer) DeleteById(ctx context.Context, req *pb.IdRequest) (*pb.InfoResponse, error) {
	err := s.blobMetadataService.DeleteByID(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to delete blob: %v", err)
	}

	return &pb.InfoResponse{
		Message: fmt.Sprintf("blob with id %s deleted successfully", req.Id),
	}, nil
}

// NewCryptoKeyUploadServer creates a new instance of CryptoKeyUploadServer.
func NewCryptoKeyUploadServer(cryptoKeyUploadService *services.CryptoKeyUploadService) (*CryptoKeyUploadServer, error) {
	return &CryptoKeyUploadServer{
		cryptoKeyUploadService: cryptoKeyUploadService,
	}, nil
}

// UploadKeys generates and uploads cryptographic keys
func (s *CryptoKeyUploadServer) Upload(req *pb.UploadKeyRequest, stream pb.CryptoKeyUpload_UploadServer) error {
	userId := uuid.New().String() // TBD: extract user id from JWT
	cryptoKeyMetas, err := s.cryptoKeyUploadService.Upload(userId, req.Algorithm, uint(req.KeySize))
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

// NewCryptoKeyDownloadServer creates a new instance of CryptoKeyDownloadServer.
func NewCryptoKeyDownloadServer(cryptoKeyDownloadService *services.CryptoKeyDownloadService) (*CryptoKeyDownloadServer, error) {
	return &CryptoKeyDownloadServer{
		cryptoKeyDownloadService: cryptoKeyDownloadService,
	}, nil
}

// DownloadById downloads a key by its ID
func (s *CryptoKeyDownloadServer) DownloadById(req *pb.KeyDownloadRequest, stream pb.CryptoKeyDownload_DownloadByIdServer) error {
	bytes, err := s.cryptoKeyDownloadService.Download(req.Id)
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
		chunk := &pb.KeyContent{
			Content: bytes[i:end],
		}

		// Send the chunk
		if err := stream.Send(chunk); err != nil {
			return fmt.Errorf("failed to send chunk: %v", err)
		}
	}

	return nil
}

// NewCryptoKeyMetadataServer creates a new instance of CryptoKeyMetadataServer.
func NewCryptoKeyMetadataServer(cryptoKeyMetadataService *services.CryptoKeyMetadataService) (*CryptoKeyMetadataServer, error) {
	return &CryptoKeyMetadataServer{
		cryptoKeyMetadataService: cryptoKeyMetadataService,
	}, nil
}

// ListMetadata lists cryptographic key metadata with optional query parameters
func (s *CryptoKeyMetadataServer) ListMetadata(req *pb.KeyMetadataQuery, stream pb.CryptoKeyMetadata_ListMetadataServer) error {
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
	if req.Limit > 0 {
		query.Limit = int(req.Limit)
	}
	if req.Offset > -1 {
		query.Offset = int(req.Offset)
	}

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
func (s *CryptoKeyMetadataServer) GetMetadataById(ctx context.Context, req *pb.IdRequest) (*pb.CryptoKeyMetaResponse, error) {
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
func (s *CryptoKeyMetadataServer) DeleteCryptoKeyById(ctx context.Context, req *pb.IdRequest) (*pb.InfoResponse, error) {
	err := s.cryptoKeyMetadataService.DeleteByID(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete crypto key: %v", err)
	}

	return &pb.InfoResponse{
		Message: fmt.Sprintf("crypto key with id %s deleted successfully", req.Id),
	}, nil
}

// Register the gRPC handlers for each service

func RegisterBlobUploadServer(server *grpc.Server, blobUploadServer *BlobUploadServer) {
	pb.RegisterBlobUploadServer(server, blobUploadServer)
}

func RegisterBlobDownloadServer(server *grpc.Server, blobDownloadServer *BlobDownloadServer) {
	pb.RegisterBlobDownloadServer(server, blobDownloadServer)
}

func RegisterBlobMetadataServer(server *grpc.Server, blobMetadataServer *BlobMetadataServer) {
	pb.RegisterBlobMetadataServer(server, blobMetadataServer)
}

func RegisterCryptoKeyUploadServer(server *grpc.Server, cryptoKeyUploadServer *CryptoKeyUploadServer) {
	pb.RegisterCryptoKeyUploadServer(server, cryptoKeyUploadServer)
}

func RegisterCryptoKeyDownloadServer(server *grpc.Server, cryptoKeyDownloadServer *CryptoKeyDownloadServer) {
	pb.RegisterCryptoKeyDownloadServer(server, cryptoKeyDownloadServer)
}

func RegisterCryptoKeyMetadataServer(server *grpc.Server, cryptoKeyMetadataServer *CryptoKeyMetadataServer) {
	pb.RegisterCryptoKeyMetadataServer(server, cryptoKeyMetadataServer)
}

// Register the gRPC-Gateway handlers for each service

// Multipart file uploads are not supported with grpc-gateway. For more details,
// see: https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/binary_file_uploads/. As a result, subsequent code can be commented.
// func RegisterBlobUploadGateway(ctx context.Context, gatewayTarget string, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
// 	// Register the handler from the endpoint (this works with gRPC-Gateway)
// 	err := pb.RegisterBlobUploadHandlerFromEndpoint(ctx, gwmux, gatewayTarget, []grpc.DialOption{grpc.WithInsecure()})
// 	if err != nil {
// 		log.Fatalf("Failed to register blob upload gateway: %v", err)
// 		return err
// 	}
// 	return nil
// }

func RegisterBlobDownloadGateway(ctx context.Context, gatewayTarget string, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := pb.RegisterBlobDownloadHandlerFromEndpoint(ctx, gwmux, gatewayTarget, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("Failed to register blob download gateway: %v", err)
		return err
	}
	return nil
}

func RegisterBlobMetadataGateway(ctx context.Context, gatewayTarget string, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := pb.RegisterBlobMetadataHandlerFromEndpoint(ctx, gwmux, gatewayTarget, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("Failed to register blob metadata gateway: %v", err)
		return err
	}
	return nil
}

func RegisterCryptoKeyUploadGateway(ctx context.Context, gatewayTarget string, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := pb.RegisterCryptoKeyUploadHandlerFromEndpoint(ctx, gwmux, gatewayTarget, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("Failed to register crypto key upload gateway: %v", err)
		return err
	}
	return nil
}

func RegisterCryptoKeyDownloadGateway(ctx context.Context, gatewayTarget string, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := pb.RegisterCryptoKeyDownloadHandlerFromEndpoint(ctx, gwmux, gatewayTarget, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("Failed to register crypto key download gateway: %v", err)
		return err
	}
	return nil
}

func RegisterCryptoKeyMetadataGateway(ctx context.Context, gatewayTarget string, gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
	err := pb.RegisterCryptoKeyMetadataHandlerFromEndpoint(ctx, gwmux, gatewayTarget, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("Failed to register crypto key metadata gateway: %v", err)
		return err
	}
	return nil
}
