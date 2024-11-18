package repository

import (
	"crypto_vault_service/internal/domain/model"
	"crypto_vault_service/internal/persistence/repository"
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestContext struct to hold DB and repositories for each test
type TestContext struct {
	DB            *gorm.DB
	BlobRepo      *repository.BlobRepositoryImpl
	CryptoKeyRepo *repository.CryptographicKeyRepositoryImpl
}

// Setup function to initialize the test DB and repositories
func setupTestDB(t *testing.T) *TestContext {
	var err error
	// Set up an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to setup DB: %v", err)
	}

	// Migrate the schema for Blob and CryptographicKey
	err = db.AutoMigrate(&model.Blob{}, &model.CryptographicKey{})
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	// Initialize the repositories with the DB instance
	blobRepo := &repository.BlobRepositoryImpl{DB: db}
	cryptoKeyRepo := &repository.CryptographicKeyRepositoryImpl{DB: db}

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
