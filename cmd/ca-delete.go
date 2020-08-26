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
	"os"
)

func init() {
	CaCmd.AddCommand(CaDeleteCmd)
}

var CaDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a CA (Certificate Authority)",
	Long: `
Delete a Certificate Authority
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

		if !roottoken {
			fmt.Print("Cannot delete a CA without using the root token.  Get the root token from 1password, and please be sure to remove it from your filesystem as soon as you're done.\n\n")
			os.Exit(1)
		}

		err = c.DeleteCA(caName)
		if err != nil {
			log.Fatalf("error deleting CA %s: %s", caName, err)
		}

		fmt.Printf("CA %s deleted.\n", caName)
	},
}
