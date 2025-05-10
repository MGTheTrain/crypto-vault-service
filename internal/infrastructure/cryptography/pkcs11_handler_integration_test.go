//go:build integration
// +build integration

package cryptography

import (
	"os"
	"testing"

	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	SlotId     = "0x0"
	ModulePath = "/usr/lib/softhsm/libsofthsm2.so"
	Label      = "MyToken"
	SOPin      = "123456"
	UserPin    = "234567"
)

type PKCS11HandlerTests struct {
	objectLabel   string
	pkcs11Handler PKCS11Handler
}

func NewPKCS11HandlerTests(t *testing.T, objectLabel string) *PKCS11HandlerTests {
	pkcs11Settings := &settings.PKCS11Settings{
		ModulePath: ModulePath,
		SOPin:      SOPin,
		UserPin:    UserPin,
		SlotId:     SlotId,
	}

	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
	}

	logInstance, err := logger.GetLogger(loggerSettings)
	require.NoError(t, err, "Failed to initialize logger")

	handler, err := NewPKCS11Handler(pkcs11Settings, logInstance)
	require.NoError(t, err, "Failed to initialize PKCS#11 handler")

	return &PKCS11HandlerTests{
		objectLabel:   objectLabel,
		pkcs11Handler: handler,
	}
}

func (p *PKCS11HandlerTests) InitializeToken(t *testing.T) {
	err := p.pkcs11Handler.InitializeToken(Label)
	require.NoError(t, err, "Failed to initialize PKCS#11 token")
}

func (p *PKCS11HandlerTests) DeleteKeyFromToken(t *testing.T) {
	for _, objType := range []string{"privkey", "pubkey", "secrkey"} {
		err := p.pkcs11Handler.DeleteObject(Label, objType, p.objectLabel)
		if err != nil {
			t.Logf("Warning: Failed to delete existing %s: %v", objType, err)
		}
	}
}

func (p *PKCS11HandlerTests) AddKeyToToken(t *testing.T, keyType string, keySize uint) {
	err := p.pkcs11Handler.AddKey(Label, p.objectLabel, keyType, keySize)
	assert.NoError(t, err, "Failed to add key to token")
}

// ---------- Tests ----------

func TestListTokens(t *testing.T) {
	test := NewPKCS11HandlerTests(t, "TestRSAKey")
	test.InitializeToken(t)

	tokens, err := test.pkcs11Handler.ListTokenSlots()
	require.NoError(t, err)
	require.NotEmpty(t, tokens)

	token := tokens[0]
	assert.NotEmpty(t, token.SlotId)
	assert.NotEmpty(t, token.Label)
	assert.NotEmpty(t, token.Manufacturer)
	assert.NotEmpty(t, token.Model)
	assert.NotEmpty(t, token.SerialNumber)
}

func TestAddRSAKey(t *testing.T) {
	test := NewPKCS11HandlerTests(t, "TestRSAKey")
	test.InitializeToken(t)
	test.AddKeyToToken(t, "RSA", 2048)
	test.DeleteKeyFromToken(t)
}

func TestAddECDSAKey(t *testing.T) {
	test := NewPKCS11HandlerTests(t, "TestECDSAKey")
	test.InitializeToken(t)
	test.AddKeyToToken(t, "ECDSA", 256)
	test.DeleteKeyFromToken(t)
}

func TestListObjects(t *testing.T) {
	test := NewPKCS11HandlerTests(t, "TestRSAKey2")
	test.InitializeToken(t)
	test.AddKeyToToken(t, "RSA", 2048)

	objects, err := test.pkcs11Handler.ListObjects(Label)
	require.NoError(t, err)
	require.NotEmpty(t, objects)

	object := objects[0]
	assert.NotEmpty(t, object.Label)
	assert.NotEmpty(t, object.Type)
	assert.NotEmpty(t, object.Usage)

	test.DeleteKeyFromToken(t)
}

func TestEncryptDecrypt(t *testing.T) {
	test := NewPKCS11HandlerTests(t, "TestRSAKey")
	test.InitializeToken(t)
	test.AddKeyToToken(t, "RSA", 2048)

	inputFile := "plain-text.txt"
	err := os.WriteFile(inputFile, []byte("This is some data to encrypt."), 0644)
	require.NoError(t, err)

	encryptedFile := "encrypted.bin"
	decryptedFile := "decrypted.txt"

	err = test.pkcs11Handler.Encrypt(Label, test.objectLabel, inputFile, encryptedFile, "RSA")
	assert.NoError(t, err)

	encryptedData, err := os.ReadFile(encryptedFile)
	require.NoError(t, err)
	assert.NotEmpty(t, encryptedData)

	err = test.pkcs11Handler.Decrypt(Label, test.objectLabel, encryptedFile, decryptedFile, "RSA")
	assert.NoError(t, err)

	decryptedData, err := os.ReadFile(decryptedFile)
	require.NoError(t, err)

	originalData, err := os.ReadFile(inputFile)
	require.NoError(t, err)
	assert.Equal(t, originalData, decryptedData)

	test.DeleteKeyFromToken(t)
	os.Remove(inputFile)
	os.Remove(encryptedFile)
	os.Remove(decryptedFile)
}

func TestSignAndVerify(t *testing.T) {
	test := NewPKCS11HandlerTests(t, "TestRSAKey")
	test.InitializeToken(t)
	test.AddKeyToToken(t, "RSA", 2048)

	dataFile := "data-to-sign.txt"
	sigFile := "data.sig"
	err := os.WriteFile(dataFile, []byte("This is some data to sign."), 0644)
	require.NoError(t, err)

	err = test.pkcs11Handler.Sign(Label, test.objectLabel, dataFile, sigFile, "RSA")
	assert.NoError(t, err)

	sigData, err := os.ReadFile(sigFile)
	require.NoError(t, err)
	assert.NotEmpty(t, sigData)

	valid, err := test.pkcs11Handler.Verify(Label, test.objectLabel, dataFile, sigFile, "RSA")
	assert.NoError(t, err)
	assert.True(t, valid)

	test.DeleteKeyFromToken(t)
	os.Remove(dataFile)
	os.Remove(sigFile)
}
