package commands

import (
	"crypto/rsa"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/utils"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// RSA Command
func EncryptRSACmd(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input")
	outputFile, _ := cmd.Flags().GetString("output")
	keyDir, _ := cmd.Flags().GetString("keyDir") // Directory to save keys

	// Validate input arguments
	if inputFile == "" || outputFile == "" || keyDir == "" {
		log.Fatalf("Error: input, output and keyDir flags are required\n")
	}

	// Generate RSA keys if no public key is provided
	var publicKey *rsa.PublicKey
	var err error
	rsa := &cryptography.RSAImpl{}

	uniqueID := uuid.New()
	// Generate RSA keys

	privateKey, publicKey, genErr := rsa.GenerateKeys(2048)
	if genErr != nil {
		log.Fatalf("Error generating RSA keys: %v\n", genErr)
	}

	privateKeyFilePath := fmt.Sprintf("%s/%s-private_key.pem", keyDir, uniqueID.String())
	// Optionally save the private and public keys
	err = rsa.SavePrivateKeyToFile(privateKey, privateKeyFilePath)
	if err != nil {
		log.Fatalf("Error saving private key: %v\n", err)
	}

	publicKeyFilePath := fmt.Sprintf("%s/%s-public_key.pem", keyDir, uniqueID.String())
	err = rsa.SavePublicKeyToFile(publicKey, publicKeyFilePath)
	if err != nil {
		log.Fatalf("Error saving public key: %v\n", err)
	}
	fmt.Println("Generated and saved RSA keys.")
	fmt.Println("Private key path:", privateKeyFilePath)
	fmt.Println("Public key path:", publicKeyFilePath)

	// Encrypt the file
	plainText, err := utils.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	encryptedData, err := rsa.Encrypt(plainText, publicKey)
	if err != nil {
		log.Fatalf("Error encrypting data: %v\n", err)
	}

	// Save encrypted file
	err = utils.WriteFile(outputFile, encryptedData)
	if err != nil {
		log.Fatalf("Error writing encrypted file: %v\n", err)
	}
	fmt.Printf("Encrypted data saved to %s\n", outputFile)
}

func DecryptRSACmd(cmd *cobra.Command, args []string) {
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
		err = rsa.SavePrivateKeyToFile(privateKey, "private_key.pem")
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
	encryptedData, err := utils.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading encrypted file: %v\n", err)
	}

	decryptedData, err := rsa.Decrypt(encryptedData, privateKey)
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
