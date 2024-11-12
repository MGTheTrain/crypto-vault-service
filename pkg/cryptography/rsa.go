package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

// RSA Interface
type RSA interface {
	Encrypt(plainText []byte, publicKey *rsa.PublicKey) ([]byte, error)
	Decrypt(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error)
	GenerateKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error)
}

// RSAImpl struct that implements the RSA interface
type RSAImpl struct{}

// GenerateRSAKeys generates RSA key pair
func (r *RSAImpl) GenerateKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA keys: %v", err)
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

// Encrypt data using RSA public key
func (r *RSAImpl) Encrypt(plainText []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return encryptedData, nil
}

// Decrypt data using RSA private key
func (r *RSAImpl) Decrypt(ciphertext []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decryptedData, nil
}
