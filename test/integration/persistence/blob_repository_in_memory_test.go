package repository_test

import (
	"crypto_vault_service/internal/domain/model"
	"crypto_vault_service/internal/persistence/repository"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Global variable for DB to avoid reinitialization
var db *gorm.DB
var repo *repository.BlobRepositoryImpl

// Setup function to initialize the test DB and repository
func setupTestDB(t *testing.T) {
	var err error
	// Set up an in-memory SQLite database for testing
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to setup DB: %v", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&model.Blob{}, &model.CryptographicKey{}) // Assuming you need both migrations
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}

	// Initialize the repository with the DB instance
	repo = &repository.BlobRepositoryImpl{DB: db}
}

// Teardown function to clean up after tests (optional, for DB cleanup)
func teardownTestDB(t *testing.T) {
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get DB connection: %v", err)
	}
	sqlDB.Close()
}

// Test setup and teardown functions for the entire test suite
func TestMain(m *testing.M) {
	// Setup before tests
	setupTestDB(nil)
	// Run tests
	code := m.Run()
	// Teardown after tests
	teardownTestDB(nil)
	// Exit with the test result code
	if code != 0 {
		fmt.Println("Tests failed.")
	}
}

// TestBlobRepository_Create tests the Create method of BlobRepositoryImpl
func TestBlobRepository_Create(t *testing.T) {
	// Create a valid CryptographicKey object
	cryptographicKey := model.CryptographicKey{
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &model.Blob{
		BlobID:           uuid.New().String(), // Generate valid UUID
		BlobStoragePath:  "/path/to/blob",
		UploadTime:       time.Now(),
		UserID:           uuid.New().String(), // Generate valid UUID for UserID
		BlobName:         "test-blob",
		BlobSize:         1024,
		BlobType:         "text",
		IsEncrypted:      false,
		IsSigned:         true,
		CryptographicKey: cryptographicKey,       // Set the CryptographicKey
		KeyID:            cryptographicKey.KeyID, // Ensure KeyID is set
	}

	// Call the Create method
	err := repo.Create(blob)
	assert.NoError(t, err, "Create should not return an error")

	// Verify the blob is created and exists in DB
	var createdBlob model.Blob
	err = db.First(&createdBlob, "blob_id = ?", blob.BlobID).Error
	assert.NoError(t, err, "Failed to find created blob")
	assert.Equal(t, blob.BlobID, createdBlob.BlobID, "BlobID should match")
	assert.Equal(t, blob.BlobName, createdBlob.BlobName, "BlobName should match")
}

// TestBlobRepository_GetById tests the GetById method of BlobRepositoryImpl
func TestBlobRepository_GetById(t *testing.T) {
	// Create a valid CryptographicKey object
	cryptographicKey := model.CryptographicKey{
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &model.Blob{
		BlobID:           uuid.New().String(), // Generate valid UUID
		BlobStoragePath:  "/path/to/blob",
		UploadTime:       time.Now(),
		UserID:           cryptographicKey.UserID, // Link to valid UserID from CryptographicKey
		BlobName:         "test-blob",
		BlobSize:         1024,
		BlobType:         "text",
		IsEncrypted:      false,
		IsSigned:         true,
		CryptographicKey: cryptographicKey,       // Set the CryptographicKey
		KeyID:            cryptographicKey.KeyID, // Ensure KeyID is set
	}

	// Create the blob in DB
	err := repo.Create(blob)
	assert.NoError(t, err, "Create should not return an error")

	// Get the blob by ID
	fetchedBlob, err := repo.GetById(blob.BlobID)
	assert.NoError(t, err, "GetById should not return an error")
	assert.NotNil(t, fetchedBlob, "Fetched blob should not be nil")
	assert.Equal(t, blob.BlobID, fetchedBlob.BlobID, "BlobID should match")
}

// TestBlobRepository_UpdateById tests the UpdateById method of BlobRepositoryImpl
func TestBlobRepository_UpdateById(t *testing.T) {
	// Create a valid CryptographicKey object
	cryptographicKey := model.CryptographicKey{
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &model.Blob{
		BlobID:           uuid.New().String(), // Generate valid UUID
		BlobStoragePath:  "/path/to/blob",
		UploadTime:       time.Now(),
		UserID:           uuid.New().String(), // Generate valid UUID for UserID
		BlobName:         "test-blob",
		BlobSize:         1024,
		BlobType:         "text",
		IsEncrypted:      false,
		IsSigned:         true,
		CryptographicKey: cryptographicKey,       // Set the CryptographicKey
		KeyID:            cryptographicKey.KeyID, // Ensure KeyID is set
	}

	// Create the blob in DB
	err := repo.Create(blob)
	assert.NoError(t, err)

	// Update the blob's name
	blob.BlobName = "updated-blob-name"
	err = repo.UpdateById(blob)
	assert.NoError(t, err)

	// Verify the blob is updated
	var updatedBlob model.Blob
	err = db.First(&updatedBlob, "blob_id = ?", blob.BlobID).Error
	assert.NoError(t, err)
	assert.Equal(t, "updated-blob-name", updatedBlob.BlobName, "BlobName should be updated")
}

// TestBlobRepository_DeleteById tests the DeleteById method of BlobRepositoryImpl
func TestBlobRepository_DeleteById(t *testing.T) {
	// Create a valid CryptographicKey object
	cryptographicKey := model.CryptographicKey{
		KeyID:     uuid.New().String(),            // Generate valid UUID for KeyID
		KeyType:   "AES",                          // Example key type
		CreatedAt: time.Now(),                     // Valid CreatedAt time
		ExpiresAt: time.Now().Add(24 * time.Hour), // Valid ExpiresAt time
		UserID:    uuid.New().String(),            // Generate valid UUID for UserID
	}

	// Create a test Blob object with valid UUIDs and required fields
	blob := &model.Blob{
		BlobID:           uuid.New().String(), // Generate valid UUID
		BlobStoragePath:  "/path/to/blob",
		UploadTime:       time.Now(),
		UserID:           uuid.New().String(), // Generate valid UUID for UserID
		BlobName:         "test-blob",
		BlobSize:         1024,
		BlobType:         "text",
		IsEncrypted:      false,
		IsSigned:         true,
		CryptographicKey: cryptographicKey,       // Set the CryptographicKey
		KeyID:            cryptographicKey.KeyID, // Ensure KeyID is set
	}

	// Create the blob in DB
	err := repo.Create(blob)
	assert.NoError(t, err)

	// Delete the blob
	err = repo.DeleteById(blob.BlobID)
	assert.NoError(t, err)

	// Verify the blob is deleted
	var deletedBlob model.Blob
	err = db.First(&deletedBlob, "blob_id = ?", blob.BlobID).Error
	assert.Error(t, err, "Blob should be deleted")
	assert.Equal(t, gorm.ErrRecordNotFound, err, "Error should be 'record not found'")
}
