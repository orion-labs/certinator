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
