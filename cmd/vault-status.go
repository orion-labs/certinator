package cmd

import (
	"fmt"
	"github.com/orion-labs/certinator/pkg/certinator"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	VaultCmd.AddCommand(VaultStatusCmd)
}

var VaultStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Fetch Vault Status",
	Long: `
Fetch Vault Status
`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := certinator.NewCertinator(verbose)
		if err != nil {
			log.Fatalf("Error creating Certinator: %s", err)
		}

		status, err := c.VaultStatus()
		if err != nil {
			log.Fatalf("Error getting status: %s", err)
		}

		fmt.Printf("Vault Status:\n")
		fmt.Printf("  Seal Type    : %s\n", status.Type)
		fmt.Printf("  Initialized  : %v\n", status.Initialized)
		fmt.Printf("  Sealed       : %v\n", status.Sealed)
		fmt.Printf("  Total Shares : %v\n", status.T)
		fmt.Printf("  Threshold    : %v\n", status.N)
		fmt.Printf("  Version      : %v\n", status.Version)
		fmt.Printf("  Cluster Name : %v\n", status.ClusterName)
		fmt.Printf("  Cluster ID   : %v\n", status.ClusterID)
	},
}
