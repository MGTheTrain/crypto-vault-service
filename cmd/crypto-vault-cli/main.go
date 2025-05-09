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
