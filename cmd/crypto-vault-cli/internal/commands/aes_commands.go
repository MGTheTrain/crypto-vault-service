package commands

import (
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type AESCommandHandler struct {
	aes    *cryptography.AES
	Logger logger.Logger
}

func NewAESCommandHandler() *AESCommandHandler {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	factory := &logger.LoggerFactory{}

	logger, err := factory.NewLogger(loggerSettings)
	if err != nil {
		log.Panicf("Error creating logger: %v", err)
		return nil
	}

	aes, err := cryptography.NewAES(logger)
	if err != nil {
		log.Panicf("%v\n", err)
		return nil
	}

	return &AESCommandHandler{
		aes:    aes,
		Logger: logger,
	}
}

// EncryptAESCmd encrypts a file using AES and saves the symmetric key with a UUID prefix
func (commandHandler *AESCommandHandler) EncryptAESCmd(cmd *cobra.Command, args []string) {

	inputFilePath, err := cmd.Flags().GetString("input-file")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	outputFilePath, err := cmd.Flags().GetString("output-file")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	keySize, err := cmd.Flags().GetInt("key-size")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	keyDir, _ := cmd.Flags().GetString("key-dir")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	key, err := commandHandler.aes.GenerateKey(keySize)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	plainText, err := os.ReadFile(inputFilePath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	encryptedData, err := commandHandler.aes.Encrypt(plainText, key)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	err = os.WriteFile(outputFilePath, encryptedData, 0644)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	uniqueID := uuid.New().String()

	keyFilePath := filepath.Join(keyDir, fmt.Sprintf("%s-symmetric-key.bin", uniqueID))
	err = os.WriteFile(keyFilePath, key, 0644)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	commandHandler.Logger.Info(fmt.Sprintf("Encrypted data saved to %s. AES key saved to %s", outputFilePath, keyFilePath))
}

// DecryptAESCmd decrypts a file using AES and reads the corresponding symmetric key with a UUID prefix
func (commandHandler *AESCommandHandler) DecryptAESCmd(cmd *cobra.Command, args []string) {
	inputFilePath, err := cmd.Flags().GetString("input-file")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	outputFilePath, err := cmd.Flags().GetString("output-file")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	symmetricKey, err := cmd.Flags().GetString("symmetric-key")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	key, err := os.ReadFile(symmetricKey)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	encryptedData, err := os.ReadFile(inputFilePath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	decryptedData, err := commandHandler.aes.Decrypt(encryptedData, key)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	err = os.WriteFile(outputFilePath, decryptedData, 0644)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	commandHandler.Logger.Info(fmt.Sprintf("Decrypted data saved to %s.", outputFilePath))
}

func InitAESCommands(rootCmd *cobra.Command) {
	handler := NewAESCommandHandler()

	var encryptAESFileCmd = &cobra.Command{
		Use:   "encrypt-aes",
		Short: "Encrypt a file using AES",
		Run:   handler.EncryptAESCmd,
	}
	encryptAESFileCmd.Flags().StringP("input-file", "", "", "Input file path")
	encryptAESFileCmd.Flags().StringP("output-file", "", "", "Output encrypted file path")
	encryptAESFileCmd.Flags().IntP("key-size", "", 16, "AES key size (default 16 bytes for AES-128)")
	encryptAESFileCmd.Flags().StringP("key-dir", "", "", "Directory to store the encryption key")
	rootCmd.AddCommand(encryptAESFileCmd)

	var decryptAESFileCmd = &cobra.Command{
		Use:   "decrypt-aes",
		Short: "Decrypt a file using AES",
		Run:   handler.DecryptAESCmd,
	}
	decryptAESFileCmd.Flags().StringP("input-file", "i", "", "Input encrypted file path")
	decryptAESFileCmd.Flags().StringP("output-file", "o", "", "Output decrypted file path")
	decryptAESFileCmd.Flags().StringP("symmetric-key", "k", "", "Path to the symmetric key")
	rootCmd.AddCommand(decryptAESFileCmd)
}
