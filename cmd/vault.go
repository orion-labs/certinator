package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(VaultCmd)
}

var VaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Vault management commands",
	Long: `
Vault management commands
`,
}
