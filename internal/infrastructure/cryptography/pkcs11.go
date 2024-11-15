package cryptography

import (
	"fmt"
	"os/exec"
	"strings"
)

// PKCS11TokenInterface defines the operations for working with a PKCS#11 token
type PKCS11TokenInterface interface {
	Pkcs11SlotSetup() error
	IsTokenSet() (bool, error)
	IsObjectSet() (bool, error)
	InitializeToken() error
	GetFreeSlot() (string, error)
	AddKey() error
	AddECDSASignKey() error
	AddRSASignKey() error
	DeleteObject(objectType, objectLabel string) error // Added method for deleting keys
}

// PKCS11Token represents the parameters and operations for interacting with a PKCS#11 token
type PKCS11Token struct {
	ModulePath  string
	TokenLabel  string
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

// Pkcs11SlotSetup sets up the PKCS#11 token, initializes it, and adds keys (ECDSA or RSA)
func (token *PKCS11Token) Pkcs11SlotSetup() error {
	// Check if OpenSC is installed
	if _, err := exec.LookPath("pkcs11-tool"); err != nil {
		return fmt.Errorf("OpenSC is not installed: %v", err)
	}

	// Validate the required environment variables and pins
	if token.ModulePath == "" {
		return fmt.Errorf("PKCS11_MODULE_PATH must be set")
	}

	if len(token.UserPin) < 4 || len(token.SOPin) < 4 {
		return fmt.Errorf("PINs must be at least 4 characters")
	}

	// Initialize token if necessary
	if err := token.InitializeToken(); err != nil {
		return err
	}

	// Add the key to the token (either ECDSA or RSA)
	if err := token.AddKey(); err != nil {
		return err
	}

	// List all token slots
	fmt.Println("### List all slots")
	if output, err := token.executePKCS11ToolCommand([]string{"-L", "--module", token.ModulePath}); err != nil {
		return err
	} else {
		fmt.Println(output)
	}

	// List all objects on the selected token
	fmt.Printf("### List all objects on Token '%s'\n", token.TokenLabel)
	if output, err := token.executePKCS11ToolCommand([]string{"-O", "--module", token.ModulePath, "--token-label", token.TokenLabel, "--pin", token.UserPin}); err != nil {
		return err
	} else {
		fmt.Println(output)
	}

	return nil
}

// IsTokenSet checks if the token exists in the given module path
func (token *PKCS11Token) IsTokenSet() (bool, error) {
	if token.ModulePath == "" || token.TokenLabel == "" {
		return false, fmt.Errorf("missing module path or token label")
	}

	args := []string{"--module", token.ModulePath, "-T"}
	output, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return false, err
	}

	if strings.Contains(output, token.TokenLabel) && strings.Contains(output, "token initialized") {
		fmt.Printf("Token with label '%s' exists.\n", token.TokenLabel)
		return true, nil
	}

	fmt.Printf("Error: Token with label '%s' does not exist.\n", token.TokenLabel)
	return false, nil
}

// IsObjectSet checks if the specified object exists on the given token
func (token *PKCS11Token) IsObjectSet() (bool, error) {
	if token.ModulePath == "" || token.TokenLabel == "" || token.ObjectLabel == "" || token.UserPin == "" {
		return false, fmt.Errorf("missing required arguments")
	}

	args := []string{"-O", "--module", token.ModulePath, "--token-label", token.TokenLabel, "--pin", token.UserPin}
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
func (token *PKCS11Token) InitializeToken() error {
	if token.ModulePath == "" || token.TokenLabel == "" || token.SOPin == "" || token.UserPin == "" {
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

	// Get the next available slot
	nextSlot, err := token.GetFreeSlot()
	if err != nil {
		return err
	}

	// Initialize the token
	args := []string{"--module", token.ModulePath, "--init-token", "--label", token.TokenLabel, "--so-pin", token.SOPin, "--init-pin", "--pin", token.UserPin, "--slot", nextSlot}
	_, err = token.executePKCS11ToolCommand(args)
	if err != nil {
		return fmt.Errorf("failed to initialize token with label '%s': %v", token.TokenLabel, err)
	}

	fmt.Printf("Token with label '%s' initialized successfully.\n", token.TokenLabel)
	return nil
}

// DeleteObject deletes a key or object from the token
func (token *PKCS11Token) DeleteObject(objectType, objectLabel string) error {
	if token.ModulePath == "" || token.TokenLabel == "" || objectLabel == "" || token.UserPin == "" {
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
		"--token-label", token.TokenLabel,
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

// GetFreeSlot finds the next available uninitialized slot
func (token *PKCS11Token) GetFreeSlot() (string, error) {
	if token.ModulePath == "" {
		return "", fmt.Errorf("missing module path")
	}

	args := []string{"-L", "--module", token.ModulePath}
	output, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return "", err
	}

	// Extract the first available uninitialized slot
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "uninitialized") {
			// Extract slot number from the line
			parts := strings.Fields(line)
			if len(parts) > 2 {
				return parts[2], nil
			}
		}
	}

	return "", fmt.Errorf("no uninitialized slot found")
}

// AddKey adds the selected key (ECDSA or RSA) to the token
func (token *PKCS11Token) AddKey() error {
	if token.ModulePath == "" || token.TokenLabel == "" || token.ObjectLabel == "" || token.UserPin == "" {
		return fmt.Errorf("missing required arguments")
	}

	// Determine key type and call the appropriate function to generate the key
	if token.KeyType == "ECDSA" {
		return token.AddECDSASignKey()
	} else if token.KeyType == "RSA" {
		return token.AddRSASignKey()
	} else {
		return fmt.Errorf("unsupported key type: %s", token.KeyType)
	}
}

// AddECDSASignKey adds an ECDSA signing key to the token
func (token *PKCS11Token) AddECDSASignKey() error {
	if token.KeySize != 256 && token.KeySize != 384 && token.KeySize != 521 {
		return fmt.Errorf("ECDSA key size must be one of 256, 384, or 521 bits, but got %d", token.KeySize)
	}

	// Generate the key pair (example using secp256r1)
	args := []string{
		"--module", token.ModulePath,
		"--token-label", token.TokenLabel,
		"--keypairgen",
		"--key-type", fmt.Sprintf("EC:secp256r1"), // Choose secp256r1 for simplicity
		"--label", token.ObjectLabel,
		"--pin", token.UserPin,
		"--usage-sign",
	}
	_, err := token.executePKCS11ToolCommand(args)
	if err != nil {
		return fmt.Errorf("failed to add ECDSA key to token: %v", err)
	}

	fmt.Printf("ECDSA key with label '%s' added to token '%s'.\n", token.ObjectLabel, token.TokenLabel)
	return nil
}

// AddRSASignKey adds an RSA signing key to the token
func (token *PKCS11Token) AddRSASignKey() error {
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
		"--token-label", token.TokenLabel,
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

	fmt.Printf("RSA key with label '%s' added to token '%s'.\n", token.ObjectLabel, token.TokenLabel)
	return nil
}
