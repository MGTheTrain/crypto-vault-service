package main

import (
	"fmt"
	"os"

	commands "crypto_vault_service/cmd/crypto-vault-cli/internal/commands"

	"github.com/spf13/cobra"
)

// Main function
func main() {
	var rootCmd = &cobra.Command{Use: "crypto-cli"}

	// AES Commands
	var encryptAESFileCmd = &cobra.Command{
		Use:   "encrypt-aes",
		Short: "Encrypt a file using AES",
		Run:   commands.EncryptAESCmd,
	}
	encryptAESFileCmd.Flags().StringP("input", "i", "", "Input file path")
	encryptAESFileCmd.Flags().StringP("output", "o", "", "Output encrypted file path")
	encryptAESFileCmd.Flags().IntP("keySize", "k", 16, "AES key size (default 16 bytes for AES-128)")
	encryptAESFileCmd.Flags().StringP("keyDir", "d", "", "Directory to store the encryption key")
	rootCmd.AddCommand(encryptAESFileCmd)

	var decryptAESFileCmd = &cobra.Command{
		Use:   "decrypt-aes",
		Short: "Decrypt a file using AES",
		Run:   commands.DecryptAESCmd,
	}
	decryptAESFileCmd.Flags().StringP("input", "i", "", "Input encrypted file path")
	decryptAESFileCmd.Flags().StringP("output", "o", "", "Output decrypted file path")
	decryptAESFileCmd.Flags().StringP("keyDir", "d", "", "Directory to read the encryption key from")
	rootCmd.AddCommand(decryptAESFileCmd)

	// RSA Commands
	var encryptRSAFileCmd = &cobra.Command{
		Use:   "encrypt-rsa",
		Short: "Encrypt a file using RSA",
		Run:   commands.EncryptRSACmd,
	}
	encryptRSAFileCmd.Flags().StringP("input", "i", "", "Input file path")
	encryptRSAFileCmd.Flags().StringP("output", "o", "", "Output encrypted file path")
	encryptRSAFileCmd.Flags().StringP("publicKey", "p", "", "Path to RSA public key")
	rootCmd.AddCommand(encryptRSAFileCmd)

	var decryptRSAFileCmd = &cobra.Command{
		Use:   "decrypt-rsa",
		Short: "Decrypt a file using RSA",
		Run:   commands.DecryptRSACmd,
	}
	decryptRSAFileCmd.Flags().StringP("input", "i", "", "Input encrypted file path")
	decryptRSAFileCmd.Flags().StringP("output", "o", "", "Output decrypted file path")
	decryptRSAFileCmd.Flags().StringP("privateKey", "r", "", "Path to RSA private key")
	rootCmd.AddCommand(decryptRSAFileCmd)

	// ECDSA Command
	var signECCMessageCmd = &cobra.Command{
		Use:   "sign-ecc",
		Short: "Sign a message using ECC",
		Run:   commands.SignECCCmd,
	}

	// Rename the input flag to messageFile for clarity
	signECCMessageCmd.Flags().StringP("input", "i", "", "Path to file that needs to be signed")
	signECCMessageCmd.Flags().StringP("keyDir", "d", "", "Directory to save generated keys (optional)")
	rootCmd.AddCommand(signECCMessageCmd)

	var verifyECCSignatureCmd = &cobra.Command{
		Use:   "verify-ecc",
		Short: "Verify a signature using ECC",
		Run:   commands.VerifyECCCmd,
	}
	verifyECCSignatureCmd.Flags().StringP("input", "i", "", "Path to ECC public key")
	verifyECCSignatureCmd.Flags().StringP("publicKey", "p", "", "The public key used to verify the signature")
	verifyECCSignatureCmd.Flags().StringP("signature", "s", "", "Signature to verify (hex format)")
	rootCmd.AddCommand(verifyECCSignatureCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
