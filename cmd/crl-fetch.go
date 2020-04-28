package cmd

import "github.com/spf13/cobra"

func init() {
	CrlCmd.AddCommand(CrlFetchCmd)
}

var CrlFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch CRL (Certificate Revocation List)",
	Long: `
Fetch Certificate Revocation List
`,
}
