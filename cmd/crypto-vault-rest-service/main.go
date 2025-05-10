package main

import (
	"context"
	v1 "crypto_vault_service/internal/api/rest/v1"
	"crypto_vault_service/internal/app/services"
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/connector"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/persistence/repository"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"crypto_vault_service/cmd/crypto-vault-rest-service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// parseArgs parses the --config flag and returns its value.
func parseArgs() (string, error) {
	configPath := flag.String("config", "../../configs/rest-app.yaml", "Path to the config YAML file")
	flag.Parse()

	if *configPath == "" {
		return "", fmt.Errorf("missing required --config argument")
	}
	return *configPath, nil
}

// @title CryptoVault Service API
// @version v1
// @description Service capable of managing cryptographic keys and securing data at rest (metadata, BLOB)
// @termsOfService TBD
// @contact.name MGTheTrain
// @contact.url TBD
// @contact.email TBD
// @license.name LGPL-2.1 license
// @license.url https://github.com/MGTheTrain/crypto-vault-service/blob/main/LICENSE
// @BasePath /api/v1/cvs
// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information
// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationUrl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information
// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information
// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationUrl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information
func main() {
	r := gin.Default()

	path, err := parseArgs()
	if err != nil {
		fmt.Printf("Warning: Could not parse arguments: %v", err)
	}

	config, err := settings.InitializeRestConfig(path)
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	logger, err := logger.GetLogger(&config.Logger)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
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
		query := "SELECT 1 FROM pg_database WHERE datname = $1"
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

	blobRepo, err := repository.NewGormBlobRepository(db, logger)
	if err != nil {
		log.Fatalf("Error creating blob repository instance: %v", err)
	}
	cryptoKeyRepo, err := repository.NewGormCryptoKeyRepository(db, logger)
	if err != nil {
		log.Fatalf("Error creating crypto key repository instance: %v", err)
	}

	ctx := context.Background()
	var blobConnector connector.BlobConnector
	if config.BlobConnector.CloudProvider == "azure" {
		blobConnector, err = connector.NewAzureBlobConnector(ctx, &config.BlobConnector, logger)
		if err != nil {
			log.Fatalf("%v", err)
			return
		}
	}

	var vaultConnector connector.VaultConnector
	if config.BlobConnector.CloudProvider == "azure" {
		vaultConnector, err = connector.NewAzureVaultConnector(ctx, &config.KeyConnector, logger)
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

	docs.SwaggerInfo.Version = v1.Version
	docs.SwaggerInfo.BasePath = v1.BasePath
	swaggerRoute := fmt.Sprintf("/api/" + v1.Version + "/cvs/swagger/*any")
	r.GET(swaggerRoute, ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":" + config.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
