package cmd

import (
	"fmt"
	"github.com/orion-labs/certinator/pkg/certinator"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var caCAAddress string
var caCN string
var caExportKey bool
var caIssuingDomains string

func init() {
	CaCmd.AddCommand(CaCreateCmd)
	CaCreateCmd.Flags().StringVarP(&caCAAddress, "address", "a", "", "Public-ish address for the CA (used for configuring CRL's)")
	CaCreateCmd.Flags().StringVarP(&caCN, "name", "n", "", "Common Name for the CA Certificate")
	CaCreateCmd.Flags().BoolVarP(&caExportKey, "export", "e", false, "Export CA Private Key")
	CaCreateCmd.Flags().StringVarP(&caIssuingDomains, "issuedomains", "d", "", "Allowed domains for certificate issuance.  Comma separated list.  e.g.: 'foo.com,bar.com'")
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

		domains := make([]string, 0)

		if caIssuingDomains == "" {
			domains = append(domains, "orion.svc.cluster.local")
		} else {
			commaRegex := regexp.MustCompile(`.+,.+`)
			if commaRegex.MatchString(caIssuingDomains) {
				domains = strings.Split(caIssuingDomains, ",")
			} else {
				domains = append(domains, caIssuingDomains)
			}
		}

		// TODO provide for configuring issuing role via dialog
		// this works as a crude start, but should be configurable in the future.
		standardIssuingRole := certinator.CertificateIssuingRole{
			Name:       "cert-issuer",
			Domains:    domains,
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

		roottoken, err := c.UsingRootToken()
		if err != nil {
			log.Fatalf("failed checking own token: %s", err)
		}

		if !roottoken {
			fmt.Print("Cannot create a CA without using the root token.  Get the root token from 1password, and please be sure to remove it from your filesystem as soon as you're done.\n\n")
			os.Exit(1)
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
		info, err := c.GenerateCaCert(caName, caCN, caExportKey)
		if err != nil {
			log.Fatalf("failed to create CA certificate for %s", caName)
		}

		err = c.CreateIssuingRole(caName, standardIssuingRole)
		if err != nil {
			log.Fatalf("failed to create issuing role in CA %s", caName)
		}

		caFile := fmt.Sprintf("%s-ca.crt", caName)
		caKeyFile := fmt.Sprintf("%s-ca.key", caName)

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

				if caExportKey {
					caKeyData, ok := info.Data["private_key"].(string)
					if !ok {
						log.Fatalf("Can't unmarshal ca private key data.")
					}

					err = ioutil.WriteFile(caKeyFile, []byte(caKeyData), 0600)
					if err != nil {
						log.Fatalf("failed to write file %s", caKeyFile)
					}
				}

				fmt.Printf("CA %s created.\n  CA Certificate written to %s\n", caName, caFile)

				if caExportKey {
					fmt.Printf("  CA Private Key writen to %s\n", caKeyFile)
				}

				return
			}
		}

	},
}
