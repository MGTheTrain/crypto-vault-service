package commands

import (
	"crypto/rsa"

	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type RSACommandHandler struct {
	rsa    *cryptography.RSA
	Logger logger.Logger
}

func NewRSACommandHandler() *RSACommandHandler {
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

	rsa, err := cryptography.NewRSA(logger)
	if err != nil {
		log.Panicf("%v\n", err)
		return nil
	}

	return &RSACommandHandler{
		rsa:    rsa,
		Logger: logger,
	}
}

// EncryptRSACmd encrypts a file using RSA and saves asymmetric key pairs
func (commandHandler *RSACommandHandler) EncryptRSACmd(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	outputFile, _ := cmd.Flags().GetString("output-file")
	keyDir, _ := cmd.Flags().GetString("key-dir")

	var publicKey *rsa.PublicKey

	uniqueID := uuid.New()

	privateKey, publicKey, err := commandHandler.rsa.GenerateKeys(2048)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	privateKeyFilePath := fmt.Sprintf("%s/%s-private-key.pem", keyDir, uniqueID.String())

	err = commandHandler.rsa.SavePrivateKeyToFile(privateKey, privateKeyFilePath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	publicKeyFilePath := fmt.Sprintf("%s/%s-public-key.pem", keyDir, uniqueID.String())
	err = commandHandler.rsa.SavePublicKeyToFile(publicKey, publicKeyFilePath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	plainText, err := os.ReadFile(inputFile)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	encryptedData, err := commandHandler.rsa.Encrypt(plainText, publicKey)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	err = os.WriteFile(outputFile, encryptedData, 0644)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	commandHandler.Logger.Info(fmt.Sprintf("Encrypted data path %s", outputFile))
}

// DecryptRSACmd decrypts a file using RSA
func (commandHandler *RSACommandHandler) DecryptRSACmd(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	outputFile, _ := cmd.Flags().GetString("output-file")
	privateKeyPath, _ := cmd.Flags().GetString("private-key")

	var privateKey *rsa.PrivateKey

	if privateKeyPath == "" {

		privKey, _, err := commandHandler.rsa.GenerateKeys(2048)
		if err != nil {
			commandHandler.Logger.Error(fmt.Sprintf("%v", err))
			return
		}
		privateKey = privKey

		err = commandHandler.rsa.SavePrivateKeyToFile(privateKey, "private-key.pem")
		if err != nil {
			commandHandler.Logger.Error(fmt.Sprintf("%v", err))
			return
		}

	}
	privateKey, err := commandHandler.rsa.ReadPrivateKey(privateKeyPath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	encryptedData, err := os.ReadFile(inputFile)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	decryptedData, err := commandHandler.rsa.Decrypt(encryptedData, privateKey)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	err = os.WriteFile(outputFile, decryptedData, 0644)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	commandHandler.Logger.Info(fmt.Sprintf("Decrypted data path %s", outputFile))
}

func InitRSACommands(rootCmd *cobra.Command) {
	handler := NewRSACommandHandler()

	var encryptRSAFileCmd = &cobra.Command{
		Use:   "encrypt-rsa",
		Short: "Encrypt a file using RSA",
		Run:   handler.EncryptRSACmd,
	}
	encryptRSAFileCmd.Flags().StringP("input-file", "", "", "Input file path")
	encryptRSAFileCmd.Flags().StringP("output-file", "", "", "Output encrypted file path")
	encryptRSAFileCmd.Flags().StringP("key-dir", "", "", "Directory to store the encryption key")
	rootCmd.AddCommand(encryptRSAFileCmd)

	var decryptRSAFileCmd = &cobra.Command{
		Use:   "decrypt-rsa",
		Short: "Decrypt a file using RSA",
		Run:   handler.DecryptRSACmd,
	}
	decryptRSAFileCmd.Flags().StringP("input-file", "", "", "Input encrypted file path")
	decryptRSAFileCmd.Flags().StringP("output-file", "", "", "Output decrypted file path")
	decryptRSAFileCmd.Flags().StringP("private-key", "", "", "Path to RSA private key")
	rootCmd.AddCommand(decryptRSAFileCmd)
}
