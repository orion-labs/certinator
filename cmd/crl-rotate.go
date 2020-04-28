package cmd

import "github.com/spf13/cobra"

func init() {
	CrlCmd.AddCommand(CrlRotateCmd)
}

var CrlRotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "Rotate CRL (Certificate Revocation List)",
	Long: `
Rotate Certificate Revocation List
`,
}
