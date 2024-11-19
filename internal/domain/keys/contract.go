package keys

// Define KeyType as a custom type (based on int)
type KeyType int

// Enum-like values using iota
const (
	AsymmetricPublic KeyType = iota
	AsymmetricPrivate
	Symmetric
)

// KeyManagement defines methods for managing cryptographic key operations.
type KeyManagement interface {
	// Upload handles the upload of blobs from file paths.
	// Returns the created Blobs metadata and any error encountered.
	Upload(filePath []string) ([]*CryptographicKey, error)

	// Download retrieves a cryptographic key by its ID and key type, returning the metadata and key data.
	// Returns the key metadata, key data as a byte slice, and any error.
	Download(keyId string, keyType KeyType) (*CryptographicKey, []byte, error)

	// DeleteByID removes a cryptographic key by its ID.
	// Returns any error encountered.
	DeleteByID(keyId string) error
}

// CryptographicKeyMetadataManagement defines the methods for managing CryptographicKey metadata
type CryptographicKeyMetadataManagement interface {
	// Create creates a new cryptographic key
	Create(key *CryptographicKey) (*CryptographicKey, error)
	// GetByID retrieves cryptographic key by ID
	GetByID(keyID string) (*CryptographicKey, error)
	// UpdateByID updates cryptographic key metadata
	UpdateByID(keyID string, updates *CryptographicKey) (*CryptographicKey, error)
	// DeleteByID deletes a cryptographic key by ID
	DeleteByID(keyID string) error
}

// KeyOperations defines methods for key management, encryption, signing, and PKCS#11 operations.
type KeyOperations interface {

	// ---Key generation---

	// GenerateKey generates keys for specified type and size (e.g., AES, RSA, ECDSA)
	GenerateKey(keyType string, keySize int) ([]byte, error)

	// ---Key storage and retrieval---

	// SaveKey saves a key to a file
	SaveKey(key []byte, filename string) error
	// LoadKey loads a key from a file
	LoadKey(filename string) ([]byte, error)

	// ---Encryption and Decryption (Symmetric algorithms like AES)---

	// EncryptWithSymmetricKey encrypts data with symmetric keys (e.g. AES)
	EncryptWithSymmetricKey(plainText []byte, key []byte) ([]byte, error)
	// DecryptWithSymmetricKey decrypts data with symmetric keys (e.g. AES)
	DecryptWithSymmetricKey(cipherText []byte, key []byte) ([]byte, error)

	// ---Asymmetric Encryption (RSA, ECDSA, PKCS#11)---

	// EncryptWithPublicKey encrypts with public key using asymmetric algorithms (RSA, ECDSA) and optionally a PKCS#11 interface
	EncryptWithPublicKey(plainText []byte, publicKey interface{}) ([]byte, error)
	// DecryptWithPrivateKey decrypt with private key using asymmetric algorithms (RSA, ECDSA) and optionally a PKCS#11 interface
	DecryptWithPrivateKey(cipherText []byte, privateKey interface{}) ([]byte, error)

	// ---Signing and Verification (For RSA, ECDSA)---

	// SignWithPrivateKey signs message with private key using asymmetric algorithms (RSA, ECDSA) and optionally a PKCS#11 interface
	SignWithPrivateKey(message []byte, privateKey interface{}) ([]byte, error)
	// VerifyWithPublicKey verifies signatures with public key using asymmetric algorithms (RSA, ECDSA) and optionally a PKCS#11 interface
	VerifyWithPublicKey(message []byte, signature []byte, publicKey interface{}) (bool, error)

	// ---PKCS#11 Operations---

	// InitializeToken initializes PKCS#11 token in the specified slot
	InitializeToken(slot string) error
	// AddKeyToToken adds key to the PKCS#11 token
	AddKeyToToken() error
	// DeleteKeyFromToken deletes key from PKCS#11 token by type and label
	DeleteKeyFromToken(objectType, objectLabel string) error
}
