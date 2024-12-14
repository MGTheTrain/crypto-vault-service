package main

import (
	"fmt"
	"os"

	commands "crypto_vault_service/cmd/crypto-vault-cli/internal/commands"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "crypto-vault-cli"}

	// AES Commands
	commands.InitAESCommands(rootCmd)

	// RSA Commands
	commands.InitRSACommands(rootCmd)

	// ECDSA Commands
	commands.InitECDSACommands(rootCmd)

	// PKCS11 Token Commands
	commands.InitPKCS11Commands(rootCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
