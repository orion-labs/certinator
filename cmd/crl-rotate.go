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

		err = c.RotateCRL(caName)
		if err != nil {
			log.Fatalf("Error fetching CRL: %s", err)
		}

		fmt.Printf("CRL for %s rotated.\n", caName)
	},
}
