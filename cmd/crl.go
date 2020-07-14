package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(CrlCmd)
}

var CrlCmd = &cobra.Command{
	Use:   "crl",
	Short: "CRL (Certificate Revocation List) commands",
	Long: `
Certificate Revocation List commands
`,
}
