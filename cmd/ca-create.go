package cmd

import (
	"fmt"
	"github.com/orion-labs/certinator/pkg/certinator"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var caCAAddress string
var caCN string

func init() {
	CaCmd.AddCommand(CaCreateCmd)
	CaCreateCmd.Flags().StringVarP(&caCAAddress, "address", "a", "", "Public-ish address for the CA (used for configuring CRL's)")
	CaCreateCmd.Flags().StringVarP(&caCN, "name", "n", "", "Common Name for the CA Certificate")
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

		// TODO provide for configuring issuing role via dialog
		// this works as a crude start, but should be configurable in the future.
		standardIssuingRole := certinator.CertificateIssuingRole{
			Name:       "cert-issuer",
			Domains:    []string{"orion.svc.cluster.local"},
			Subdomains: true,
			IpSans:     false,
			Localhost:  false,
			MaxTTL:     "8760h",
			Ttl:        "8760h",
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

		err = c.TuneCA(caName)
		if err != nil {
			log.Fatalf("failed to tune CA %s", caName)
		}

		// TODO provide for configuring CA Address by flag and dialog
		if caCAAddress == "" {
			// crl only works from within the cluster
			caCAAddress = fmt.Sprintf("http://vault.orion.svc.cluster.local/v1/%s/config/crl", caName)
		}

		err = c.ConfigureCRL(caName, caCAAddress)
		if err != nil {
			log.Fatalf("failed to configure CRL for CA %s", caName)
		}

		// TODO provide for configuraing CA CN via flag and dialog
		if caCN == "" {
			// default CN based on CA name
			caCN = fmt.Sprintf("%s.orion-ptt.orion.svc.cluster.local", caName)
		}

		// TODO add option to export CA private key
		info, err := c.GenerateCaCert(caName, caCN, false)
		if err != nil {
			log.Fatalf("failed to create CA certificate for %s", caName)
		}

		err = c.CreateIssuingRole(caName, standardIssuingRole)
		if err != nil {
			log.Fatalf("failed to create issuing role in CA %s", caName)
		}

		caFile := fmt.Sprintf("%s-ca.crt", caName)
		// TODO Provide for exporting private keys

		if info != nil {
			if info.Data != nil {
				caData, ok := info.Data["issuing_ca"].(string)
				if !ok {
					log.Fatalf("Can't unmarshal ca certificate data.")
				}

				err = ioutil.WriteFile(caFile, []byte(caData), 0644)
				if err != nil {
					log.Fatalf("failed to write file %s", caFile)
				}

				fmt.Printf("CA %s created.\n", caName)

				return
			}
		}

	},
}
