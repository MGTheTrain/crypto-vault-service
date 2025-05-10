//go:build unit
// +build unit

package cryptography

import (
	"crypto/rsa"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RSAProcessorTests encapsulates RSAProcessor test cases
type RSAProcessorTests struct {
	processor RSAProcessor
}

// NewRSAProcessorTests creates a new instance of RSAProcessorTests
func NewRSAProcessorTests(t *testing.T) *RSAProcessorTests {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logInstance, err := logger.GetLogger(loggerSettings)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	processor, err := NewRSAProcessor(logInstance)
	if err != nil {
		t.Fatalf("Failed to create RSA processor: %v", err)
	}

	return &RSAProcessorTests{
		processor: processor,
	}
}

func (rt *RSAProcessorTests) TestGenerateRSAKeys(t *testing.T) {
	privateKey, publicKey, err := rt.processor.GenerateKeys(2048)
	assert.NoError(t, err)
	assert.NotNil(t, privateKey)
	assert.NotNil(t, publicKey)
	assert.IsType(t, &rsa.PublicKey{}, publicKey)
	assert.Equal(t, 2048, privateKey.N.BitLen())
}

func (rt *RSAProcessorTests) TestEncryptDecrypt(t *testing.T) {
	privateKey, publicKey, err := rt.processor.GenerateKeys(2048)
	assert.NoError(t, err)

	plainText := []byte("This is a secret message")
	encrypted, err := rt.processor.Encrypt(plainText, publicKey)
	assert.NoError(t, err)
	decrypted, err := rt.processor.Decrypt(encrypted, privateKey)
	assert.NoError(t, err)
	assert.Equal(t, plainText, decrypted)
}

func (rt *RSAProcessorTests) TestSaveAndReadKeys(t *testing.T) {
	privateKey, publicKey, err := rt.processor.GenerateKeys(2048)
	assert.NoError(t, err)

	privFile := "private.pem"
	pubFile := "public.pem"

	assert.NoError(t, rt.processor.SavePrivateKeyToFile(privateKey, privFile))
	assert.NoError(t, rt.processor.SavePublicKeyToFile(publicKey, pubFile))

	readPriv, err := rt.processor.ReadPrivateKey(privFile)
	assert.NoError(t, err)
	assert.Equal(t, privateKey.N, readPriv.N)
	assert.Equal(t, privateKey.E, readPriv.E)

	readPub, err := rt.processor.ReadPublicKey(pubFile)
	assert.NoError(t, err)
	assert.Equal(t, publicKey.N, readPub.N)
	assert.Equal(t, publicKey.E, readPub.E)

	os.Remove(privFile)
	os.Remove(pubFile)
}

func (rt *RSAProcessorTests) TestEncryptWithInvalidKey(t *testing.T) {
	_, publicKey, err := rt.processor.GenerateKeys(2048)
	assert.NoError(t, err)

	plainText := []byte("This should fail encryption")
	encrypted, err := rt.processor.Encrypt(plainText, publicKey)
	assert.NoError(t, err)

	wrongPrivKey, _, err := rt.processor.GenerateKeys(2048)
	assert.NoError(t, err)

	_, err = rt.processor.Decrypt(encrypted, wrongPrivKey)
	assert.Error(t, err)
}

func (rt *RSAProcessorTests) TestSavePrivateKeyInvalidPath(t *testing.T) {
	privateKey, _, err := rt.processor.GenerateKeys(2048)
	assert.NoError(t, err)

	err = rt.processor.SavePrivateKeyToFile(privateKey, "/invalid/path/private.pem")
	assert.Error(t, err)
}

func (rt *RSAProcessorTests) TestSavePublicKeyInvalidPath(t *testing.T) {
	_, publicKey, err := rt.processor.GenerateKeys(2048)
	assert.NoError(t, err)

	err = rt.processor.SavePublicKeyToFile(publicKey, "/invalid/path/public.pem")
	assert.Error(t, err)
}

func (rt *RSAProcessorTests) TestSignAndVerify(t *testing.T) {
	privateKey, publicKey, err := rt.processor.GenerateKeys(2048)
	assert.NoError(t, err)

	data := []byte("This is a test message")
	signature, err := rt.processor.Sign(data, privateKey)
	assert.NoError(t, err)
	assert.NotNil(t, signature)

	valid, err := rt.processor.Verify(data, signature, publicKey)
	assert.NoError(t, err)
	assert.True(t, valid)

	tampered := []byte("This is a tampered message")
	valid, err = rt.processor.Verify(tampered, signature, publicKey)
	assert.Error(t, err)
	assert.False(t, valid)
}

func TestRSAProcessor(t *testing.T) {
	rt := NewRSAProcessorTests(t)

	t.Run("TestGenerateRSAKeys", rt.TestGenerateRSAKeys)
	t.Run("TestEncryptDecrypt", rt.TestEncryptDecrypt)
	t.Run("TestSaveAndReadKeys", rt.TestSaveAndReadKeys)
	t.Run("TestEncryptWithInvalidKey", rt.TestEncryptWithInvalidKey)
	t.Run("TestSavePrivateKeyInvalidPath", rt.TestSavePrivateKeyInvalidPath)
	t.Run("TestSavePublicKeyInvalidPath", rt.TestSavePublicKeyInvalidPath)
	t.Run("TestSignAndVerify", rt.TestSignAndVerify)
}
