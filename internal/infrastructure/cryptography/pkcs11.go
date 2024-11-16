package cryptography

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
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

// Encrypt encrypts data using the cryptographic capabilities of the PKCS#11 token.
func (token *PKCS11Token) Encrypt(inputFilePath, outputFilePath string) error {
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

// Decrypt decrypts data using the cryptographic capabilities of the PKCS#11 token.
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

// Sign signs data using the cryptographic capabilities of the PKCS#11 token.
func (token *PKCS11Token) Sign(inputFilePath, outputFilePath string) error {
	// Validate required parameters
	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
		return fmt.Errorf("missing required arguments for signing")
	}

	if token.KeyType != "RSA" {
		return fmt.Errorf("only RSA keys are supported for decryption")
	}

	// Check if the input file exists
	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %v", err)
	}

	uniqueID := uuid.New()
	// Step 1: Hash the data using OpenSSL (SHA-256)
	hashFile := fmt.Sprintf("%s-data.hash", uniqueID)
	hashCmd := exec.Command("openssl", "dgst", "-sha256", "-binary", inputFilePath)

	// Redirect the output of the hash command to a file
	hashOut, err := os.Create(hashFile)
	if err != nil {
		return fmt.Errorf("failed to create hash output file: %v", err)
	}
	defer hashOut.Close()

	hashCmd.Stdout = hashOut
	hashCmd.Stderr = os.Stderr

	// Execute the hashing command
	if err := hashCmd.Run(); err != nil {
		return fmt.Errorf("failed to hash data: %v", err)
	}
	fmt.Println("Data hashed successfully.")

	// Step 2: Sign the hashed data using pkcs11-tool
	signCmd := exec.Command("pkcs11-tool",
		"--module", token.ModulePath,
		"--token-label", token.Label,
		"--pin", token.UserPin,
		"--sign",
		"--mechanism", "RSA-PKCS-PSS", // Using RSA-PKCS-PSS for signing
		"--hash-algorithm", "SHA256", // Hash algorithm to match the hashing step
		"--input-file", hashFile, // Input file containing the hashed data
		"--output-file", outputFilePath, // Output signature file
		"--signature-format", "openssl", // Use OpenSSL signature format
		"--label", token.ObjectLabel, // Key label used for signing
	)

	// Run the signing command
	signOutput, err := signCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to sign data: %v\nOutput: %s", err, signOutput)
	}

	// Step 3: Remove the hash file from the filesystem
	os.Remove(hashFile)

	fmt.Printf("Signing successful. Signature written to %s\n", outputFilePath)
	return nil
}

// Verify verifies the signature of data using the cryptographic capabilities of the PKCS#11 token.
func (token *PKCS11Token) Verify(dataFilePath, signatureFilePath string) (bool, error) {
	valid := false
	// Validate required parameters
	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
		return valid, fmt.Errorf("missing required arguments for verification")
	}

	if token.KeyType != "RSA" {
		return valid, fmt.Errorf("only RSA keys are supported for decryption")
	}

	// Check if the input files exist
	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
		return valid, fmt.Errorf("data file does not exist: %v", err)
	}
	if _, err := os.Stat(signatureFilePath); os.IsNotExist(err) {
		return valid, fmt.Errorf("signature file does not exist: %v", err)
	}

	uniqueID := uuid.New()
	// Step 1: Retrieve the public key from the PKCS#11 token and save it as public.der
	publicKeyFile := fmt.Sprintf("%s-public.der", uniqueID)

	args := []string{
		"--module", token.ModulePath,
		"--token-label", token.Label,
		"--pin", token.UserPin,
		"--read-object",
		"--label", token.ObjectLabel,
		"--type", "pubkey", // Extract public key
		"--output-file", publicKeyFile, // Output file for public key in DER format
	}

	cmd := exec.Command("pkcs11-tool", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return valid, fmt.Errorf("failed to retrieve public key: %v\nOutput: %s", err, output)
	}
	fmt.Println("Public key retrieved successfully.")

	// Step 2: Verify the signature using OpenSSL and the retrieved public key
	verifyCmd := exec.Command(
		"openssl", "dgst", "-keyform", "DER", "-verify", publicKeyFile, "-sha256", // Use SHA256 for hash
		"-sigopt", "rsa_padding_mode:pss", // Use PSS padding
		"-sigopt", "rsa_pss_saltlen:-1", // Set salt length to default (-1 for auto)
		"-signature", signatureFilePath, // Path to the signature file
		"-binary", dataFilePath, // Path to the data file
	)

	// Run the verification command
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

	// Step 3: Remove the public key from the filesystem
	os.Remove(publicKeyFile)

	return valid, nil
}
