package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	v1 "crypto_vault_service/internal/api/grpc/v1"
	"crypto_vault_service/internal/app/services"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/persistence/repository"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	path := "../../configs/grpc-app.yaml"

	config, err := settings.InitializeGrpcConfig(path)
	if err != nil {
		fmt.Printf("failed to initialize config: %v", err)
	}

	logger, err := logger.GetLogger(&config.Logger)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	// Database connection and migrations
	var db *gorm.DB
	switch config.Database.Type {
	case "postgres":
		dsn := config.Database.DSN
		if dsn == "" {
			log.Fatalf("POSTGRES_DSN environment variable is not set")
		}

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}

		// Check if database exists
		var dbExists bool
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get raw DB connection: %v", err)
		}
		query := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname='%s'", config.Database.Name)
		err = sqlDB.QueryRow(query).Scan(&dbExists)

		if err != nil && err.Error() != "sql: no rows in result set" {
			log.Fatalf("Failed to check if database '%s' exists: %v", config.Database.Name, err)
		}

		if !dbExists {
			_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", config.Database.Name))
			if err != nil {
				log.Fatalf("Failed to create database '%s': %v", config.Database.Name, err)
			}
			fmt.Printf("Database '%s' created successfully.\n", config.Database.Name)
		} else {
			fmt.Printf("Database '%s' already exists. Skipping creation.\n", config.Database.Name)
		}

		// Reconnect to the newly created database
		dsn = fmt.Sprintf(config.Database.DSN+" dbname=%s", config.Database.Name)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL database '%s': %v", config.Database.Name, err)
		}
	default:
		log.Fatalf("Unsupported database type: %s", config.Database.Type)
	}

	// Migrate the schema for Blob and CryptoKey
	err = db.AutoMigrate(&blobs.BlobMeta{}, &keys.CryptoKeyMeta{})
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	// setup infrastructure instances
	blobRepo, err := repository.NewGormBlobRepository(db, logger)
	if err != nil {
		log.Fatalf("Error creating blob repository instance: %v", err)
	}
	cryptoKeyRepo, err := repository.NewGormCryptoKeyRepository(db, logger)
	if err != nil {
		log.Fatalf("Error creating crypto key repository instance: %v", err)
	}

	blobConnector, err := connector.NewAzureBlobConnector(&config.BlobConnector, logger)
	if err != nil {
		log.Fatalf("%v", err)
	}

	vaultConnector, err := connector.NewAzureVaultConnector(&config.KeyConnector, logger)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Initialize services
	blobUploadService, err := services.NewBlobUploadService(blobConnector, blobRepo, vaultConnector, cryptoKeyRepo, logger)
	if err != nil {
		log.Fatalf("%v", err)
	}
	blobDownloadService, err := services.NewBlobDownloadService(blobConnector, blobRepo, vaultConnector, cryptoKeyRepo, logger)
	if err != nil {
		log.Fatalf("%v", err)
	}
	blobMetadataService, err := services.NewBlobMetadataService(blobRepo, blobConnector, logger)
	if err != nil {
		log.Fatalf("%v", err)
	}
	cryptoKeyUploadService, err := services.NewCryptoKeyUploadService(vaultConnector, cryptoKeyRepo, logger)
	if err != nil {
		log.Fatalf("%v", err)
	}
	cryptoKeyDownloadService, err := services.NewCryptoKeyDownloadService(vaultConnector, cryptoKeyRepo, logger)
	if err != nil {
		log.Fatalf("%v", err)
	}
	cryptoKeyMetadataService, err := services.NewCryptoKeyMetadataService(vaultConnector, cryptoKeyRepo, logger)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Create gRPC server and register the gRPC services
	blobUploadServer, err := v1.NewBlobUploadServer(blobUploadService)
	if err != nil {
		log.Fatalf("failed to create blob upload server: %v", err)
	}

	blobDownloadServer, err := v1.NewBlobDownloadServer(blobDownloadService)
	if err != nil {
		log.Fatalf("failed to create blob download server: %v", err)
	}

	blobMetadataServer, err := v1.NewBlobMetadataServer(blobMetadataService)
	if err != nil {
		log.Fatalf("failed to create blob metadata server: %v", err)
	}

	cryptoKeyUploadServer, err := v1.NewCryptoKeyUploadServer(cryptoKeyUploadService)
	if err != nil {
		log.Fatalf("failed to create crypto key upload server: %v", err)
	}

	cryptoKeyDownloadServer, err := v1.NewCryptoKeyDownloadServer(cryptoKeyDownloadService)
	if err != nil {
		log.Fatalf("failed to create crypto key download server: %v", err)
	}

	cryptoKeyMetadataServer, err := v1.NewCryptoKeyMetadataServer(cryptoKeyMetadataService)
	if err != nil {
		log.Fatalf("failed to create crypto key metadata server: %v", err)
	}

	grpcServer := grpc.NewServer()

	v1.RegisterBlobUploadServer(grpcServer, blobUploadServer)
	v1.RegisterBlobDownloadServer(grpcServer, blobDownloadServer)
	v1.RegisterBlobMetadataServer(grpcServer, blobMetadataServer)
	v1.RegisterCryptoKeyUploadServer(grpcServer, cryptoKeyUploadServer)
	v1.RegisterCryptoKeyDownloadServer(grpcServer, cryptoKeyDownloadServer)
	v1.RegisterCryptoKeyMetadataServer(grpcServer, cryptoKeyMetadataServer)

	// Enable reflection in order to list services via `grpcurl -plaintext localhost:50051 list`
	reflection.Register(grpcServer)

	// Set up listener for gRPC server
	lis, err := net.Listen("tcp", ":"+config.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("gRPC server started at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Set up gRPC-Gateway mux
	gwmux := runtime.NewServeMux()

	gatewayTarget := "0.0.0.0:" + config.Port
	// Create a client connection to the gRPC server
	conn, err := grpc.NewClient(gatewayTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	// Register all services for the gRPC-Gateway mux

	creds := insecure.NewCredentials()
	// Multipart file uploads are not supported with grpc-gateway. For more details
	// checkout: https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/binary_file_uploads/. As a result, subsequent code can be commented.
	// err = v1.RegisterBlobUploadGateway(context.Background(), gatewayTarget, gwmux, conn, creds)
	// if err != nil {
	// 	log.Fatalf("Failed to register blob upload gateway: %v", err)
	// }
	err = v1.RegisterBlobDownloadGateway(context.Background(), gatewayTarget, gwmux, conn, creds)
	if err != nil {
		log.Fatalf("Failed to register blob download gateway: %v", err)
	}
	err = v1.RegisterBlobMetadataGateway(context.Background(), gatewayTarget, gwmux, conn, creds)
	if err != nil {
		log.Fatalf("Failed to register blob metadata gateway: %v", err)
	}
	err = v1.RegisterCryptoKeyUploadGateway(context.Background(), gatewayTarget, gwmux, conn, creds)
	if err != nil {
		log.Fatalf("Failed to register crypto key upload gateway: %v", err)
	}
	err = v1.RegisterCryptoKeyDownloadGateway(context.Background(), gatewayTarget, gwmux, conn, creds)
	if err != nil {
		log.Fatalf("Failed to register crypto key download gateway: %v", err)
	}
	err = v1.RegisterCryptoKeyMetadataGateway(context.Background(), gatewayTarget, gwmux, conn, creds)
	if err != nil {
		log.Fatalf("Failed to register crypto key metadata gateway: %v", err)
	}

	gatewayPort := config.GatewayPort
	// Set up the HTTP server to serve the Gateway
	gwServer := &http.Server{
		Addr:    ":" + gatewayPort,
		Handler: gwmux,
	}

	log.Printf("gRPC-Gateway server started at http://0.0.0.0:%v", gatewayPort)
	log.Fatalln(gwServer.ListenAndServe())
}
