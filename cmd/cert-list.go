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

		roottoken, err := c.UsingRootToken()
		if err != nil {
			log.Fatalf("failed checking own token: %s", err)
		}

		if roottoken {
			fmt.Print("You are currently using the root token.  You should not be doing this unless it's really necessary.\n\n")
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
