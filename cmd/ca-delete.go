package cmd

import (
	"fmt"
	"github.com/orion-labs/certinator/pkg/certinator"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	CaCmd.AddCommand(CaDeleteCmd)
}

var CaDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a CA (Certificate Authority)",
	Long: `
Delete a Certificate Authority
`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := certinator.NewCertinator(verbose)
		if err != nil {
			log.Fatalf("Error creating Certinator: %s", err)
		}

		if len(args) > 0 {
			if caName == "" {
				caName = args[0]
			}
		}

		err = c.DeleteCA(caName)
		if err != nil {
			log.Fatalf("error deleting CA %s: %s", caName, err)
		}

		fmt.Printf("CA %s deleted.\n", caName)
	},
}
