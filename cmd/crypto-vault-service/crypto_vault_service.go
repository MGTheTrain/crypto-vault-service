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

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	path := "../../configs/app.yaml"

	config, err := settings.Initialize(path)
	if err != nil {
		fmt.Printf("failed to initialize config: %v", err)
	}

	logger, err := logger.GetLogger(&config.Logger)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

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

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get raw DB connection: %v", err)
		}

		// Check if the database exists
		var dbExists bool
		query := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname='%s'", config.Database.Name)
		err = sqlDB.QueryRow(query).Scan(&dbExists)

		if err != nil && err.Error() != "sql: no rows in result set" {
			log.Fatalf("Failed to check if database '%s' exists: %v", config.Database.Name, err)
		}

		if !dbExists {
			// If the database doesn't exist, create it
			_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", config.Database.Name))
			if err != nil {
				log.Fatalf("Failed to create database '%s': %v", config.Database.Name, err)
			}
			fmt.Printf("Database '%s' created successfully.\n", config.Database.Name)
		} else {
			fmt.Printf("Database '%s' already exists. Skipping creation.\n", config.Database.Name)
		}

		dsn = fmt.Sprintf("user=postgres password=postgres host=localhost port=5432 dbname=%s sslmode=disable", config.Database.Name)
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

	blobRepo, err := repository.NewGormBlobRepository(db, logger)
	if err != nil {
		log.Fatalf("Error creating blob repository instance: %v", err)
	}
	cryptoKeyRepo, err := repository.NewGormCryptoKeyRepository(db, logger)
	if err != nil {
		log.Fatalf("Error creating crypto key repository instance: %v", err)
	}

	var blobConnector connector.BlobConnector
	if config.BlobConnector.CloudProvider == "azure" {
		blobConnector, err = connector.NewAzureBlobConnector(&config.BlobConnector, logger)
		if err != nil {
			log.Fatalf("%v", err)
			return
		}
	}

	var vaultConnector connector.VaultConnector
	if config.BlobConnector.CloudProvider == "azure" {
		vaultConnector, err = connector.NewAzureVaultConnector(&config.KeyConnector, logger)
		if err != nil {
			log.Fatalf("%v", err)
			return
		}
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

	if err := r.Run(":" + config.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
