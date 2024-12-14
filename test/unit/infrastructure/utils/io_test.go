package utils

import (
	"os"
	"testing"

	"crypto_vault_service/internal/infrastructure/utils"

	"github.com/stretchr/testify/assert"
)

// TestWriteFileAndReadFile tests the WriteFile and ReadFile functions
func TestWriteFileAndReadFile(t *testing.T) {
	// Prepare a test file path
	testFilePath := "testfile.txt"
	testData := []byte("This is a test message.")

	// Write data to the file
	err := utils.WriteFile(testFilePath, testData)
	assert.NoError(t, err, "Error writing to file")

	// Read data from the file
	readData, err := utils.ReadFile(testFilePath)
	assert.NoError(t, err, "Error reading from file")

	// Ensure that the written and read data are the same
	assert.Equal(t, testData, readData, "The data read from the file should match the data written")

	// Clean up the test file
	err = os.Remove(testFilePath)
	assert.NoError(t, err, "Error cleaning up the test file")
}

// TestReadFileWithNonExistentFile tests the ReadFile function with a non-existent file
func TestReadFileWithNonExistentFile(t *testing.T) {
	// Use a file path that doesn't exist
	testFilePath := "non_existent_file.txt"

	// Attempt to read from the non-existent file
	data, err := utils.ReadFile(testFilePath)

	// We expect an error because the file does not exist
	assert.Error(t, err, "Reading from a non-existent file should return an error")
	assert.Nil(t, data, "No data should be returned for a non-existent file")
}
