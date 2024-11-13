package cryptography

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	cryptography "crypto_vault_service/internal/infrastructure/cryptography"
	"encoding/hex"
)

// TestGenerateKeys tests the key generation functionality
func TestGenerateKeys(t *testing.T) {
	ecc := &cryptography.ECDSAImpl{}

	// Generate ECDSA keys using P256 curve
	privateKey, publicKey, err := ecc.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.NotNil(t, publicKey)
	assert.Equal(t, elliptic.P256(), privateKey.PublicKey.Curve)
	assert.Equal(t, elliptic.P256(), publicKey.Curve)
}

// TestSignVerify tests signing and verifying functionality
func TestSignVerify(t *testing.T) {
	ecc := &cryptography.ECDSAImpl{}

	// Generate ECDSA keys
	privateKey, publicKey, err := ecc.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	// Message to sign
	message := []byte("This is a test message.")

	// Sign the message
	signature, err := ecc.Sign(message, privateKey)
	assert.NoError(t, err)
	assert.NotNil(t, signature)

	// Verify the signature
	valid, err := ecc.Verify(message, signature, publicKey)
	assert.NoError(t, err)
	assert.True(t, valid, "The signature should be valid")

	// Modify the message and try verifying the signature
	modifiedMessage := []byte("This is a modified message.")
	valid, err = ecc.Verify(modifiedMessage, signature, publicKey)
	assert.NoError(t, err)
	assert.False(t, valid, "The signature should not be valid for a modified message")
}

// TestSaveAndReadKeys tests saving and reading the private and public keys from PEM files
func TestSaveAndReadKeys(t *testing.T) {
	ecc := &cryptography.ECDSAImpl{}

	// Generate ECDSA keys
	privateKey, publicKey, err := ecc.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	// Save private and public keys to files
	privateKeyFile := "private.pem"
	publicKeyFile := "public.pem"
	err = ecc.SavePrivateKeyToFile(privateKey, privateKeyFile)
	assert.NoError(t, err)

	err = ecc.SavePublicKeyToFile(publicKey, publicKeyFile)
	assert.NoError(t, err)

	// Read the private and public keys from the files
	readPrivateKey, err := ecc.ReadPrivateKey(privateKeyFile)
	assert.NoError(t, err)
	assert.Equal(t, privateKey.D, readPrivateKey.D)
	assert.Equal(t, privateKey.PublicKey.X, readPrivateKey.PublicKey.X)
	assert.Equal(t, privateKey.PublicKey.Y, readPrivateKey.PublicKey.Y)

	readPublicKey, err := ecc.ReadPublicKey(publicKeyFile)
	assert.NoError(t, err)
	assert.Equal(t, publicKey.X, readPublicKey.X)
	assert.Equal(t, publicKey.Y, readPublicKey.Y)

	// Clean up the generated files
	os.Remove(privateKeyFile)
	os.Remove(publicKeyFile)
}

// TestSaveSignatureToFile tests saving a signature to a file
func TestSaveSignatureToFile(t *testing.T) {
	ecc := &cryptography.ECDSAImpl{}

	// Generate ECDSA keys
	privateKey, _, err := ecc.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	// Message to sign
	message := []byte("This is a test message.")

	// Sign the message
	signature, err := ecc.Sign(message, privateKey)
	assert.NoError(t, err)
	assert.NotNil(t, signature)

	// Save the signature to a file
	signatureFile := "signature.hex"
	err = ecc.SaveSignatureToFile(signatureFile, signature)
	assert.NoError(t, err)

	// Read the saved signature from the file
	hexData, err := ioutil.ReadFile(signatureFile)
	assert.NoError(t, err)

	// Decode the hex signature
	decodedSignature, err := hex.DecodeString(string(hexData))
	assert.NoError(t, err)
	assert.Equal(t, signature, decodedSignature)

	// Clean up the generated signature file
	os.Remove(signatureFile)
}

// TestSignWithInvalidPrivateKey tests signing with an invalid private key
func TestSignWithInvalidPrivateKey(t *testing.T) {
	ecc := &cryptography.ECDSAImpl{}

	// Generate ECDSA keys (valid ones)
	_, _, err := ecc.GenerateKeys(elliptic.P256())
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
	_, err = ecc.Sign(message, invalidPrivateKey)
	assert.Error(t, err, "Signing with an invalid private key should fail")
}

// TestVerifyWithInvalidPublicKey tests verifying with an invalid public key
func TestVerifyWithInvalidPublicKey(t *testing.T) {
	ecc := &cryptography.ECDSAImpl{}

	// Generate ECDSA keys
	privateKey, _, err := ecc.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	// Sign the message
	message := []byte("This is a test message.")
	signature, err := ecc.Sign(message, privateKey)
	assert.NoError(t, err)

	// Create an invalid public key (e.g., public key X = 0)
	invalidPublicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     new(big.Int).SetInt64(0),
		Y:     new(big.Int).SetInt64(0),
	}

	// Attempt to verify the signature with the invalid public key
	valid, err := ecc.Verify(message, signature, invalidPublicKey)
	assert.NoError(t, err)
	assert.False(t, valid, "Verification with an invalid public key should fail")
}
