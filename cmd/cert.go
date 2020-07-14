package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(CertCmd)
}

var CertCmd = &cobra.Command{
	Use:   "cert",
	Short: "Certificate commands",
	Long: `
Certificate commands
`,
}
