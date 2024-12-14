package commands

import (
	"crypto/rsa"
	"crypto_vault_service/cmd/crypto-vault-cli/internal/status"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// RSA Command
func EncryptRSACmd(cmd *cobra.Command, args []string) {
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

	keyDir, err := cmd.Flags().GetString("key-dir")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	var publicKey *rsa.PublicKey
	rsa := &cryptography.RSA{}

	uniqueID := uuid.New()

	privateKey, publicKey, err := rsa.GenerateKeys(2048)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	privateKeyFilePath := fmt.Sprintf("%s/%s-private-key.pem", keyDir, uniqueID.String())

	err = rsa.SavePrivateKeyToFile(privateKey, privateKeyFilePath)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	publicKeyFilePath := fmt.Sprintf("%s/%s-public-key.pem", keyDir, uniqueID.String())
	err = rsa.SavePublicKeyToFile(publicKey, publicKeyFilePath)
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

	encryptedData, err := rsa.Encrypt(plainText, publicKey)
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
	info := status.NewInfo(fmt.Sprintf("Generated and saved RSA keys. Private key path: %s. Public key path: %s. Encrypted data saved to %s", privateKeyFilePath, publicKeyFilePath, outputFile))
	info.PrintJsonInfo(false)
}

func DecryptRSACmd(cmd *cobra.Command, args []string) {
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

	privateKeyPath, err := cmd.Flags().GetString("private-key")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	var privateKey *rsa.PrivateKey
	rsa := &cryptography.RSA{}
	if privateKeyPath == "" {

		privKey, _, err := rsa.GenerateKeys(2048)
		if err != nil {
			e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
			e.PrintJsonError()
			return
		}
		privateKey = privKey

		err = rsa.SavePrivateKeyToFile(privateKey, "private-key.pem")
		if err != nil {
			e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
			e.PrintJsonError()
			return
		}
		info := status.NewInfo("Generated and saved private key")
		info.PrintJsonInfo(false)
	} else {

		privateKey, err = rsa.ReadPrivateKey(privateKeyPath)
		if err != nil {
			e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
			e.PrintJsonError()
			return
		}
	}

	encryptedData, err := os.ReadFile(inputFile)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	decryptedData, err := rsa.Decrypt(encryptedData, privateKey)
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
