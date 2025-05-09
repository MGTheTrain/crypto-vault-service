package commands

import (
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// PKCS11CommandsHandler holds settings and methods for managing PKCS#11 token operations
type PKCS11CommandsHandler struct {
	pkcs11Handler cryptography.IPKCS11Handler
	Logger        logger.Logger
}

func NewPKCS11CommandsHandler() *PKCS11CommandsHandler {
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

	pkcs11Settings, err := ReadPkcs11SettingsFromEnv()

	if err != nil {
		logger.Error(fmt.Sprintf("%v", err))
		return nil
	}

	pkcs11Handler, err := cryptography.NewPKCS11Handler(pkcs11Settings, logger)
	if err != nil {
		log.Panicf("%v\n", err)
		return nil
	}

	return &PKCS11CommandsHandler{
		pkcs11Handler: pkcs11Handler,
		Logger:        logger,
	}
}

func ReadPkcs11SettingsFromEnv() (*settings.PKCS11Settings, error) {
	modulePath := os.Getenv("PKCS11_MODULE_PATH")
	if modulePath == "" {
		return nil, fmt.Errorf("environment variable PKCS11_MODULE_PATH is not set")
	}

	soPin := os.Getenv("PKCS11_SO_PIN")
	if soPin == "" {
		return nil, fmt.Errorf("environment variable PKCS11_SO_PIN is not set")
	}

	userPin := os.Getenv("PKCS11_USER_PIN")
	if userPin == "" {
		return nil, fmt.Errorf("environment variable PKCS11_USER_PIN is not set")
	}

	slotId := os.Getenv("PKCS11_SLOT_ID")
	if slotId == "" {
		return nil, fmt.Errorf("environment variable PKCS11_SLOT_ID is not set")
	}

	// Construct and return the settings struct
	settings := &settings.PKCS11Settings{
		ModulePath: modulePath,
		SOPin:      soPin,
		UserPin:    userPin,
		SlotId:     slotId,
	}

	return settings, nil
}

// ListTokenSlotsCmd lists PKCS#11 tokens
func (commandHandler *PKCS11CommandsHandler) ListTokenSlotsCmd(cmd *cobra.Command, args []string) {

	tokens, err := commandHandler.pkcs11Handler.ListTokenSlots()
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	tokensJSON, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	commandHandler.Logger.Info(string(tokensJSON))
}

// ListObjectsSlotsCmd lists PKCS#11 token objects
func (commandHandler *PKCS11CommandsHandler) ListObjectsSlotsCmd(cmd *cobra.Command, args []string) {
	tokenLabel, _ := cmd.Flags().GetString("token-label")

	objects, err := commandHandler.pkcs11Handler.ListObjects(tokenLabel)
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	objectsJSON, err := json.MarshalIndent(objects, "", "  ")
	if err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}

	commandHandler.Logger.Info(string(objectsJSON))
}

// InitializeTokenCmd initializes a PKCS#11 token
func (commandHandler *PKCS11CommandsHandler) InitializeTokenCmd(cmd *cobra.Command, args []string) {
	tokenLabel, _ := cmd.Flags().GetString("token-label")

	if err := commandHandler.pkcs11Handler.InitializeToken(tokenLabel); err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}
}

// AddKeyCmd adds a key to the PKCS#11 token
func (commandHandler *PKCS11CommandsHandler) AddKeyCmd(cmd *cobra.Command, args []string) {
	tokenLabel, _ := cmd.Flags().GetString("token-label")
	objectLabel, _ := cmd.Flags().GetString("object-label")
	keyType, _ := cmd.Flags().GetString("key-type")
	keySize, _ := cmd.Flags().GetUint("key-size")

	if err := commandHandler.pkcs11Handler.AddKey(tokenLabel, objectLabel, keyType, keySize); err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}
}

// DeleteObjectCmd deletes an object (key) from the PKCS#11 token
func (commandHandler *PKCS11CommandsHandler) DeleteObjectCmd(cmd *cobra.Command, args []string) {
	tokenLabel, _ := cmd.Flags().GetString("token-label")
	objectType, _ := cmd.Flags().GetString("object-type")
	objectLabel, _ := cmd.Flags().GetString("object-label")

	if err := commandHandler.pkcs11Handler.DeleteObject(tokenLabel, objectType, objectLabel); err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}
}

// EncryptCmd encrypts data using the PKCS#11 token
func (commandHandler *PKCS11CommandsHandler) EncryptCmd(cmd *cobra.Command, args []string) {
	tokenLabel, _ := cmd.Flags().GetString("token-label")
	objectLabel, _ := cmd.Flags().GetString("object-label")
	inputFilePath, _ := cmd.Flags().GetString("input-file")
	outputFilePath, _ := cmd.Flags().GetString("output-file")
	keyType, _ := cmd.Flags().GetString("key-type")

	if err := commandHandler.pkcs11Handler.Encrypt(tokenLabel, objectLabel, inputFilePath, outputFilePath, keyType); err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}
}

// DecryptCmd decrypts data using the PKCS#11 token
func (commandHandler *PKCS11CommandsHandler) DecryptCmd(cmd *cobra.Command, args []string) {
	tokenLabel, _ := cmd.Flags().GetString("token-label")
	objectLabel, _ := cmd.Flags().GetString("object-label")
	inputFilePath, _ := cmd.Flags().GetString("input-file")
	outputFilePath, _ := cmd.Flags().GetString("output-file")
	keyType, _ := cmd.Flags().GetString("key-type")

	if err := commandHandler.pkcs11Handler.Decrypt(tokenLabel, objectLabel, inputFilePath, outputFilePath, keyType); err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}
}

// SignCmd signs data using the PKCS#11 token
func (commandHandler *PKCS11CommandsHandler) SignCmd(cmd *cobra.Command, args []string) {
	tokenLabel, _ := cmd.Flags().GetString("token-label")
	objectLabel, _ := cmd.Flags().GetString("object-label")
	dataFilePath, _ := cmd.Flags().GetString("data-file")
	signatureFilePath, _ := cmd.Flags().GetString("signature-file")
	keyType, _ := cmd.Flags().GetString("key-type")

	if err := commandHandler.pkcs11Handler.Sign(tokenLabel, objectLabel, dataFilePath, signatureFilePath, keyType); err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}
}

// VerifyCmd verifies the signature using the PKCS#11 token
func (commandHandler *PKCS11CommandsHandler) VerifyCmd(cmd *cobra.Command, args []string) {
	tokenLabel, _ := cmd.Flags().GetString("token-label")
	objectLabel, _ := cmd.Flags().GetString("object-label")
	dataFilePath, _ := cmd.Flags().GetString("data-file")
	signatureFilePath, _ := cmd.Flags().GetString("signature-file")
	keyType, _ := cmd.Flags().GetString("key-type")

	if _, err := commandHandler.pkcs11Handler.Verify(tokenLabel, objectLabel, dataFilePath, signatureFilePath, keyType); err != nil {
		commandHandler.Logger.Error(fmt.Sprintf("%v", err))
		return
	}
}

// InitPKCS11Commands initializes all the PKCS#11 commands
func InitPKCS11Commands(rootCmd *cobra.Command) {
	handler := NewPKCS11CommandsHandler()

	var pkcs11InitializeTokenCmd = &cobra.Command{
		Use:   "initialize-token",
		Short: "Initialize a PKCS#11 token",
		Run:   handler.InitializeTokenCmd,
	}
	pkcs11InitializeTokenCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	rootCmd.AddCommand(pkcs11InitializeTokenCmd)

	var listTokenSlotsCmd = &cobra.Command{
		Use:   "list-slots",
		Short: "List PKCS#11 token slots",
		Run:   handler.ListTokenSlotsCmd,
	}
	rootCmd.AddCommand(listTokenSlotsCmd)

	var listObjectsSlotsCmd = &cobra.Command{
		Use:   "list-objects",
		Short: "List PKCS#11 token objects",
		Run:   handler.ListObjectsSlotsCmd,
	}
	listObjectsSlotsCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	rootCmd.AddCommand(listObjectsSlotsCmd)

	var pkcs11AddKeyCmd = &cobra.Command{
		Use:   "add-key",
		Short: "Add key (ECDSA or RSA) to the PKCS#11 token",
		Run:   handler.AddKeyCmd,
	}
	pkcs11AddKeyCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11AddKeyCmd.Flags().String("object-label", "", "Label of the object (key)")
	pkcs11AddKeyCmd.Flags().String("key-type", "", "Type of the key (ECDSA or RSA)")
	pkcs11AddKeyCmd.Flags().Uint("key-size", 0, "Key size in bits (2048 for RSA, 256 for ECDSA)")
	rootCmd.AddCommand(pkcs11AddKeyCmd)

	var pkcs11DeleteObjectCmd = &cobra.Command{
		Use:   "delete-object",
		Short: "Delete an object (key) from the PKCS#11 token",
		Run:   handler.DeleteObjectCmd,
	}
	pkcs11DeleteObjectCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11DeleteObjectCmd.Flags().String("object-label", "", "Label of the object to delete")
	pkcs11DeleteObjectCmd.Flags().String("object-type", "", "Type of the object (e.g., privkey, pubkey, cert)")
	rootCmd.AddCommand(pkcs11DeleteObjectCmd)

	// --------------------------- Encryption Command ---------------------------
	var pkcs11EncryptCmd = &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt data using a PKCS#11 token",
		Run:   handler.EncryptCmd,
	}
	pkcs11EncryptCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11EncryptCmd.Flags().String("object-label", "", "Label of the object (key) for encryption")
	pkcs11EncryptCmd.Flags().String("key-type", "", "Type of the key (ECDSA or RSA)")
	pkcs11EncryptCmd.Flags().String("input-file", "", "Path to the unencrypted input file")
	pkcs11EncryptCmd.Flags().String("output-file", "", "Path to encrypted output file")
	rootCmd.AddCommand(pkcs11EncryptCmd)

	// --------------------------- Decryption Command ---------------------------
	var pkcs11DecryptCmd = &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt data using a PKCS#11 token",
		Run:   handler.DecryptCmd,
	}
	pkcs11DecryptCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11DecryptCmd.Flags().String("object-label", "", "Label of the object (key) for decryption")
	pkcs11DecryptCmd.Flags().String("key-type", "", "Type of the key (ECDSA or RSA)")
	pkcs11DecryptCmd.Flags().String("input-file", "", "Path to the encrypted input file")
	pkcs11DecryptCmd.Flags().String("output-file", "", "Path to decrypted output file")
	rootCmd.AddCommand(pkcs11DecryptCmd)

	// --------------------------- Signature Command ---------------------------
	var pkcs11SignCmd = &cobra.Command{
		Use:   "sign",
		Short: "Sign data using a PKCS#11 token",
		Run:   handler.SignCmd,
	}
	pkcs11SignCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11SignCmd.Flags().String("object-label", "", "Label of the object (key) for signing")
	pkcs11SignCmd.Flags().String("key-type", "", "Type of the key (ECDSA or RSA)")
	pkcs11SignCmd.Flags().String("data-file", "", "Path to the input file to be signed")
	pkcs11SignCmd.Flags().String("signature-file", "", "Path to store the signature output file")
	rootCmd.AddCommand(pkcs11SignCmd)

	// --------------------------- Verify Command ---------------------------
	var pkcs11VerifyCmd = &cobra.Command{
		Use:   "verify",
		Short: "Verify the signature using a PKCS#11 token",
		Run:   handler.VerifyCmd,
	}
	pkcs11VerifyCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11VerifyCmd.Flags().String("object-label", "", "Label of the object (key) for signature verification")
	pkcs11VerifyCmd.Flags().String("key-type", "", "Type of the key (ECDSA or RSA)")
	pkcs11VerifyCmd.Flags().String("data-file", "", "Path to the input file to verify the signature")
	pkcs11VerifyCmd.Flags().String("signature-file", "", "Path to the signature file used for signature verifying")
	rootCmd.AddCommand(pkcs11VerifyCmd)
}
