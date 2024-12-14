package commands

import (
	"crypto_vault_service/internal/infrastructure/cryptography"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid" // Import UUID package
	"github.com/spf13/cobra"
)

// Encrypts a file using AES and saves the symmetric key with a UUID prefix
func EncryptAESCmd(cmd *cobra.Command, args []string) {

	inputFile, err := cmd.Flags().GetString("input-file")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	outputFile, err := cmd.Flags().GetString("output-file")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	keySize, err := cmd.Flags().GetInt("key-size")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	keyDir, _ := cmd.Flags().GetString("key-dir")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	aes := &cryptography.AES{}

	// Generate AES Key
	key, err := aes.GenerateKey(keySize)
	if err != nil {
		log.Fatalf("Error generating AES key: %v\n", err)
	}

	// Encrypt the file
	plainText, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	encryptedData, err := aes.Encrypt(plainText, key)
	if err != nil {
		log.Fatalf("Error encrypting data: %v\n", err)
	}

	// Save encrypted file
	err = os.WriteFile(outputFile, encryptedData, 0644)
	if err != nil {
		log.Fatalf("Error writing encrypted file: %v\n", err)
	}
	fmt.Printf("Encrypted data saved to %s\n", outputFile)

	// Generate a UUID for the key filename
	uniqueID := uuid.New().String()

	// Save the AES key with the UUID prefix in the specified key directory
	keyFilePath := filepath.Join(keyDir, fmt.Sprintf("%s-symmetric-key.bin", uniqueID))
	err = os.WriteFile(keyFilePath, key, 0644)
	if err != nil {
		log.Fatalf("Error writing AES key to file: %v\n", err)
	}
	fmt.Printf("AES key saved to %s\n", keyFilePath)
}

// Decrypts a file using AES and reads the corresponding symmetric key with a UUID prefix
func DecryptAESCmd(cmd *cobra.Command, args []string) {
	inputFile, err := cmd.Flags().GetString("input-file")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	outputFile, err := cmd.Flags().GetString("output-file")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	symmetricKey, err := cmd.Flags().GetString("symmetric-key")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

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
	encryptedData, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading encrypted file: %v\n", err)
	}

	aes := &cryptography.AES{}

	decryptedData, err := aes.Decrypt(encryptedData, key)
	if err != nil {
		log.Fatalf("Error decrypting data: %v\n", err)
	}

	// Save decrypted file
	err = os.WriteFile(outputFile, decryptedData, 0644)
	if err != nil {
		log.Fatalf("Error writing decrypted file: %v\n", err)
	}
	fmt.Printf("Decrypted data saved to %s\n", outputFile)
}

func InitAESCommands(rootCmd *cobra.Command) {
	// AES Commands
	var encryptAESFileCmd = &cobra.Command{
		Use:   "encrypt-aes",
		Short: "Encrypt a file using AES",
		Run:   EncryptAESCmd,
	}
	encryptAESFileCmd.Flags().StringP("input-file", "", "", "Input file path")
	encryptAESFileCmd.Flags().StringP("output-file", "", "", "Output encrypted file path")
	encryptAESFileCmd.Flags().IntP("key-size", "", 16, "AES key size (default 16 bytes for AES-128)")
	encryptAESFileCmd.Flags().StringP("key-dir", "", "", "Directory to store the encryption key")
	rootCmd.AddCommand(encryptAESFileCmd)

	var decryptAESFileCmd = &cobra.Command{
		Use:   "decrypt-aes",
		Short: "Decrypt a file using AES",
		Run:   DecryptAESCmd,
	}
	decryptAESFileCmd.Flags().StringP("input-file", "i", "", "Input encrypted file path")
	decryptAESFileCmd.Flags().StringP("output-file", "o", "", "Output decrypted file path")
	decryptAESFileCmd.Flags().StringP("symmetric-key", "k", "", "Path to the symmetric key")
	rootCmd.AddCommand(decryptAESFileCmd)
}
