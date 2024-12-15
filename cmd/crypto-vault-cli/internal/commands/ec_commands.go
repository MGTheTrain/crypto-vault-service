package commands

import (
	"crypto/ecdsa"
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

	factory := &logger.LoggerFactory{}

	logger, err := factory.NewLogger(loggerSettings)
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

// SignECCCmd signs the contents of a file with ECDSA
func (commandHandler *ECCommandHandler) SignECCCmd(cmd *cobra.Command, args []string) {

	inputFile, err := cmd.Flags().GetString("input-file")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	keyDir, err := cmd.Flags().GetString("key-dir")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	var privateKey *ecdsa.PrivateKey
	var publicKey *ecdsa.PublicKey

	privateKey, publicKey, err = commandHandler.ec.GenerateKeys(elliptic.P256())
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	fileContent, err := os.ReadFile(inputFile)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	signature, err := commandHandler.ec.Sign(fileContent, privateKey)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	uniqueID := uuid.New()

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

	signatureFilePath := fmt.Sprintf("%s/%s-signature.sig", keyDir, uniqueID.String())
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
	inputFilePath, err := cmd.Flags().GetString("input-file")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	publicKeyPath, err := cmd.Flags().GetString("public-key")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	signatureFile, err := cmd.Flags().GetString("signature-file")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	var publicKey *ecdsa.PublicKey

	if publicKeyPath == "" {
		commandHandler.Logger.Error("Public key is required for ECC signature verification")
		return
	} else {
		publicKey, err = commandHandler.ec.ReadPublicKey(publicKeyPath, elliptic.P256())
		if err != nil {
			commandHandler.Logger.Error(fmt.Sprintf("%v", err))
			return
		}
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

	var signECCMessageCmd = &cobra.Command{
		Use:   "sign-ecc",
		Short: "Sign a message using ECC",
		Run:   handler.SignECCCmd,
	}
	signECCMessageCmd.Flags().StringP("input-file", "", "", "Path to file that needs to be signed")
	signECCMessageCmd.Flags().StringP("key-dir", "", "", "Directory to save generated keys (optional)")
	rootCmd.AddCommand(signECCMessageCmd)

	var verifyECCSignatureCmd = &cobra.Command{
		Use:   "verify-ecc",
		Short: "Verify a signature using ECC",
		Run:   handler.VerifyECCCmd,
	}
	verifyECCSignatureCmd.Flags().StringP("input-file", "", "", "Path to ECC public key")
	verifyECCSignatureCmd.Flags().StringP("public-key", "", "", "The public key used to verify the signature")
	verifyECCSignatureCmd.Flags().StringP("signature-file", "", "", "Signature to verify (hex format)")
	rootCmd.AddCommand(verifyECCSignatureCmd)
}
