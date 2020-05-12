package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	VaultCmd.AddCommand(VaultInitCmd)
}

var VaultInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Vault instance",
	Long: `
Initialize a Vault instance
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Vault init has not been implemented yet.\nIn the mean time, try running `vault operator init` directly against the vault server.")
	},
}
