package cmd

import "github.com/spf13/cobra"

func init() {
	CaCmd.AddCommand(CaListCmd)
}

var CaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List CA's (Certificate Authorities)",
	Long: `
List Certificate Authorities
`,
}
