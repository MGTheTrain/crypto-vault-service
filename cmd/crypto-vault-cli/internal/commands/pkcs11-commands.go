package commands

import (
	"crypto_vault_service/internal/infrastructure/cryptography"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// Command to check if token is set
func IsTokenSetCmd(cmd *cobra.Command, args []string) {
	modulePath, _ := cmd.Flags().GetString("module")
	Label, _ := cmd.Flags().GetString("token-label")

	token := &cryptography.PKCS11Token{
		ModulePath: modulePath,
		Label:      Label,
	}

	isSet, err := token.IsTokenSet()
	if err != nil {
		log.Fatalf("Error checking token: %v", err)
	}
	if isSet {
		fmt.Println("Token is set.")
	} else {
		fmt.Println("Token is not set.")
	}
}

// Command to check if object is set
func IsObjectSetCmd(cmd *cobra.Command, args []string) {
	modulePath, _ := cmd.Flags().GetString("module")
	Label, _ := cmd.Flags().GetString("token-label")
	objectLabel, _ := cmd.Flags().GetString("object-label")
	userPin, _ := cmd.Flags().GetString("user-pin")

	token := &cryptography.PKCS11Token{
		ModulePath:  modulePath,
		Label:       Label,
		ObjectLabel: objectLabel,
		UserPin:     userPin,
	}

	isSet, err := token.IsObjectSet()
	if err != nil {
		log.Fatalf("Error checking object: %v", err)
	}
	if isSet {
		fmt.Println("Object is set.")
	} else {
		fmt.Println("Object is not set.")
	}
}

// Command to initialize a PKCS#11 token
func InitializeTokenCmd(cmd *cobra.Command, args []string) {
	slot, _ := cmd.Flags().GetString("slot")
	modulePath, _ := cmd.Flags().GetString("module")
	Label, _ := cmd.Flags().GetString("token-label")
	soPin, _ := cmd.Flags().GetString("so-pin")
	userPin, _ := cmd.Flags().GetString("user-pin")

	token := &cryptography.PKCS11Token{
		ModulePath: modulePath,
		Label:      Label,
		SOPin:      soPin,
		UserPin:    userPin,
	}

	if err := token.InitializeToken(slot); err != nil {
		log.Fatalf("Error initializing token: %v", err)
	}
}

// Command to add key to PKCS#11 token
func AddKeyCmd(cmd *cobra.Command, args []string) {
	modulePath, _ := cmd.Flags().GetString("module")
	Label, _ := cmd.Flags().GetString("token-label")
	objectLabel, _ := cmd.Flags().GetString("object-label")
	keyType, _ := cmd.Flags().GetString("key-type") // "ECDSA" or "RSA"
	keySize, _ := cmd.Flags().GetInt("key-size")
	userPin, _ := cmd.Flags().GetString("user-pin")

	token := &cryptography.PKCS11Token{
		ModulePath:  modulePath,
		Label:       Label,
		ObjectLabel: objectLabel,
		KeyType:     keyType,
		KeySize:     keySize,
		UserPin:     userPin,
	}

	if err := token.AddKey(); err != nil {
		log.Fatalf("Error adding key: %v", err)
	}
	fmt.Println("Key added successfully!")
}

// Command to delete an object (key) from the PKCS#11 token
func DeleteObjectCmd(cmd *cobra.Command, args []string) {
	modulePath, _ := cmd.Flags().GetString("module")
	Label, _ := cmd.Flags().GetString("token-label")
	objectLabel, _ := cmd.Flags().GetString("object-label")
	objectType, _ := cmd.Flags().GetString("object-type")
	userPin, _ := cmd.Flags().GetString("user-pin")

	token := &cryptography.PKCS11Token{
		ModulePath:  modulePath,
		Label:       Label,
		ObjectLabel: objectLabel,
		UserPin:     userPin,
	}

	if err := token.DeleteObject(objectType, objectLabel); err != nil {
		log.Fatalf("Error deleting object: %v", err)
	}
	fmt.Println("Object deleted successfully!")
}

func InitPKCS11Commands(rootCmd *cobra.Command) {
	var pkcs11IsTokenSetCmd = &cobra.Command{
		Use:   "is-token-set",
		Short: "Check if PKCS#11 token is set",
		Run:   IsTokenSetCmd,
	}
	pkcs11IsTokenSetCmd.Flags().String("module", "", "Path to the PKCS#11 module")
	pkcs11IsTokenSetCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	rootCmd.AddCommand(pkcs11IsTokenSetCmd)

	var pkcs11IsObjectSetCmd = &cobra.Command{
		Use:   "is-object-set",
		Short: "Check if object exists in the PKCS#11 token",
		Run:   IsObjectSetCmd,
	}
	pkcs11IsObjectSetCmd.Flags().String("module", "", "Path to the PKCS#11 module")
	pkcs11IsObjectSetCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11IsObjectSetCmd.Flags().String("object-label", "", "Label of the object")
	pkcs11IsObjectSetCmd.Flags().String("user-pin", "", "User PIN")
	rootCmd.AddCommand(pkcs11IsObjectSetCmd)

	var pkcs11InitializeTokenCmd = &cobra.Command{
		Use:   "initialize-token",
		Short: "Initialize a PKCS#11 token",
		Run:   InitializeTokenCmd,
	}
	pkcs11InitializeTokenCmd.Flags().String("slot", "", "The token slot id")
	pkcs11InitializeTokenCmd.Flags().String("module", "", "Path to the PKCS#11 module")
	pkcs11InitializeTokenCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11InitializeTokenCmd.Flags().String("so-pin", "", "Security Officer PIN")
	pkcs11InitializeTokenCmd.Flags().String("user-pin", "", "User PIN")
	rootCmd.AddCommand(pkcs11InitializeTokenCmd)

	var pkcs11AddKeyCmd = &cobra.Command{
		Use:   "add-key",
		Short: "Add key (ECDSA or RSA) to the PKCS#11 token",
		Run:   AddKeyCmd,
	}
	pkcs11AddKeyCmd.Flags().String("module", "", "Path to the PKCS#11 module")
	pkcs11AddKeyCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11AddKeyCmd.Flags().String("object-label", "", "Label of the object (key)")
	pkcs11AddKeyCmd.Flags().String("key-type", "", "Type of the key (ECDSA or RSA)")
	pkcs11AddKeyCmd.Flags().Int("key-size", 0, "Key size in bits (2048 for RSA, 256 for ECDSA)")
	pkcs11AddKeyCmd.Flags().String("user-pin", "", "User PIN")
	rootCmd.AddCommand(pkcs11AddKeyCmd)

	var pkcs11DeleteObjectCmd = &cobra.Command{
		Use:   "delete-object",
		Short: "Delete an object (key) from the PKCS#11 token",
		Run:   DeleteObjectCmd,
	}
	pkcs11DeleteObjectCmd.Flags().String("module", "", "Path to the PKCS#11 module")
	pkcs11DeleteObjectCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11DeleteObjectCmd.Flags().String("object-label", "", "Label of the object to delete")
	pkcs11DeleteObjectCmd.Flags().String("object-type", "", "Type of the object (e.g., privkey, pubkey, cert)")
	pkcs11DeleteObjectCmd.Flags().String("user-pin", "", "User PIN")
	rootCmd.AddCommand(pkcs11DeleteObjectCmd)
}
