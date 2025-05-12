// Package main is the entry point for the crypto-vault-cli application.
// It initializes the root command and registers various sub-commands (AES, RSA, ECDSA, PKCS#11)
// for the CLI, then executes the command-line interface.
package main

import (
	"log"

	commands "crypto_vault_service/cmd/crypto-vault-cli/internal/commands"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "crypto-vault-cli"}

	commands.InitAESCommands(rootCmd)
	commands.InitRSACommands(rootCmd)
	commands.InitECDSACommands(rootCmd)

	_, err := commands.ReadPkcs11SettingsFromEnv()
	if err == nil {
		commands.InitPKCS11Commands(rootCmd)
	} else {
		log.Println("Skipping PKCS#11 command registration: ", err.Error())
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}
