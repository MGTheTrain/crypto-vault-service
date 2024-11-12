package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

// AES Functions
// GenerateRandomAESKey generates a random AES key of the specified size
func GenerateRandomAESKey(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %v", err)
	}
	return key, nil
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

// Encrypts data using AES in CBC mode
func EncryptAES(plainText, key []byte) ([]byte, error) {
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

// Decrypts data using AES in CBC mode
func DecryptAES(ciphertext, key []byte) ([]byte, error) {
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

// RSA Functions
func GenerateRSAKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA keys: %v", err)
	}

	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

func SavePrivateKeyToFile(privateKey *rsa.PrivateKey, filename string) error {
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

	return nil
}

func SavePublicKeyToFile(publicKey *rsa.PublicKey, filename string) error {
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

	return nil
}

// Encrypt data with RSA public key
func EncryptWithRSA(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return encryptedData, nil
}

// Decrypt data with RSA private key
func DecryptWithRSA(privateKey *rsa.PrivateKey, encryptedData []byte) ([]byte, error) {
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decryptedData, nil
}

// File Operations
func ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func WriteFile(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0644)
}

// AES Command

// Read RSA private key from PEM file
func ReadPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	privKeyPEM, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key file: %v", err)
	}

	block, _ := pem.Decode(privKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	return privateKey, nil
}

// Read RSA public key from PEM file
func ReadPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	pubKeyPEM, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key file: %v", err)
	}

	block, _ := pem.Decode(pubKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse public key: %v", err)
	}

	return publicKey, nil
}
