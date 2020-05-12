package cmd

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(CaCmd)
}

var CaCmd = &cobra.Command{
	Use:   "ca",
	Short: "CA (Certificate Authority) commands",
	Long: `
Certificate Authority commands
`,
}
