package commands

import (
	"crypto_vault_service/internal/infrastructure/cryptography"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// Command for Slot Setup
func Pkcs11SlotSetupCmd(cmd *cobra.Command, args []string) {
	// Read flags for the required parameters
	modulePath, _ := cmd.Flags().GetString("module")
	tokenLabel, _ := cmd.Flags().GetString("token-label")
	soPin, _ := cmd.Flags().GetString("so-pin")
	userPin, _ := cmd.Flags().GetString("user-pin")

	// Create an instance of PKCS11Token with the provided flags
	token := &cryptography.PKCS11Token{
		ModulePath: modulePath,
		TokenLabel: tokenLabel,
		SOPin:      soPin,
		UserPin:    userPin,
	}

	// Call the method to set up the token slot
	if err := token.Pkcs11SlotSetup(); err != nil {
		log.Fatalf("Error setting up slot: %v", err)
	}
	fmt.Println("Token slot setup successfully!")
}

// Command to check if token is set
func IsTokenSetCmd(cmd *cobra.Command, args []string) {
	modulePath, _ := cmd.Flags().GetString("module")
	tokenLabel, _ := cmd.Flags().GetString("token-label")

	token := &cryptography.PKCS11Token{
		ModulePath: modulePath,
		TokenLabel: tokenLabel,
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
	tokenLabel, _ := cmd.Flags().GetString("token-label")
	objectLabel, _ := cmd.Flags().GetString("object-label")
	userPin, _ := cmd.Flags().GetString("user-pin")

	token := &cryptography.PKCS11Token{
		ModulePath:  modulePath,
		TokenLabel:  tokenLabel,
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
	modulePath, _ := cmd.Flags().GetString("module")
	tokenLabel, _ := cmd.Flags().GetString("token-label")
	soPin, _ := cmd.Flags().GetString("so-pin")
	userPin, _ := cmd.Flags().GetString("user-pin")

	token := &cryptography.PKCS11Token{
		ModulePath: modulePath,
		TokenLabel: tokenLabel,
		SOPin:      soPin,
		UserPin:    userPin,
	}

	if err := token.InitializeToken(); err != nil {
		log.Fatalf("Error initializing token: %v", err)
	}
	fmt.Println("Token initialized successfully!")
}

// Command to get free slot for the PKCS#11 token
func GetFreeSlotCmd(cmd *cobra.Command, args []string) {
	modulePath, _ := cmd.Flags().GetString("module")

	token := &cryptography.PKCS11Token{
		ModulePath: modulePath,
	}

	slot, err := token.GetFreeSlot()
	if err != nil {
		log.Fatalf("Error getting free slot: %v", err)
	}
	fmt.Printf("Free slot: %s\n", slot)
}

// Command to add key to PKCS#11 token
func AddKeyCmd(cmd *cobra.Command, args []string) {
	modulePath, _ := cmd.Flags().GetString("module")
	tokenLabel, _ := cmd.Flags().GetString("token-label")
	objectLabel, _ := cmd.Flags().GetString("object-label")
	keyType, _ := cmd.Flags().GetString("key-type") // "ECDSA" or "RSA"
	keySize, _ := cmd.Flags().GetInt("key-size")
	userPin, _ := cmd.Flags().GetString("user-pin")

	token := &cryptography.PKCS11Token{
		ModulePath:  modulePath,
		TokenLabel:  tokenLabel,
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
	tokenLabel, _ := cmd.Flags().GetString("token-label")
	objectLabel, _ := cmd.Flags().GetString("object-label")
	objectType, _ := cmd.Flags().GetString("object-type")
	userPin, _ := cmd.Flags().GetString("user-pin")

	token := &cryptography.PKCS11Token{
		ModulePath:  modulePath,
		TokenLabel:  tokenLabel,
		ObjectLabel: objectLabel,
		UserPin:     userPin,
	}

	if err := token.DeleteObject(objectType, objectLabel); err != nil {
		log.Fatalf("Error deleting object: %v", err)
	}
	fmt.Println("Object deleted successfully!")
}

func InitPKCS11Commands(rootCmd *cobra.Command) {
	var pkcs11SlotSetupCmd = &cobra.Command{
		Use:   "setup-slot",
		Short: "Setup PKCS#11 token slot",
		Run:   Pkcs11SlotSetupCmd,
	}
	pkcs11SlotSetupCmd.Flags().String("module", "", "Path to the PKCS#11 module")
	pkcs11SlotSetupCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11SlotSetupCmd.Flags().String("so-pin", "", "Security Officer PIN")
	pkcs11SlotSetupCmd.Flags().String("user-pin", "", "User PIN")
	rootCmd.AddCommand(pkcs11SlotSetupCmd)

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
	pkcs11InitializeTokenCmd.Flags().String("module", "", "Path to the PKCS#11 module")
	pkcs11InitializeTokenCmd.Flags().String("token-label", "", "Label of the PKCS#11 token")
	pkcs11InitializeTokenCmd.Flags().String("so-pin", "", "Security Officer PIN")
	pkcs11InitializeTokenCmd.Flags().String("user-pin", "", "User PIN")
	rootCmd.AddCommand(pkcs11InitializeTokenCmd)

	var pkcs11GetFreeSlotCmd = &cobra.Command{
		Use:   "get-free-slot",
		Short: "Get the next available uninitialized slot",
		Run:   GetFreeSlotCmd,
	}
	pkcs11GetFreeSlotCmd.Flags().String("module", "", "Path to the PKCS#11 module")
	rootCmd.AddCommand(pkcs11GetFreeSlotCmd)

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
