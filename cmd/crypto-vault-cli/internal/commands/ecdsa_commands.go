package commands

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto_vault_service/cmd/crypto-vault-cli/internal/status"
	"crypto_vault_service/internal/infrastructure/cryptography"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// signECCCmd signs the contents of a file with ECDSA
func SignECCCmd(cmd *cobra.Command, args []string) {

	inputFile, err := cmd.Flags().GetString("input-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	keyDir, err := cmd.Flags().GetString("key-dir")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	EC := &cryptography.EC{}
	var privateKey *ecdsa.PrivateKey
	var publicKey *ecdsa.PublicKey

	privateKey, publicKey, err = EC.GenerateKeys(elliptic.P256())
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	fileContent, err := os.ReadFile(inputFile)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	signature, err := EC.Sign(fileContent, privateKey)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	info := status.NewInfo(fmt.Sprintf("Signature: %x", signature))
	info.PrintJsonInfo(false)

	uniqueID := uuid.New()

	if privateKey != nil && keyDir != "" {
		privateKeyFilePath := fmt.Sprintf("%s/%s-private-key.pem", keyDir, uniqueID.String())

		err = EC.SavePrivateKeyToFile(privateKey, privateKeyFilePath)
		if err != nil {
			e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
			e.PrintJsonError()
			return
		}
		info := status.NewInfo(fmt.Sprintf("Private key saved to: %s", privateKeyFilePath))
		info.PrintJsonInfo(false)
	}

	if publicKey != nil && keyDir != "" {
		publicKeyFilePath := fmt.Sprintf("%s/%s-public-key.pem", keyDir, uniqueID.String())
		err = EC.SavePublicKeyToFile(publicKey, publicKeyFilePath)
		if err != nil {
			e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
			e.PrintJsonError()
			return
		}
		info := status.NewInfo(fmt.Sprintf("Public key saved to: %s", publicKeyFilePath))
		info.PrintJsonInfo(false)
	}

	if keyDir != "" {
		signatureFilePath := fmt.Sprintf("%s/%s-signature.sig", keyDir, uniqueID.String())
		err = EC.SaveSignatureToFile(signatureFilePath, signature)
		if err != nil {
			if err != nil {
				e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
				e.PrintJsonError()
				return
			}
		}
		info := status.NewInfo(fmt.Sprintf("Signature file saved to: %s", signatureFilePath))
		info.PrintJsonInfo(false)
	}
}

// verifyECCCmd verifies the signature of a file's content using ECDSA
func VerifyECCCmd(cmd *cobra.Command, args []string) {
	inputFile, err := cmd.Flags().GetString("input-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	publicKeyPath, err := cmd.Flags().GetString("public-key")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	signatureFile, err := cmd.Flags().GetString("signature-file")
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	EC := &cryptography.EC{}
	var publicKey *ecdsa.PublicKey

	if publicKeyPath == "" {
		log.Fatalf("Public key is required for ECC signature verification.")
		e := status.NewError("Public key is required for ECC signature verification", status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	} else {
		publicKey, err = EC.ReadPublicKey(publicKeyPath, elliptic.P256())
		if err != nil {
			e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
			e.PrintJsonError()
			return
		}
	}

	fileContent, err := os.ReadFile(inputFile)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	signatureHex, err := os.ReadFile(signatureFile)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	signature, err := hex.DecodeString(string(signatureHex))
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	valid, err := EC.Verify(fileContent, signature, publicKey)
	if err != nil {
		e := status.NewError(fmt.Sprintf("%v", err), status.ErrCodeInternalError)
		e.PrintJsonError()
		return
	}

	if valid {
		info := status.NewInfo("Signature is valid")
		info.PrintJsonInfo(false)
	} else {
		info := status.NewInfo("Signature is invalid")
		info.PrintJsonInfo(false)
	}
}

func InitECDSACommands(rootCmd *cobra.Command) {
	var signECCMessageCmd = &cobra.Command{
		Use:   "sign-ecc",
		Short: "Sign a message using ECC",
		Run:   SignECCCmd,
	}
	signECCMessageCmd.Flags().StringP("input-file", "", "", "Path to file that needs to be signed")
	signECCMessageCmd.Flags().StringP("key-dir", "", "", "Directory to save generated keys (optional)")
	rootCmd.AddCommand(signECCMessageCmd)

	var verifyECCSignatureCmd = &cobra.Command{
		Use:   "verify-ecc",
		Short: "Verify a signature using ECC",
		Run:   VerifyECCCmd,
	}
	verifyECCSignatureCmd.Flags().StringP("input-file", "", "", "Path to ECC public key")
	verifyECCSignatureCmd.Flags().StringP("public-key", "", "", "The public key used to verify the signature")
	verifyECCSignatureCmd.Flags().StringP("signature-file", "", "", "Signature to verify (hex format)")
	rootCmd.AddCommand(verifyECCSignatureCmd)
}
