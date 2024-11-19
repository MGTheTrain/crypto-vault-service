package keys

// KeyType defines a custom type for key types based on an integer.
type KeyType int

// Enum-like values for different key types, using iota to generate sequential values.
const (
	AsymmetricPublic  KeyType = iota // Public key in asymmetric cryptography (e.g., RSA, ECDSA)
	AsymmetricPrivate                // Private key in asymmetric cryptography (e.g., RSA, ECDSA)
	Symmetric                        // Symmetric key (e.g., AES)
)

// CryptKeyUploadService defines methods for uploading cryptographic keys.
type CryptKeyUploadService interface {
	// Upload uploads cryptographic keys from specified file paths.
	// It returns a slice of CryptoKeyMeta and any error encountered during the upload process.
	Upload(filePaths []string) ([]*CryptoKeyMeta, error)
}

// CryptoKeyMetadataService defines methods for managing cryptographic key metadata and deleting keys.
type CryptoKeyMetadataService interface {
	// List retrieves metadata for all cryptographic keys.
	// It returns a slice of CryptoKeyMeta and any error encountered during the retrieval process.
	List() ([]*CryptoKeyMeta, error)

	// GetByID retrieves the metadata of a cryptographic key by its unique ID.
	// It returns the CryptoKeyMeta and any error encountered during the retrieval process.
	GetByID(keyID string) (*CryptoKeyMeta, error)

	// DeleteByID deletes a cryptographic key and its associated metadata by ID.
	// It returns any error encountered during the deletion process.
	DeleteByID(keyID string) error
}

// CryptoKeyDownloadService defines methods for downloading cryptographic keys.
type CryptoKeyDownloadService interface {
	// Download retrieves a cryptographic key by its ID and type.
	// It returns the CryptoKeyMeta, the key data as a byte slice, and any error encountered during the download process.
	Download(keyID string, keyType KeyType) (*CryptoKeyMeta, []byte, error)
}

// KeyOperations defines methods for cryptographic key management, encryption, signing, and PKCS#11 operations.
type CryptoKeyOperationService interface {

	// --- Key Generation ---

	// GenerateKey generates a cryptographic key of the specified type and size (e.g., AES, RSA, ECDSA).
	// It returns the generated key as a byte slice and any error encountered during the key generation.
	GenerateKey(keyType string, keySize int) ([]byte, error)

	// --- Key Storage and Retrieval ---

	// SaveKey saves a cryptographic key to a specified file.
	// It returns any error encountered during the saving process.
	SaveKey(key []byte, filename string) error

	// LoadKey loads a cryptographic key from a specified file.
	// It returns the loaded key as a byte slice and any error encountered during the loading process.
	LoadKey(filename string) ([]byte, error)

	// --- Encryption and Decryption (Symmetric algorithms like AES) ---

	// EncryptWithSymmetricKey encrypts data using a symmetric key (e.g., AES).
	// It returns the encrypted data as a byte slice and any error encountered during encryption.
	EncryptWithSymmetricKey(plainText []byte, key []byte) ([]byte, error)

	// DecryptWithSymmetricKey decrypts data using a symmetric key (e.g., AES).
	// It returns the decrypted data as a byte slice and any error encountered during decryption.
	DecryptWithSymmetricKey(cipherText []byte, key []byte) ([]byte, error)

	// --- Asymmetric Encryption (RSA, ECDSA, PKCS#11) ---

	// EncryptWithPublicKey encrypts data with a public key using asymmetric encryption algorithms (e.g., RSA, ECDSA).
	// It optionally supports PKCS#11 hardware tokens for key storage.
	// It returns the encrypted data as a byte slice and any error encountered during encryption.
	EncryptWithPublicKey(plainText []byte, publicKey interface{}) ([]byte, error)

	// DecryptWithPrivateKey decrypts data with a private key using asymmetric encryption algorithms (e.g., RSA, ECDSA).
	// It optionally supports PKCS#11 hardware tokens for key storage.
	// It returns the decrypted data as a byte slice and any error encountered during decryption.
	DecryptWithPrivateKey(cipherText []byte, privateKey interface{}) ([]byte, error)

	// --- Signing and Verification (For RSA, ECDSA) ---

	// SignWithPrivateKey signs a message using a private key with asymmetric algorithms (e.g., RSA, ECDSA).
	// It optionally supports PKCS#11 hardware tokens for key storage.
	// It returns the signature and any error encountered during the signing process.
	SignWithPrivateKey(message []byte, privateKey interface{}) ([]byte, error)

	// VerifyWithPublicKey verifies a signature using a public key with asymmetric algorithms (e.g., RSA, ECDSA).
	// It optionally supports PKCS#11 hardware tokens for key storage.
	// It returns true if the signature is valid, false otherwise, and any error encountered during the verification process.
	VerifyWithPublicKey(message []byte, signature []byte, publicKey interface{}) (bool, error)

	// --- PKCS#11 Operations ---

	// InitializeToken initializes a PKCS#11 token in the specified hardware slot.
	// It returns any error encountered during the initialization.
	InitializeToken(slot string) error

	// AddKeyToToken adds a cryptographic key to a PKCS#11 token.
	// It returns any error encountered during the addition of the key.
	AddKeyToToken() error

	// DeleteKeyFromToken deletes a cryptographic key from a PKCS#11 token by its type and label.
	// It returns any error encountered during the deletion of the key.
	DeleteKeyFromToken(objectType, objectLabel string) error
}
