package cmd

import (
	"fmt"
	"github.com/orion-labs/certinator/pkg/certinator"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	CaCmd.AddCommand(CaListCmd)
}

var CaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List CA's (Certificate Authorities)",
	Long: `
List Certificate Authorities
`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := certinator.NewCertinator(verbose)
		if err != nil {
			log.Fatalf("Error creating Certinator: %s", err)
		}

		if len(args) > 0 {
			if caName == "" {
				caName = args[1]
			}
		}

		cas, err := c.ListCAs()
		if err != nil {
			log.Fatalf("error listing CA's: %s", err)
		}

		fmt.Printf("CA's:\n%s.\n", cas)
	},
}
