package commands

import (
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid" // Import UUID package
	"github.com/spf13/cobra"
)

// Encrypts a file using AES and saves the symmetric key with a UUID prefix
func EncryptAESCmd(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input")
	outputFile, _ := cmd.Flags().GetString("output")
	keySize, _ := cmd.Flags().GetInt("keySize")
	keyDir, _ := cmd.Flags().GetString("keyDir")

	// Validate input arguments
	if inputFile == "" || outputFile == "" || keyDir == "" {
		log.Fatalf("Error: input, output, and keyDir flags are required\n")
	}

	aes := &cryptography.AESImpl{}

	// Generate AES Key
	key, err := aes.GenerateKey(keySize)
	if err != nil {
		log.Fatalf("Error generating AES key: %v\n", err)
	}

	// Encrypt the file
	plainText, err := utils.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	encryptedData, err := aes.Encrypt(plainText, key)
	if err != nil {
		log.Fatalf("Error encrypting data: %v\n", err)
	}

	// Save encrypted file
	err = utils.WriteFile(outputFile, encryptedData)
	if err != nil {
		log.Fatalf("Error writing encrypted file: %v\n", err)
	}
	fmt.Printf("Encrypted data saved to %s\n", outputFile)

	// Generate a UUID for the key filename
	uniqueID := uuid.New().String() // Generate a unique UUID

	// Save the AES key with the UUID prefix in the specified key directory
	keyFilePath := filepath.Join(keyDir, fmt.Sprintf("%s-symmetric_key.bin", uniqueID))
	err = utils.WriteFile(keyFilePath, key)
	if err != nil {
		log.Fatalf("Error writing AES key to file: %v\n", err)
	}
	fmt.Printf("AES key saved to %s\n", keyFilePath)
}

// Decrypts a file using AES and reads the corresponding symmetric key with a UUID prefix
func DecryptAESCmd(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input")
	outputFile, _ := cmd.Flags().GetString("output")
	symmetricKey, _ := cmd.Flags().GetString("symmetricKey")

	// Validate input arguments
	if inputFile == "" || outputFile == "" || symmetricKey == "" {
		log.Fatalf("Error: input, output and symmetricKey flags are required\n")
	}

	// Read the symmetric key from the file
	key, err := os.ReadFile(symmetricKey)
	if err != nil {
		log.Fatalf("Error reading symmetric key from file: %v\n", err)
	}

	// Decrypt the file
	encryptedData, err := utils.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading encrypted file: %v\n", err)
	}

	aes := &cryptography.AESImpl{}

	decryptedData, err := aes.Decrypt(encryptedData, key)
	if err != nil {
		log.Fatalf("Error decrypting data: %v\n", err)
	}

	// Save decrypted file
	err = utils.WriteFile(outputFile, decryptedData)
	if err != nil {
		log.Fatalf("Error writing decrypted file: %v\n", err)
	}
	fmt.Printf("Decrypted data saved to %s\n", outputFile)
}
