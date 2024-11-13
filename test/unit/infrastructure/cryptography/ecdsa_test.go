package cryptography

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptography "crypto_vault_service/internal/infrastructure/cryptography"
	"math/big"
	"os"
	"testing"

	"encoding/hex"

	"github.com/stretchr/testify/assert"
)

type ECDSATests struct {
	ecc *cryptography.ECDSAImpl
}

// NewECDSATests is a constructor that creates a new instance of ECDSATests
func NewECDSATests() *ECDSATests {
	return &ECDSATests{
		ecc: &cryptography.ECDSAImpl{},
	}
}

// TestGenerateKeys tests the key generation functionality
func (et *ECDSATests) TestGenerateKeys(t *testing.T) {
	// Generate ECDSA keys using P256 curve
	privateKey, publicKey, err := et.ecc.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.NotNil(t, publicKey)
	assert.Equal(t, elliptic.P256(), privateKey.PublicKey.Curve)
	assert.Equal(t, elliptic.P256(), publicKey.Curve)
}

// TestSignVerify tests signing and verifying functionality
func (et *ECDSATests) TestSignVerify(t *testing.T) {
	// Generate ECDSA keys
	privateKey, publicKey, err := et.ecc.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	// Message to sign
	message := []byte("This is a test message.")

	// Sign the message
	signature, err := et.ecc.Sign(message, privateKey)
	assert.NoError(t, err)
	assert.NotNil(t, signature)

	// Verify the signature
	valid, err := et.ecc.Verify(message, signature, publicKey)
	assert.NoError(t, err)
	assert.True(t, valid, "The signature should be valid")

	// Modify the message and try verifying the signature
	modifiedMessage := []byte("This is a modified message.")
	valid, err = et.ecc.Verify(modifiedMessage, signature, publicKey)
	assert.NoError(t, err)
	assert.False(t, valid, "The signature should not be valid for a modified message")
}

// TestSaveAndReadKeys tests saving and reading the private and public keys from PEM files
func (et *ECDSATests) TestSaveAndReadKeys(t *testing.T) {
	// Generate ECDSA keys
	privateKey, publicKey, err := et.ecc.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	// Save private and public keys to files
	privateKeyFile := "private.pem"
	publicKeyFile := "public.pem"
	err = et.ecc.SavePrivateKeyToFile(privateKey, privateKeyFile)
	assert.NoError(t, err)

	err = et.ecc.SavePublicKeyToFile(publicKey, publicKeyFile)
	assert.NoError(t, err)

	// Read the private and public keys from the files
	readPrivateKey, err := et.ecc.ReadPrivateKey(privateKeyFile)
	assert.NoError(t, err)
	assert.Equal(t, privateKey.D, readPrivateKey.D)
	assert.Equal(t, privateKey.PublicKey.X, readPrivateKey.PublicKey.X)
	assert.Equal(t, privateKey.PublicKey.Y, readPrivateKey.PublicKey.Y)

	readPublicKey, err := et.ecc.ReadPublicKey(publicKeyFile)
	assert.NoError(t, err)
	assert.Equal(t, publicKey.X, readPublicKey.X)
	assert.Equal(t, publicKey.Y, readPublicKey.Y)

	// Clean up the generated files
	os.Remove(privateKeyFile)
	os.Remove(publicKeyFile)
}

// TestSaveSignatureToFile tests saving a signature to a file
func (et *ECDSATests) TestSaveSignatureToFile(t *testing.T) {
	// Generate ECDSA keys
	privateKey, _, err := et.ecc.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	// Message to sign
	message := []byte("This is a test message.")

	// Sign the message
	signature, err := et.ecc.Sign(message, privateKey)
	assert.NoError(t, err)
	assert.NotNil(t, signature)

	// Save the signature to a file
	signatureFile := "signature.hex"
	err = et.ecc.SaveSignatureToFile(signatureFile, signature)
	assert.NoError(t, err)

	// Read the saved signature from the file
	hexData, err := os.ReadFile(signatureFile)
	assert.NoError(t, err)

	// Decode the hex signature
	decodedSignature, err := hex.DecodeString(string(hexData))
	assert.NoError(t, err)
	assert.Equal(t, signature, decodedSignature)

	// Clean up the generated signature file
	os.Remove(signatureFile)
}

// TestSignWithInvalidPrivateKey tests signing with an invalid private key
func (et *ECDSATests) TestSignWithInvalidPrivateKey(t *testing.T) {
	// Generate ECDSA keys (valid ones)
	_, _, err := et.ecc.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	// Modify the private key to make it invalid (e.g., set D to 0)
	invalidPrivateKey := &ecdsa.PrivateKey{
		D: new(big.Int).SetInt64(0), // Invalid key with D = 0
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
		},
	}

	// Attempt to sign a message with the invalid private key
	message := []byte("This message will fail to sign")
	_, err = et.ecc.Sign(message, invalidPrivateKey)
	assert.Error(t, err, "Signing with an invalid private key should fail")
}

// TestVerifyWithInvalidPublicKey tests verifying with an invalid public key
func (et *ECDSATests) TestVerifyWithInvalidPublicKey(t *testing.T) {
	// Generate ECDSA keys
	privateKey, _, err := et.ecc.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	// Sign the message
	message := []byte("This is a test message.")
	signature, err := et.ecc.Sign(message, privateKey)
	assert.NoError(t, err)

	// Create an invalid public key (e.g., public key X = 0)
	invalidPublicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetInt64(0),
		Y:     new(big.Int).SetInt64(0),
	}

	// Attempt to verify the signature with the invalid public key
	valid, err := et.ecc.Verify(message, signature, invalidPublicKey)
	assert.NoError(t, err)
	assert.False(t, valid, "Verification with an invalid public key should fail")
}

func TestECDSA(t *testing.T) {
	// Create a new ECDSA test suite instance
	et := NewECDSATests()

	// Run each test method
	t.Run("TestGenerateKeys", et.TestGenerateKeys)
	t.Run("TestSignVerify", et.TestSignVerify)
	t.Run("TestSaveAndReadKeys", et.TestSaveAndReadKeys)
	t.Run("TestSaveSignatureToFile", et.TestSaveSignatureToFile)
	t.Run("TestSignWithInvalidPrivateKey", et.TestSignWithInvalidPrivateKey)
	t.Run("TestVerifyWithInvalidPublicKey", et.TestVerifyWithInvalidPublicKey)
}
