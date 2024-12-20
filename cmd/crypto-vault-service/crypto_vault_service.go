package main

import (
	v1 "crypto_vault_service/internal/api/v1"
	"crypto_vault_service/internal/app/services"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/persistence/repository"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	// TBD: consider env vars or load config yml file and
	// utilize settings objects during costruction of other objects

	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	var db *gorm.DB

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}

	switch dbType {
	case "postgres":
		dsn := "user=postgres password=postgres host=localhost port=5432 sslmode=disable"
		if dsn == "" {
			log.Fatalf("POSTGRES_DSN environment variable is not set")
		}

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}

		uniqueDBName := "blobs_" + uuid.New().String()

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get raw DB connection: %v", err)
		}

		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", uniqueDBName))
		if err != nil {
			log.Fatalf("Failed to create database '%s': %v", uniqueDBName, err)
		}
		fmt.Printf("Database '%s' created successfully.\n", uniqueDBName)

		dsn = fmt.Sprintf("user=postgres password=postgres host=localhost port=5432 dbname=%s sslmode=disable", uniqueDBName)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL database '%s': %v", uniqueDBName, err)
		}

	case "sqlite":
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to SQLite: %v", err)
		}

	default:
		log.Fatalf("Unsupported DB_TYPE value: %s", dbType)
	}

	// Migrate the schema for Blob and CryptoKey
	err = db.AutoMigrate(&blobs.BlobMeta{}, &keys.CryptoKeyMeta{})
	if err != nil {
		log.Fatalf("Failed to migrate schema: %v", err)
	}

	blobRepo, err := repository.NewGormBlobRepository(db, logger)
	if err != nil {
		log.Fatalf("Error creating blob repository instance: %v", err)
	}
	cryptoKeyRepo, err := repository.NewGormCryptoKeyRepository(db, logger)
	if err != nil {
		log.Fatalf("Error creating crypto key repository instance: %v", err)
	}

	blobConnectorSettings := &settings.BlobConnectorSettings{
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}

	blobConnector, err := connector.NewAzureBlobConnector(blobConnectorSettings, logger)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	keyConnectorSettings := &settings.KeyConnectorSettings{
		ConnectionString: "DefaultEndpointsProtocol=http;AccountName=devstoreaccount1;AccountKey=Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw==;BlobEndpoint=http://127.0.0.1:10000/devstoreaccount1;",
		ContainerName:    "testblobs",
	}
	vaultConnector, err := connector.NewAzureVaultConnector(keyConnectorSettings, logger)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	blobUploadService, err := services.NewBlobUploadService(blobConnector, blobRepo, vaultConnector, cryptoKeyRepo, logger)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	blobDownloadService, err := services.NewBlobDownloadService(blobConnector, blobRepo, vaultConnector, cryptoKeyRepo, logger)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	blobMetadataService, err := services.NewBlobMetadataService(blobRepo, blobConnector, logger)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	cryptoKeyUploadService, err := services.NewCryptoKeyUploadService(vaultConnector, cryptoKeyRepo, logger)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	cryptoKeyDownloadService, err := services.NewCryptoKeyDownloadService(vaultConnector, cryptoKeyRepo, logger)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	cryptoKeyMetadataService, err := services.NewCryptoKeyMetadataService(vaultConnector, cryptoKeyRepo, logger)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	v1.SetupRoutes(r, blobUploadService, blobDownloadService, blobMetadataService, cryptoKeyUploadService, cryptoKeyDownloadService, cryptoKeyMetadataService)

	// r.Use(v1.AuthMiddleware())

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
