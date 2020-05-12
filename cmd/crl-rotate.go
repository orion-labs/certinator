package cmd

import (
	"fmt"
	"github.com/orion-labs/certinator/pkg/certinator"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	CrlCmd.AddCommand(CrlRotateCmd)
}

var CrlRotateCmd = &cobra.Command{
	Use:   "rotate",
	Short: "RotateCRL (Certificate Revocation List)",
	Long: `
Rotate Certificate Revocation List
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

		err = c.RotateCRL(caName)
		if err != nil {
			log.Fatalf("Error fetching CRL: %s", err)
		}

		fmt.Printf("CRL for %s rotated.\n", caName)
	},
}
