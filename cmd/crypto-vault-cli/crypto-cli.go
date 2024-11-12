package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	cryptography "crypto_vault_service/pkg/cryptography"
)

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

	aes := &cryptography.AESImpl{}

	// Generate AES Key
	key, err := aes.GenerateKey(keySize)
	if err != nil {
		log.Fatalf("Error generating AES key: %v\n", err)
	}

	// Encrypt the file
	plainText, err := cryptography.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	encryptedData, err := aes.Encrypt(plainText, key)
	if err != nil {
		log.Fatalf("Error encrypting data: %v\n", err)
	}

	// Save encrypted file
	err = cryptography.WriteFile(outputFile, encryptedData)
	if err != nil {
		log.Fatalf("Error writing encrypted file: %v\n", err)
	}
	fmt.Printf("Encrypted data saved to %s\n", outputFile)

	// Save the AES key to the specified key directory
	keyFilePath := filepath.Join(keyDir, "encryption_key.bin")
	err = cryptography.WriteFile(keyFilePath, key)
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
	encryptedData, err := cryptography.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading encrypted file: %v\n", err)
	}

	aes := &cryptography.AESImpl{}

	decryptedData, err := aes.Decrypt(encryptedData, key)
	if err != nil {
		log.Fatalf("Error decrypting data: %v\n", err)
	}

	// Save decrypted file
	err = cryptography.WriteFile(outputFile, decryptedData)
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
	rsa := &cryptography.RSAImpl{}
	if publicKeyPath == "" {
		// Generate RSA keys

		privateKey, pubKey, genErr := rsa.GenerateKeys(2048)
		if genErr != nil {
			log.Fatalf("Error generating RSA keys: %v\n", genErr)
		}
		publicKey = pubKey

		// Optionally save the private and public keys
		err = cryptography.SavePrivateKeyToFile(privateKey, "data/private_key.pem")
		if err != nil {
			log.Fatalf("Error saving private key: %v\n", err)
		}
		err = cryptography.SavePublicKeyToFile(publicKey, "data/public_key.pem")
		if err != nil {
			log.Fatalf("Error saving public key: %v\n", err)
		}
		fmt.Println("Generated and saved RSA keys.")
	} else {
		// Read the provided public key
		publicKey, err = cryptography.ReadPublicKey(publicKeyPath)
		if err != nil {
			log.Fatalf("Error reading public key: %v\n", err)
		}
	}

	// Encrypt the file
	plainText, err := cryptography.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	encryptedData, err := rsa.Encrypt(plainText, publicKey)
	if err != nil {
		log.Fatalf("Error encrypting data: %v\n", err)
	}

	// Save encrypted file
	err = cryptography.WriteFile(outputFile, encryptedData)
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
	rsa := &cryptography.RSAImpl{}
	if privateKeyPath == "" {
		// Generate RSA keys
		privKey, _, genErr := rsa.GenerateKeys(2048)
		if genErr != nil {
			log.Fatalf("Error generating RSA keys: %v\n", genErr)
		}
		privateKey = privKey

		// Optionally save the private and public keys
		err = cryptography.SavePrivateKeyToFile(privateKey, "private_key.pem")
		if err != nil {
			log.Fatalf("Error saving private key: %v\n", err)
		}
		fmt.Println("Generated and saved private key.")
	} else {
		// Read the provided private key
		privateKey, err = cryptography.ReadPrivateKey(privateKeyPath)
		if err != nil {
			log.Fatalf("Error reading private key: %v\n", err)
		}
	}

	// Decrypt the file
	encryptedData, err := cryptography.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading encrypted file: %v\n", err)
	}

	decryptedData, err := rsa.Decrypt(encryptedData, privateKey)
	if err != nil {
		log.Fatalf("Error decrypting data: %v\n", err)
	}

	// Save decrypted file
	err = cryptography.WriteFile(outputFile, decryptedData)
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
