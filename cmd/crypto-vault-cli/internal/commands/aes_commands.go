package commands

import (
	"crypto_vault_service/cmd/crypto-vault-cli/internal/status"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// Encrypts a file using AES and saves the symmetric key with a UUID prefix
func EncryptAESCmd(cmd *cobra.Command, args []string) {

	inputFile, err := cmd.Flags().GetString("input-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	outputFile, err := cmd.Flags().GetString("output-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	keySize, err := cmd.Flags().GetInt("key-size")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	keyDir, _ := cmd.Flags().GetString("key-dir")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	aes := &cryptography.AES{}

	key, err := aes.GenerateKey(keySize)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	plainText, err := os.ReadFile(inputFile)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	encryptedData, err := aes.Encrypt(plainText, key)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	err = os.WriteFile(outputFile, encryptedData, 0644)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	uniqueID := uuid.New().String()

	keyFilePath := filepath.Join(keyDir, fmt.Sprintf("%s-symmetric-key.bin", uniqueID))
	err = os.WriteFile(keyFilePath, key, 0644)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	info := status.NewInfo(fmt.Sprintf("Encrypted data saved to %s. AES key saved to %s", outputFile, keyFilePath))
	info.PrintJsonInfo(false)
}

// Decrypts a file using AES and reads the corresponding symmetric key with a UUID prefix
func DecryptAESCmd(cmd *cobra.Command, args []string) {
	inputFile, err := cmd.Flags().GetString("input-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	outputFile, err := cmd.Flags().GetString("output-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	symmetricKey, err := cmd.Flags().GetString("symmetric-key")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	key, err := os.ReadFile(symmetricKey)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	encryptedData, err := os.ReadFile(inputFile)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	aes := &cryptography.AES{}

	decryptedData, err := aes.Decrypt(encryptedData, key)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	err = os.WriteFile(outputFile, decryptedData, 0644)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	info := status.NewInfo(fmt.Sprintf("Decrypted data saved to %s", outputFile))
	info.PrintJsonInfo(false)
}

func InitAESCommands(rootCmd *cobra.Command) {
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
