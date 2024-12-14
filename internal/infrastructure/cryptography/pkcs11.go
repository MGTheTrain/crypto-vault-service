package cryptography

import (
	"crypto_vault_service/internal/infrastructure/settings"
	"crypto_vault_service/internal/infrastructure/utils"
	"fmt"
	"os/exec"
	"strings"
)

// Token represents a PKCS#11 token with its label and other metadata.
type Token struct {
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

// IPKCS11TokenHandler defines the operations for working with a PKCS#11 token
type IPKCS11TokenHandler interface {
	// ListTokens lists all available tokens in the available slots
	ListTokens() ([]Token, error)
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
	Verify(label, objectLabel, keyType, dataFilePath, signatureFilePath string) (bool, error)
	// DeleteObject deletes a key or object from the token
	DeleteObject(label, objectType, objectLabel string) error
}

// PKCS11TokenHandler represents the parameters and operations for interacting with a PKCS#11 token
type PKCS11TokenHandler struct {
	Settings *settings.PKCS11Settings
}

// NewPKCS11TokenHandler creates and returns a new instance of PKCS11TokenHandler
func NewPKCS11TokenHandler(settings settings.PKCS11Settings) (*PKCS11TokenHandler, error) {
	if err := settings.Validate(); err != nil {
		return nil, err
	}

	return &PKCS11TokenHandler{
		Settings: &settings,
	}, nil
}

// Private method to execute pkcs11-tool commands and return output
func (token *PKCS11TokenHandler) executePKCS11ToolCommand(args []string) (string, error) {
	cmd := exec.Command("pkcs11-tool", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pkcs11-tool command failed: %v\nOutput: %s", err, output)
	}
	return string(output), nil
}

// ListTokens lists all available tokens in the available slots
func (token *PKCS11TokenHandler) ListTokens() ([]Token, error) {
	if err := utils.CheckNonEmptyStrings(token.Settings.ModulePath); err != nil {
		return nil, err
	}

	listCmd := exec.Command(
		"pkcs11-tool", "--module", token.Settings.ModulePath, "-L",
	)

	output, err := listCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens with pkcs11-tool: %v\nOutput: %s", err, output)
	}

	var tokens []Token
	lines := strings.Split(string(output), "\n")
	var currentToken *Token

	for _, line := range lines {

		if strings.Contains(line, "token label") {
			if currentToken != nil {
				tokens = append(tokens, *currentToken)
			}

			currentToken = &Token{}
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
func (token *PKCS11TokenHandler) ListObjects(tokenLabel string) ([]TokenObject, error) {
	//
	if err := utils.CheckNonEmptyStrings(tokenLabel, token.Settings.ModulePath); err != nil {
		return nil, err
	}

	listObjectsCmd := exec.Command(
		"pkcs11-tool", "--module", token.Settings.ModulePath, "-O", "--token-label", tokenLabel, "--pin", token.Settings.UserPin,
	)

	output, err := listObjectsCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list objects with pkcs11-tool: %v\nOutput: %s", err, output)
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
func (token *PKCS11TokenHandler) isTokenSet(label string) (bool, error) {
	if err := utils.CheckNonEmptyStrings(label); err != nil {
		return false, err
	}

	args := []string{"--module", token.Settings.ModulePath, "-T"}
	output, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return false, err
	}

	if strings.Contains(output, label) && strings.Contains(output, "token initialized") {
		fmt.Printf("Token with label '%s' exists.\n", label)
		return true, nil
	}

	fmt.Printf("Error: Token with label '%s' does not exist.\n", label)
	return false, nil
}

// InitializeToken initializes the token with the provided label and pins
func (token *PKCS11TokenHandler) InitializeToken(label string) error {
	if err := utils.CheckNonEmptyStrings(label); err != nil {
		return err
	}

	tokenExists, err := token.isTokenSet(label)
	if err != nil {
		return err
	}

	if tokenExists {
		fmt.Println("Skipping initialization. Token label exists.")
		return nil
	}

	args := []string{"--module", token.Settings.ModulePath, "--init-token", "--label", label, "--so-pin", token.Settings.SOPin, "--init-pin", "--pin", token.Settings.UserPin, "--slot", token.Settings.SlotId}
	_, err = token.executePKCS11ToolCommand(args)
	if err != nil {
		return fmt.Errorf("failed to initialize token with label '%s': %v", label, err)
	}

	fmt.Printf("Token with label '%s' initialized successfully.\n", label)
	return nil
}

// AddKey adds the selected key (ECDSA or RSA) to the token
func (token *PKCS11TokenHandler) AddKey(label, objectLabel, keyType string, keySize uint) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel, keyType); err != nil {
		return err
	}

	if keyType == "ECDSA" {
		return token.addECDSASignKey(label, objectLabel, keySize)
	} else if keyType == "RSA" {
		return token.addRSASignKey(label, objectLabel, keySize)
	} else {
		return fmt.Errorf("unsupported key type: %s", keyType)
	}
}

// addECDSASignKey adds an ECDSA signing key to the token
func (token *PKCS11TokenHandler) addECDSASignKey(label, objectLabel string, keySize uint) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel); err != nil {
		return err
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
		return fmt.Errorf("failed to add ECDSA key to token: %v", err)
	}

	fmt.Printf("ECDSA key with label '%s' added to token '%s'.\n", objectLabel, label)
	return nil
}

// addRSASignKey adds an RSA signing key to the token
func (token *PKCS11TokenHandler) addRSASignKey(label, objectLabel string, keySize uint) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel); err != nil {
		return err
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
		return fmt.Errorf("failed to add RSA key to token: %v", err)
	}

	fmt.Printf("RSA key with label '%s' added to token '%s'.\n", objectLabel, label)
	return nil
}

// Encrypt encrypts data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pkcs
func (token *PKCS11TokenHandler) Encrypt(label, objectLabel, inputFilePath, outputFilePath, keyType string) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel, inputFilePath, outputFilePath, keyType); err != nil {
		return err
	}

	if err := utils.CheckFilesExist(inputFilePath); err != nil {
		return err
	}

	if keyType != "RSA" {
		return fmt.Errorf("only RSA keys are supported for encryption")
	}

	// Prepare the URI to use PKCS#11 engine for accessing the public key
	keyURI := fmt.Sprintf("pkcs11:token=%s;object=%s;type=public;pin-value=%s", label, objectLabel, token.Settings.UserPin)

	// Run OpenSSL command to encrypt using the public key from the PKCS#11 token
	encryptCmd := exec.Command(
		"openssl", "pkeyutl", "-engine", "pkcs11", "-keyform", "engine", "-pubin", "-encrypt",
		"-inkey", keyURI, "-pkeyopt", "rsa_padding_mode:pkcs1", "-in", inputFilePath, "-out", outputFilePath,
	)

	// Execute the encryption command
	encryptOutput, err := encryptCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to encrypt data with OpenSSL: %v\nOutput: %s", err, encryptOutput)
	}

	fmt.Printf("Encryption successful. Encrypted data written to %s\n", outputFilePath)
	return nil
}

// Decrypt decrypts data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pkcs
func (token *PKCS11TokenHandler) Decrypt(label, objectLabel, inputFilePath, outputFilePath, keyType string) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel, inputFilePath, outputFilePath, keyType); err != nil {
		return err
	}

	if err := utils.CheckFilesExist(inputFilePath); err != nil {
		return err
	}

	if keyType != "RSA" {
		return fmt.Errorf("only RSA keys are supported for decryption")
	}

	// Prepare the URI to use PKCS#11 engine for accessing the private key
	keyURI := fmt.Sprintf("pkcs11:token=%s;object=%s;type=private;pin-value=%s", label, objectLabel, token.Settings.UserPin)

	// Run OpenSSL command to decrypt the data using the private key from the PKCS#11 token
	decryptCmd := exec.Command(
		"openssl", "pkeyutl", "-engine", "pkcs11", "-keyform", "engine", "-decrypt",
		"-inkey", keyURI, "-pkeyopt", "rsa_padding_mode:pkcs1", "-in", inputFilePath, "-out", outputFilePath,
	)

	// Execute the decryption command
	decryptOutput, err := decryptCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to decrypt data with OpenSSL: %v\nOutput: %s", err, decryptOutput)
	}

	fmt.Printf("Decryption successful. Decrypted data written to %s\n", outputFilePath)
	return nil
}

// Sign signs data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pss
func (token *PKCS11TokenHandler) Sign(label, objectLabel, dataFilePath, signatureFilePath, keyType string) error {
	if err := utils.CheckNonEmptyStrings(label, objectLabel, dataFilePath, signatureFilePath, keyType); err != nil {
		return err
	}

	if err := utils.CheckFilesExist(dataFilePath); err != nil {
		return err
	}

	if keyType != "RSA" && keyType != "ECDSA" {
		return fmt.Errorf("only RSA and ECDSA keys are supported for signing")
	}

	// Prepare the OpenSSL command based on key type
	var signCmd *exec.Cmd
	var signatureFormat string
	if keyType == "RSA" {
		signatureFormat = "rsa_padding_mode:pss"
		// Command for signing with RSA-PSS
		signCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-sign",
			"pkcs11:token="+label+";object="+objectLabel+";type=private;pin-value="+token.Settings.UserPin,
			"-sigopt", signatureFormat,
			"-sha384", // Use SHA-384
			"-out", signatureFilePath, dataFilePath,
		)
	} else if keyType == "ECDSA" {
		// Command for signing with ECDSA
		signCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-sign",
			"pkcs11:token="+label+";object="+objectLabel+";type=private;pin-value="+token.Settings.UserPin,
			"-sha384", // ECDSA typically uses SHA-384
			"-out", signatureFilePath, dataFilePath,
		)
	}

	// Execute the sign command
	signOutput, err := signCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to sign data: %v\nOutput: %s", err, signOutput)
	}

	fmt.Printf("Signing successful. Signature written to %s\n", signatureFilePath)
	return nil
}

// Verify verifies the signature of data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pss
func (token *PKCS11TokenHandler) Verify(label, objectLabel, keyType, dataFilePath, signatureFilePath string) (bool, error) {
	valid := false

	if err := utils.CheckNonEmptyStrings(label, objectLabel, keyType, dataFilePath, signatureFilePath); err != nil {
		return valid, err
	}

	if err := utils.CheckFilesExist(dataFilePath, signatureFilePath); err != nil {
		return valid, err
	}

	if keyType != "RSA" && keyType != "ECDSA" {
		return valid, fmt.Errorf("only RSA and ECDSA keys are supported for verification")
	}

	// Prepare the OpenSSL command based on key type
	var verifyCmd *exec.Cmd
	if keyType == "RSA" {
		// Command for verifying with RSA-PSS
		verifyCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-verify",
			"pkcs11:token="+label+";object="+objectLabel+";type=public;pin-value="+token.Settings.UserPin,
			"-sigopt", "rsa_padding_mode:pss",
			"-sha384", // Use SHA-384 for verification
			"-signature", signatureFilePath, "-binary", dataFilePath,
		)
	} else if keyType == "ECDSA" {
		// Command for verifying with ECDSA
		verifyCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-verify",
			"pkcs11:token="+label+";object="+objectLabel+";type=public;pin-value="+token.Settings.UserPin,
			"-sha384", // ECDSA typically uses SHA-384
			"-signature", signatureFilePath, "-binary", dataFilePath,
		)
	}

	// Execute the verify command
	verifyOutput, err := verifyCmd.CombinedOutput()
	if err != nil {
		return valid, fmt.Errorf("failed to verify signature: %v\nOutput: %s", err, verifyOutput)
	}

	// Check the output from OpenSSL to determine if the verification was successful
	if strings.Contains(string(verifyOutput), "Verified OK") {
		fmt.Println("Verification successful: The signature is valid.")
		valid = true
	} else {
		fmt.Println("Verification failed: The signature is invalid.")
	}

	return valid, nil
}

// DeleteObject deletes a key or object from the token
func (token *PKCS11TokenHandler) DeleteObject(label, objectType, objectLabel string) error {
	if err := utils.CheckNonEmptyStrings(label, objectType, objectLabel); err != nil {
		return err
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
		return fmt.Errorf("failed to delete object of type '%s' with label '%s': %v", objectType, objectLabel, err)
	}

	fmt.Printf("Object of type '%s' with label '%s' deleted successfully.\n", objectType, objectLabel)
	return nil
}
