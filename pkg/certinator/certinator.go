package certinator

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

const DEFAULT_CA_MAX_LEASE = "43800h0m0s"
const DEFAULT_CERTIFICATE_ROLE = "cert-issuer"

//type Config struct {
//	ServiceCaName string   `json:"service_ca_name"`
//	ClientCaName  string   `json:"client_ca_name"`
//	Services      []string `json:"services"`
//	Clients       []string `json:"clients"`
//}

type CertificateIssuingRole struct {
	Name       string
	Domains    []string
	Subdomains bool
	IpSans     bool
	Localhost  bool
	MaxTTL     string
	Ttl        string
}

// CertificateRequest  Struct for keeping track of a certificate that will be created.  Definitely NOT a CSR.
type CertificateRequest struct {
	CommonName string `json:"common_name"`
	Sans       string `json:"subject_alt_names"`
	IpSans     string `json:"ip_sans"`
	Ttl        string `json:"ttl"`
}

// CertInfo struct for holding certificate info returned from vault.
type CertInfo struct {
	IssuingCA      string
	PrivateKey     string
	PrivateKeyType string
	SerialNumber   string
	Certificate    string
	Expiration     int64
}

type CertificateBundle map[string]CertInfo

type Certinator struct {
	Client  *api.Client
	Verbose bool
}

func NewCertinator(verbose bool) (c *Certinator, err error) {
	client, err := VaultClient("", "", verbose)
	if err != nil {
		err = errors.Wrapf(err, "failed creating vault client")
		return c, err
	}

	c = &Certinator{
		Client:  client,
		Verbose: verbose,
	}

	return c, err
}

func VerboseOutput(verbose bool, message string, args ...interface{}) {
	if verbose {
		if len(args) == 0 {
			fmt.Printf("%s\n", message)
			return
		}

		msg := fmt.Sprintf(message, args...)
		fmt.Printf("%s\n", msg)
	}
}

func ExampleCertificateRequestFile() string {
	return `[
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
`
}

func (c *Certinator) UsingRootToken() (ok bool, err error) {
	secret, err := c.Client.Logical().Read("auth/token/lookup-self")
	if err != nil {
		err = errors.Wrapf(err, "failed to look up own token")
		return ok, err
	}

	if secret != nil {
		if secret.Data != nil {
			policies := secret.Data["policies"]
			policy, t := policies.([]interface{})
			if t {
				p, t := policy[0].(string)
				if t {
					if p == "root" {
						ok = true
						return ok, err
					}
				}
			}
		}
	}

	return ok, err
}

// TODO add new ca to certinator policy when a new CA is created

/*


Connect to Vault via VAULT_ADDR

Init:

    init -> read json
    unseal -> read from json

    write root token to disk
    write unseal keys to disk

Create CA's for:
        * Services
        * Clients

Create Service Certs:

    Read service list from config file or CLI.

    Create Service Certs -> Write to:
        * K8S secrets
        * PEM Files

Create Client Certs:

    Read client list from config file or CLI.

    Create Client Certs -> Write to:
        * K8S Secerts
        * PEM FIles



*/
