package cmd

import "github.com/spf13/cobra"

func init() {
	CertCmd.AddCommand(CertRevokeCmd)
}

var CertRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke Certificates",
	Long: `
Revoke Certificates
`,
}
