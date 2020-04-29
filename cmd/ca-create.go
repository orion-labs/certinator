package cmd

import (
	"fmt"
	"github.com/orion-labs/certinator/pkg/certinator"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	CaCmd.AddCommand(CaCreateCmd)
}

var CaCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a CA (Certificate Authority)",
	Long: `
Create a Certificate Authority
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

		err = c.CreateCA(caName)
		if err != nil {
			log.Fatalf("error creating CA %s: %s", caName, err)
		}

		fmt.Printf("CA %s created.\n", caName)
	},
}
