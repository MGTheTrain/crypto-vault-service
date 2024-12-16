package commands

import (
	"crypto/elliptic"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type ECCommandHandler struct {
	ec     *cryptography.EC
	Logger logger.Logger
}

func NewECCommandHandler() *ECCommandHandler {
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

	ec, err := cryptography.NewEC(logger)
	if err != nil {
		log.Panicf("%v\n", err)
		return nil
	}

	return &ECCommandHandler{
		ec:     ec,
		Logger: logger,
	}
}

// GenerateECKeysCmd generates EC key pairs and persists those in a selected directory
func (commandHandler *ECCommandHandler) GenerateECKeysCmd(cmd *cobra.Command, args []string) {
	keySize, _ := cmd.Flags().GetInt("key-size")
	keyDir, _ := cmd.Flags().GetString("key-dir")

	uniqueID := uuid.New()

	var curve elliptic.Curve
	if keySize == 224 {
		curve = elliptic.P224()
	} else if keySize == 256 {
		curve = elliptic.P256()
	} else if keySize == 384 {
		curve = elliptic.P384()
	} else if keySize == 521 {
		curve = elliptic.P521()
	} else {
		commandHandler.Logger.Error(fmt.Sprintf("key size %v not supported", keySize))
		return
	}
	privateKey, publicKey, err := commandHandler.ec.GenerateKeys(curve)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	privateKeyFilePath := fmt.Sprintf("%s/%s-private-key.pem", keyDir, uniqueID.String())
	err = commandHandler.ec.SavePrivateKeyToFile(privateKey, privateKeyFilePath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	publicKeyFilePath := fmt.Sprintf("%s/%s-public-key.pem", keyDir, uniqueID.String())
	err = commandHandler.ec.SavePublicKeyToFile(publicKey, publicKeyFilePath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}
}

// SignECCCmd signs the contents of a file with ECDSA
func (commandHandler *ECCommandHandler) SignECCCmd(cmd *cobra.Command, args []string) {
	inputFilePath, _ := cmd.Flags().GetString("input-file")
	privateKeyFilePath, _ := cmd.Flags().GetString("private-key")
	signatureFilePath, _ := cmd.Flags().GetString("output-file")

	fileContent, err := os.ReadFile(inputFilePath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	privateKey, err := commandHandler.ec.ReadPrivateKey(privateKeyFilePath, elliptic.P256())
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	signature, err := commandHandler.ec.Sign(fileContent, privateKey)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	err = commandHandler.ec.SaveSignatureToFile(signatureFilePath, signature)
	if err != nil {
		if err != nil {
			commandHandler.Logger.Error(fmt.Sprintf("%v", err))
			return
		}
	}
}

// verifyECCCmd verifies the signature of a file's content using ECDSA
func (commandHandler *ECCommandHandler) VerifyECCCmd(cmd *cobra.Command, args []string) {
	inputFilePath, _ := cmd.Flags().GetString("input-file")
	publicKeyPath, _ := cmd.Flags().GetString("public-key")
	signatureFile, _ := cmd.Flags().GetString("signature-file")

	publicKey, err := commandHandler.ec.ReadPublicKey(publicKeyPath, elliptic.P256())
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	fileContent, err := os.ReadFile(inputFilePath)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	signatureHex, err := os.ReadFile(signatureFile)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	signature, err := hex.DecodeString(string(signatureHex))
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	valid, err := commandHandler.ec.Verify(fileContent, signature, publicKey)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	if valid {
		commandHandler.Logger.Info(fmt.Sprintf("Signature valid for %s", inputFilePath))
	} else {
		commandHandler.Logger.Info(fmt.Sprintf("Signature invalid for %s", inputFilePath))
	}
}

func InitECDSACommands(rootCmd *cobra.Command) {
	handler := NewECCommandHandler()

	var generateECKeysCmd = &cobra.Command{
		Use:   "generate-ecc-keys",
		Short: "Generate ECC keys",
		Run:   handler.GenerateECKeysCmd,
	}
	generateECKeysCmd.Flags().IntP("key-size", "", 256, "ECC key size (default 256 bytes for ECC-256)")
	generateECKeysCmd.Flags().StringP("key-dir", "", "", "Directory to store the ECC keys")
	rootCmd.AddCommand(generateECKeysCmd)

	var signECCMessageCmd = &cobra.Command{
		Use:   "sign-ecc",
		Short: "Sign a message using ECC",
		Run:   handler.SignECCCmd,
	}
	signECCMessageCmd.Flags().StringP("input-file", "", "", "Path to file that needs to be signed")
	signECCMessageCmd.Flags().StringP("private-key", "", "", "Path to ECC private key")
	signECCMessageCmd.Flags().StringP("output-file", "", "", "Path to signature output file")
	rootCmd.AddCommand(signECCMessageCmd)

	var verifyECCSignatureCmd = &cobra.Command{
		Use:   "verify-ecc",
		Short: "Verify a signature using ECC",
		Run:   handler.VerifyECCCmd,
	}
	verifyECCSignatureCmd.Flags().StringP("input-file", "", "", "Path to file which needs to be validated")
	verifyECCSignatureCmd.Flags().StringP("public-key", "", "", "Path to ECC public key")
	verifyECCSignatureCmd.Flags().StringP("signature-file", "", "", "Path to signature input file")
	rootCmd.AddCommand(verifyECCSignatureCmd)
}
