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

// AESProcessorTests encapsulates AES-related test cases
type AESProcessorTests struct {
	processor AESProcessor
}

// NewAESProcessorTests creates a new instance of AESProcessorTests
func NewAESProcessorTests(t *testing.T) *AESProcessorTests {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logInstance, err := logger.GetLogger(loggerSettings)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	aes, err := NewAESProcessor(logInstance)
	if err != nil {
		t.Fatalf("Failed to create AES processor: %v", err)
	}

	return &AESProcessorTests{
		processor: aes,
	}
}

func (at *AESProcessorTests) TestEncryptDecrypt(t *testing.T) {
	key, err := at.processor.GenerateKey(16)
	assert.NoError(t, err)

	plainText := []byte("This is a test message.")

	ciphertext, err := at.processor.Encrypt(plainText, key)
	assert.NoError(t, err)
	assert.NotNil(t, ciphertext)
	assert.Greater(t, len(ciphertext), 0)

	decryptedText, err := at.processor.Decrypt(ciphertext, key)
	assert.NoError(t, err)
	assert.NotNil(t, decryptedText)
	assert.Equal(t, plainText, decryptedText)
}

func (at *AESProcessorTests) TestEncryptionWithInvalidKey(t *testing.T) {
	key := []byte("shortkey")
	plainText := []byte("This is a test.")

	_, err := at.processor.Encrypt(plainText, key)
	assert.Error(t, err)
}

func (at *AESProcessorTests) TestGenerateKey(t *testing.T) {
	key, err := at.processor.GenerateKey(16)
	assert.NoError(t, err)
	assert.Equal(t, 16, len(key))

	key256, err := at.processor.GenerateKey(32)
	assert.NoError(t, err)
	assert.Equal(t, 32, len(key256))
}

func (at *AESProcessorTests) TestDecryptWithWrongKey(t *testing.T) {
	key, err := at.processor.GenerateKey(16)
	assert.NoError(t, err)

	plainText := []byte("Test decryption with wrong key.")
	ciphertext, err := at.processor.Encrypt(plainText, key)
	assert.NoError(t, err)

	wrongKey, err := at.processor.GenerateKey(16)
	assert.NoError(t, err)

	decrypted, err := at.processor.Decrypt(ciphertext, wrongKey)

	// Accept either an error or bad output (not equal to original plaintext)
	if err == nil {
		assert.NotEqual(t, plainText, decrypted, "Decryption with wrong key should not return original message")
	} else {
		assert.Error(t, err, "Expected an error when decrypting with wrong key")
	}
}

func (at *AESProcessorTests) TestDecryptShortCiphertext(t *testing.T) {
	key, err := at.processor.GenerateKey(16)
	assert.NoError(t, err)

	_, err = at.processor.Decrypt([]byte("short"), key)
	assert.Error(t, err)
}

// Entry point to run AESProcessorTests
func TestAESProcessor(t *testing.T) {
	tests := NewAESProcessorTests(t)

	t.Run("TestEncryptDecrypt", tests.TestEncryptDecrypt)
	t.Run("TestEncryptionWithInvalidKey", tests.TestEncryptionWithInvalidKey)
	t.Run("TestGenerateKey", tests.TestGenerateKey)
	t.Run("TestDecryptWithWrongKey", tests.TestDecryptWithWrongKey)
	t.Run("TestDecryptShortCiphertext", tests.TestDecryptShortCiphertext)
}
