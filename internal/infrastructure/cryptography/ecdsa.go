package cryptography

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
)

// ECDSA Interface
type ECDSA interface {
	GenerateKeys(curve elliptic.Curve) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error)
	Sign(message []byte, privateKey *ecdsa.PrivateKey) ([]byte, error)
	Verify(message, signature []byte, publicKey *ecdsa.PublicKey) (bool, error)
	SaveSignatureToFile(filename string, data []byte) error
	SavePrivateKeyToFile(privateKey *ecdsa.PrivateKey, filename string) error
	SavePublicKeyToFile(publicKey *ecdsa.PublicKey, filename string) error
	ReadPrivateKey(privateKeyPath string) (*ecdsa.PrivateKey, error)
	ReadPublicKey(publicKeyPath string) (*ecdsa.PublicKey, error)
}

// ECDSAImpl struct that implements the ECDSA interface
type ECDSAImpl struct{}

// GenerateKeys generates an elliptic curve key pair
func (e *ECDSAImpl) GenerateKeys(curve elliptic.Curve) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate elliptic curve keys: %v", err)
	}

	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

// Sign signs a message with the private key
func (e *ECDSAImpl) Sign(message []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	// Check if the private key is valid (D should not be zero)
	if privateKey.D.Sign() == 0 {
		return nil, fmt.Errorf("invalid private key: D cannot be zero")
	}

	// Hash the message before signing it
	hash := sha256.Sum256(message)
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %v", err)
	}

	// Encode the signature as r and s
	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

// Verify verifies the signature of a message with the public key
func (e *ECDSAImpl) Verify(message, signature []byte, publicKey *ecdsa.PublicKey) (bool, error) {
	// Hash the message before verifying it
	hash := sha256.Sum256(message)

	// Split the signature into r and s
	r, s := signature[:len(signature)/2], signature[len(signature)/2:]
	rInt := new(big.Int).SetBytes(r)
	sInt := new(big.Int).SetBytes(s)

	// Verify the signature
	valid := ecdsa.Verify(publicKey, hash[:], rInt, sInt)
	return valid, nil
}

// SavePrivateKeyToFile saves the private key to a PEM file using encoding/pem
func (e *ECDSAImpl) SavePrivateKeyToFile(privateKey *ecdsa.PrivateKey, filename string) error {
	// Marshal private key components (private key 'D' and public key components 'X' and 'Y')
	privKeyBytes := append(privateKey.D.Bytes(), privateKey.PublicKey.X.Bytes()...)
	privKeyBytes = append(privKeyBytes, privateKey.PublicKey.Y.Bytes()...)

	// Prepare the PEM block
	privKeyPem := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privKeyBytes,
	}

	// Write the PEM block to a file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer file.Close()

	// Encode and write the private key in PEM format
	err = pem.Encode(file, privKeyPem)
	if err != nil {
		return fmt.Errorf("failed to encode private key: %v", err)
	}

	return nil
}

// SavePublicKeyToFile saves the public key to a PEM file using encoding/pem
func (e *ECDSAImpl) SavePublicKeyToFile(publicKey *ecdsa.PublicKey, filename string) error {
	pubKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)

	// Prepare the PEM block for the public key
	pubKeyPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}

	// Write the PEM block to a file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %v", err)
	}
	defer file.Close()

	// Encode and write the public key in PEM format
	err = pem.Encode(file, pubKeyPem)
	if err != nil {
		return fmt.Errorf("failed to encode public key: %v", err)
	}

	return nil
}

// SaveSignatureToFile can be used for storing signature files in hex format
func (e *ECDSAImpl) SaveSignatureToFile(filename string, data []byte) error {
	hexData := hex.EncodeToString(data)
	err := ioutil.WriteFile(filename, []byte(hexData), 0644)
	if err != nil {
		return fmt.Errorf("failed to write data to file %s: %v", filename, err)
	}
	return nil
}

// ReadPrivateKey reads an ECDSA private key from a PEM file using encoding/pem
func (e *ECDSAImpl) ReadPrivateKey(privateKeyPath string) (*ecdsa.PrivateKey, error) {
	privKeyPEM, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key file: %v", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(privKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	// Extract the private key (first part is 'D', followed by 'X' and 'Y' of the public key)
	privKeyBytes := block.Bytes
	privKeyD := new(big.Int).SetBytes(privKeyBytes[:32])  // first 32 bytes for D
	pubKeyX := new(big.Int).SetBytes(privKeyBytes[32:64]) // next 32 bytes for X
	pubKeyY := new(big.Int).SetBytes(privKeyBytes[64:96]) // last 32 bytes for Y

	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(), // Assume P256 curve, you can make it dynamic
		X:     pubKeyX,
		Y:     pubKeyY,
	}

	privateKey := &ecdsa.PrivateKey{
		D:         privKeyD,
		PublicKey: *publicKey,
	}

	return privateKey, nil
}

// ReadPublicKey reads an ECDSA public key from a PEM file using encoding/pem
func (e *ECDSAImpl) ReadPublicKey(publicKeyPath string) (*ecdsa.PublicKey, error) {
	pubKeyPEM, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key file: %v", err)
	}

	// Decode the PEM block
	block, _ := pem.Decode(pubKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	// Extract the public key (first 32 bytes for 'X' and next 32 bytes for 'Y')
	pubKeyBytes := block.Bytes
	pubKeyX := new(big.Int).SetBytes(pubKeyBytes[:32])   // first 32 bytes for X
	pubKeyY := new(big.Int).SetBytes(pubKeyBytes[32:64]) // next 32 bytes for Y

	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P256(), // Assume P256 curve, you can make it dynamic
		X:     pubKeyX,
		Y:     pubKeyY,
	}

	return publicKey, nil
}
