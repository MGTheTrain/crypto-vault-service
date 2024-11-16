package cryptography

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PKCS11TokenInterface defines the operations for working with a PKCS#11 token
type PKCS11TokenInterface interface {
	IsTokenSet() (bool, error)
	IsObjectSet() (bool, error)
	InitializeToken() error
	AddKey() error
	Encrypt() error
	Decrypt() error
	Sign() error
	Verify() error
	DeleteObject(objectType, objectLabel string) error
}

// PKCS11Token represents the parameters and operations for interacting with a PKCS#11 token
type PKCS11Token struct {
	ModulePath  string
	Label       string
	SOPin       string
	UserPin     string
	ObjectLabel string
	KeyType     string // "ECDSA" or "RSA"
	KeySize     int    // Key size in bits for RSA or ECDSA (e.g., 256 for ECDSA, 2048 for RSA)
}

// Public method to execute pkcs11-tool commands and return output
func (token *PKCS11Token) executePKCS11ToolCommand(args []string) (string, error) {
	cmd := exec.Command("pkcs11-tool", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pkcs11-tool command failed: %v\nOutput: %s", err, output)
	}
	return string(output), nil
}

// IsTokenSet checks if the token exists in the given module path
func (token *PKCS11Token) IsTokenSet() (bool, error) {
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
func (token *PKCS11Token) IsObjectSet() (bool, error) {
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
func (token *PKCS11Token) InitializeToken(slot string) error {
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
func (token *PKCS11Token) DeleteObject(objectType, objectLabel string) error {
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
func (token *PKCS11Token) AddKey() error {
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
func (token *PKCS11Token) addECDSASignKey() error {
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
func (token *PKCS11Token) addRSASignKey() error {
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

func (token *PKCS11Token) Encrypt(inputFilePath, outputFilePath string) error {
	// Validate required parameters
	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
		return fmt.Errorf("missing required arguments for encryption")
	}

	if token.KeyType != "RSA" {
		return fmt.Errorf("only RSA keys are supported for encryption")
	}

	// Temporary file to store the public key in DER format
	publicKeyFile := "public.der"

	// Step 1: Retrieve the public key from the PKCS#11 token using pkcs11-tool
	args := []string{
		"--module", token.ModulePath,
		"--token-label", token.Label,
		"--pin", token.UserPin,
		"--read-object",
		"--label", token.ObjectLabel,
		"--type", "pubkey", // Retrieve public key
		"--output-file", publicKeyFile, // Store public key in DER format
	}

	cmd := exec.Command("pkcs11-tool", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to retrieve public key: %v\nOutput: %s", err, output)
	}
	fmt.Println("Public key retrieved successfully.")

	// Check if the public key file was generated
	if _, err := os.Stat(publicKeyFile); os.IsNotExist(err) {
		return fmt.Errorf("public key file not found: %s", publicKeyFile)
	}

	// Step 2: Encrypt the data using OpenSSL and the retrieved public key
	encryptCmd := exec.Command("openssl", "pkeyutl", "-encrypt", "-pubin", "-inkey", publicKeyFile, "-keyform", "DER", "-in", inputFilePath, "-out", outputFilePath)
	encryptOutput, err := encryptCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to encrypt data with OpenSSL: %v\nOutput: %s", err, encryptOutput)
	}

	// Step 3: Remove the public key from the filesystem
	os.Remove("public.der")

	fmt.Printf("Encryption successful. Encrypted data written to %s\n", outputFilePath)
	return nil
}

func (token *PKCS11Token) Decrypt(inputFilePath, outputFilePath string) error {
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

	// Create or validate the output file (will be overwritten)
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create or open output file: %v", err)
	}
	defer outputFile.Close()

	// Step 1: Prepare the command to decrypt the data using pkcs11-tool
	args := []string{
		"--module", token.ModulePath,
		"--token-label", token.Label,
		"--pin", token.UserPin,
		"--decrypt",
		"--label", token.ObjectLabel,
		"--mechanism", "RSA-PKCS", // Specify the RSA-PKCS mechanism
		"--input-file", inputFilePath, // Input file with encrypted data
	}

	// Run the decryption command
	cmd := exec.Command("pkcs11-tool", args...)
	output, err := cmd.CombinedOutput()

	// Capture the decrypted data and filter out any extra output
	if err != nil {
		return fmt.Errorf("decryption failed: %v\nOutput: %s", err, output)
	}

	// Split the output into lines and filter out unwanted lines
	lines := strings.Split(string(output), "\n")
	var decryptedData []byte
	for i, line := range lines {
		if !strings.Contains(line, "Using decrypt algorithm RSA-PKCS") {
			// If this is not the last line, append a newline
			decryptedData = append(decryptedData, []byte(line)...)
			if i < len(lines)-1 {
				decryptedData = append(decryptedData, '\n')
			}
		}
	}

	// Write the actual decrypted data (without extra info) to the output file
	_, err = outputFile.Write(decryptedData)
	if err != nil {
		return fmt.Errorf("failed to write decrypted data to output file: %v", err)
	}

	fmt.Printf("Decryption successful. Decrypted data written to %s\n", outputFilePath)
	return nil
}
