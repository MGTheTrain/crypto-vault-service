package cryptography

// For reference only, to demonstrate experiments with pkcs11-tool for encryption, decryption, signing, and verification.

// // Encrypt encrypts data using the cryptographic capabilities of the PKCS#11 token. Currently only supports RSA keys. Refer to: https://docs.nitrokey.com/nethsm/pkcs11-tool#pkcs11-tool
// func (token *PKCS11TokenImpl) Encrypt(inputFilePath, outputFilePath string) error {
// 	// Validate required parameters
// 	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
// 		return fmt.Errorf("missing required arguments for encryption")
// 	}

// 	if token.KeyType != "RSA" {
// 		return fmt.Errorf("only RSA keys are supported for encryption")
// 	}

// 	uniqueID := uuid.New()
// 	// Temporary file to store the public key in DER format
// 	publicKeyFile := fmt.Sprintf("%s-public.der", uniqueID)

// 	// Step 1: Retrieve the public key from the PKCS#11 token using pkcs11-tool
// 	args := []string{
// 		"--module", token.ModulePath,
// 		"--token-label", token.Label,
// 		"--pin", token.UserPin,
// 		"--read-object",
// 		"--label", token.ObjectLabel,
// 		"--type", "pubkey", // Retrieve public key
// 		"--output-file", publicKeyFile, // Store public key in DER format
// 	}

// 	cmd := exec.Command("pkcs11-tool", args...)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to retrieve public key: %v\nOutput: %s", err, output)
// 	}
// 	fmt.Println("Public key retrieved successfully.")

// 	// Check if the public key file was generated
// 	if _, err := os.Stat(publicKeyFile); os.IsNotExist(err) {
// 		return fmt.Errorf("public key file not found: %s", publicKeyFile)
// 	}

// 	// Step 2: Encrypt the data using OpenSSL and the retrieved public key
// 	encryptCmd := exec.Command("openssl", "pkeyutl", "-encrypt", "-pubin", "-inkey", publicKeyFile, "-keyform", "DER", "-in", inputFilePath, "-out", outputFilePath)
// 	encryptOutput, err := encryptCmd.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to encrypt data with OpenSSL: %v\nOutput: %s", err, encryptOutput)
// 	}

// 	// Step 3: Remove the public key from the filesystem
// 	os.Remove(publicKeyFile)

// 	fmt.Printf("Encryption successful. Encrypted data written to %s\n", outputFilePath)
// 	return nil
// }

// // Decrypt decrypts data using the cryptographic capabilities of the PKCS#11 token. Currently only supports RSA keys. Refer to: https://docs.nitrokey.com/nethsm/pkcs11-tool#pkcs11-tool
// func (token *PKCS11TokenImpl) Decrypt(inputFilePath, outputFilePath string) error {
// 	// Validate required parameters
// 	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
// 		return fmt.Errorf("missing required arguments for decryption")
// 	}

// 	if token.KeyType != "RSA" {
// 		return fmt.Errorf("only RSA keys are supported for decryption")
// 	}

// 	// Check if input file exists
// 	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
// 		return fmt.Errorf("input file does not exist: %v", err)
// 	}

// 	// Create or validate the output file (will be overwritten)
// 	outputFile, err := os.Create(outputFilePath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create or open output file: %v", err)
// 	}
// 	defer outputFile.Close()

// 	// Step 1: Prepare the command to decrypt the data using pkcs11-tool
// 	args := []string{
// 		"--module", token.ModulePath,
// 		"--token-label", token.Label,
// 		"--pin", token.UserPin,
// 		"--decrypt",
// 		"--label", token.ObjectLabel,
// 		"--mechanism", "RSA-PKCS", // Specify the RSA-PKCS mechanism
// 		"--input-file", inputFilePath, // Input file with encrypted data
// 	}

// 	// Run the decryption command
// 	cmd := exec.Command("pkcs11-tool", args...)
// 	output, err := cmd.CombinedOutput()

// 	// Capture the decrypted data and filter out any extra output
// 	if err != nil {
// 		return fmt.Errorf("decryption failed: %v\nOutput: %s", err, output)
// 	}

// 	// Split the output into lines and filter out unwanted lines
// 	lines := strings.Split(string(output), "\n")
// 	var decryptedData []byte
// 	for i, line := range lines {
// 		if !strings.Contains(line, "Using decrypt algorithm RSA-PKCS") {
// 			// If this is not the last line, append a newline
// 			decryptedData = append(decryptedData, []byte(line)...)
// 			if i < len(lines)-1 {
// 				decryptedData = append(decryptedData, '\n')
// 			}
// 		}
// 	}

// 	// Write the actual decrypted data (without extra info) to the output file
// 	_, err = outputFile.Write(decryptedData)
// 	if err != nil {
// 		return fmt.Errorf("failed to write decrypted data to output file: %v", err)
// 	}

// 	fmt.Printf("Decryption successful. Decrypted data written to %s\n", outputFilePath)
// 	return nil
// }

// // Sign signs data using the cryptographic capabilities of the PKCS#11 token. Currently only supports RSA keys. Refer to: https://docs.nitrokey.com/nethsm/pkcs11-tool#pkcs11-tool
// func (token *PKCS11TokenImpl) Sign(inputFilePath, outputFilePath string) error {
// 	// Validate required parameters
// 	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
// 		return fmt.Errorf("missing required arguments for signing")
// 	}

// 	if token.KeyType != "RSA" {
// 		return fmt.Errorf("only RSA keys are supported for decryption")
// 	}

// 	// Check if the input file exists
// 	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
// 		return fmt.Errorf("input file does not exist: %v", err)
// 	}

// 	uniqueID := uuid.New()
// 	// Step 1: Hash the data using OpenSSL (SHA-256)
// 	hashFile := fmt.Sprintf("%s-data.hash", uniqueID)
// 	hashCmd := exec.Command("openssl", "dgst", "-sha256", "-binary", inputFilePath)

// 	// Redirect the output of the hash command to a file
// 	hashOut, err := os.Create(hashFile)
// 	if err != nil {
// 		return fmt.Errorf("failed to create hash output file: %v", err)
// 	}
// 	defer hashOut.Close()

// 	hashCmd.Stdout = hashOut
// 	hashCmd.Stderr = os.Stderr

// 	// Execute the hashing command
// 	if err := hashCmd.Run(); err != nil {
// 		return fmt.Errorf("failed to hash data: %v", err)
// 	}
// 	fmt.Println("Data hashed successfully.")

// 	// Step 2: Sign the hashed data using pkcs11-tool
// 	signCmd := exec.Command("pkcs11-tool",
// 		"--module", token.ModulePath,
// 		"--token-label", token.Label,
// 		"--pin", token.UserPin,
// 		"--sign",
// 		"--mechanism", "RSA-PKCS-PSS", // Using RSA-PKCS-PSS for signing
// 		"--hash-algorithm", "SHA256", // Hash algorithm to match the hashing step
// 		"--input-file", hashFile, // Input file containing the hashed data
// 		"--output-file", outputFilePath, // Output signature file
// 		"--signature-format", "openssl", // Use OpenSSL signature format
// 		"--label", token.ObjectLabel, // Key label used for signing
// 	)

// 	// Run the signing command
// 	signOutput, err := signCmd.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to sign data: %v\nOutput: %s", err, signOutput)
// 	}

// 	// Step 3: Remove the hash file from the filesystem
// 	os.Remove(hashFile)

// 	fmt.Printf("Signing successful. Signature written to %s\n", outputFilePath)
// 	return nil
// }

// // Verify verifies the signature of data using the cryptographic capabilities of the PKCS#11 token.
// func (token *PKCS11TokenImpl) Verify(dataFilePath, signatureFilePath string) (bool, error) {
// 	valid := false
// 	// Validate required parameters
// 	if token.ModulePath == "" || token.Label == "" || token.ObjectLabel == "" || token.UserPin == "" {
// 		return valid, fmt.Errorf("missing required arguments for verification")
// 	}

// 	if token.KeyType != "RSA" {
// 		return valid, fmt.Errorf("only RSA keys are supported for decryption")
// 	}

// 	// Check if the input files exist
// 	if _, err := os.Stat(dataFilePath); os.IsNotExist(err) {
// 		return valid, fmt.Errorf("data file does not exist: %v", err)
// 	}
// 	if _, err := os.Stat(signatureFilePath); os.IsNotExist(err) {
// 		return valid, fmt.Errorf("signature file does not exist: %v", err)
// 	}

// 	uniqueID := uuid.New()
// 	// Step 1: Retrieve the public key from the PKCS#11 token and save it as public.der
// 	publicKeyFile := fmt.Sprintf("%s-public.der", uniqueID)

// 	args := []string{
// 		"--module", token.ModulePath,
// 		"--token-label", token.Label,
// 		"--pin", token.UserPin,
// 		"--read-object",
// 		"--label", token.ObjectLabel,
// 		"--type", "pubkey", // Extract public key
// 		"--output-file", publicKeyFile, // Output file for public key in DER format
// 	}

// 	cmd := exec.Command("pkcs11-tool", args...)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		return valid, fmt.Errorf("failed to retrieve public key: %v\nOutput: %s", err, output)
// 	}
// 	fmt.Println("Public key retrieved successfully.")

// 	// Step 2: Verify the signature using OpenSSL and the retrieved public key
// 	verifyCmd := exec.Command(
// 		"openssl", "dgst", "-keyform", "DER", "-verify", publicKeyFile, "-sha256", // Use SHA256 for hash
// 		"-sigopt", "rsa_padding_mode:pss", // Use PSS padding
// 		"-sigopt", "rsa_pss_saltlen:-1", // Set salt length to default (-1 for auto)
// 		"-signature", signatureFilePath, // Path to the signature file
// 		"-binary", dataFilePath, // Path to the data file
// 	)

// 	// Run the verification command
// 	verifyOutput, err := verifyCmd.CombinedOutput()
// 	if err != nil {
// 		return valid, fmt.Errorf("failed to verify signature: %v\nOutput: %s", err, verifyOutput)
// 	}

// 	// Check the output from OpenSSL to determine if the verification was successful
// 	if strings.Contains(string(verifyOutput), "Verified OK") {
// 		fmt.Println("Verification successful: The signature is valid.")
// 		valid = true
// 	} else {
// 		fmt.Println("Verification failed: The signature is invalid.")
// 	}

// 	// Step 3: Remove the public key from the filesystem
// 	os.Remove(publicKeyFile)

// 	return valid, nil
// }
