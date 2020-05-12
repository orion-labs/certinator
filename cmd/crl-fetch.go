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
Fetch Certificate Revocation List.

Writes it to STDOUT for easy scripting or piping to openssl.  e.g

	certinator crl fetch service | openssl crl -text -noout
	
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

		if caName == "" {
			log.Fatalf("Must specify CA name to fetch CRL.\nTry again with -c <ca name> or `certinator crl fetch <ca name>`.\n")
		}

		crlPem, err := c.FetchCRL(caName)
		if err != nil {
			log.Fatalf("Error fetching CRL: %s", err)
		}

		fmt.Printf("%s\n", crlPem)
	},
}
