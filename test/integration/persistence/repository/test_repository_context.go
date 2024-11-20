package repository

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

// TestRepositoryContext struct to hold DB and repositories for each test
type TestRepositoryContext struct {
	DB            *gorm.DB
	BlobRepo      *repository.GormBlobRepository
	CryptoKeyRepo *repository.GormCryptoKeyRepository
}

// Setup function to initialize the test DB and repositories
func setupTestDB(t *testing.T) *TestRepositoryContext {
	var err error
	var db *gorm.DB

	// Check for the DB type to use (SQLite in-memory or PostgreSQL)
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		// Default to SQLite in-memory if DB_TYPE is not set
		dbType = "sqlite"
	}

	switch dbType {
	case "postgres":
		// Setup PostgreSQL connection
		dsn := "user=postgres password=postgres host=localhost port=5432 sslmode=disable"
		if dsn == "" {
			t.Fatalf("POSTGRES_DSN environment variable is not set")
		}

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}

		// Check if the `blobs` database exists
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

		// Now, open a connection to the `blobs` database
		dsn = "user=postgres password=postgres host=localhost port=5432 dbname=blobs sslmode=disable"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to PostgreSQL database 'blobs': %v", err)
		}

	case "sqlite":
		// Setup SQLite in-memory connection
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

	// Return the test context that holds the DB and repositories
	return &TestRepositoryContext{
		DB:            db,
		BlobRepo:      blobRepo,
		CryptoKeyRepo: cryptoKeyRepo,
	}
}

// Teardown function to clean up after tests (optional, for DB cleanup)
func teardownTestDB(t *testing.T, ctx *TestRepositoryContext) {
	sqlDB, err := ctx.DB.DB()
	if err != nil {
		t.Fatalf("Failed to get DB connection: %v", err)
	}
	sqlDB.Close()
}

// TestMain setup and teardown for the entire test suite
func TestMain(m *testing.M) {
	// Set up test context
	ctx := setupTestDB(nil)
	// Run tests
	code := m.Run()
	// Clean up after tests
	teardownTestDB(nil, ctx)
	// Exit with the test result code
	if code != 0 {
		fmt.Println("Tests failed.")
	}
}
