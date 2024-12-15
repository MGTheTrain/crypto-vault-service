package helpers

import (
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/persistence/repository"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestDBContext struct {
	DB            *gorm.DB
	BlobRepo      *repository.GormBlobRepository
	CryptoKeyRepo *repository.GormCryptoKeyRepository
}

// SetupTestDB initializes the test database and repositories based on the DB_TYPE environment variable
func SetupTestDB(t *testing.T) *TestDBContext {
	var err error
	var db *gorm.DB

	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite" // Default to SQLite in-memory if DB_TYPE is not set
	}

	switch dbType {
	case "postgres":
		// PostgreSQL setup
		dsn := "user=postgres password=postgres host=localhost port=5432 sslmode=disable"
		if dsn == "" {
			t.Fatalf("POSTGRES_DSN environment variable is not set")
		}

		// Connect to PostgreSQL without specifying a database (so we can create one if necessary)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}

		// Generate a unique database name using UUID
		uniqueDBName := "blobs_" + uuid.New().String()

		// Ensure the unique `blobs` database exists, create if necessary
		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("Failed to get raw DB connection: %v", err)
		}

		// Create the new database
		_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", uniqueDBName))
		if err != nil {
			t.Fatalf("Failed to create database '%s': %v", uniqueDBName, err)
		}
		fmt.Printf("Database '%s' created successfully.\n", uniqueDBName)

		// Now that the unique `blobs` database is created, connect to it
		dsn = fmt.Sprintf("user=postgres password=postgres host=localhost port=5432 dbname=%s sslmode=disable", uniqueDBName)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to PostgreSQL database '%s': %v", uniqueDBName, err)
		}

	case "sqlite":
		db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to SQLite: %v", err)
		}

	default:
		t.Fatalf("Unsupported DB_TYPE value: %s", dbType)
	}

	// Migrate the schema for Blob and CryptoKey
	err = db.AutoMigrate(&blobs.BlobMeta{}, &keys.CryptoKeyMeta{})
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	// Initialize the repositories with the DB instance
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	factory := &logger.LoggerFactory{}

	logger, err := factory.NewLogger(loggerSettings)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	blobRepo, err := repository.NewGormBlobRepository(db, logger)
	if err != nil {
		log.Fatalf("Error creating blob repository instance: %v", err)
	}
	cryptoKeyRepo, err := repository.NewGormCryptoKeyRepository(db, logger)
	if err != nil {
		log.Fatalf("Error creating crypto key repository instance: %v", err)
	}

	return &TestDBContext{
		DB:            db,
		BlobRepo:      blobRepo,
		CryptoKeyRepo: cryptoKeyRepo,
	}
}

// TeardownTestDB closes the DB connection and cleans up the database after the test
func TeardownTestDB(t *testing.T, ctx *TestDBContext, dbType string) {
	sqlDB, err := ctx.DB.DB()
	if err != nil {
		t.Fatalf("Failed to get DB connection: %v", err)
	}

	// If using PostgreSQL, drop the unique database created during the test
	if dbType == "postgres" {
		// Get the database name from the DSN or context (you might store it during DB setup)
		databaseName := ctx.DB.Migrator().CurrentDatabase()

		// Close the current DB connection before dropping the database
		sqlDB.Close()

		// Connect again to PostgreSQL without specifying a database (connect to the default one)
		dsn := "user=postgres password=postgres host=localhost port=5432 sslmode=disable"
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to reconnect to PostgreSQL: %v", err)
		}

		// Drop the unique database
		tx := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", databaseName))
		if tx.Error != nil {
			t.Fatalf("Failed to drop database '%s': %v", databaseName, tx.Error)
		}
		fmt.Printf("Database '%s' dropped successfully.\n", databaseName)
	} else {
		// For SQLite, no need to drop the in-memory database, just close the connection
		sqlDB.Close()
	}
}
