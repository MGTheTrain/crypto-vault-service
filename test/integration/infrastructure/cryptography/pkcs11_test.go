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
func NewPKCS11Test(modulePath, tokenLabel, soPin, userPin, objectLabel, keyType string, keySize int) *PKCS11Test {
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

// Setup initializes the PKCS#11 token, sets the token label, and generates the keys.
func (p *PKCS11Test) Setup(t *testing.T) {
	// Initialize the PKCS#11 token
	err := p.Token.InitializeToken()
	require.NoError(t, err, "Failed to initialize PKCS#11 token")

	// Check if the token is properly initialized
	isTokenSet, err := p.Token.IsTokenSet()
	require.NoError(t, err, "Error checking if token is set")
	assert.True(t, isTokenSet, "The token should be initialized and set")
}

// AddKeyToToken is a helper function to add a key to a token
func (p *PKCS11Test) AddKeyToToken(t *testing.T) {
	var err error
	err = p.Token.AddKey()

	assert.NoError(t, err, "Failed to add key to the token")

	// Verify the key was added to the token
	isObjectSet, err := p.Token.IsObjectSet()
	assert.NoError(t, err, "Error checking if key is set")
	assert.True(t, isObjectSet, fmt.Sprintf("The %s key should be added to the token", p.Token.KeyType))
}

// TestAddRSAKey tests adding an RSA key to a PKCS#11 token
func TestAddRSAKey(t *testing.T) {
	test := NewPKCS11Test("/usr/lib/softhsm/libsofthsm2.so", "TestToken", "123456", "234567", "TestRSAKey", "RSA", 2048)

	// Set up the token
	test.Setup(t)

	// Add an RSA key to the token
	test.AddKeyToToken(t)
}

// TestAddECDSAKey tests adding an ECDSA key to a PKCS#11 token
func TestAddECDSAKey(t *testing.T) {
	test := NewPKCS11Test("/usr/lib/softhsm/libsofthsm2.so", "TestToken", "123456", "234567", "TestECDSAKey", "ECDSA", 256)

	// Set up the token
	test.Setup(t)

	// Add an ECDSA key to the token
	test.AddKeyToToken(t)
}
