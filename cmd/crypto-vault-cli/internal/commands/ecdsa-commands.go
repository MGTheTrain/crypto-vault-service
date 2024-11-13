package commands

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/utils"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// ECDSA command
// signECCCmd signs the contents of a file with ECDSA
func SignECCCmd(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input") // File to sign
	keyDir, _ := cmd.Flags().GetString("keyDir")   // Directory to save keys

	// ECC implementation
	ecdsaImpl := &cryptography.ECDSAImpl{}
	var privateKey *ecdsa.PrivateKey
	var publicKey *ecdsa.PublicKey
	var err error

	// Generate new ECC keys if no private key is provided
	privateKey, publicKey, err = ecdsaImpl.GenerateKeys(elliptic.P256())
	if err != nil {
		log.Fatalf("Error generating ECC keys: %v\n", err)
	}

	// Read the file content
	fileContent, err := utils.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	// Sign the file content (hash the content before signing)
	signature, err := ecdsaImpl.Sign(fileContent, privateKey)
	if err != nil {
		log.Fatalf("Error signing file content: %v\n", err)
	}

	// Output the signature
	fmt.Printf("Signature: %x\n", signature)

	// Save the private and public keys to files (if they were generated)
	if privateKey != nil && keyDir != "" {
		privateKeyFilePath := fmt.Sprintf("%s/private_key.pem", keyDir)
		err = ecdsaImpl.SavePrivateKeyToFile(privateKey, privateKeyFilePath)
		if err != nil {
			log.Fatalf("Error saving private key: %v\n", err)
		}
		fmt.Printf("Private key saved to: %s\n", privateKeyFilePath)
	}

	if publicKey != nil && keyDir != "" {
		publicKeyFilePath := fmt.Sprintf("%s/public_key.pem", keyDir)
		err = ecdsaImpl.SavePublicKeyToFile(publicKey, publicKeyFilePath)
		if err != nil {
			log.Fatalf("Error saving public key: %v\n", err)
		}
		fmt.Printf("Public key saved to: %s\n", publicKeyFilePath)
	}

	// Save the signature to a file in the data folder (optional, based on the input file)
	if keyDir != "" {
		signatureFilePath := fmt.Sprintf("%s/signature.sig", keyDir)
		err = ecdsaImpl.SaveSignatureToFile(signatureFilePath, signature)
		if err != nil {
			log.Fatalf("Error saving signature: %v\n", err)
		}
		fmt.Printf("Signature saved to: %s\n", signatureFilePath)
	}
}

// verifyECCCmd verifies the signature of a file's content using ECDSA
func VerifyECCCmd(cmd *cobra.Command, args []string) {
	publicKeyPath, _ := cmd.Flags().GetString("publicKey") // Path to public key
	inputFile, _ := cmd.Flags().GetString("input")         // Input file to verify
	signatureFile, _ := cmd.Flags().GetString("signature") // Path to signature file

	// ECC implementation
	ecdsaImpl := &cryptography.ECDSAImpl{}
	var publicKey *ecdsa.PublicKey
	var err error

	// Read the public key
	if publicKeyPath == "" {
		log.Fatalf("Public key is required for ECC signature verification.\n")
	} else {
		publicKey, err = ecdsaImpl.ReadPublicKey(publicKeyPath)
		if err != nil {
			log.Fatalf("Error reading public key: %v\n", err)
		}
	}

	// Read the file content (optional: you can also hash the content before verifying)
	fileContent, err := utils.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	// Read the signature (from hex file)
	signatureHex, err := os.ReadFile(signatureFile)
	if err != nil {
		log.Fatalf("Error reading signature file: %v\n", err)
	}

	// Decode the hex string back to bytes
	signature, err := hex.DecodeString(string(signatureHex))
	if err != nil {
		log.Fatalf("Error decoding signature hex: %v\n", err)
	}

	// Verify the signature
	valid, err := ecdsaImpl.Verify(fileContent, signature, publicKey)
	if err != nil {
		log.Fatalf("Error verifying signature: %v\n", err)
	}

	if valid {
		fmt.Println("Signature is valid.")
	} else {
		fmt.Println("Signature is invalid.")
	}
}
