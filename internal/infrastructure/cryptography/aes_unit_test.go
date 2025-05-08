//go:build unit
// +build unit

package cryptography

import (
	"log"
	"testing"

	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"

	"github.com/stretchr/testify/assert"
)

// AES struct to encapsulate AES-related test cases
type AESTests struct {
	aes *AES
}

// NewAESTests is a constructor that creates a new instance of AESTests
func NewAESTests(t *testing.T) *AESTests {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	aes, err := NewAES(logger)
	if err != nil {
		t.Logf("%v\n", err)
	}

	return &AESTests{
		aes: aes,
	}
}

// TestEncryptDecrypt tests the encryption and decryption functionality
func (at *AESTests) TestEncryptDecrypt(t *testing.T) {
	// Generate a random key of 16 bytes (128-bit AES)
	key, err := at.aes.GenerateKey(16)
	assert.NoError(t, err)

	// Define a plaintext to encrypt and decrypt
	plainText := []byte("This is a test message.")

	// Encrypt the plaintext
	ciphertext, err := at.aes.Encrypt(plainText, key)
	assert.NoError(t, err)
	assert.NotNil(t, ciphertext)
	assert.Greater(t, len(ciphertext), 0, "Ciphertext should be longer than 0")

	// Decrypt the ciphertext
	decryptedText, err := at.aes.Decrypt(ciphertext, key)
	assert.NoError(t, err)
	assert.NotNil(t, decryptedText)

	// Assert the decrypted text is the same as the original plaintext
	assert.Equal(t, plainText, decryptedText)
}

// TestEncryptionWithInvalidKey tests encryption with invalid key sizes
func (at *AESTests) TestEncryptionWithInvalidKey(t *testing.T) {
	// Try generating an invalid key (e.g., 8 bytes instead of a standard AES size)
	key := []byte("shortkey")
	plainText := []byte("This is a test.")

	// Try encrypting with an invalid key
	_, err := at.aes.Encrypt(plainText, key)
	assert.Error(t, err)
}

// TestGenerateKey tests key generation functionality
func (at *AESTests) TestGenerateKey(t *testing.T) {
	// Generate a random AES key with 16 bytes (128-bit AES)
	key, err := at.aes.GenerateKey(16)
	assert.NoError(t, err)
	assert.Equal(t, len(key), 16)

	// Try generating a 32-byte AES key (256-bit AES)
	key256, err := at.aes.GenerateKey(32)
	assert.NoError(t, err)
	assert.Equal(t, len(key256), 32)
}

// TestDecryptWithWrongKey tests decryption with a wrong key
func (at *AESTests) TestDecryptWithWrongKey(t *testing.T) {
	// Generate a random 16-byte AES key
	key, err := at.aes.GenerateKey(16)
	assert.NoError(t, err)

	// Encrypt the data
	plainText := []byte("Test decryption with wrong key.")
	ciphertext, err := at.aes.Encrypt(plainText, key)
	assert.NoError(t, err)

	// Generate a new, different key for decryption
	anotherKey, err := at.aes.GenerateKey(16)
	assert.NoError(t, err)

	// Try to decrypt with the wrong key
	_, err = at.aes.Decrypt(ciphertext, anotherKey)
	assert.Error(t, err, "Decryption with the wrong key should fail")
}

// TestDecryptShortCiphertext tests the case where the ciphertext is too short
func (at *AESTests) TestDecryptShortCiphertext(t *testing.T) {
	// Generate a random key
	key, err := at.aes.GenerateKey(16)
	assert.NoError(t, err)

	// Attempt to decrypt a too-short ciphertext
	_, err = at.aes.Decrypt([]byte("short"), key)
	assert.Error(t, err, "Decrypting a ciphertext that's too short should fail")
}

// TestAES is the entry point to run the AES tests
func TestAES(t *testing.T) {
	// Create a new AES test suite instance
	at := NewAESTests(t)

	// Run each test method
	t.Run("TestEncryptDecrypt", at.TestEncryptDecrypt)
	t.Run("TestEncryptionWithInvalidKey", at.TestEncryptionWithInvalidKey)
	t.Run("TestGenerateKey", at.TestGenerateKey)
	t.Run("TestDecryptWithWrongKey", at.TestDecryptWithWrongKey)
	t.Run("TestDecryptShortCiphertext", at.TestDecryptShortCiphertext)
}
