package cmd

import (
	"fmt"
	"github.com/orion-labs/certinator/pkg/certinator"
	"github.com/spf13/cobra"
	"log"
)

var revokeCN string

func init() {
	CertCmd.AddCommand(CertRevokeCmd)
	CertRevokeCmd.Flags().StringVarP(&revokeCN, "commonname", "n", "", "common name of certificate to revoke.")
}

var CertRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke Certificates",
	Long: `
Revoke Certificates
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

		err = c.RevokeCert(revokeCN, caName)
		if err != nil {
			log.Fatalf("error revoking certificate %s in CA %s: %s", revokeCN, caName, err)
		}

		fmt.Printf("Certificate %s deleted.\n", revokeCN)
	},
}
