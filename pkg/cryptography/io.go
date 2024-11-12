package cryptography

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

// File Operations
func ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func WriteFile(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0644)
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
func ReadPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	pubKeyPEM, err := ioutil.ReadFile(publicKeyPath)
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
