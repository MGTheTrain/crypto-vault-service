package commands

import (
	"crypto_vault_service/cmd/crypto-vault-cli/internal/status"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/infrastructure/utils"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// PKCS11CommandsHandler holds settings and methods for managing PKCS#11 token operations
type PKCS11CommandsHandler struct{}

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
	modulePath, err := cmd.Flags().GetString("module")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	soPin, err := cmd.Flags().GetString("so-pin")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	userPin, err := cmd.Flags().GetString("user-pin")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	slotId, err := cmd.Flags().GetString("slot-id")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	err = h.writePkcs11ConfigFile(modulePath, soPin, userPin, slotId)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	info := status.NewInfo("created pkcs11-settings.json")
	info.PrintJsonInfo(false)
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
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	tokens, err := tokenHandler.ListTokenSlots()
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	tokensJSON, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	info := status.NewInfo(string(tokensJSON))
	info.PrintJsonInfo(true)
}

// ListObjectsSlotsCmd lists PKCS#11 token objects
func (h *PKCS11CommandsHandler) ListObjectsSlotsCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	tokenLabel, err := cmd.Flags().GetString("token-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	objects, err := tokenHandler.ListObjects(tokenLabel)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	objectsJSON, err := json.MarshalIndent(objects, "", "  ")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	info := status.NewInfo(string(objectsJSON))
	info.PrintJsonInfo(true)
}

// InitializeTokenCmd initializes a PKCS#11 token
func (h *PKCS11CommandsHandler) InitializeTokenCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	tokenLabel, err := cmd.Flags().GetString("token-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	if err := tokenHandler.InitializeToken(tokenLabel); err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
}

// AddKeyCmd adds a key to the PKCS#11 token
func (h *PKCS11CommandsHandler) AddKeyCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	tokenLabel, err := cmd.Flags().GetString("token-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	objectLabel, err := cmd.Flags().GetString("object-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	keyType, err := cmd.Flags().GetString("key-type")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	keySize, err := cmd.Flags().GetUint("key-size")

	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	if err := tokenHandler.AddKey(tokenLabel, objectLabel, keyType, keySize); err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
}

// DeleteObjectCmd deletes an object (key) from the PKCS#11 token
func (h *PKCS11CommandsHandler) DeleteObjectCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	tokenLabel, err := cmd.Flags().GetString("token-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	objectType, err := cmd.Flags().GetString("object-type")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	objectLabel, err := cmd.Flags().GetString("object-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	if err := tokenHandler.DeleteObject(tokenLabel, objectType, objectLabel); err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
}

// EncryptCmd encrypts data using the PKCS#11 token
func (h *PKCS11CommandsHandler) EncryptCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	tokenLabel, err := cmd.Flags().GetString("token-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	objectLabel, err := cmd.Flags().GetString("object-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	inputFilePath, err := cmd.Flags().GetString("input-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	outputFilePath, err := cmd.Flags().GetString("output-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	keyType, err := cmd.Flags().GetString("key-type")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	if err := tokenHandler.Encrypt(tokenLabel, objectLabel, inputFilePath, outputFilePath, keyType); err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
}

// DecryptCmd decrypts data using the PKCS#11 token
func (h *PKCS11CommandsHandler) DecryptCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	tokenLabel, err := cmd.Flags().GetString("token-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	objectLabel, err := cmd.Flags().GetString("object-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	inputFilePath, err := cmd.Flags().GetString("input-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	outputFilePath, err := cmd.Flags().GetString("output-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	keyType, err := cmd.Flags().GetString("key-type")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	if err := tokenHandler.Decrypt(tokenLabel, objectLabel, inputFilePath, outputFilePath, keyType); err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
}

// SignCmd signs data using the PKCS#11 token
func (h *PKCS11CommandsHandler) SignCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	tokenLabel, err := cmd.Flags().GetString("token-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	objectLabel, err := cmd.Flags().GetString("object-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	dataFilePath, err := cmd.Flags().GetString("data-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	signatureFilePath, err := cmd.Flags().GetString("signature-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	keyType, err := cmd.Flags().GetString("key-type")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	if err := tokenHandler.Sign(tokenLabel, objectLabel, dataFilePath, signatureFilePath, keyType); err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
}

// VerifyCmd verifies the signature using the PKCS#11 token
func (h *PKCS11CommandsHandler) VerifyCmd(cmd *cobra.Command, args []string) {
	tokenHandler, err := h.getTokenHandler()
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	tokenLabel, err := cmd.Flags().GetString("token-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	objectLabel, err := cmd.Flags().GetString("object-label")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	dataFilePath, err := cmd.Flags().GetString("data-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	signatureFilePath, err := cmd.Flags().GetString("signature-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}
	keyType, err := cmd.Flags().GetString("key-type")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	if _, err := tokenHandler.Verify(tokenLabel, objectLabel, dataFilePath, signatureFilePath, keyType); err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
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