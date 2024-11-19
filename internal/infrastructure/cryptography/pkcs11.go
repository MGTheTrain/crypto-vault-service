package cryptography

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PKCS11Token defines the operations for working with a PKCS#11 token
type PKCS11Token interface {
	IsTokenSet() (bool, error)
	IsObjectSet() (bool, error)
	InitializeToken(slot string) error
	AddKey() error
	Encrypt(inputFilePath, outputFilePath string) error
	Decrypt(inputFilePath, outputFilePath string) error
	Sign(inputFilePath, outputFilePath string) error
	Verify(dataFilePath, signatureFilePath string) (bool, error)
	DeleteObject(objectType, objectLabel string) error
}

// PKCS11TokenImpl represents the parameters and operations for interacting with a PKCS#11 token
type PKCS11TokenImpl struct {
	ModulePath  string
	Label       string
	SOPin       string
	UserPin     string
	ObjectLabel string
	KeyType     string // "ECDSA" or "RSA"
	KeySize     int    // Key size in bits for RSA or ECDSA (e.g., 256 for ECDSA, 2048 for RSA)
}

// Public method to execute pkcs11-tool commands and return output
func (token *PKCS11TokenImpl) executePKCS11ToolCommand(args []string) (string, error) {
	cmd := exec.Command("pkcs11-tool", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pkcs11-tool command failed: %v\nOutput: %s", err, output)
	}
	return string(output), nil
}

// IsTokenSet checks if the token exists in the given module path
func (token *PKCS11TokenImpl) IsTokenSet() (bool, error) {
	if token.ModulePath == "" || token.Label == "" {
		return false, fmt.Errorf("missing module path or token label")
	}

	args := []string{"--module", token.ModulePath, "-T"}
	output, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return false, err
	}

	if strings.Contains(output, token.Label) && strings.Contains(output, "token initialized") {
		fmt.Printf("Token with label '%s' exists.\n", token.Label)
		return true, nil
	}

	fmt.Printf("Error: Token with label '%s' does not exist.\n", token.Label)
	return false, nil
}

// IsObjectSet checks if the specified object exists on the given token
func (token *PKCS11TokenImpl) IsObjectSet() (bool, error) {
	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
		return false, fmt.Errorf("missing required arguments")
	}

	args := []string{"-O", "--module", token.ModulePath, "--token-label", token.Label, "--pin", token.UserPin}
	output, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return false, err
	}

	if strings.Contains(output, token.ObjectLabel) {
		fmt.Printf("Object with label '%s' exists.\n", token.ObjectLabel)
		return true, nil
	}

	fmt.Printf("Error: Object with label '%s' does not exist.\n", token.ObjectLabel)
	return false, nil
}

// InitializeToken initializes the token with the provided label and pins
func (token *PKCS11TokenImpl) InitializeToken(slot string) error {
	if token.ModulePath == "" || token.Label == "" || token.SOPin == "" || token.UserPin == "" || slot == "" {
		return fmt.Errorf("missing required parameters for token initialization")
	}

	// Check if the token is already initialized
	tokenExists, err := token.IsTokenSet()
	if err != nil {
		return err
	}

	if tokenExists {
		fmt.Println("Skipping initialization. Token label exists.")
		return nil
	}

	// Initialize the token
	args := []string{"--module", token.ModulePath, "--init-token", "--label", token.Label, "--so-pin", token.SOPin, "--init-pin", "--pin", token.UserPin, "--slot", slot}
	_, err = token.executePKCS11ToolCommand(args)
	if err != nil {
		return fmt.Errorf("failed to initialize token with label '%s': %v", token.Label, err)
	}

	fmt.Printf("Token with label '%s' initialized successfully.\n", token.Label)
	return nil
}

// DeleteObject deletes a key or object from the token
func (token *PKCS11TokenImpl) DeleteObject(objectType, objectLabel string) error {
	if token.ModulePath == "" || token.Label == "" || objectLabel == "" || token.UserPin == "" {
		return fmt.Errorf("missing required arguments to delete object")
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

	// Execute the pkcs11-tool command to delete the object
	args := []string{
		"--module", token.ModulePath,
		"--token-label", token.Label,
		"--pin", token.UserPin,
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

// AddKey adds the selected key (ECDSA or RSA) to the token
func (token *PKCS11TokenImpl) AddKey() error {
	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
		return fmt.Errorf("missing required arguments")
	}

	// Determine key type and call the appropriate function to generate the key
	if token.KeyType == "ECDSA" {
		return token.addECDSASignKey()
	} else if token.KeyType == "RSA" {
		return token.addRSASignKey()
	} else {
		return fmt.Errorf("unsupported key type: %s", token.KeyType)
	}
}

// addECDSASignKey adds an ECDSA signing key to the token
func (token *PKCS11TokenImpl) addECDSASignKey() error {
	if token.KeySize != 256 && token.KeySize != 384 && token.KeySize != 521 {
		return fmt.Errorf("ECDSA key size must be one of 256, 384, or 521 bits, but got %d", token.KeySize)
	}

	// Generate the key pair (example using secp256r1)
	// Supported ECDSA key sizes and their corresponding elliptic curves
	ecdsaCurves := map[int]string{
		256: "secp256r1",
		384: "secp384r1",
		521: "secp521r1",
	}

	curve, supported := ecdsaCurves[token.KeySize]
	if !supported {
		return fmt.Errorf("ECDSA key size must be one of 256, 384, or 521 bits, but got %d", token.KeySize)
	}

	// Generate the key pair using the correct elliptic curve
	args := []string{
		"--module", token.ModulePath,
		"--token-label", token.Label,
		"--keypairgen",
		"--key-type", fmt.Sprintf("EC:%s", curve), // Use the dynamically selected curve
		"--label", token.ObjectLabel,
		"--pin", token.UserPin,
		"--usage-sign",
	}

	_, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return fmt.Errorf("failed to add ECDSA key to token: %v", err)
	}

	fmt.Printf("ECDSA key with label '%s' added to token '%s'.\n", token.ObjectLabel, token.Label)
	return nil
}

// addRSASignKey adds an RSA signing key to the token
func (token *PKCS11TokenImpl) addRSASignKey() error {
	// Supported RSA key sizes (for example, 2048, 3072, and 4096)
	supportedRSASizes := []int{2048, 3072, 4096}

	validKeySize := false
	for _, size := range supportedRSASizes {
		if token.KeySize == size {
			validKeySize = true
			break
		}
	}

	if !validKeySize {
		return fmt.Errorf("RSA key size must be one of %v bits, but got %d", supportedRSASizes, token.KeySize)
	}

	args := []string{
		"--module", token.ModulePath,
		"--token-label", token.Label,
		"--keypairgen",
		"--key-type", fmt.Sprintf("RSA:%d", token.KeySize),
		"--label", token.ObjectLabel,
		"--pin", token.UserPin,
		"--usage-sign",
	}
	_, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return fmt.Errorf("failed to add RSA key to token: %v", err)
	}

	fmt.Printf("RSA key with label '%s' added to token '%s'.\n", token.ObjectLabel, token.Label)
	return nil
}

// Encrypt encrypts data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pkcs
func (token *PKCS11TokenImpl) Encrypt(inputFilePath, outputFilePath string) error {
	// Validate required parameters
	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
		return fmt.Errorf("missing required arguments for encryption")
	}

	if token.KeyType != "RSA" {
		return fmt.Errorf("only RSA keys are supported for encryption")
	}

	// Step 1: Prepare the URI to use PKCS#11 engine for accessing the public key
	keyURI := fmt.Sprintf("pkcs11:token=%s;object=%s;type=public;pin-value=%s", token.Label, token.ObjectLabel, token.UserPin)

	// Step 2: Run OpenSSL command to encrypt using the public key from the PKCS#11 token
	encryptCmd := exec.Command(
		"openssl", "pkeyutl", "-engine", "pkcs11", "-keyform", "engine", "-pubin", "-encrypt",
		"-inkey", keyURI, "-pkeyopt", "rsa_padding_mode:pkcs1", "-in", inputFilePath, "-out", outputFilePath,
	)

	// Execute the encryption command
	encryptOutput, err := encryptCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to encrypt data with OpenSSL: %v\nOutput: %s", err, encryptOutput)
	}

	// Output success message
	fmt.Printf("Encryption successful. Encrypted data written to %s\n", outputFilePath)
	return nil
}

// Decrypt decrypts data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pkcs
func (token *PKCS11TokenImpl) Decrypt(inputFilePath, outputFilePath string) error {
	// Validate required parameters
	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
		return fmt.Errorf("missing required arguments for decryption")
	}

	if token.KeyType != "RSA" {
		return fmt.Errorf("only RSA keys are supported for decryption")
	}

	// Check if input file exists
	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %v", err)
	}

	// Step 1: Prepare the URI to use PKCS#11 engine for accessing the private key
	keyURI := fmt.Sprintf("pkcs11:token=%s;object=%s;type=private;pin-value=%s", token.Label, token.ObjectLabel, token.UserPin)

	// Step 2: Run OpenSSL command to decrypt the data using the private key from the PKCS#11 token
	decryptCmd := exec.Command(
		"openssl", "pkeyutl", "-engine", "pkcs11", "-keyform", "engine", "-decrypt",
		"-inkey", keyURI, "-pkeyopt", "rsa_padding_mode:pkcs1", "-in", inputFilePath, "-out", outputFilePath,
	)

	// Execute the decryption command
	decryptOutput, err := decryptCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to decrypt data with OpenSSL: %v\nOutput: %s", err, decryptOutput)
	}

	// Output success message
	fmt.Printf("Decryption successful. Decrypted data written to %s\n", outputFilePath)
	return nil
}

// Sign signs data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pss
func (token *PKCS11TokenImpl) Sign(inputFilePath, outputFilePath string) error {
	// Validate required parameters
	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
		return fmt.Errorf("missing required arguments for signing")
	}

	if token.KeyType != "RSA" && token.KeyType != "ECDSA" {
		return fmt.Errorf("only RSA and ECDSA keys are supported for signing")
	}

	// Check if the input file exists
	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %v", err)
	}

	// Step 1: Prepare the OpenSSL command based on key type
	var signCmd *exec.Cmd
	var signatureFormat string
	if token.KeyType == "RSA" {
		signatureFormat = "rsa_padding_mode:pss"
		// Command for signing with RSA-PSS
		signCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-sign",
			"pkcs11:token="+token.Label+";object="+token.ObjectLabel+";type=private;pin-value="+token.UserPin,
			"-sigopt", signatureFormat,
			"-sha384", // Use SHA-384
			"-out", outputFilePath, inputFilePath,
		)
	} else if token.KeyType == "ECDSA" {
		// Command for signing with ECDSA
		signCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-sign",
			"pkcs11:token="+token.Label+";object="+token.ObjectLabel+";type=private;pin-value="+token.UserPin,
			"-sha384", // ECDSA typically uses SHA-384
			"-out", outputFilePath, inputFilePath,
		)
	}

	// Execute the sign command
	signOutput, err := signCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to sign data: %v\nOutput: %s", err, signOutput)
	}

	fmt.Printf("Signing successful. Signature written to %s\n", outputFilePath)
	return nil
}

// Verify verifies the signature of data using the cryptographic capabilities of the PKCS#11 token. Refer to: https://docs.yubico.com/hardware/yubihsm-2/hsm-2-user-guide/hsm2-openssl-libp11.html#rsa-pss
func (token *PKCS11TokenImpl) Verify(dataFilePath, signatureFilePath string) (bool, error) {
	valid := false

	// Validate required parameters
	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
		return valid, fmt.Errorf("missing required arguments for verification")
	}

	if token.KeyType != "RSA" && token.KeyType != "ECDSA" {
		return valid, fmt.Errorf("only RSA and ECDSA keys are supported for verification")
	}

	// Check if the input files exist
	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		return valid, fmt.Errorf("data file does not exist: %v", err)
	}
	if _, err := os.Stat(signatureFilePath); os.IsNotExist(err) {
		return valid, fmt.Errorf("signature file does not exist: %v", err)
	}

	// Step 1: Prepare the OpenSSL command based on key type
	var verifyCmd *exec.Cmd
	if token.KeyType == "RSA" {
		// Command for verifying with RSA-PSS
		verifyCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-verify",
			"pkcs11:token="+token.Label+";object="+token.ObjectLabel+";type=public;pin-value="+token.UserPin,
			"-sigopt", "rsa_padding_mode:pss",
			"-sha384", // Use SHA-384 for verification
			"-signature", signatureFilePath, "-binary", dataFilePath,
		)
	} else if token.KeyType == "ECDSA" {
		// Command for verifying with ECDSA
		verifyCmd = exec.Command(
			"openssl", "dgst", "-engine", "pkcs11", "-keyform", "engine", "-verify",
			"pkcs11:token="+token.Label+";object="+token.ObjectLabel+";type=public;pin-value="+token.UserPin,
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
