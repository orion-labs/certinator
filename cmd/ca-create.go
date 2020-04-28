package cmd

import "github.com/spf13/cobra"

func init() {
	CaCmd.AddCommand(CaCreateCmd)
}

var CaCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a CA (Certificate Authority)",
	Long: `
Create a Certificate Authority
`,
}
