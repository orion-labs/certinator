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
var certRequestCommonName string
var certRequestSans string
var certRequestIpSans string
var certRequestTtl string

func init() {
	CertCmd.AddCommand(CertCreateCmd)
	CertCreateCmd.Flags().StringVarP(&certRequestFile, "file", "f", "", "JSON file containing info on certificates to create")
	CertCreateCmd.Flags().StringVarP(&certRequestCommonName, "name", "n", "", "common name of certificate to be created")
	CertCreateCmd.Flags().StringVarP(&certRequestSans, "sans", "s", "", "Subject Alternate Names for certificate.")
	CertCreateCmd.Flags().StringVarP(&certRequestIpSans, "ipsans", "i", "", "IP sans for certificate.")
	CertCreateCmd.Flags().StringVarP(&certRequestTtl, "ttl", "t", "", "ttl for certificate")

}

var CertCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Certificates",
	Long: `
Create Certificates.

Create a single certificate with:

	certinator cert create -c <ca name> -n <common name> -s <comma separated list of Subject Alternate Names> -i <comma separated list of IP SAN's> -t <ttl of certificate>

Create multiple certificates with a config file of format:
	
	[
		{
			"common_name": "some.cert.name.com",
			"subject_alternate_names": "alt.name.one,alt.name.two,alt.name.three",
			"ip_sans": "192.168.0.1,192.168.0.2,192.168.0.3",
			"ttl": "8760h"
		},
		{
			"common_name": "some.other.cert.name.com",
			"subject_alternate_names": "alt.name.one,alt.name.two,alt.name.three",
			"ip_sans": "192.168.0.1,192.168.0.2,192.168.0.3",
			"ttl": "8760h"
		}
	]

Via:
	certinator cert create -f <file name>

NB: TTL's are in seconds unless a unit suffix is supplied.  i.e. 8760 seconds.  8760h is 8760 hours or roughly one year.  Forgetting to add the units will lead to surprising behavior.

`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			if caName == "" {
				caName = args[0]
			}
		}

		if caName == "" {
			log.Fatalf("Cannot issue certificates without a CA name.  Try again with `certinator cert create -c <ca name>` or `certinator cert create <ca name>`")
		}

		// at minimum, a common name and a ttl are required
		if certRequestCommonName != "" && certRequestTtl == "" {
			// or if not, you need a config file
			if certRequestFile == "" {
				log.Fatalf("Cannot create certifiates without either a common name and a ttl or a properly formatted json request file.\nRequest file should look like:\n%s", certinator.ExampleCertificateRequestFile())
			}
		}

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

		certRequests := make([]certinator.CertificateRequest, 0)

		if certRequestFile != "" {
			reqBytes, err := ioutil.ReadFile(certRequestFile)
			if err != nil {
				log.Fatalf("error reading request file for CA %q: %s", certRequestFile, err)
			}

			err = json.Unmarshal(reqBytes, &certRequests)
			if err != nil {
				log.Fatalf("failed unmarshalling JSON in %s: %s", certRequestFile, err)
			}
		} else {
			certRequests = append(certRequests, certinator.CertificateRequest{
				CommonName: certRequestCommonName,
				Sans:       certRequestSans,
				IpSans:     certRequestIpSans,
				Ttl:        certRequestTtl,
			})
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

			fmt.Printf("Files Written:\n  %s\n  %s\n\n", certFileName, keyFileName)
		}
	},
}
