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

type RSACommandHandler struct {
	rsaProcessor cryptography.RSAProcessor
	Logger       logger.Logger
}

func NewRSACommandHandler() *RSACommandHandler {
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

	rsaProcessor, err := cryptography.NewRSAProcessor(logger)
	if err != nil {
		log.Panicf("%v\n", err)
		return nil
	}

	return &RSACommandHandler{
		rsaProcessor: rsaProcessor,
		Logger:       logger,
	}
}

// GenerateRSAKeysCmd generates RSA key pairs and persists those in a selected directory
func (commandHandler *RSACommandHandler) GenerateRSAKeysCmd(cmd *cobra.Command, args []string) {
	keySize, _ := cmd.Flags().GetInt("key-size")
	keyDir, _ := cmd.Flags().GetString("key-dir")

	uniqueID := uuid.New()

	privateKey, publicKey, err := commandHandler.rsaProcessor.GenerateKeys(keySize)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	privateKeyFilePath := fmt.Sprintf("%s/%s-private-key.pem", keyDir, uniqueID.String())

	err = commandHandler.rsaProcessor.SavePrivateKeyToFile(privateKey, privateKeyFilePath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	publicKeyFilePath := fmt.Sprintf("%s/%s-public-key.pem", keyDir, uniqueID.String())
	err = commandHandler.rsaProcessor.SavePublicKeyToFile(publicKey, publicKeyFilePath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}
}

// EncryptRSACmd encrypts a file using RSA and saves asymmetric key pairs
func (commandHandler *RSACommandHandler) EncryptRSACmd(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input-file")
	outputFile, _ := cmd.Flags().GetString("output-file")
	publicKeyPath, _ := cmd.Flags().GetString("public-key")

	publicKey, err := commandHandler.rsaProcessor.ReadPublicKey(publicKeyPath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	plainText, err := os.ReadFile(filepath.Clean(inputFile))
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	encryptedData, err := commandHandler.rsaProcessor.Encrypt(plainText, publicKey)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	err = os.WriteFile(outputFile, encryptedData, 0600)
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

	privateKey, err := commandHandler.rsaProcessor.ReadPrivateKey(privateKeyPath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	encryptedData, err := os.ReadFile(filepath.Clean(inputFile))
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	decryptedData, err := commandHandler.rsaProcessor.Decrypt(encryptedData, privateKey)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	err = os.WriteFile(outputFile, decryptedData, 0600)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	commandHandler.Logger.Info(fmt.Sprintf("Decrypted data path %s", outputFile))
}

// SignRSACmd signs a file using RSA and saves the signature
func (commandHandler *RSACommandHandler) SignRSACmd(cmd *cobra.Command, args []string) {
	inputFilePath, _ := cmd.Flags().GetString("input-file")
	signatureFilePath, _ := cmd.Flags().GetString("output-file")
	privateKeyPath, _ := cmd.Flags().GetString("private-key")

	// Read private key
	privateKey, err := commandHandler.rsaProcessor.ReadPrivateKey(privateKeyPath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	// Read data to sign
	data, err := os.ReadFile(filepath.Clean(inputFilePath))
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	// Sign the data
	signature, err := commandHandler.rsaProcessor.Sign(data, privateKey)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	// Save the signature to a file
	err = os.WriteFile(signatureFilePath, signature, 0600)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	commandHandler.Logger.Info(fmt.Sprintf("Signature saved at %s", signatureFilePath))
}

// VerifyRSACmd verifies a signature using RSA
func (commandHandler *RSACommandHandler) VerifyRSACmd(cmd *cobra.Command, args []string) {
	inputFilePath, _ := cmd.Flags().GetString("input-file")
	signatureFilePath, _ := cmd.Flags().GetString("signature-file")
	publicKeyPath, _ := cmd.Flags().GetString("public-key")

	// Read public key
	publicKey, err := commandHandler.rsaProcessor.ReadPublicKey(publicKeyPath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	// Read data and signature
	data, err := os.ReadFile(filepath.Clean(inputFilePath))
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	signature, err := os.ReadFile(filepath.Clean(signatureFilePath))
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	// Verify the signature
	valid, err := commandHandler.rsaProcessor.Verify(data, signature, publicKey)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	if valid {
		commandHandler.Logger.Info("Signature is valid")
	} else {
		commandHandler.Logger.Error("Signature is invalid")
	}
}

func InitRSACommands(rootCmd *cobra.Command) {
	handler := NewRSACommandHandler()

	var generateRSAKeysCmd = &cobra.Command{
		Use:   "generate-rsa-keys",
		Short: "Generate RSA keys",
		Run:   handler.GenerateRSAKeysCmd,
	}
	generateRSAKeysCmd.Flags().IntP("key-size", "", 2048, "RSA key size (default 2048 bytes for RSA-2048)")
	generateRSAKeysCmd.Flags().StringP("key-dir", "", "", "Directory to store the RSA keys")
	rootCmd.AddCommand(generateRSAKeysCmd)

	var encryptRSAFileCmd = &cobra.Command{
		Use:   "encrypt-rsa",
		Short: "Encrypt a file using RSA",
		Run:   handler.EncryptRSACmd,
	}
	encryptRSAFileCmd.Flags().StringP("input-file", "", "", "Path to input file which needs to be encrypted")
	encryptRSAFileCmd.Flags().StringP("output-file", "", "", "Path to encrypted output file")
	encryptRSAFileCmd.Flags().StringP("public-key", "", "", "Path to RSA public private key")
	rootCmd.AddCommand(encryptRSAFileCmd)

	var decryptRSAFileCmd = &cobra.Command{
		Use:   "decrypt-rsa",
		Short: "Decrypt a file using RSA",
		Run:   handler.DecryptRSACmd,
	}
	decryptRSAFileCmd.Flags().StringP("input-file", "", "", "Path to encrypted file")
	decryptRSAFileCmd.Flags().StringP("output-file", "", "", "Path to decrypted output file")
	decryptRSAFileCmd.Flags().StringP("private-key", "", "", "Path to RSA private key")
	rootCmd.AddCommand(decryptRSAFileCmd)

	var signRSAFileCmd = &cobra.Command{
		Use:   "sign-rsa",
		Short: "Sign a file using RSA",
		Run:   handler.SignRSACmd,
	}

	signRSAFileCmd.Flags().StringP("input-file", "", "", "Path to file which needs to be signed")
	signRSAFileCmd.Flags().StringP("output-file", "", "", "Path to signature output file")
	signRSAFileCmd.Flags().StringP("private-key", "", "", "Path to RSA private key")
	rootCmd.AddCommand(signRSAFileCmd)

	var verifyRSAFileCmd = &cobra.Command{
		Use:   "verify-rsa",
		Short: "Verify a file is valid using RSA",
		Run:   handler.VerifyRSACmd,
	}

	verifyRSAFileCmd.Flags().StringP("input-file", "", "", "Path to file which needs to be validated")
	verifyRSAFileCmd.Flags().StringP("signature-file", "", "", "Path to signature input file")
	verifyRSAFileCmd.Flags().StringP("public-key", "", "", "Path to RSA public key")
	rootCmd.AddCommand(verifyRSAFileCmd)
}
