package cmd

import "github.com/spf13/cobra"

func init() {
	VaultCmd.AddCommand(VaultUnsealCmd)
}

var VaultUnsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseal Vault Instances",
	Long: `
Unseal Vault Instances
`,
}
