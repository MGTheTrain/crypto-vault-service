package cryptography

import (
	"crypto/rsa"
	cryptography "crypto_vault_service/internal/infrastructure/cryptography"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RSATests struct to encapsulate RSA-related test cases
type RSATests struct {
	rsaImpl *cryptography.RSAImpl
}

// NewRSATests is a constructor that creates a new instance of RSATests
func NewRSATests() *RSATests {
	return &RSATests{
		rsaImpl: &cryptography.RSAImpl{},
	}
}

// TestGenerateRSAKeys tests the generation of RSA keys
func (rt *RSATests) TestGenerateRSAKeys(t *testing.T) {
	// Generate RSA keys with 2048-bit size
	privateKey, publicKey, err := rt.rsaImpl.GenerateKeys(2048)
	assert.NoError(t, err, "Error generating RSA keys")
	assert.NotNil(t, privateKey, "Private key should not be nil")
	assert.NotNil(t, publicKey, "Public key should not be nil")

	// Ensure the public key's type is *rsa.PublicKey
	assert.IsType(t, &rsa.PublicKey{}, publicKey)
	assert.Equal(t, 2048, privateKey.N.BitLen(), "Private key bit length should be 2048")
}

// TestEncryptDecrypt tests the encryption and decryption methods of RSA
func (rt *RSATests) TestEncryptDecrypt(t *testing.T) {
	// Generate RSA keys
	privateKey, publicKey, err := rt.rsaImpl.GenerateKeys(2048)
	assert.NoError(t, err)

	// Message to encrypt
	plainText := []byte("This is a secret message")

	// Encrypt the message
	encryptedData, err := rt.rsaImpl.Encrypt(plainText, publicKey)
	assert.NoError(t, err, "Error encrypting data")

	// Decrypt the message
	decryptedData, err := rt.rsaImpl.Decrypt(encryptedData, privateKey)
	assert.NoError(t, err, "Error decrypting data")

	// Ensure the decrypted data matches the original message
	assert.Equal(t, plainText, decryptedData, "Decrypted data should match the original plaintext")
}

// TestSaveAndReadKeys tests saving and reading RSA keys to and from files
func (rt *RSATests) TestSaveAndReadKeys(t *testing.T) {
	// Generate RSA keys
	privateKey, publicKey, err := rt.rsaImpl.GenerateKeys(2048)
	assert.NoError(t, err)

	// Save keys to files
	privateKeyFile := "private.pem"
	publicKeyFile := "public.pem"

	err = rt.rsaImpl.SavePrivateKeyToFile(privateKey, privateKeyFile)
	assert.NoError(t, err, "Error saving private key to file")

	err = rt.rsaImpl.SavePublicKeyToFile(publicKey, publicKeyFile)
	assert.NoError(t, err, "Error saving public key to file")

	// Read the keys back from the files
	readPrivateKey, err := rt.rsaImpl.ReadPrivateKey(privateKeyFile)
	assert.NoError(t, err, "Error reading private key from file")
	assert.Equal(t, privateKey.N, readPrivateKey.N, "Private key N component should match")
	assert.Equal(t, privateKey.E, readPrivateKey.E, "Private key E component should match")

	readPublicKey, err := rt.rsaImpl.ReadPublicKey(publicKeyFile)
	assert.NoError(t, err, "Error reading public key from file")
	assert.Equal(t, publicKey.N, readPublicKey.N, "Public key N component should match")
	assert.Equal(t, publicKey.E, readPublicKey.E, "Public key E component should match")

	// Clean up the generated files
	os.Remove(privateKeyFile)
	os.Remove(publicKeyFile)
}

// TestEncryptWithInvalidKey tests encryption with an invalid public key
func (rt *RSATests) TestEncryptWithInvalidKey(t *testing.T) {
	// Generate RSA keys
	_, _, err := rt.rsaImpl.GenerateKeys(2048)
	assert.NoError(t, err)

	// Attempt to encrypt with a nil public key (invalid case)
	plainText := []byte("This should fail encryption")
	_, err = rt.rsaImpl.Encrypt(plainText, nil)
	assert.Error(t, err, "Encryption should fail with an invalid public key")

	// Attempt to decrypt with a nil private key (invalid case)
	_, err = rt.rsaImpl.Decrypt(plainText, nil)
	assert.Error(t, err, "Decryption should fail with an invalid private key")

	// Attempt to decrypt with a different private key (invalid case)
	_, err = rt.rsaImpl.Decrypt(plainText, &rsa.PrivateKey{})
	assert.Error(t, err, "Decryption should fail with an invalid private key")
}

// TestSavePrivateKeyInvalidPath tests saving a private key to an invalid path
func (rt *RSATests) TestSavePrivateKeyInvalidPath(t *testing.T) {
	// Generate RSA keys
	privateKey, _, err := rt.rsaImpl.GenerateKeys(2048)
	assert.NoError(t, err)

	// Try saving the private key to an invalid file path
	err = rt.rsaImpl.SavePrivateKeyToFile(privateKey, "/invalid/path/private.pem")
	assert.Error(t, err, "Saving private key to an invalid path should return an error")
}

// TestSavePublicKeyInvalidPath tests saving a public key to an invalid path
func (rt *RSATests) TestSavePublicKeyInvalidPath(t *testing.T) {
	// Generate RSA keys
	_, publicKey, err := rt.rsaImpl.GenerateKeys(2048)
	assert.NoError(t, err)

	// Try saving the public key to an invalid file path
	err = rt.rsaImpl.SavePublicKeyToFile(publicKey, "/invalid/path/public.pem")
	assert.Error(t, err, "Saving public key to an invalid path should return an error")
}

// TestRSA is the entry point to run the RSA tests
func TestRSA(t *testing.T) {
	// Create a new RSA test suite instance
	rt := NewRSATests()

	// Run each test method
	t.Run("TestGenerateRSAKeys", rt.TestGenerateRSAKeys)
	t.Run("TestEncryptDecrypt", rt.TestEncryptDecrypt)
	t.Run("TestSaveAndReadKeys", rt.TestSaveAndReadKeys)
	t.Run("TestEncryptWithInvalidKey", rt.TestEncryptWithInvalidKey)
	t.Run("TestSavePrivateKeyInvalidPath", rt.TestSavePrivateKeyInvalidPath)
	t.Run("TestSavePublicKeyInvalidPath", rt.TestSavePublicKeyInvalidPath)
}
