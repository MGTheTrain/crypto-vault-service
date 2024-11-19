package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// AES Interface
type AES interface {
	Encrypt(plainText, key []byte) ([]byte, error)
	Decrypt(ciphertext, key []byte) ([]byte, error)
	GenerateKey(keySize int) ([]byte, error)
}

// AESImpl struct that implements the AES interface
type AESImpl struct{}

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
func (a *AESImpl) GenerateKey(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %v", err)
	}
	return key, nil
}

// Encrypt data using AES in CBC mode
func (a *AESImpl) Encrypt(plainText, key []byte) ([]byte, error) {
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

	return ciphertext, nil
}

// Decrypt data using AES in CBC mode
func (a *AESImpl) Decrypt(ciphertext, key []byte) ([]byte, error) {
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

	return pkcs7Unpad(ciphertext, aes.BlockSize)
}
