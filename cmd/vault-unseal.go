package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	VaultCmd.AddCommand(VaultUnsealCmd)
}

var VaultUnsealCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseal Vault Instances",
	Long: `
Unseal Vault Instances
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Vault unseal has not been implemented yet.\nIn the mean time, try running `vault operator unseal <key>` directly against the vault server.")
	},
}
