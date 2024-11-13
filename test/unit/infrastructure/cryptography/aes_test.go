package cryptography

import (
	"testing"

	cryptography "crypto_vault_service/internal/infrastructure/cryptography"

	"github.com/stretchr/testify/assert"
)

// TestEncryptDecrypt tests the encryption and decryption functionality
func TestEncryptDecrypt(t *testing.T) {
	aesImpl := &cryptography.AESImpl{}

	// Generate a random key of 16 bytes (128-bit AES)
	key, err := aesImpl.GenerateKey(16)
	assert.NoError(t, err)

	// Define a plaintext to encrypt and decrypt
	plainText := []byte("This is a test message.")

	// Encrypt the plaintext
	ciphertext, err := aesImpl.Encrypt(plainText, key)
	assert.NoError(t, err)
	assert.NotNil(t, ciphertext)
	assert.Greater(t, len(ciphertext), 0, "Ciphertext should be longer than 0")

	// Decrypt the ciphertext
	decryptedText, err := aesImpl.Decrypt(ciphertext, key)
	assert.NoError(t, err)
	assert.NotNil(t, decryptedText)

	// Assert the decrypted text is the same as the original plaintext
	assert.Equal(t, plainText, decryptedText)
}

// TestEncryptionWithInvalidKey tests encryption with invalid key sizes
func TestEncryptionWithInvalidKey(t *testing.T) {
	aesImpl := &cryptography.AESImpl{}

	// Try generating an invalid key (e.g., 8 bytes instead of a standard AES size)
	key := []byte("shortkey")
	plainText := []byte("This is a test.")

	// Try encrypting with an invalid key
	_, err := aesImpl.Encrypt(plainText, key)
	assert.Error(t, err)
}

// TestGenerateKey tests key generation functionality
func TestGenerateKey(t *testing.T) {
	aesImpl := &cryptography.AESImpl{}

	// Generate a random AES key with 16 bytes (128-bit AES)
	key, err := aesImpl.GenerateKey(16)
	assert.NoError(t, err)
	assert.Equal(t, len(key), 16)

	// Try generating a 32-byte AES key (256-bit AES)
	key256, err := aesImpl.GenerateKey(32)
	assert.NoError(t, err)
	assert.Equal(t, len(key256), 32)
}

// TestDecryptWithWrongKey tests decryption with a wrong key
func TestDecryptWithWrongKey(t *testing.T) {
	aesImpl := &cryptography.AESImpl{}

	// Generate a random 16-byte AES key
	key, err := aesImpl.GenerateKey(16)
	assert.NoError(t, err)

	// Encrypt the data
	plainText := []byte("Test decryption with wrong key.")
	ciphertext, err := aesImpl.Encrypt(plainText, key)
	assert.NoError(t, err)

	// Generate a new, different key for decryption
	anotherKey, err := aesImpl.GenerateKey(16)
	assert.NoError(t, err)

	// Try to decrypt with the wrong key
	_, err = aesImpl.Decrypt(ciphertext, anotherKey)
	assert.Error(t, err, "Decryption with the wrong key should fail")
}

// TestDecryptShortCiphertext tests the case where the ciphertext is too short
func TestDecryptShortCiphertext(t *testing.T) {
	aesImpl := &cryptography.AESImpl{}

	// Generate a random key
	key, err := aesImpl.GenerateKey(16)
	assert.NoError(t, err)

	// Attempt to decrypt a too-short ciphertext
	_, err = aesImpl.Decrypt([]byte("short"), key)
	assert.Error(t, err, "Decrypting a ciphertext that's too short should fail")
}
