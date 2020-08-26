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
	CaCmd.AddCommand(CaListCmd)
}

var CaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List CA's (Certificate Authorities)",
	Long: `
List Certificate Authorities
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

		roottoken, err := c.UsingRootToken()
		if err != nil {
			log.Fatalf("failed checking own token: %s", err)
		}

		if roottoken {
			fmt.Print("You are currently using the root token.  You should not be doing this unless it's really necessary.\n\n")
		}

		cas, err := c.ListCAs()
		if err != nil {
			log.Fatalf("error listing CA's: %s", err)
		}

		if len(cas) == 0 {
			fmt.Printf("No configured Certificate Authorities.\n")
			return
		}

		fmt.Printf("Certificate Authorities:\n")
		for _, ca := range cas {
			fmt.Printf("  %s\n", ca)
		}
	},
}
