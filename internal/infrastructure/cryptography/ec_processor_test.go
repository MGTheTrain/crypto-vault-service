//go:build unit
// +build unit

package cryptography

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"encoding/hex"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ecProcessorTests struct {
	processor ECProcessor
}

// NewECProcessorTests initializes the ECProcessor test suite
func NewECProcessorTests(t *testing.T) *ecProcessorTests {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logInstance, err := logger.GetLogger(loggerSettings)
	if err != nil {
		t.Fatalf("Error creating logger: %v", err)
	}

	ec, err := NewECProcessor(logInstance)
	if err != nil {
		t.Fatalf("Error creating EC processor: %v", err)
	}

	return &ecProcessorTests{
		processor: ec,
	}
}

func (et *ecProcessorTests) TestGenerateKeys(t *testing.T) {
	priv, pub, err := et.processor.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)
	assert.NotNil(t, priv)
	assert.NotNil(t, pub)
	assert.Equal(t, elliptic.P256(), priv.PublicKey.Curve)
	assert.Equal(t, elliptic.P256(), pub.Curve)
}

func (et *ecProcessorTests) TestSignVerify(t *testing.T) {
	priv, pub, err := et.processor.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	msg := []byte("This is a test message.")
	sig, err := et.processor.Sign(msg, priv)
	assert.NoError(t, err)
	assert.NotNil(t, sig)

	valid, err := et.processor.Verify(msg, sig, pub)
	assert.NoError(t, err)
	assert.True(t, valid)

	invalidMsg := []byte("Modified message.")
	valid, err = et.processor.Verify(invalidMsg, sig, pub)
	assert.NoError(t, err)
	assert.False(t, valid)
}

func (et *ecProcessorTests) TestSaveAndReadKeys(t *testing.T) {
	priv, pub, err := et.processor.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	privateFile := "private_test.pem"
	publicFile := "public_test.pem"

	err = et.processor.SavePrivateKeyToFile(priv, privateFile)
	assert.NoError(t, err)

	err = et.processor.SavePublicKeyToFile(pub, publicFile)
	assert.NoError(t, err)

	readPriv, err := et.processor.ReadPrivateKey(privateFile, elliptic.P256())
	assert.NoError(t, err)
	assert.Equal(t, priv.D, readPriv.D)
	assert.Equal(t, priv.PublicKey.X, readPriv.PublicKey.X)
	assert.Equal(t, priv.PublicKey.Y, readPriv.PublicKey.Y)

	readPub, err := et.processor.ReadPublicKey(publicFile, elliptic.P256())
	assert.NoError(t, err)
	assert.Equal(t, pub.X, readPub.X)
	assert.Equal(t, pub.Y, readPub.Y)

	os.Remove(privateFile)
	os.Remove(publicFile)
}

func (et *ecProcessorTests) TestSaveSignatureToFile(t *testing.T) {
	priv, _, err := et.processor.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	msg := []byte("This is a test message.")
	sig, err := et.processor.Sign(msg, priv)
	assert.NoError(t, err)

	sigFile := "signature_test.hex"
	err = et.processor.SaveSignatureToFile(sigFile, sig)
	assert.NoError(t, err)

	hexData, err := os.ReadFile(sigFile)
	assert.NoError(t, err)

	decoded, err := hex.DecodeString(string(hexData))
	assert.NoError(t, err)
	assert.Equal(t, sig, decoded)

	os.Remove(sigFile)
}

func (et *ecProcessorTests) TestSignWithInvalidPrivateKey(t *testing.T) {
	invalidPriv := &ecdsa.PrivateKey{
		D: new(big.Int).SetInt64(0),
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
		},
	}
	_, err := et.processor.Sign([]byte("Invalid signing"), invalidPriv)
	assert.Error(t, err)
}

func (et *ecProcessorTests) TestVerifyWithInvalidPublicKey(t *testing.T) {
	priv, _, err := et.processor.GenerateKeys(elliptic.P256())
	assert.NoError(t, err)

	msg := []byte("Test message")
	sig, err := et.processor.Sign(msg, priv)
	assert.NoError(t, err)

	invalidPub := &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     big.NewInt(0),
		Y:     big.NewInt(0),
	}

	valid, err := et.processor.Verify(msg, sig, invalidPub)
	assert.NoError(t, err)
	assert.False(t, valid)
}

// TestECDSA runs all ECProcessor tests
func TestECDSA(t *testing.T) {
	suite := NewECProcessorTests(t)

	t.Run("GenerateKeys", suite.TestGenerateKeys)
	t.Run("SignVerify", suite.TestSignVerify)
	t.Run("SaveAndReadKeys", suite.TestSaveAndReadKeys)
	t.Run("SaveSignatureToFile", suite.TestSaveSignatureToFile)
	t.Run("SignWithInvalidPrivateKey", suite.TestSignWithInvalidPrivateKey)
	t.Run("VerifyWithInvalidPublicKey", suite.TestVerifyWithInvalidPublicKey)
}
