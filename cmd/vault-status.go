// Copyright © Copyright 2020 Orion Labs, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
