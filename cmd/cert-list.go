package cmd

import (
	"fmt"
	"github.com/orion-labs/certinator/pkg/certinator"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	CertCmd.AddCommand(CertListCmd)
}

var CertListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Certificates",
	Long: `
List Certificates
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

		certs, err := c.ListCerts(caName)
		if err != nil {
			log.Fatalf("error listing certificates on CA %s: %s", caName, err)
		}

		if len(certs) == 0 {
			fmt.Printf("No Certificates created in CA %s\n", caName)
			return
		}

		fmt.Printf("Certificates in CA %s:\n", caName)
		for _, c := range certs {
			fmt.Printf("  %s\n", c)
		}
	},
}
