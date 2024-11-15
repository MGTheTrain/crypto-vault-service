package cryptography

import (
	"fmt"
	"testing"

	cryptography "crypto_vault_service/internal/infrastructure/cryptography"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PKCS11Test struct {
	Token *cryptography.PKCS11Token
}

// NewPKCS11Test sets up the test environment for PKCS#11 integration tests
func NewPKCS11Test(slot, modulePath, tokenLabel, soPin, userPin, objectLabel, keyType string, keySize int) *PKCS11Test {
	return &PKCS11Test{
		Token: &cryptography.PKCS11Token{
			ModulePath:  modulePath,
			TokenLabel:  tokenLabel,
			SOPin:       soPin,
			UserPin:     userPin,
			ObjectLabel: objectLabel,
			KeyType:     keyType,
			KeySize:     keySize,
		},
	}
}
func (p *PKCS11Test) Setup(t *testing.T) {
	// Initialize the PKCS#11 token
	tokenSlot := "0x1"
	err := p.Token.InitializeToken(tokenSlot)
	require.NoError(t, err, "Failed to initialize PKCS#11 token")

	isTokenSet, err := p.Token.IsTokenSet()
	require.NoError(t, err, "Error checking if token is set")
	assert.True(t, isTokenSet, "The token should be initialized and set")
}

// DeleteKeyFromToken deletes any existing key with the same label before adding a new key.
func (p *PKCS11Test) DeleteKeyFromToken(t *testing.T) {
	err := p.Token.DeleteObject("privkey", p.Token.ObjectLabel)
	if err != nil {
		t.Logf("Warning: Failed to delete existing private key: %v\n", err)
	}

	err = p.Token.DeleteObject("pubkey", p.Token.ObjectLabel)
	if err != nil {
		t.Logf("Warning: Failed to delete existing public key: %v\n", err)
	}

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

	p.DeleteKeyFromToken(t)
}

// TestAddRSAKey tests adding an RSA key to a PKCS#11 token
func TestAddRSAKey(t *testing.T) {
	test := NewPKCS11Test("0x1", "/usr/lib/softhsm/libsofthsm2.so", "my-token", "123456", "234567", "TestRSAKey", "RSA", 2048)

	test.Setup(t)

	// Add an RSA key to the token
	test.AddKeyToToken(t)
}

// TestAddECDSAKey tests adding an ECDSA key to a PKCS#11 token
func TestAddECDSAKey(t *testing.T) {
	test := NewPKCS11Test("0x1", "/usr/lib/softhsm/libsofthsm2.so", "my-token", "123456", "234567", "TestECDSAKey", "ECDSA", 256)

	test.Setup(t)

	// Add an ECDSA key to the token
	test.AddKeyToToken(t)
}
