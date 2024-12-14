package commands

import (
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/infrastructure/utils"
	"encoding/json"
	"fmt"
	"os"

	"log"

	"github.com/spf13/cobra"
)

// PKCS11CommandsHandler holds settings and methods for managing PKCS#11 token operations
type PKCS11CommandsHandler struct{}

// GetFlagString retrieves a flag value and logs an error if it is missing or invalid
func GetFlagString(cmd *cobra.Command, flagName, errMessage string) (string, error) {
	value, err := cmd.Flags().GetString(flagName)
	if err != nil || value == "" {
		return "", fmt.Errorf("%s: %v", errMessage, err)
	}
	return value, nil
}

// validatePKCS11Settings checks if the PKCS#11 settings (ModulePath, SOPin, UserPin, and SlotId)
func (h *PKCS11CommandsHandler) validatePKCS11Settings(tokenHandler *cryptography.PKCS11Handler) error {
	if err := utils.CheckNonEmptyStrings(
		tokenHandler.Settings.ModulePath,
		tokenHandler.Settings.SOPin,
		tokenHandler.Settings.UserPin,
		tokenHandler.Settings.SlotId); err != nil {
		return fmt.Errorf("ensure PKCS#11 settings have been configured trough `configure-pkcs11-settings` command: %v", err)
	}
	return nil
}

// readPkcs11ConfigFile reads the pkcs11-settings.json file and create the settings object
func (h *PKCS11CommandsHandler) readPkcs11ConfigFile() (*settings.PKCS11Settings, error) {
	plainText, err := os.ReadFile("pkcs11-settings.json")
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %s", err)
	}

	var settings settings.PKCS11Settings
	err = json.Unmarshal(plainText, &settings)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON into struct: %s", err)
	}

	return &settings, nil
}

// writePkcs11ConfigFile writes the pkcs11-settings.json config file
func (h *PKCS11CommandsHandler) writePkcs11ConfigFile(modulePath, soPin, userPin, slotId string) error {
	settings := map[string]string{
		"modulePath": modulePath,
		"soPin":      soPin,
		"userPin":    userPin,
		"slotId":     slotId,
	}

	settingsJSON, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling settings to JSON: %v", err)
	}

	file, err := os.Create("pkcs11-settings.json")
	if err != nil {
		return fmt.Errorf("error creating JSON file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(settingsJSON)
	if err != nil {
		return fmt.Errorf("error writing to JSON file: %v", err)
	}
	return nil
}

// storePKCS11SettingsCmd command saves the PKCS#11 settings to a JSON configuration file
func (h *PKCS11CommandsHandler) storePKCS11SettingsCmd(cmd *cobra.Command, args []string) {
	modulePath, err := GetFlagString(cmd, "module", "Error fetching module path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	soPin, err := GetFlagString(cmd, "so-pin", "Error fetching SO Pin flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	userPin, err := GetFlagString(cmd, "user-pin", "Error fetching user pin flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	slotId, err := GetFlagString(cmd, "slot-id", "Error fetching slot id flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	h.writePkcs11ConfigFile(modulePath, soPin, userPin, slotId)
	fmt.Println("created pkcs11-settings.json")
}

// getTokenHandler reads the PKCS#11 config file and validates the settings.
func (h *PKCS11CommandsHandler) getTokenHandler() (*cryptography.PKCS11Handler, error) {

	pkcs11Settings, err := h.readPkcs11ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("error reading PKCS#11 config file: %v", err)
	}

	tokenHandler := &cryptography.PKCS11Handler{
		Settings: pkcs11Settings,
	}

	err = h.validatePKCS11Settings(tokenHandler)
	if err != nil {
		return nil, fmt.Errorf("error validating PKCS#11 settings: %v", err)
	}

	return tokenHandler, nil
}

// ListTokenSlotsCmd lists PKCS#11 tokens
func (h *PKCS11CommandsHandler) ListTokenSlotsCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	tokens, err := tokenHandler.ListTokenSlots()
	if err != nil {
		log.Fatalf("Error initializing token: %v", err)
		return
	}

	tokensJSON, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling tokens to JSON: %v", err)
		return
	}

	fmt.Println(string(tokensJSON))
}

// ListObjectsSlotsCmd lists PKCS#11 token objects
func (h *PKCS11CommandsHandler) ListObjectsSlotsCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	tokenLabel, err := GetFlagString(cmd, "token-label", "Error fetching token-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	objects, err := tokenHandler.ListObjects(tokenLabel)
	if err != nil {
		log.Fatalf("Error initializing token: %v", err)
		return
	}

	objectsJSON, err := json.MarshalIndent(objects, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling objects to JSON: %v", err)
		return
	}

	fmt.Println(string(objectsJSON))
}

// InitializeTokenCmd initializes a PKCS#11 token
func (h *PKCS11CommandsHandler) InitializeTokenCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	tokenLabel, err := GetFlagString(cmd, "token-label", "Error fetching token-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	if err := tokenHandler.InitializeToken(tokenLabel); err != nil {
		log.Fatalf("Error initializing token: %v", err)
	}
}

// AddKeyCmd adds a key to the PKCS#11 token
func (h *PKCS11CommandsHandler) AddKeyCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	tokenLabel, err := GetFlagString(cmd, "token-label", "Error fetching token-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	objectLabel, err := GetFlagString(cmd, "object-label", "Error fetching object-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	keyType, err := GetFlagString(cmd, "key-type", "Error fetching key-type path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	keySize, err := cmd.Flags().GetUint("key-size")

	if err != nil {
		log.Panicf("Error fetching key-size path flag: %v", err)
	}

	if err := tokenHandler.AddKey(tokenLabel, objectLabel, keyType, keySize); err != nil {
		log.Fatalf("Error adding key: %v", err)
	}
}

// DeleteObjectCmd deletes an object (key) from the PKCS#11 token
func (h *PKCS11CommandsHandler) DeleteObjectCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	tokenLabel, err := GetFlagString(cmd, "token-label", "Error fetching token-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	objectType, err := GetFlagString(cmd, "object-type", "Error fetching object-type path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	objectLabel, err := GetFlagString(cmd, "object-label", "Error fetching object-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	if err := tokenHandler.DeleteObject(tokenLabel, objectType, objectLabel); err != nil {
		log.Fatalf("Error deleting object: %v", err)
	}
}

// EncryptCmd encrypts data using the PKCS#11 token
func (h *PKCS11CommandsHandler) EncryptCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	tokenLabel, err := GetFlagString(cmd, "token-label", "Error fetching token-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	objectLabel, err := GetFlagString(cmd, "object-label", "Error fetching object-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	inputFilePath, err := GetFlagString(cmd, "input-file", "Error input-file path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	outputFilePath, err := GetFlagString(cmd, "output-file", "Error output-file path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	keyType, err := GetFlagString(cmd, "key-type", "Error key-type path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	if err := tokenHandler.Encrypt(tokenLabel, objectLabel, inputFilePath, outputFilePath, keyType); err != nil {
		log.Fatalf("Error encrypting data: %v", err)
	}
}

// DecryptCmd decrypts data using the PKCS#11 token
func (h *PKCS11CommandsHandler) DecryptCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	tokenLabel, err := GetFlagString(cmd, "token-label", "Error fetching token-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	objectLabel, err := GetFlagString(cmd, "object-label", "Error fetching object-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	inputFilePath, err := GetFlagString(cmd, "input-file", "Error input-file path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	outputFilePath, err := GetFlagString(cmd, "output-file", "Error output-file path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	keyType, err := GetFlagString(cmd, "key-type", "Error key-type path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	if err := tokenHandler.Decrypt(tokenLabel, objectLabel, inputFilePath, outputFilePath, keyType); err != nil {
		log.Fatalf("Error decrypting data: %v", err)
	}
}

// SignCmd signs data using the PKCS#11 token
func (h *PKCS11CommandsHandler) SignCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	tokenLabel, err := GetFlagString(cmd, "token-label", "Error fetching token-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	objectLabel, err := GetFlagString(cmd, "object-label", "Error fetching object-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	dataFilePath, err := GetFlagString(cmd, "data-file", "Error data-file path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	signatureFilePath, err := GetFlagString(cmd, "signature-file", "Error signature-file path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	keyType, err := GetFlagString(cmd, "key-type", "Error key-type path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	if err := tokenHandler.Sign(tokenLabel, objectLabel, dataFilePath, signatureFilePath, keyType); err != nil {
		log.Fatalf("Error signing data: %v", err)
	}
}

// VerifyCmd verifies the signature using the PKCS#11 token
func (h *PKCS11CommandsHandler) VerifyCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	tokenLabel, err := GetFlagString(cmd, "token-label", "Error fetching token-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	objectLabel, err := GetFlagString(cmd, "object-label", "Error fetching object-label path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	dataFilePath, err := GetFlagString(cmd, "data-file", "Error data-file path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	signatureFilePath, err := GetFlagString(cmd, "signature-file", "Error signature-file path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	keyType, err := GetFlagString(cmd, "key-type", "Error key-type path flag")
	if err != nil {
		log.Fatalf("%v", err)
		return
	}

	if _, err := tokenHandler.Verify(tokenLabel, objectLabel, dataFilePath, signatureFilePath, keyType); err != nil {
		log.Fatalf("Error verifying signature: %v", err)
	}
}

// InitPKCS11Commands initializes all the PKCS#11 commands
func InitPKCS11Commands(rootCmd *cobra.Command) {
	handler := &PKCS11CommandsHandler{}

	var storePKCS11SettingsCmd = &cobra.Command{
		Use:   "store-pkcs11-settings",
		Short: "Stores PKCS#11 settings locally in the pkcs11-settings.json file",
		Run:   handler.storePKCS11SettingsCmd,
	}
	storePKCS11SettingsCmd.Flags().String("slot-id", "", "The token slot id")
	storePKCS11SettingsCmd.Flags().String("module", "", "Path to the PKCS#11 module")
	storePKCS11SettingsCmd.Flags().String("so-pin", "", "Security Officer PIN")
	storePKCS11SettingsCmd.Flags().String("user-pin", "", "User PIN")
	rootCmd.AddCommand(storePKCS11SettingsCmd)

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
