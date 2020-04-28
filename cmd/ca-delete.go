package cmd

import "github.com/spf13/cobra"

func init() {
	CaCmd.AddCommand(CaDeleteCmd)
}

var CaDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a CA (Certificate Authority)",
	Long: `
Delete a Certificate Authority
`,
}
