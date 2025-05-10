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
	aesProcessor cryptography.AESProcessor
	Logger       logger.Logger
}

func NewAESCommandHandler() *AESCommandHandler {
	loggerSettings := &settings.LoggerSettings{
		LogLevel: "info",
		LogType:  "console",
		FilePath: "",
	}

	logger, err := logger.GetLogger(loggerSettings)
	if err != nil {
		log.Panicf("Error creating logger: %v", err)
		return nil
	}

	aesProcessor, err := cryptography.NewAESProcessor(logger)
	if err != nil {
		log.Panicf("%v\n", err)
		return nil
	}

	return &AESCommandHandler{
		aesProcessor: aesProcessor,
		Logger:       logger,
	}
}

// GenerateAESKeysCmd generates AES key pairs and persists those in a selected directory
func (commandHandler *AESCommandHandler) GenerateAESKeysCmd(cmd *cobra.Command, args []string) {
	keySize, _ := cmd.Flags().GetInt("key-size")
	keyDir, _ := cmd.Flags().GetString("key-dir")

	uniqueID := uuid.New()

	secretKey, err := commandHandler.aesProcessor.GenerateKey(keySize)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	keyFilePath := filepath.Join(keyDir, fmt.Sprintf("%s-symmetric-key.bin", uniqueID))
	err = os.WriteFile(keyFilePath, secretKey, 0644)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}
	commandHandler.Logger.Info(fmt.Sprintf("AES key saved to %s", keyFilePath))
}

// EncryptAESCmd encrypts a file using AES and saves the symmetric key with a UUID prefix
func (commandHandler *AESCommandHandler) EncryptAESCmd(cmd *cobra.Command, args []string) {
	inputFilePath, _ := cmd.Flags().GetString("input-file")
	outputFilePath, _ := cmd.Flags().GetString("output-file")
	symmetricKey, _ := cmd.Flags().GetString("symmetric-key")

	plainText, err := os.ReadFile(inputFilePath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	key, err := os.ReadFile(symmetricKey)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	encryptedData, err := commandHandler.aesProcessor.Encrypt(plainText, key)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	err = os.WriteFile(outputFilePath, encryptedData, 0644)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	commandHandler.Logger.Info(fmt.Sprintf("Encrypted data saved to %s", outputFilePath))
}

// DecryptAESCmd decrypts a file using AES and reads the corresponding symmetric key with a UUID prefix
func (commandHandler *AESCommandHandler) DecryptAESCmd(cmd *cobra.Command, args []string) {
	inputFilePath, _ := cmd.Flags().GetString("input-file")
	outputFilePath, _ := cmd.Flags().GetString("output-file")
	symmetricKey, _ := cmd.Flags().GetString("symmetric-key")

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

	decryptedData, err := commandHandler.aesProcessor.Decrypt(encryptedData, key)
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

	var generateAESKeysCmd = &cobra.Command{
		Use:   "generate-aes-keys",
		Short: "Generate AES keys",
		Run:   handler.GenerateAESKeysCmd,
	}
	generateAESKeysCmd.Flags().IntP("key-size", "", 16, "AES key size (default 16 bytes for AES-128)")
	generateAESKeysCmd.Flags().StringP("key-dir", "", "", "Directory to store the encryption key")
	rootCmd.AddCommand(generateAESKeysCmd)

	var encryptAESFileCmd = &cobra.Command{
		Use:   "encrypt-aes",
		Short: "Encrypt a file using AES",
		Run:   handler.EncryptAESCmd,
	}
	encryptAESFileCmd.Flags().StringP("input-file", "", "", "Path to input file that needs to be encrypted")
	encryptAESFileCmd.Flags().StringP("output-file", "", "", "Path to encrypted output file")
	encryptAESFileCmd.Flags().StringP("symmetric-key", "", "", "Path to the symmetric key")
	rootCmd.AddCommand(encryptAESFileCmd)

	var decryptAESFileCmd = &cobra.Command{
		Use:   "decrypt-aes",
		Short: "Decrypt a file using AES",
		Run:   handler.DecryptAESCmd,
	}
	decryptAESFileCmd.Flags().StringP("input-file", "", "", "Input encrypted file path")
	decryptAESFileCmd.Flags().StringP("output-file", "", "", "Path to decrypted output file")
	decryptAESFileCmd.Flags().StringP("symmetric-key", "", "", "Path to the symmetric key")
	rootCmd.AddCommand(decryptAESFileCmd)
}
