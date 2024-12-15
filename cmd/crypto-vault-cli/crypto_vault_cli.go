package main

import (
	"fmt"
	"os"

	commands "crypto_vault_service/cmd/crypto-vault-cli/internal/commands"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "crypto-vault-cli"}

	commands.InitAESCommands(rootCmd)

	commands.InitRSACommands(rootCmd)

	commands.InitECDSACommands(rootCmd)

	commands.InitPKCS11Commands(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
