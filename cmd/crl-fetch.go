package cmd

import (
	"fmt"
	"github.com/orion-labs/certinator/pkg/certinator"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	CrlCmd.AddCommand(CrlFetchCmd)

}

var CrlFetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch CRL (Certificate Revocation List)",
	Long: `
Fetch Certificate Revocation List
`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := certinator.NewCertinator(verbose)
		if err != nil {
			log.Fatalf("Error creating Certinator: %s", err)
		}

		crlPem, err := c.FetchCRL(caName)
		if err != nil {
			log.Fatalf("Error fetching CRL: %s", err)
		}

		fmt.Printf("%s\n", crlPem)
	},
}
