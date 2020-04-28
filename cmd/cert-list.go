package cmd

import "github.com/spf13/cobra"

func init() {
	CertCmd.AddCommand(CertListCmd)
}

var CertListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Certificates",
	Long: `
List Certificates
`,
}
