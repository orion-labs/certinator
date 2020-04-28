package cmd

import "github.com/spf13/cobra"

func init() {
	VaultCmd.AddCommand(VaultInitCmd)
}

var VaultInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Vault instance",
	Long: `
Initialize a Vault instance
`,
}
