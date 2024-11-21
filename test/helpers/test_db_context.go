package helpers

import (
	"crypto_vault_service/internal/domain/blobs"
	"crypto_vault_service/internal/domain/keys"
	"crypto_vault_service/internal/persistence/repository"
	"fmt"
	"os"
	"testing"

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

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}

		// Ensure the `blobs` database exists
		sqlDB, err := db.DB()
		if err != nil {
			t.Fatalf("Failed to get raw DB connection: %v", err)
		}

		// Query to check if the `blobs` database exists
		var dbExists bool
		err = sqlDB.QueryRow("SELECT 1 FROM pg_database WHERE datname = 'blobs'").Scan(&dbExists)
		if err != nil {
			if err.Error() != "sql: no rows in result set" {
				t.Fatalf("Failed to check if database exists: %v", err)
			}
			// The database does not exist, create it
			_, err = sqlDB.Exec("CREATE DATABASE blobs")
			if err != nil {
				t.Fatalf("Failed to create database: %v", err)
			}
			fmt.Println("Database 'blobs' created successfully.")
		}

		// Open the connection to `blobs` database
		dsn = "user=postgres password=postgres host=localhost port=5432 dbname=blobs sslmode=disable"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to PostgreSQL database 'blobs': %v", err)
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
	blobRepo := &repository.GormBlobRepository{DB: db}
	cryptoKeyRepo := &repository.GormCryptoKeyRepository{DB: db}

	return &TestDBContext{
		DB:            db,
		BlobRepo:      blobRepo,
		CryptoKeyRepo: cryptoKeyRepo,
	}
}

// TeardownTestDB closes the DB connection after the test
func TeardownTestDB(t *testing.T, ctx *TestDBContext) {
	sqlDB, err := ctx.DB.DB()
	if err != nil {
		t.Fatalf("Failed to get DB connection: %v", err)
	}
	sqlDB.Close()
}
