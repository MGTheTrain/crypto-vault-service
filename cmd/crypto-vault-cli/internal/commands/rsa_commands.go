package commands

import (
	"crypto/rsa"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// RSA Command
func EncryptRSACmd(cmd *cobra.Command, args []string) {
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

	keyDir, err := cmd.Flags().GetString("key-dir")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	// Validate input arguments
	if inputFile == "" || outputFile == "" || keyDir == "" {
		log.Fatalf("Error: input, output and keyDir flags are required\n")
	}

	// Generate RSA keys if no public key is provided
	var publicKey *rsa.PublicKey
	rsa := &cryptography.RSA{}

	uniqueID := uuid.New()
	// Generate RSA keys

	privateKey, publicKey, err := rsa.GenerateKeys(2048)
	if err != nil {
		log.Fatalf("Error generating RSA keys: %v\n", err)
	}

	privateKeyFilePath := fmt.Sprintf("%s/%s-private-key.pem", keyDir, uniqueID.String())
	// Optionally save the private and public keys
	err = rsa.SavePrivateKeyToFile(privateKey, privateKeyFilePath)
	if err != nil {
		log.Fatalf("Error saving private key: %v\n", err)
	}

	publicKeyFilePath := fmt.Sprintf("%s/%s-public-key.pem", keyDir, uniqueID.String())
	err = rsa.SavePublicKeyToFile(publicKey, publicKeyFilePath)
	if err != nil {
		log.Fatalf("Error saving public key: %v\n", err)
	}
	fmt.Println("Generated and saved RSA keys.")
	fmt.Println("Private key path:", privateKeyFilePath)
	fmt.Println("Public key path:", publicKeyFilePath)

	// Encrypt the file
	plainText, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	encryptedData, err := rsa.Encrypt(plainText, publicKey)
	if err != nil {
		log.Fatalf("Error encrypting data: %v\n", err)
	}

	// Save encrypted file
	err = os.WriteFile(outputFile, encryptedData, 0644)
	if err != nil {
		log.Fatalf("Error writing encrypted file: %v\n", err)
	}
	fmt.Printf("Encrypted data saved to %s\n", outputFile)
}

func DecryptRSACmd(cmd *cobra.Command, args []string) {
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

	privateKeyPath, err := cmd.Flags().GetString("private-key")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	// Generate RSA keys if no private key is provided
	var privateKey *rsa.PrivateKey
	rsa := &cryptography.RSA{}
	if privateKeyPath == "" {
		// Generate RSA keys
		privKey, _, err := rsa.GenerateKeys(2048)
		if err != nil {
			log.Fatalf("Error generating RSA keys: %v\n", err)
		}
		privateKey = privKey

		// Optionally save the private and public keys
		err = rsa.SavePrivateKeyToFile(privateKey, "private-key.pem")
		if err != nil {
			log.Fatalf("Error saving private key: %v\n", err)
		}
		fmt.Println("Generated and saved private key.")
	} else {
		// Read the provided private key
		privateKey, err = rsa.ReadPrivateKey(privateKeyPath)
		if err != nil {
			log.Fatalf("Error reading private key: %v\n", err)
		}
	}

	// Decrypt the file
	encryptedData, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading encrypted file: %v\n", err)
	}

	decryptedData, err := rsa.Decrypt(encryptedData, privateKey)
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

func InitRSACommands(rootCmd *cobra.Command) {
	var encryptRSAFileCmd = &cobra.Command{
		Use:   "encrypt-rsa",
		Short: "Encrypt a file using RSA",
		Run:   EncryptRSACmd,
	}
	encryptRSAFileCmd.Flags().StringP("input-file", "", "", "Input file path")
	encryptRSAFileCmd.Flags().StringP("output-file", "", "", "Output encrypted file path")
	encryptRSAFileCmd.Flags().StringP("key-dir", "", "", "Directory to store the encryption key")
	rootCmd.AddCommand(encryptRSAFileCmd)

	var decryptRSAFileCmd = &cobra.Command{
		Use:   "decrypt-rsa",
		Short: "Decrypt a file using RSA",
		Run:   DecryptRSACmd,
	}
	decryptRSAFileCmd.Flags().StringP("input-file", "", "", "Input encrypted file path")
	decryptRSAFileCmd.Flags().StringP("output-file", "", "", "Output decrypted file path")
	decryptRSAFileCmd.Flags().StringP("private-key", "", "", "Path to RSA private key")
	rootCmd.AddCommand(decryptRSAFileCmd)
}
