package cmd

import "github.com/spf13/cobra"

func init() {
	CertCmd.AddCommand(CertCreateCmd)
}

var CertCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Certificates",
	Long: `
Create Certificates
`,
}
