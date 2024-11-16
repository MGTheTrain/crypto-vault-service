package cryptography

import (
	"fmt"
	"os"
	"testing"

	cryptography "crypto_vault_service/internal/infrastructure/cryptography"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PKCS11Test struct {
	Token *cryptography.PKCS11Token
}

// NewPKCS11Test sets up the test environment for PKCS#11 integration tests
func NewPKCS11Test(slot, modulePath, Label, soPin, userPin, objectLabel, keyType string, keySize int) *PKCS11Test {
	return &PKCS11Test{
		Token: &cryptography.PKCS11Token{
			ModulePath:  modulePath,
			Label:       Label,
			SOPin:       soPin,
			UserPin:     userPin,
			ObjectLabel: objectLabel,
			KeyType:     keyType,
			KeySize:     keySize,
		},
	}
}

// Setup initializes the PKCS#11 token
func (p *PKCS11Test) Setup(t *testing.T) {
	tokenSlot := "0x1"
	err := p.Token.InitializeToken(tokenSlot)
	require.NoError(t, err, "Failed to initialize PKCS#11 token")

	isTokenSet, err := p.Token.IsTokenSet()
	require.NoError(t, err, "Error checking if token is set")
	assert.True(t, isTokenSet, "The token should be initialized and set")
}

// DeleteKeyFromToken deletes any existing key with the same label before adding a new key.
func (p *PKCS11Test) DeleteKeyFromToken(t *testing.T) {
	// Deleting the private key
	err := p.Token.DeleteObject("privkey", p.Token.ObjectLabel)
	if err != nil {
		t.Logf("Warning: Failed to delete existing private key: %v\n", err)
	}

	// Deleting the public key
	err = p.Token.DeleteObject("pubkey", p.Token.ObjectLabel)
	if err != nil {
		t.Logf("Warning: Failed to delete existing public key: %v\n", err)
	}

	// Deleting the secret key (only if it exists)
	err = p.Token.DeleteObject("secrkey", p.Token.ObjectLabel)
	if err != nil {
		t.Logf("Warning: Failed to delete existing secret key: %v\n", err)
	}
}

// AddKeyToToken is a helper function to add a key to a token
func (p *PKCS11Test) AddKeyToToken(t *testing.T) {
	err := p.Token.AddKey()
	assert.NoError(t, err, "Failed to add key to the token")

	isObjectSet, err := p.Token.IsObjectSet()
	assert.NoError(t, err, "Error checking if key is set")
	assert.True(t, isObjectSet, fmt.Sprintf("The %s key should be added to the token", p.Token.KeyType))
}

// TestAddRSAKey tests adding an RSA key to a PKCS#11 token
func TestAddRSAKey(t *testing.T) {
	test := NewPKCS11Test("0x1", "/usr/lib/softhsm/libsofthsm2.so", "MyToken", "123456", "234567", "TestRSAKey", "RSA", 2048)

	test.Setup(t)

	// Add an RSA key to the token
	test.AddKeyToToken(t)

	test.DeleteKeyFromToken(t)
}

// TestAddECDSAKey tests adding an ECDSA key to a PKCS#11 token
func TestAddECDSAKey(t *testing.T) {
	test := NewPKCS11Test("0x1", "/usr/lib/softhsm/libsofthsm2.so", "MyToken", "123456", "234567", "TestECDSAKey", "ECDSA", 256)

	test.Setup(t)

	// Add an ECDSA key to the token
	test.AddKeyToToken(t)

	test.DeleteKeyFromToken(t)
}

// TestEncryptDecrypt tests the encryption ad decryption functionality of the PKCS#11 token
func TestEncryptDecrypt(t *testing.T) {
	// Prepare the test PKCS#11Token
	test := NewPKCS11Test("0x1", "/usr/lib/softhsm/libsofthsm2.so", "MyToken", "123456", "234567", "TestRSAKey", "RSA", 2048)
	test.Setup(t)

	// Add an RSA key to the token
	test.AddKeyToToken(t)

	// Sample input file with plaintext data (for testing purposes)
	inputFilePath := "plain-text.txt"
	err := os.WriteFile(inputFilePath, []byte("This is some data to encrypt."), 0644)
	require.NoError(t, err, "Failed to write plaintext data to input file")

	// Output file path where encrypted data will be stored
	outputFilePath := "encrypted.bin"

	// Encrypt the data using the Encrypt method
	err = test.Token.Encrypt(inputFilePath, outputFilePath)
	assert.NoError(t, err, "Failed to encrypt data using the PKCS#11 token")

	// Try reading the encrypted data from the output file
	encryptedData, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read encrypted data from output file: %v", err)
	}

	// Ensure the encrypted data is non-empty
	assert.NotEmpty(t, encryptedData, "Encrypted data should not be empty")

	// Decrypt the data using the Decrypt method
	decryptedFilePath := "decrypted.txt"
	err = test.Token.Decrypt(outputFilePath, decryptedFilePath)
	assert.NoError(t, err, "Failed to decrypt data using the PKCS#11 token")

	// Try reading the decrypted data from the file
	decryptedData, err := os.ReadFile(decryptedFilePath)
	if err != nil {
		t.Fatalf("Failed to read decrypted data from output file: %v", err)
	}

	// Ensure the decrypted data matches the original plaintext data
	originalData, err := os.ReadFile(inputFilePath)
	require.NoError(t, err, "Failed to read original input file")

	assert.Equal(t, originalData, decryptedData, "Decrypted data should match the original plaintext")

	// Clean up by deleting the key from the token
	test.DeleteKeyFromToken(t)

	// Optionally, delete the files after the test
	err = os.Remove(inputFilePath)
	require.NoError(t, err, "Failed to remove input file")
	err = os.Remove(outputFilePath)
	require.NoError(t, err, "Failed to remove encrypted file")
	err = os.Remove(decryptedFilePath)
	require.NoError(t, err, "Failed to remove decrypted file")
}
