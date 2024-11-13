package commands

import (
	"crypto/rsa"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/utils"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// RSA Command
func EncryptRSACmd(cmd *cobra.Command, args []string) {
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
		err = rsa.SavePrivateKeyToFile(privateKey, "data/private_key.pem")
		if err != nil {
			log.Fatalf("Error saving private key: %v\n", err)
		}
		err = rsa.SavePublicKeyToFile(publicKey, "data/public_key.pem")
		if err != nil {
			log.Fatalf("Error saving public key: %v\n", err)
		}
		fmt.Println("Generated and saved RSA keys.")
	} else {
		// Read the provided public key
		publicKey, err = rsa.ReadPublicKey(publicKeyPath)
		if err != nil {
			log.Fatalf("Error reading public key: %v\n", err)
		}
	}

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
