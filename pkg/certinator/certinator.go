package certinator

import (
	"github.com/hashicorp/vault/api"
)

const DEFAULT_SERVICE_CA = "service"
const DEFAULT_CLIENT_CA = "client"
const DEFAULT_CA_MAX_LEASE = "43800h0m0s"
const DEFAULT_CERTIFICATE_ROLE = "cert-issuer"

type Config struct {
	ServiceCaName string   `json:"service_ca_name"`
	ClientCaName  string   `json:"client_ca_name"`
	Services      []string `json:"services"`
	Clients       []string `json:"clients"`
}

type CertificateIssuingRole struct {
	Name       string
	Domains    []string
	Subdomains bool
	IpSans     bool
	Localhost  bool
	MaxTTL     string
	Ttl        string
}

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
	Client *api.Client
	Config *Config
}

// CertificateRequest  Struct for keeping track of a certificate that will be created.  Definitely NOT a CSR.
type CertificateRequest struct {
	CommonName string
	Sans       string
	IpSans     string
	Ttl        string
}

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
