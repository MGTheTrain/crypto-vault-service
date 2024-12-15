package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto_vault_service/internal/infrastructure/logger"
	"fmt"
)

// IAES Interface
type IAES interface {
	Encrypt(plainText, key []byte) ([]byte, error)
	Decrypt(ciphertext, key []byte) ([]byte, error)
	GenerateKey(keySize int) ([]byte, error)
}

// AES struct that implements the IAES interface
type AES struct {
	Logger logger.Logger
}

// NewAES creates and returns a new instance of AES
func NewAES(logger logger.Logger) (*AES, error) {
	return &AES{
		Logger: logger,
	}, nil
}

// Pad data to make it a multiple of AES block size
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padded := append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
	return padded
}

// Unpad data after decryption
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	padding := int(data[length-1])

	if padding > length || padding > blockSize {
		return nil, fmt.Errorf("invalid padding size")
	}
	return data[:length-padding], nil
}

// GenerateRandomAESKey generates a random AES key of the specified size
func (a *AES) GenerateKey(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %v", err)
	}

	a.Logger.Info("Generated AES key pairs.")
	return key, nil
}

// Encrypt data using AES in CBC mode
func (a *AES) Encrypt(plainText, key []byte) ([]byte, error) {
	if key == nil {
		return nil, fmt.Errorf("key key cannot be nil")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plainText = pkcs7Pad(plainText, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(plainText))
	iv := ciphertext[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainText)

	a.Logger.Info("AES encryption succeeded.")
	return ciphertext, nil
}

// Decrypt data using AES in CBC mode
func (a *AES) Decrypt(ciphertext, key []byte) ([]byte, error) {
	if key == nil {
		return nil, fmt.Errorf("key key cannot be nil")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	a.Logger.Info("AES decryption succeeded.")
	return pkcs7Unpad(ciphertext, aes.BlockSize)
}
