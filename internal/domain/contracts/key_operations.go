package contracts

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
