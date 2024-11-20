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

// TestContext struct to hold DB and repositories for each test
type TestContext struct {
	DB            *gorm.DB
	BlobRepo      *repository.GormBlobRepository
	CryptoKeyRepo *repository.GormCryptoKeyRepository
}

// Setup function to initialize the test DB and repositories
func setupTestDB(t *testing.T) *TestContext {
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
		dsn := os.Getenv("POSTGRES_DSN") // Example: "user=username dbname=test sslmode=disable"
		if dsn == "" {
			t.Fatalf("POSTGRES_DSN environment variable is not set")
		}
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("Failed to connect to PostgreSQL: %v", err)
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
	return &TestContext{
		DB:            db,
		BlobRepo:      blobRepo,
		CryptoKeyRepo: cryptoKeyRepo,
	}
}

// Teardown function to clean up after tests (optional, for DB cleanup)
func teardownTestDB(t *testing.T, ctx *TestContext) {
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
