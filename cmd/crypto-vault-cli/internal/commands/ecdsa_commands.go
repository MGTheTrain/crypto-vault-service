package commands

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// ECDSA command
// signECCCmd signs the contents of a file with ECDSA
func SignECCCmd(cmd *cobra.Command, args []string) {

	inputFile, err := cmd.Flags().GetString("input-file")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	keyDir, err := cmd.Flags().GetString("key-dir")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	// ECC implementation
	EC := &cryptography.EC{}
	var privateKey *ecdsa.PrivateKey
	var publicKey *ecdsa.PublicKey

	// Generate new ECC keys if no private key is provided
	privateKey, publicKey, err = EC.GenerateKeys(elliptic.P256())
	if err != nil {
		log.Fatalf("Error generating ECC keys: %v\n", err)
	}

	// Read the file content
	fileContent, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v\n", err)
	}

	// Sign the file content (hash the content before signing)
	signature, err := EC.Sign(fileContent, privateKey)
	if err != nil {
		log.Fatalf("Error signing file content: %v\n", err)
	}

	// Output the signature
	fmt.Printf("Signature: %x\n", signature)

	uniqueID := uuid.New()
	// Save the private and public keys to files (if they were generated)
	if privateKey != nil && keyDir != "" {
		privateKeyFilePath := fmt.Sprintf("%s/%s-private-key.pem", keyDir, uniqueID.String())

		err = EC.SavePrivateKeyToFile(privateKey, privateKeyFilePath)
		if err != nil {
			log.Fatalf("Error saving private key: %v\n", err)
		}
		fmt.Printf("Private key saved to: %s\n", privateKeyFilePath)
	}

	if publicKey != nil && keyDir != "" {
		publicKeyFilePath := fmt.Sprintf("%s/%s-public-key.pem", keyDir, uniqueID.String())
		err = EC.SavePublicKeyToFile(publicKey, publicKeyFilePath)
		if err != nil {
			log.Fatalf("Error saving public key: %v\n", err)
		}
		fmt.Printf("Public key saved to: %s\n", publicKeyFilePath)
	}

	// Save the signature to a file in the data folder (optional, based on the input file)
	if keyDir != "" {
		signatureFilePath := fmt.Sprintf("%s/%s-signature.sig", keyDir, uniqueID.String())
		err = EC.SaveSignatureToFile(signatureFilePath, signature)
		if err != nil {
			log.Fatalf("Error saving signature: %v\n", err)
		}
		fmt.Printf("Signature saved to: %s\n", signatureFilePath)
	}
}

// verifyECCCmd verifies the signature of a file's content using ECDSA
func VerifyECCCmd(cmd *cobra.Command, args []string) {
	inputFile, err := cmd.Flags().GetString("input-file")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	publicKeyPath, err := cmd.Flags().GetString("public-key")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	signatureFile, err := cmd.Flags().GetString("signature-file")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	// ECC implementation
	EC := &cryptography.EC{}
	var publicKey *ecdsa.PublicKey

	// Read the public key
	if publicKeyPath == "" {
		log.Fatalf("Public key is required for ECC signature verification.\n")
	} else {
		publicKey, err = EC.ReadPublicKey(publicKeyPath, elliptic.P256())
		if err != nil {
			log.Fatalf("Error reading public key: %v\n", err)
		}
	}

	// Read the file content (optional: you can also hash the content before verifying)
	fileContent, err := os.ReadFile(inputFile)
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
	valid, err := EC.Verify(fileContent, signature, publicKey)
	if err != nil {
		log.Fatalf("Error verifying signature: %v\n", err)
	}

	if valid {
		fmt.Println("Signature is valid.")
	} else {
		fmt.Println("Signature is invalid.")
	}
}

func InitECDSACommands(rootCmd *cobra.Command) {
	var signECCMessageCmd = &cobra.Command{
		Use:   "sign-ecc",
		Short: "Sign a message using ECC",
		Run:   SignECCCmd,
	}
	signECCMessageCmd.Flags().StringP("input-file", "", "", "Path to file that needs to be signed")
	signECCMessageCmd.Flags().StringP("key-dir", "", "", "Directory to save generated keys (optional)")
	rootCmd.AddCommand(signECCMessageCmd)

	var verifyECCSignatureCmd = &cobra.Command{
		Use:   "verify-ecc",
		Short: "Verify a signature using ECC",
		Run:   VerifyECCCmd,
	}
	verifyECCSignatureCmd.Flags().StringP("input-file", "", "", "Path to ECC public key")
	verifyECCSignatureCmd.Flags().StringP("public-key", "", "", "The public key used to verify the signature")
	verifyECCSignatureCmd.Flags().StringP("signature-file", "", "", "Signature to verify (hex format)")
	rootCmd.AddCommand(verifyECCSignatureCmd)
}
