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

		roottoken, err := c.UsingRootToken()
		if err != nil {
			log.Fatalf("failed checking own token: %s", err)
		}

		if roottoken {
			fmt.Print("You are currently using the root token.  You should not be doing this unless it's really necessary.\n\n")
		}

		if len(args) > 0 {
			if caName == "" {
				caName = args[0]
			}
		}

		if len(args) > 1 {
			if revokeCN == "" {
				revokeCN = args[1]
			}
		}

		err = c.RevokeCert(revokeCN, caName)
		if err != nil {
			log.Fatalf("error revoking certificate %s in CA %s: %s", revokeCN, caName, err)
		}

		fmt.Printf("Certificate %s deleted.\n", revokeCN)
	},
}
