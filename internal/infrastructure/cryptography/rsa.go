package cryptography

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto_vault_service/internal/infrastructure/logger"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// IRSA Interface
type IRSA interface {
	Encrypt(plainText []byte, publicKey *rsa.PublicKey) ([]byte, error)
	Decrypt(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error)
	Sign(data []byte, privateKey *rsa.PrivateKey) ([]byte, error)
	Verify(data []byte, signature []byte, publicKey *rsa.PublicKey) (bool, error)
	GenerateKeys(keySize int) (*rsa.PrivateKey, *rsa.PublicKey, error)
	SavePrivateKeyToFile(privateKey *rsa.PrivateKey, filename string) error
	SavePublicKeyToFile(publicKey *rsa.PublicKey, filename string) error
	ReadPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error)
	ReadPublicKey(publicKeyPath string) (*rsa.PublicKey, error)
}

// RSA struct that implements the IRSA interface
type RSA struct {
	logger logger.Logger
}

// NewRSA creates and returns a new instance of RSA
func NewRSA(logger logger.Logger) (*RSA, error) {
	return &RSA{
		logger: logger,
	}, nil
}

// GenerateRSAKeys generates RSA key pair
func (r *RSA) GenerateKeys(keySize int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA keys: %v", err)
	}
	publicKey := &privateKey.PublicKey
	r.logger.Info("Generated RSA key pairs")
	return privateKey, publicKey, nil
}

// Encrypt data using RSA public key
// Encrypt encrypts the given plaintext using RSA encryption.
// If the plaintext is too large, it will split it into smaller chunks and encrypt each one separately.
func (r *RSA) Encrypt(plainText []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("public key cannot be nil")
	}

	// Maximum size for the plaintext that can be encrypted with the RSA key
	// For a 2048-bit RSA key, it's approximately 245 bytes after accounting for padding
	maxSize := publicKey.Size() - 11 // PKCS#1 v1.5 padding size

	// If the plaintext is too large, split it into smaller chunks
	var encryptedData []byte
	for len(plainText) > 0 {
		// Determine the chunk size
		chunkSize := maxSize
		if len(plainText) < chunkSize {
			chunkSize = len(plainText)
		}

		// Encrypt the current chunk
		encryptedChunk, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText[:chunkSize])
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt data: %v", err)
		}

		// Append the encrypted chunk to the result
		encryptedData = append(encryptedData, encryptedChunk...)

		// Move to the next chunk
		plainText = plainText[chunkSize:]
	}

	r.logger.Info("RSA encryption succeeded")
	return encryptedData, nil
}

// Decrypt data using RSA private key. It handles multiple chunks of encrypted data.
func (r *RSA) Decrypt(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("private key cannot be nil")
	}

	// Maximum size for the decrypted data, which is the RSA key size
	maxSize := privateKey.Size()

	var decryptedData []byte
	for len(ciphertext) > 0 {
		// Determine the chunk size
		chunkSize := maxSize
		if len(ciphertext) < chunkSize {
			chunkSize = len(ciphertext)
		}

		// Decrypt the current chunk
		decryptedChunk, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext[:chunkSize])
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt data: %v", err)
		}

		// Append the decrypted chunk to the result
		decryptedData = append(decryptedData, decryptedChunk...)

		// Move to the next chunk
		ciphertext = ciphertext[chunkSize:]
	}

	r.logger.Info("RSA decryption succeeded")
	return decryptedData, nil
}

// Sign data using RSA private key
func (r *RSA) Sign(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("private key cannot be nil")
	}

	// Use the SHA-256 hash algorithm for signing
	hashed := sha256.Sum256(data)

	// Sign the hashed data with the private key
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign data: %v", err)
	}

	r.logger.Info("RSA signing succeeded")
	return signature, nil
}

// Verify RSA signature with public key
func (r *RSA) Verify(data []byte, signature []byte, publicKey *rsa.PublicKey) (bool, error) {
	if publicKey == nil {
		return false, fmt.Errorf("public key cannot be nil")
	}

	// Use the SHA-256 hash algorithm for verification
	hashed := sha256.Sum256(data)

	// Verify the signature using the public key
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return false, fmt.Errorf("failed to verify signature: %v", err)
	}

	r.logger.Info("RSA signature verified successfully")
	return true, nil
}

// SavePrivateKeyToFile saves the private key to a PEM file using encoding/pem
func (r *RSA) SavePrivateKeyToFile(privateKey *rsa.PrivateKey, filename string) error {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privKeyPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer file.Close()

	err = pem.Encode(file, privKeyPem)
	if err != nil {
		return fmt.Errorf("failed to encode private key: %v", err)
	}

	r.logger.Info(fmt.Sprintf("Saved RSA private key %s", filename))
	return nil
}

// SavePublicKeyToFile saves the public key to a PEM file using encoding/pem
func (r *RSA) SavePublicKeyToFile(publicKey *rsa.PublicKey, filename string) error {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}

	pubKeyPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %v", err)
	}
	defer file.Close()

	err = pem.Encode(file, pubKeyPem)
	if err != nil {
		return fmt.Errorf("failed to encode public key: %v", err)
	}

	r.logger.Info(fmt.Sprintf("Saved RSA public key %s", filename))

	return nil
}

// Read RSA private key from PEM file
func (r *RSA) ReadPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	privKeyPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key file: %v", err)
	}

	block, _ := pem.Decode(privKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	// First try to parse as PKCS#1 format
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privateKey, nil
	}

	// If PKCS#1 parsing fails, try parsing as PKCS#8 format
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key in either PKCS#1 or PKCS#8 format: %v", err)
	}

	// Type assertion to *rsa.PrivateKey if it is indeed an RSA key
	privateKey, ok := privateKeyInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not of type RSA")
	}

	return privateKey, nil
}

// Read RSA public key from PEM file
func (r *RSA) ReadPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	pubKeyPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key file: %v", err)
	}

	block, _ := pem.Decode(pubKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	// Try to parse as PKCS#1 format first
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return publicKey, nil
	}

	// If PKCS#1 parsing fails, try parsing as PKCS#8 format
	pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse public key in either PKCS#1 or PKCS#8 format: %v", err)
	}

	// Type assertion to *rsa.PublicKey if it is indeed an RSA key
	publicKey, ok := pubKeyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type RSA")
	}

	return publicKey, nil
}
