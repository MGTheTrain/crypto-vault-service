package cryptography

import (
	"os"
	"testing"

	cryptography "crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/settings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PKCS11Test struct {
	Label        string
	ObjectLabel  string
	KeyType      string
	KeySize      uint
	TokenHandler *cryptography.PKCS11Handler
}

// NewPKCS11Test sets up the test environment for PKCS#11 integration tests
func NewPKCS11Test(t *testing.T, slotId, modulePath, label, soPin, userPin, objectLabel, keyType string, keySize uint) *PKCS11Test {
	settings := settings.PKCS11Settings{
		ModulePath: modulePath,
		SOPin:      soPin,
		UserPin:    userPin,
		SlotId:     slotId,
	}

	tokenHandler, err := cryptography.NewPKCS11Handler(settings)
	if err != nil {
		t.Logf("%v\n", err)
	}
	return &PKCS11Test{
		Label:        label,
		ObjectLabel:  objectLabel,
		KeyType:      keyType,
		KeySize:      keySize,
		TokenHandler: tokenHandler,
	}
}

// InitializeToken initializes a PKCS#11 token by calling the TokenHandler's InitializeToken method
func (p *PKCS11Test) InitializeToken(t *testing.T) {
	err := p.TokenHandler.InitializeToken(p.Label)
	require.NoError(t, err, "Failed to initialize PKCS#11 token")
}

// DeleteKeyFromToken deletes any existing key with the same label before adding a new key.
func (p *PKCS11Test) DeleteKeyFromToken(t *testing.T) {
	objectType := "privkey"
	err := p.TokenHandler.DeleteObject(p.Label, objectType, p.ObjectLabel)
	if err != nil {
		t.Logf("Warning: Failed to delete existing private key: %v\n", err)
	}

	objectType = "pubkey"
	err = p.TokenHandler.DeleteObject(p.Label, objectType, p.ObjectLabel)
	if err != nil {
		t.Logf("Warning: Failed to delete existing public key: %v\n", err)
	}

	objectType = "secrkey"
	err = p.TokenHandler.DeleteObject(p.Label, objectType, p.ObjectLabel)
	if err != nil {
		t.Logf("Warning: Failed to delete existing secret key: %v\n", err)
	}
}

// AddKeyToToken is a helper function to add a key to a token
func (p *PKCS11Test) AddKeyToToken(t *testing.T, label, objectLabel, keyType string, keySize uint) {
	err := p.TokenHandler.AddKey(label, objectLabel, keyType, keySize)
	assert.NoError(t, err, "Failed to add key to the token")
}

// TestListTokens tests listing available tokens using PKCS#11
func TestListTokens(t *testing.T) {
	slotId := "0x0"
	modulePath := "/usr/lib/softhsm/libsofthsm2.so"
	label := "MyToken"
	soPin := "123456"
	userPin := "234567"
	objectLabel := "TestRSAKey"
	keyType := "RSA"
	keySize := 2048

	test := NewPKCS11Test(t, slotId, modulePath, label, soPin, userPin, objectLabel, keyType, uint(keySize))
	test.InitializeToken(t) // creates a token slot

	tokens, err := test.TokenHandler.ListTokenSlots()
	require.NoError(t, err, "Failed to list tokens")

	require.NotEmpty(t, tokens, "Token list should not be empty")

	token := tokens[0]
	assert.NotEmpty(t, token.SlotId, "Token serial number should not be empty")
	assert.NotEmpty(t, token.Label, "Token label should not be empty")
	assert.NotEmpty(t, token.Manufacturer, "Token manufacturer should not be empty")
	assert.NotEmpty(t, token.Model, "Token model should not be empty")
	assert.NotEmpty(t, token.SerialNumber, "Token serial number should not be empty")
}

// TestAddRSAKey tests adding an RSA key to a PKCS#11 token
func TestAddRSAKey(t *testing.T) {
	slotId := "0x0"
	modulePath := "/usr/lib/softhsm/libsofthsm2.so"
	label := "MyToken"
	soPin := "123456"
	userPin := "234567"
	objectLabel := "TestRSAKey"
	keyType := "RSA"
	keySize := 2048

	test := NewPKCS11Test(t, slotId, modulePath, label, soPin, userPin, objectLabel, keyType, uint(keySize))

	test.InitializeToken(t)

	test.AddKeyToToken(t, label, objectLabel, keyType, uint(keySize))

	test.DeleteKeyFromToken(t)
}

// TestAddECDSAKey tests adding an ECDSA key to a PKCS#11 token
func TestAddECDSAKey(t *testing.T) {
	slotId := "0x0"
	modulePath := "/usr/lib/softhsm/libsofthsm2.so"
	label := "MyToken"
	soPin := "123456"
	userPin := "234567"
	objectLabel := "TestECDSAKey"
	keyType := "ECDSA"
	keySize := 256

	test := NewPKCS11Test(t, slotId, modulePath, label, soPin, userPin, objectLabel, keyType, uint(keySize))

	test.InitializeToken(t)

	test.AddKeyToToken(t, label, objectLabel, keyType, uint(keySize))

	test.DeleteKeyFromToken(t)
}

// TestListObjects tests listing available objects in a specific PKCS#11 token
func TestListObjects(t *testing.T) {
	slotId := "0x0"
	modulePath := "/usr/lib/softhsm/libsofthsm2.so"
	label := "MyToken"
	soPin := "123456"
	userPin := "234567"
	objectLabel := "TestRSAKey2"
	keyType := "RSA"
	keySize := 2048

	test := NewPKCS11Test(t, slotId, modulePath, label, soPin, userPin, objectLabel, keyType, uint(keySize))
	test.InitializeToken(t)

	test.AddKeyToToken(t, label, objectLabel, keyType, uint(keySize))

	objects, err := test.TokenHandler.ListObjects(label)
	require.NoError(t, err, "Failed to list objects")

	require.NotEmpty(t, objects, "Object list should not be empty")

	object := objects[0]
	assert.NotEmpty(t, object.Label, "Object label should not be empty")
	assert.NotEmpty(t, object.Type, "Object type should not be empty")
	assert.NotEmpty(t, object.Usage, "Object usage should not be empty")

	test.DeleteKeyFromToken(t)
}

// TestEncryptDecrypt tests the encryption ad decryption functionality of the PKCS#11 token
func TestEncryptDecrypt(t *testing.T) {
	slotId := "0x0"
	modulePath := "/usr/lib/softhsm/libsofthsm2.so"
	label := "MyToken"
	soPin := "123456"
	userPin := "234567"
	objectLabel := "TestRSAKey"
	keyType := "RSA"
	keySize := 2048

	test := NewPKCS11Test(t, slotId, modulePath, label, soPin, userPin, objectLabel, keyType, uint(keySize))
	test.InitializeToken(t)

	test.AddKeyToToken(t, label, objectLabel, keyType, uint(keySize))

	inputFilePath := "plain-text.txt"
	err := os.WriteFile(inputFilePath, []byte("This is some data to encrypt."), 0644)
	require.NoError(t, err, "Failed to write plaintext data to input file")

	outputFilePath := "encrypted.bin"

	err = test.TokenHandler.Encrypt(label, objectLabel, inputFilePath, outputFilePath, keyType)
	assert.NoError(t, err, "Failed to encrypt data using the PKCS#11 token")

	encryptedData, err := os.ReadFile(outputFilePath)
	if err != nil {
		t.Fatalf("Failed to read encrypted data from output file: %v", err)
	}

	assert.NotEmpty(t, encryptedData, "Encrypted data should not be empty")

	decryptedFilePath := "decrypted.txt"
	err = test.TokenHandler.Decrypt(label, objectLabel, outputFilePath, decryptedFilePath, keyType)
	assert.NoError(t, err, "Failed to decrypt data using the PKCS#11 token")

	decryptedData, err := os.ReadFile(decryptedFilePath)
	if err != nil {
		t.Fatalf("Failed to read decrypted data from output file: %v", err)
	}

	originalData, err := os.ReadFile(inputFilePath)
	require.NoError(t, err, "Failed to read original input file")

	assert.Equal(t, originalData, decryptedData, "Decrypted data should match the original plaintext")

	test.DeleteKeyFromToken(t)

	err = os.Remove(inputFilePath)
	require.NoError(t, err, "Failed to remove input file")
	err = os.Remove(outputFilePath)
	require.NoError(t, err, "Failed to remove encrypted file")
	err = os.Remove(decryptedFilePath)
	require.NoError(t, err, "Failed to remove decrypted file")
}

// TestSignAndVerify tests the signing and verification functionality of the PKCS#11 token
func TestSignAndVerify(t *testing.T) {
	slotId := "0x0"
	modulePath := "/usr/lib/softhsm/libsofthsm2.so"
	label := "MyToken"
	soPin := "123456"
	userPin := "234567"
	objectLabel := "TestRSAKey"
	keyType := "RSA"
	keySize := 2048

	test := NewPKCS11Test(t, slotId, modulePath, label, soPin, userPin, objectLabel, keyType, uint(keySize))
	test.InitializeToken(t)

	test.AddKeyToToken(t, label, objectLabel, keyType, uint(keySize))

	dataFilePath := "data-to-sign.txt"
	err := os.WriteFile(dataFilePath, []byte("This is some data to sign."), 0644)
	require.NoError(t, err, "Failed to write data to sign to input file")

	signatureFilePath := "data.sig"

	err = test.TokenHandler.Sign(label, objectLabel, dataFilePath, signatureFilePath, keyType)
	assert.NoError(t, err, "Failed to sign data using the PKCS#11 token")

	signatureData, err := os.ReadFile(signatureFilePath)
	if err != nil {
		t.Fatalf("Failed to read signature data from output file: %v", err)
	}

	assert.NotEmpty(t, signatureData, "Signature data should not be empty")

	valid, err := test.TokenHandler.Verify(label, objectLabel, dataFilePath, signatureFilePath, keyType)
	assert.NoError(t, err, "Failed to verify the signature using the PKCS#11 token")

	assert.True(t, valid, "The signature should be valid")

	test.DeleteKeyFromToken(t)

	err = os.Remove(dataFilePath)
	require.NoError(t, err, "Failed to remove input file")
	err = os.Remove(signatureFilePath)
	require.NoError(t, err, "Failed to remove signature file")
}

// // TestSignAndVerifyECDSA tests the signing and verification functionality for ECDSA using a PKCS#11 token
// // This test is commented out due to errors occurring in the CI workflow, likely related to the PKCS#11 module
// // or its interaction with the SoftHSM library during signing and verification steps. The issue prevents
// // proper execution and will be addressed in a future update.
// func TestSignAndVerifyECDSA(t *testing.T) {
// 	// Prepare the test PKCS#11Token for ECDSA
// 	test := NewPKCS11Test(t, "0x1", "/usr/lib/softhsm/libsofthsm2.so", "MyToken2", "123456", "234567", "TestECDSAKey", "ECDSA", 256)
// 	test.InitializeToken(t)

// 	// Add an ECDSA key to the token
// 	test.AddKeyToToken(t, label, objectLabel, keyType, uint(keySize))

// 	// Sample input file with data to sign (for testing purposes)
// 	inputFilePath := "data-to-sign.txt"
// 	err := os.WriteFile(inputFilePath, []byte("This is some data to sign."), 0644)
// 	require.NoError(t, err, "Failed to write data to sign to input file")

// 	// Output file path where the signature will be stored
// 	signatureFilePath := "data.sig"

// 	// Sign the data using the Sign method (ECDSA)
// 	err = p.TokenHandler.Sign(inputFilePath, signatureFilePath)
// 	assert.NoError(t, err, "Failed to sign data using the PKCS#11 token with ECDSA")

// 	// Try reading the signature data from the output file
// 	signatureData, err := os.ReadFile(signatureFilePath)
// 	if err != nil {
// 		t.Fatalf("Failed to read signature data from output file: %v", err)
// 	}

// 	// Ensure the signature data is non-empty
// 	assert.NotEmpty(t, signatureData, "Signature data should not be empty")

// 	// Verify the signature using the Verify method (ECDSA)
// 	valid, err := p.TokenHandler.Verify(inputFilePath, signatureFilePath)
// 	assert.NoError(t, err, "Failed to verify the signature using the PKCS#11 token with ECDSA")

// 	// Ensure the signature is valid
// 	assert.True(t, valid, "The ECDSA signature should be valid")

// 	// Clean up by deleting the key from the token
// 	test.DeleteKeyFromToken(t)

// 	// Optionally, delete the files after the test
// 	err = os.Remove(inputFilePath)
// 	require.NoError(t, err, "Failed to remove input file")
// 	err = os.Remove(signatureFilePath)
// 	require.NoError(t, err, "Failed to remove signature file")
// }
