package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

func runCommand(t *testing.T, cmd string, args []string) (string, error) {
	command := exec.Command(cmd, args...)
	var out bytes.Buffer
	command.Stdout = &out
	command.Stderr = &out

	err := command.Run()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %v", err, out.String())
		return "", err
	}

	return out.String(), nil
}

func TestAESEncryptionAndDecryption(t *testing.T) {
	uuid := "test-uuid-1234"
	inputFile := "data/input.txt"

	encOutputFile := fmt.Sprintf("data/%s-output.enc", uuid)
	cmdEncrypt := "go"
	argsEncrypt := []string{"run", "crypto_vault_cli.go", "encrypt-aes", "--input-file", inputFile, "--output-file", encOutputFile, "--symmetric-key", "your-generated-symmetric-key"}

	_, err := runCommand(t, cmdEncrypt, argsEncrypt)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decOutputFile := fmt.Sprintf("data/%s-decrypted.txt", uuid)
	cmdDecrypt := "go"
	argsDecrypt := []string{"run", "crypto_vault_cli.go", "decrypt-aes", "--input-file", encOutputFile, "--output-file", decOutputFile, "--symmetric-key", "your-generated-symmetric-key"}

	_, err = runCommand(t, cmdDecrypt, argsDecrypt)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}
}

func TestRSAEncryptionAndDecryption(t *testing.T) {
	uuid := "test-uuid-5678"
	inputFile := "data/input.txt"

	encOutputFile := fmt.Sprintf("data/%s-encrypted.txt", uuid)
	cmdEncryptRSA := "go"
	argsEncryptRSA := []string{"run", "crypto_vault_cli.go", "encrypt-rsa", "--input-file", inputFile, "--output-file", encOutputFile, "--public-key", "your-generated-public-key"}

	_, err := runCommand(t, cmdEncryptRSA, argsEncryptRSA)
	if err != nil {
		t.Fatalf("RSA Encryption failed: %v", err)
	}

	decOutputFile := fmt.Sprintf("data/%s-decrypted.txt", uuid)
	cmdDecryptRSA := "go"
	argsDecryptRSA := []string{"run", "crypto_vault_cli.go", "decrypt-rsa", "--input-file", encOutputFile, "--output-file", decOutputFile, "--private-key", "your-generated-private-key"}

	_, err = runCommand(t, cmdDecryptRSA, argsDecryptRSA)
	if err != nil {
		t.Fatalf("RSA Decryption failed: %v", err)
	}
}

func TestRSASignAndVerify(t *testing.T) {
	uuid := "test-uuid-5678"
	inputFile := "data/input.txt"
	signatureFile := fmt.Sprintf("data/%s-signature.bin", uuid)

	// Sign
	cmdSignRSA := "go"
	argsSignRSA := []string{"run", "crypto_vault_cli.go", "sign-rsa", "--input-file", inputFile, "--output-file", signatureFile, "--private-key", "your-generated-private-key"}

	_, err := runCommand(t, cmdSignRSA, argsSignRSA)
	if err != nil {
		t.Fatalf("RSA Signing failed: %v", err)
	}

	// Verify
	cmdVerifyRSA := "go"
	argsVerifyRSA := []string{
		"run", "crypto_vault_cli.go", "verify-rsa", "--input-file", inputFile, "--signature-file", signatureFile, "--public-key", "your-generated-public-key",
	}

	_, err = runCommand(t, cmdVerifyRSA, argsVerifyRSA)
	if err != nil {
		t.Fatalf("RSA Verification failed: %v", err)
	}
}

func TestSigningAndVerificationECDSA(t *testing.T) {
	uuid := "test-uuid-ecc"
	inputFile := "data/input.txt"
	signatureFile := fmt.Sprintf("data/%s-signature.bin", uuid)

	// Sign
	cmdSignECDSA := "go"
	argsSignECDSA := []string{"run", "crypto_vault_cli.go", "sign-ecc", "--input-file", inputFile, "--output-file", signatureFile, "--private-key", "your-generated-private-key"}

	_, err := runCommand(t, cmdSignECDSA, argsSignECDSA)
	if err != nil {
		t.Fatalf("ECDSA Signing failed: %v", err)
	}

	// Verify
	cmdVerifyECDSA := "go"
	argsVerifyECDSA := []string{
		"run", "crypto_vault_cli.go", "verify-ecc", "--input-file", inputFile, "--signature-file", signatureFile, "--public-key", "your-generated-public-key",
	}

	_, err = runCommand(t, cmdVerifyECDSA, argsVerifyECDSA)
	if err != nil {
		t.Fatalf("ECDSA Verification failed: %v", err)
	}
}

func TestPKCS11EncryptionAndDecryption(t *testing.T) {
	uuid := "test-uuid-pkcs11"
	inputFile := "data/input.txt"

	encOutputFile := fmt.Sprintf("data/%s-encrypted-output.enc", uuid)
	cmdEncryptPKCS11 := "go"
	argsEncryptPKCS11 := []string{
		"run", "crypto_vault_cli.go", "encrypt", "--token-label", "my-token", "--object-label", "my-rsa-key", "--key-type", "RSA", "--input-file", inputFile, "--output-file", encOutputFile,
	}

	_, err := runCommand(t, cmdEncryptPKCS11, argsEncryptPKCS11)
	if err != nil {
		t.Fatalf("PKCS11 Encryption failed: %v", err)
	}

	decOutputFile := fmt.Sprintf("data/%s-decrypted-output.txt", uuid)
	cmdDecryptPKCS11 := "go"
	argsDecryptPKCS11 := []string{
		"run", "crypto_vault_cli.go", "decrypt", "--token-label", "my-token", "--object-label", "my-rsa-key", "--key-type", "RSA", "--input-file", encOutputFile, "--output-file", decOutputFile,
	}

	_, err = runCommand(t, cmdDecryptPKCS11, argsDecryptPKCS11)
	if err != nil {
		t.Fatalf("PKCS11 Decryption failed: %v", err)
	}
}

func TestPKCS11KeyManagement(t *testing.T) {

	// Store PKCS#11 settings
	cmdStorePKCS11 := "go"
	argsStorePKCS11 := []string{
		"run", "crypto_vault_cli.go", "store-pkcs11-settings", "--module", "/usr/lib/softhsm/libsofthsm2.so", "--so-pin", "1234", "--user-pin", "5678", "--slot-id", "0x0",
	}

	_, err := runCommand(t, cmdStorePKCS11, argsStorePKCS11)
	if err != nil {
		t.Fatalf("Storing PKCS#11 settings failed: %v", err)
	}

	// Add RSA Key
	cmdAddRSAKey := "go"
	argsAddRSAKey := []string{"run", "crypto_vault_cli.go", "add-key", "--token-label", "my-token", "--object-label", "my-rsa-key", "--key-type", "RSA", "--key-size", "2048"}

	_, err = runCommand(t, cmdAddRSAKey, argsAddRSAKey)
	if err != nil {
		t.Fatalf("Adding RSA Key to PKCS#11 failed: %v", err)
	}

	// List Objects
	cmdListObjects := "go"
	argsListObjects := []string{"run", "crypto_vault_cli.go", "list-objects", "--token-label", "my-token"}

	_, err = runCommand(t, cmdListObjects, argsListObjects)
	if err != nil {
		t.Fatalf("Listing PKCS#11 objects failed: %v", err)
	}

	// Delete RSA Key
	cmdDeleteRSAKey := "go"
	argsDeleteRSAKey := []string{"run", "crypto_vault_cli.go", "delete-object", "--token-label", "my-token", "--object-label", "my-rsa-key", "--object-type", "pubkey"}

	_, err = runCommand(t, cmdDeleteRSAKey, argsDeleteRSAKey)
	if err != nil {
		t.Fatalf("Deleting RSA Key from PKCS#11 failed: %v", err)
	}
}
