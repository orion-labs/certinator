package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/orion-labs/certinator/pkg/certinator"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var certRequestFile string

func init() {
	CertCmd.AddCommand(CertCreateCmd)
}

var CertCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Certificates",
	Long: `
Create Certificates
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			if certRequestFile == "" {
				certRequestFile = args[0]
			}
		}

		if certRequestFile == "" {
			log.Fatalf("Cannot create certifiates without a properly formatted json request file.\nRequest file should look like:\n%s", certinator.ExampleCertificateRequestFile())
		}

		c, err := certinator.NewCertinator(verbose)
		if err != nil {
			log.Fatalf("Error creating Certinator: %s", err)
		}

		certRequests := make([]certinator.CertificateRequest, 0)

		reqBytes, err := ioutil.ReadFile(certRequestFile)
		if err != nil {
			log.Fatalf("error reading request file %s: %s", certRequestFile, err)
		}

		err = json.Unmarshal(reqBytes, &certRequests)
		if err != nil {
			log.Fatalf("failed unmarshalling JSON in %s: %s", certRequestFile, err)
		}

		bundle, err := c.CreateCerts(caName, certRequests)
		if err != nil {
			log.Fatalf("error creating certificates: %s", err)
		}

		for cn, data := range bundle {
			certFileName := fmt.Sprintf("%s.crt", cn)
			keyFileName := fmt.Sprintf("%s.key", cn)

			err := ioutil.WriteFile(certFileName, []byte(data.Certificate), 0644)
			if err != nil {
				log.Fatalf("err writing certificate file %s: %s", certFileName, err)
			}

			err = ioutil.WriteFile(keyFileName, []byte(data.PrivateKey), 0600)
			if err != nil {
				log.Fatalf("err writing private key file %s: %s", keyFileName, err)
			}
		}
	},
}
