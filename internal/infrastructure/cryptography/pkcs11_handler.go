package cryptography

import (
	"crypto_vault_service/internal/infrastructure/logger"
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/infrastructure/utils"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// Token represents a PKCS#11 token with its label and other metadata.
type Token struct {
	SlotID       string
	Label        string
	Manufacturer string
	Model        string
	SerialNumber string
}

// TokenObject represents a PKCS#11 object (e.g. public or private key) with metadata.
type TokenObject struct {
	Label  string
	Type   string // The type of the object (e.g. RSA, ECDSA)
	Usage  string // The usage of the object (e.g. encrypt, sign, decrypt)
	Access string // Access controls for the object (e.g. sensitive, always sensitive)
}

// PKCS11Handler defines the operations for working with a PKCS#11 token
type PKCS11Handler interface {
	// ListTokenSlots lists all available tokens in the available slots
	ListTokenSlots() ([]Token, error)
	// ListObjects lists all objects (e.g. keys) in a specific token based on the token label
	ListObjects(tokenLabel string) ([]TokenObject, error)
	// InitializeToken initializes the token with the provided label and pins
	InitializeToken(label string) error
	// AddKey adds the selected key (ECDSA or RSA) to the token
	AddKey(label, objectLabel, keyType string, keySize uint) error
	// Encrypt encrypts data using the cryptographic capabilities of the PKCS#11 token
	Encrypt(label, objectLabel, inputFilePath, outputFilePath, keyType string) error
	// Decrypt decrypts data using the cryptographic capabilities of the PKCS#11 token
	Decrypt(label, objectLabel, inputFilePath, outputFilePath, keyType string) error
	// Sign signs data using the cryptographic capabilities of the PKCS#11 token
	Sign(label, objectLabel, dataFilePath, signatureFilePath, keyType string) error
	// Verify verifies the signature of data using the cryptographic capabilities of the PKCS#11 token
	Verify(label, objectLabel, dataFilePath, signatureFilePath, keyType string) (bool, error)
	// DeleteObject deletes a key or object from the token
	DeleteObject(label, objectType, objectLabel string) error
}

// pkcs11Handler represents the parameters and operations for interacting with a PKCS#11 token
type pkcs11Handler struct {
	Settings *settings.PKCS11Settings
	Logger   logger.Logger
}

// NewPKCS11Handler creates and returns a new instance of PKCS11Handler
func NewPKCS11Handler(settings *settings.PKCS11Settings, logger logger.Logger) (PKCS11Handler, error) {
	if err := settings.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate settings: %w", err)
	}

	return &pkcs11Handler{
		Settings: settings,
		Logger:   logger,
	}, nil
}

// Private method to execute pkcs11-tool commands and return output
func (token *pkcs11Handler) executePKCS11ToolCommand(args []string) (string, error) {
	cmd := exec.Command("pkcs11-tool", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pkcs11-tool command failed: %w\nOutput: %s", err, output)
	}
	return string(output), nil
}

// ListTokenSlots lists all available tokens in the available slots
func (token *pkcs11Handler) ListTokenSlots() ([]Token, error) {
	if err := utils.CheckNonEmptyStrings(token.Settings.ModulePath); err != nil {
		return nil, fmt.Errorf("failed to check non-empty string for ModulePath='%s': %w", token.Settings.ModulePath, err)
	}

	// #nosec G204 -- TODO(MGTheTrain) validate all inputs used in exec.Command
	listCmd := exec.Command(
		"pkcs11-tool", "--module", token.Settings.ModulePath, "-L",
	)

	output, err := listCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens with pkcs11-tool: %w\nOutput: %s", err, output)
	}

	var tokens []Token
	lines := strings.Split(string(output), "\n")
	var currentToken *Token

	for _, line := range lines {

		if strings.Contains(line, "Slot") {
			if currentToken != nil {
				tokens = append(tokens, *currentToken)
			}

			currentToken = &Token{
				SlotID:       "",
				Label:        "",
				Manufacturer: "",
				Model:        "",
				SerialNumber: "",
			}

			re := regexp.MustCompile(`\((0x[0-9a-fA-F]+)\)`) // e.g. `(0x39e9d82d)` in `Slot 1 (0x39e9d82d): SoftHSM slot ID 0x39e9d82d`
			matches := re.FindStringSubmatch(line)
			currentToken.SlotID = matches[1]
		}

		if strings.Contains(line, "token label") {
			currentToken.Label = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if currentToken != nil {
			if strings.Contains(line, "token manufacturer") {
				currentToken.Manufacturer = strings.TrimSpace(strings.Split(line, ":")[1])
			}
			if strings.Contains(line, "token model") {
				currentToken.Model = strings.TrimSpace(strings.Split(line, ":")[1])
			}
			if strings.Contains(line, "serial num") {
				currentToken.SerialNumber = strings.TrimSpace(strings.Split(line, ":")[1])
			}
		}
	}

	if currentToken != nil {
		tokens = append(tokens, *currentToken)
	}

	return tokens, nil
}

// ListObjects lists all objects (e.g. keys) in a specific token based on the token label.
func (token *pkcs11Handler) ListObjects(tokenLabel string) ([]TokenObject, error) {
	if err := utils.CheckNonEmptyStrings(tokenLabel, token.Settings.ModulePath); err != nil {
		return nil, fmt.Errorf("failed to check non-empty strings for tokenLabel='%s' and ModulePath='%s': %w", tokenLabel, token.Settings.ModulePath, err)
	}

	// #nosec G204 -- TODO(MGTheTrain) validate all inputs used in exec.Command
	listObjectsCmd := exec.Command(
		"pkcs11-tool", "--module", token.Settings.ModulePath, "-O", "--token-label", tokenLabel, "--pin", token.Settings.UserPin,
	)

	output, err := listObjectsCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list objects with pkcs11-tool: %w\nOutput: %s", err, output)
	}

	var objects []TokenObject
	lines := strings.Split(string(output), "\n")
	var currentObject *TokenObject

	for _, line := range lines {

		if strings.Contains(line, "Private") || strings.Contains(line, "Public") || strings.Contains(line, "Secret") {
			if currentObject != nil {
				objects = append(objects, *currentObject)
			}

			currentObject = &TokenObject{
				Label:  "",
				Type:   "",
				Usage:  "",
				Access: "",
			}
			currentObject.Type = line
		}
		if strings.Contains(line, "label:") {
			currentObject.Label = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Usage:") {
			currentObject.Usage = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Access:") {
			currentObject.Access = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	if currentObject != nil {
		objects = append(objects, *currentObject)
	}

	return objects, nil
}

// isTokenSet checks if the token exists in the given module path
func (token *pkcs11Handler) isTokenSet(label string) (bool, error) {
	if err := utils.CheckNonEmptyStrings(label); err != nil {
		return false, fmt.Errorf("failed to check non-empty string for label='%s': %w", label, err)
	}

	args := []string{"--module", token.Settings.ModulePath, "-T"}
	output, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return false, fmt.Errorf("failed to execute PKCS#11 tool command with args=%v: %w", args, err)
	}

	if strings.Contains(output, label) && strings.Contains(output, "token initialized") {
		token.Logger.Info(fmt.Sprintf("Token with label '%s' exists.\n", label))
		return true, nil
	}

	token.Logger.Info(fmt.Sprintf("Token with label '%s' does not exist.\n", label))
	return false, nil
}

// InitializeToken initializes the token with the provided label and pins
func (token *pkcs11Handler) InitializeToken(label string) error {
	if err := utils.CheckNonEmptyStrings(label); err != nil {
		return fmt.Errorf("failed to check non-empty string for label='%s': %w", label, err)
	}

	tokenExists, err := token.isTokenSet(label)
	if err != nil {
		return fmt.Errorf("failed to check if token is set for label='%s': %w", label, err)
	}

	if tokenExists {
		return nil
	}

	args := []string{"--module", token.Settings.ModulePath, "--init-token", "--label", label, "--so-pin", token.Settings.SOPin, "--init-pin", "--pin", token.Settings.UserPin, "--slot", token.Settings.SlotID}
	_, err = token.executePKCS11ToolCommand(args)
	if err != nil {
		return fmt.Errorf("failed to initialize token with label '%s': %w", label, err)
	}

	token.Logger.Info(fmt.Sprintf("Token with label '%s' initialized successfully.\n", label))
	return nil
}

// AddKey adds the selected key (ECDSA or RSA) to the token
func (token *pkcs11Handler) AddKey(label, objectLabel, keyType string, keySize uint) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel, keyType); err != nil {
		return fmt.Errorf("failed to check non-empty strings for label='%s', objectLabel='%s', keyType='%s': %w", label, objectLabel, keyType, err)
	}

	switch keyType {
	case "ECDSA":
		return token.addECDSASignKey(label, objectLabel, keySize)
	case "RSA":
		return token.addRSASignKey(label, objectLabel, keySize)
	default:
		return fmt.Errorf("unsupported key type: %s", keyType)
	}
}

// addECDSASignKey adds an ECDSA signing key to the token
func (token *pkcs11Handler) addECDSASignKey(label, objectLabel string, keySize uint) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel); err != nil {
		return fmt.Errorf("failed to check non-empty strings for label='%s' and objectLabel='%s': %w", label, objectLabel, err)
	}

	// Generate the key pair (example using secp256r1)
	// Supported ECDSA key sizes and their corresponding elliptic curves
	ecdsaCurves := map[uint]string{
		256: "secp256r1",
		384: "secp384r1",
		521: "secp521r1",
	}

	curve, supported := ecdsaCurves[keySize]
	if !supported {
		return fmt.Errorf("ECDSA key size must be one of 256, 384, or 521 bits, but got %d", keySize)
	}

	args := []string{
		"--module", token.Settings.ModulePath,
		"--token-label", label,
		"--keypairgen",
		"--key-type", fmt.Sprintf("EC:%s", curve), // Use the dynamically selected curve
		"--label", objectLabel,
		"--pin", token.Settings.UserPin,
		"--usage-sign",
	}

	_, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return fmt.Errorf("failed to add ECDSA key to token: %w", err)
	}

	token.Logger.Info(fmt.Sprintf("ECDSA key with label '%s' added to token '%s'", objectLabel, label))
	return nil
}

// addRSASignKey adds an RSA signing key to the token
func (token *pkcs11Handler) addRSASignKey(label, objectLabel string, keySize uint) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel); err != nil {
		return fmt.Errorf("failed to check non-empty strings for label='%s' and objectLabel='%s': %w", label, objectLabel, err)
	}

	// Supported RSA key sizes (for example, 2048, 3072, and 4096)
	supportedRSASizes := []uint{2048, 3072, 4096}

	validKeySize := false
	for _, size := range supportedRSASizes {
		if keySize == size {
			validKeySize = true
			break
		}
	}

	if !validKeySize {
		return fmt.Errorf("RSA key size must be one of %v bits, but got %d", supportedRSASizes, keySize)
	}

	args := []string{
		"--module", token.Settings.ModulePath,
		"--token-label", label,
		"--keypairgen",
		"--key-type", fmt.Sprintf("RSA:%d", keySize),
		"--label", objectLabel,
		"--pin", token.Settings.UserPin,
		"--usage-sign",
	}
	_, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return fmt.Errorf("failed to add RSA key to token: %w", err)
	}

	token.Logger.Info(fmt.Sprintf("RSA key with label '%s' added to token '%s'", objectLabel, label))
	return nil
}

// Encrypt encrypts data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pkcs
func (token *pkcs11Handler) Encrypt(label, objectLabel, inputFilePath, outputFilePath, keyType string) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel, inputFilePath, outputFilePath, keyType); err != nil {
		return fmt.Errorf("failed to check non-empty strings for label='%s', objectLabel='%s', inputFilePath='%s', outputFilePath='%s', keyType='%s': %w",
			label, objectLabel, inputFilePath, outputFilePath, keyType, err)
	}

	if err := utils.CheckFilesExist(inputFilePath); err != nil {
		return fmt.Errorf("failed to check if input file exists (inputFilePath='%s'): %w", inputFilePath, err)
	}

	if keyType != "RSA" {
		return fmt.Errorf("only RSA keys are supported for encryption")
	}

	// Prepare the URI to use PKCS#11 engine for accessing the public key
	keyURI := fmt.Sprintf("pkcs11:token=%s;object=%s;type=public;pin-value=%s", label, objectLabel, token.Settings.UserPin)

	// Run OpenSSL command to encrypt using the public key from the PKCS#11 token
	// #nosec G204 -- TODO(MGTheTrain) validate all inputs used in exec.Command
	encryptCmd := exec.Command(
		"openssl", "pkeyutl", "-engine", "pkcs11", "-keyform", "engine", "-pubin", "-encrypt",
		"-inkey", keyURI, "-pkeyopt", "rsa_padding_mode:pkcs1", "-in", inputFilePath, "-out", outputFilePath,
	)

	// Execute the encryption command
	encryptOutput, err := encryptCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to encrypt data with OpenSSL: %w\nOutput: %s", err, encryptOutput)
	}

	token.Logger.Info(fmt.Sprintf("Encryption successful. Encrypted data written to %s", outputFilePath))
	return nil
}

// Decrypt decrypts data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pkcs
func (token *pkcs11Handler) Decrypt(label, objectLabel, inputFilePath, outputFilePath, keyType string) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel, inputFilePath, outputFilePath, keyType); err != nil {
		return fmt.Errorf("failed to check non-empty strings for label='%s', objectLabel='%s', inputFilePath='%s', outputFilePath='%s', keyType='%s': %w",
			label, objectLabel, inputFilePath, outputFilePath, keyType, err)
	}

	if err := utils.CheckFilesExist(inputFilePath); err != nil {
		return fmt.Errorf("failed to check if input file exists (inputFilePath='%s'): %w", inputFilePath, err)
	}

	if keyType != "RSA" {
		return fmt.Errorf("only RSA keys are supported for decryption")
	}

	// Prepare the URI to use PKCS#11 engine for accessing the private key
	keyURI := fmt.Sprintf("pkcs11:token=%s;object=%s;type=private;pin-value=%s", label, objectLabel, token.Settings.UserPin)

	// Run OpenSSL command to decrypt the data using the private key from the PKCS#11 token
	// #nosec G204 -- TODO(MGTheTrain) validate all inputs used in exec.Command
	decryptCmd := exec.Command(
		"openssl", "pkeyutl", "-engine", "pkcs11", "-keyform", "engine", "-decrypt",
		"-inkey", keyURI, "-pkeyopt", "rsa_padding_mode:pkcs1", "-in", inputFilePath, "-out", outputFilePath,
	)

	// Execute the decryption command
	decryptOutput, err := decryptCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to decrypt data with OpenSSL: %w\nOutput: %s", err, decryptOutput)
	}

	token.Logger.Info(fmt.Sprintf("Decryption successful. Decrypted data written to %s", outputFilePath))
	return nil
}

// Sign signs data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pss
func (token *pkcs11Handler) Sign(label, objectLabel, dataFilePath, signatureFilePath, keyType string) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel, dataFilePath, signatureFilePath, keyType); err != nil {
		return fmt.Errorf("failed to check non-empty strings for label='%s', objectLabel='%s', dataFilePath='%s', signatureFilePath='%s', keyType='%s': %w",
			label, objectLabel, dataFilePath, signatureFilePath, keyType, err)
	}

	if err := utils.CheckFilesExist(dataFilePath); err != nil {
		return fmt.Errorf("failed to check if file exists (dataFilePath='%s'): %w", dataFilePath, err)
	}

	if keyType != "RSA" && keyType != "ECDSA" {
		return fmt.Errorf("only RSA and ECDSA keys are supported for signing")
	}

	// Prepare the OpenSSL command based on key type
	var signCmd *exec.Cmd
	var signatureFormat string

	switch keyType {
	case "RSA":
		signatureFormat = "rsa_padding_mode:pss"
		// Command for signing with RSA-PSS
		// #nosec G204 -- TODO(MGTheTrain) validate all inputs used in exec.Command
		signCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-sign",
			"pkcs11:token="+label+";object="+objectLabel+";type=private;pin-value="+token.Settings.UserPin,
			"-sigopt", signatureFormat,
			"-sha384", // Use SHA-384
			"-out", signatureFilePath, dataFilePath,
		)
	case "ECDSA":
		// Command for signing with ECDSA
		// #nosec G204 -- TODO(MGTheTrain) validate all inputs used in exec.Command
		signCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-sign",
			"pkcs11:token="+label+";object="+objectLabel+";type=private;pin-value="+token.Settings.UserPin,
			"-sha384", // ECDSA typically uses SHA-384
			"-out", signatureFilePath, dataFilePath,
		)
	default:
		return fmt.Errorf("unsupported key type: %s", keyType)
	}

	// Execute the sign command
	signOutput, err := signCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to sign data: %w\nOutput: %s", err, signOutput)
	}

	token.Logger.Info(fmt.Sprintf("Signing successful. Signature written to %s", signatureFilePath))
	return nil
}

// Verify verifies the signature of data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pss
func (token *pkcs11Handler) Verify(label, objectLabel, dataFilePath, signatureFilePath, keyType string) (bool, error) {

	if err := utils.CheckNonEmptyStrings(label, objectLabel, keyType, dataFilePath, signatureFilePath); err != nil {
		return false, fmt.Errorf("failed to check non-empty strings for label='%s', objectLabel='%s', keyType='%s', dataFilePath='%s', signatureFilePath='%s': %w",
			label, objectLabel, keyType, dataFilePath, signatureFilePath, err)
	}

	if err := utils.CheckFilesExist(dataFilePath, signatureFilePath); err != nil {
		return false, fmt.Errorf("failed to check if files exist (dataFilePath='%s', signatureFilePath='%s'): %w",
			dataFilePath, signatureFilePath, err)
	}

	if keyType != "RSA" && keyType != "ECDSA" {
		return false, fmt.Errorf("only RSA and ECDSA keys are supported for verification")
	}

	// Prepare the OpenSSL command based on key type
	var verifyCmd *exec.Cmd

	switch keyType {
	case "RSA":
		// Command for verifying with RSA-PSS
		// #nosec G204 -- TODO(MGTheTrain) validate all inputs used in exec.Command
		verifyCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-verify",
			"pkcs11:token="+label+";object="+objectLabel+";type=public;pin-value="+token.Settings.UserPin,
			"-sigopt", "rsa_padding_mode:pss",
			"-sha384", // Use SHA-384 for verification
			"-signature", signatureFilePath, "-binary", dataFilePath,
		)
	case "ECDSA":
		// Command for verifying with ECDSA
		// #nosec G204 -- TODO(MGTheTrain) validate all inputs used in exec.Command
		verifyCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-verify",
			"pkcs11:token="+label+";object="+objectLabel+";type=public;pin-value="+token.Settings.UserPin,
			"-sha384", // ECDSA typically uses SHA-384
			"-signature", signatureFilePath, "-binary", dataFilePath,
		)
	default:
		return false, fmt.Errorf("unsupported key type: %s", keyType)
	}

	// Execute the verify command
	verifyOutput, err := verifyCmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to verify signature: %w\nOutput: %s", err, verifyOutput)
	}

	// Check the output from OpenSSL to determine if the verification was successful
	if strings.Contains(string(verifyOutput), "Verified OK") {
		token.Logger.Info("The signature is valid")
		return true, nil
	}
	token.Logger.Info("The signature is invalid")
	return false, nil
}

// DeleteObject deletes a key or object from the token
func (token *pkcs11Handler) DeleteObject(label, objectType, objectLabel string) error {
	if err := utils.CheckNonEmptyStrings(label, objectType, objectLabel); err != nil {
		return fmt.Errorf("failed to check non-empty strings for label='%s', objectType='%s', objectLabel='%s': %w", label, objectType, objectLabel, err)
	}

	// Ensure the object type is valid (privkey, pubkey, secrkey, cert, data)
	validObjectTypes := map[string]bool{
		"privkey": true,
		"pubkey":  true,
		"secrkey": true,
		"cert":    true,
		"data":    true,
	}

	if !validObjectTypes[objectType] {
		return fmt.Errorf("invalid object type '%s'. Valid types are privkey, pubkey, secrkey, cert, data", objectType)
	}

	args := []string{
		"--module", token.Settings.ModulePath,
		"--token-label", label,
		"--pin", token.Settings.UserPin,
		"--delete-object",
		"--type", objectType,
		"--label", objectLabel,
	}

	_, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return fmt.Errorf("failed to delete object of type '%s' with label '%s': %w", objectType, objectLabel, err)
	}

	token.Logger.Info(fmt.Sprintf("Object of type '%s' with label '%s' deleted successfully.\n", objectType, objectLabel))
	return nil
}
