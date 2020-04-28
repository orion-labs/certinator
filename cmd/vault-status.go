package cmd

import "github.com/spf13/cobra"

func init() {
	VaultCmd.AddCommand(VaultStatusCmd)
}

var VaultStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Fetch Vault Status",
	Long: `
Fetch Vault Status
`,
}
