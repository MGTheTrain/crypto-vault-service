package main

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
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// AES Functions
// GenerateRandomAESKey generates a random AES key of the specified size
func generateRandomAESKey(keySize int) ([]byte, error) {
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
func encryptAES(plainText, key []byte) ([]byte, error) {
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
func decryptAES(ciphertext, key []byte) ([]byte, error) {
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
func generateRSAKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA keys: %v", err)
	}

	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

func savePrivateKeyToFile(privateKey *rsa.PrivateKey, filename string) error {
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

func savePublicKeyToFile(publicKey *rsa.PublicKey, filename string) error {
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
func encryptWithRSA(publicKey *rsa.PublicKey, data []byte) ([]byte, error) {
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}
	return encryptedData, nil
}

// Decrypt data with RSA private key
func decryptWithRSA(privateKey *rsa.PrivateKey, encryptedData []byte) ([]byte, error) {
	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %v", err)
	}
	return decryptedData, nil
}

// File Operations
func readFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

func writeFile(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0644)
}

// AES Command

// Read RSA private key from PEM file
func readPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
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
func readPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
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

// Encrypts a file using AES and saves the encryption key
func encryptAESCmd(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input")
	outputFile, _ := cmd.Flags().GetString("output")
	keySize, _ := cmd.Flags().GetInt("keySize")
	keyDir, _ := cmd.Flags().GetString("keyDir")

	// Validate input arguments
	if inputFile == "" || outputFile == "" || keyDir == "" {
		log.Fatalf("Error: input, output, and keyDir flags are required\n")
	}

	// Generate AES Key
	key, err := generateRandomAESKey(keySize)
	if err != nil {
		log.Fatalf("Error generating AES key: %v\n", err)
	}

	// Encrypt the file
	plainText, err := readFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	encryptedData, err := encryptAES(plainText, key)
	if err != nil {
		log.Fatalf("Error encrypting data: %v\n", err)
	}

	// Save encrypted file
	err = writeFile(outputFile, encryptedData)
	if err != nil {
		log.Fatalf("Error writing encrypted file: %v\n", err)
	}
	fmt.Printf("Encrypted data saved to %s\n", outputFile)

	// Save the AES key to the specified key directory
	keyFilePath := filepath.Join(keyDir, "encryption_key.bin")
	err = writeFile(keyFilePath, key)
	if err != nil {
		log.Fatalf("Error writing AES key to file: %v\n", err)
	}
	fmt.Printf("AES key saved to %s\n", keyFilePath)
}

func decryptAESCmd(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input")
	outputFile, _ := cmd.Flags().GetString("output")
	keyDir, _ := cmd.Flags().GetString("keyDir")

	// Validate input arguments
	if inputFile == "" || outputFile == "" || keyDir == "" {
		log.Fatalf("Error: input, output, and keyDir flags are required\n")
	}

	// Read the encryption key from the specified directory
	keyFilePath := filepath.Join(keyDir, "encryption_key.bin")
	key, err := ioutil.ReadFile(keyFilePath)
	if err != nil {
		log.Fatalf("Error reading encryption key from file: %v\n", err)
	}

	// Decrypt the file
	encryptedData, err := readFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading encrypted file: %v\n", err)
	}

	decryptedData, err := decryptAES(encryptedData, key)
	if err != nil {
		log.Fatalf("Error decrypting data: %v\n", err)
	}

	// Save decrypted file
	err = writeFile(outputFile, decryptedData)
	if err != nil {
		log.Fatalf("Error writing decrypted file: %v\n", err)
	}
	fmt.Printf("Decrypted data saved to %s\n", outputFile)
}

// RSA Command
func encryptRSACmd(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input")
	outputFile, _ := cmd.Flags().GetString("output")
	publicKeyPath, _ := cmd.Flags().GetString("publicKey")

	// Generate RSA keys if no public key is provided
	var publicKey *rsa.PublicKey
	var err error
	if publicKeyPath == "" {
		// Generate RSA keys
		privateKey, pubKey, genErr := generateRSAKeys(2048)
		if genErr != nil {
			log.Fatalf("Error generating RSA keys: %v\n", genErr)
		}
		publicKey = pubKey

		// Optionally save the private and public keys
		err = savePrivateKeyToFile(privateKey, "data/private_key.pem")
		if err != nil {
			log.Fatalf("Error saving private key: %v\n", err)
		}
		err = savePublicKeyToFile(publicKey, "data/public_key.pem")
		if err != nil {
			log.Fatalf("Error saving public key: %v\n", err)
		}
		fmt.Println("Generated and saved RSA keys.")
	} else {
		// Read the provided public key
		publicKey, err = readPublicKey(publicKeyPath)
		if err != nil {
			log.Fatalf("Error reading public key: %v\n", err)
		}
	}

	// Encrypt the file
	plainText, err := readFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	encryptedData, err := encryptWithRSA(publicKey, plainText)
	if err != nil {
		log.Fatalf("Error encrypting data: %v\n", err)
	}

	// Save encrypted file
	err = writeFile(outputFile, encryptedData)
	if err != nil {
		log.Fatalf("Error writing encrypted file: %v\n", err)
	}
	fmt.Printf("Encrypted data saved to %s\n", outputFile)
}

func decryptRSACmd(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input")
	outputFile, _ := cmd.Flags().GetString("output")
	privateKeyPath, _ := cmd.Flags().GetString("privateKey")

	// Generate RSA keys if no private key is provided
	var privateKey *rsa.PrivateKey
	var err error
	if privateKeyPath == "" {
		// Generate RSA keys
		privKey, _, genErr := generateRSAKeys(2048)
		if genErr != nil {
			log.Fatalf("Error generating RSA keys: %v\n", genErr)
		}
		privateKey = privKey

		// Optionally save the private and public keys
		err = savePrivateKeyToFile(privateKey, "private_key.pem")
		if err != nil {
			log.Fatalf("Error saving private key: %v\n", err)
		}
		fmt.Println("Generated and saved private key.")
	} else {
		// Read the provided private key
		privateKey, err = readPrivateKey(privateKeyPath)
		if err != nil {
			log.Fatalf("Error reading private key: %v\n", err)
		}
	}

	// Decrypt the file
	encryptedData, err := readFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading encrypted file: %v\n", err)
	}

	decryptedData, err := decryptWithRSA(privateKey, encryptedData)
	if err != nil {
		log.Fatalf("Error decrypting data: %v\n", err)
	}

	// Save decrypted file
	err = writeFile(outputFile, decryptedData)
	if err != nil {
		log.Fatalf("Error writing decrypted file: %v\n", err)
	}
	fmt.Printf("Decrypted data saved to %s\n", outputFile)
}

// Main function
func main() {
	var rootCmd = &cobra.Command{Use: "crypto-cli"}

	// AES Commands
	var encryptAESFileCmd = &cobra.Command{
		Use:   "encrypt-aes",
		Short: "Encrypt a file using AES",
		Run:   encryptAESCmd,
	}
	encryptAESFileCmd.Flags().StringP("input", "i", "", "Input file path")
	encryptAESFileCmd.Flags().StringP("output", "o", "", "Output encrypted file path")
	encryptAESFileCmd.Flags().IntP("keySize", "k", 16, "AES key size (default 16 bytes for AES-128)")
	encryptAESFileCmd.Flags().StringP("keyDir", "d", "", "Directory to store the encryption key")
	rootCmd.AddCommand(encryptAESFileCmd)

	var decryptAESFileCmd = &cobra.Command{
		Use:   "decrypt-aes",
		Short: "Decrypt a file using AES",
		Run:   decryptAESCmd,
	}
	decryptAESFileCmd.Flags().StringP("input", "i", "", "Input encrypted file path")
	decryptAESFileCmd.Flags().StringP("output", "o", "", "Output decrypted file path")
	decryptAESFileCmd.Flags().StringP("keyDir", "d", "", "Directory to read the encryption key from")
	rootCmd.AddCommand(decryptAESFileCmd)

	// RSA Commands
	var encryptRSAFileCmd = &cobra.Command{
		Use:   "encrypt-rsa",
		Short: "Encrypt a file using RSA",
		Run:   encryptRSACmd,
	}
	encryptRSAFileCmd.Flags().StringP("input", "i", "", "Input file path")
	encryptRSAFileCmd.Flags().StringP("output", "o", "", "Output encrypted file path")
	encryptRSAFileCmd.Flags().StringP("publicKey", "p", "", "Path to RSA public key")
	rootCmd.AddCommand(encryptRSAFileCmd)

	var decryptRSAFileCmd = &cobra.Command{
		Use:   "decrypt-rsa",
		Short: "Decrypt a file using RSA",
		Run:   decryptRSACmd,
	}
	decryptRSAFileCmd.Flags().StringP("input", "i", "", "Input encrypted file path")
	decryptRSAFileCmd.Flags().StringP("output", "o", "", "Output decrypted file path")
	decryptRSAFileCmd.Flags().StringP("privateKey", "r", "", "Path to RSA private key")
	rootCmd.AddCommand(decryptRSAFileCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
