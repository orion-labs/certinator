package certinator

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

const DEFAULT_SERVICE_CA = "service"
const DEFAULT_CLIENT_CA = "client"

type Config struct {
	ServiceCaName string   `json:"service_ca_name"`
	ClientCaName  string   `json:"client_ca_name"`
	Services      []string `json:"services"`
	Clients       []string `json:"clients"`
}

type Certinator struct {
	Client *api.Client
	Config *Config
}

func (c *Certinator) VaultInitialized() (ok bool, err error) {

	ok, err = c.Client.Sys().InitStatus()

	return ok, err
}

func (c *Certinator) VaultSealed() (ok bool, err error) {
	status, err := c.Client.Sys().SealStatus()
	if err != nil {
		err = errors.Wrap(err, "failed to check seal status")
		return ok, err
	}

	ok = status.Sealed

	return ok, err
}

func (c *Certinator) CaExists(name string) (ok bool, err error) {
	path := "sys/mounts"
	secret, err := c.Client.Logical().Read(path)
	if err != nil {
		err = errors.Wrapf(err, "failed to read %s", path)
	}

	mountName := fmt.Sprintf("%s/", name)

	for k, v := range secret.Data {
		if k == mountName {
			d, ok := v.(map[string]interface{})
			if ok {
				if d["type"] == "pki" {
					ok = true
					return ok, err
				}
			}
		}
	}
	return ok, err
}

func (c *Certinator) CreateCa(name string) (err error) {
	data := map[string]interface{}{
		"type":        "pki",
		"description": fmt.Sprintf("%s certificate authority", name),
	}

	path := fmt.Sprintf("sys/mounts/%s", name)

	_, err = c.Client.Logical().Write(path, data)
	if err != nil {
		err = errors.Wrapf(err, "failed creating %s CA", name)
		return err
	}

	return err
}

func (c *Certinator) DeleteCa(name string) (err error) {
	path := fmt.Sprintf("sys/mounts/%s", name)

	_, err = c.Client.Logical().Delete(path)
	if err != nil {
		err = errors.Wrapf(err, "Failed to delete CA at %s", name)
		return err
	}

	return err
}

func (c *Certinator) InitVault() (err error) {
	// TODO implement and test InitVault
	return err
}

func (c *Certinator) UnsealVault() (err error) {
	// TODO implement and test UnsealVault
	return err
}

/*

vault secrets enable -path=kafka -description="Kafka Client Certificate Authority" -max-lease-ttl=43800h pki
vault secrets tune -max-lease-ttl=43800h kafka
vault write kafka/root/generate/internal common_name=kafka-ca.scribd.com ttl=43800h | tee cacert
vault write kafka/config/urls issuing_certificates=${VAULT_ADDR}/v1/kafka/ca crl_distribution_points=${VAULT_ADDR}/v1/kafka/crl
vault write kafka/roles/cert-issuer allowed_domains="kafka.scribd.com" allow_subdomains="true" max_ttl="8760h" ttl="8760h" allow_ip_sans=true allow_localhost=true

Connect to Vault via VAULT_ADDR

Init:

    init -> read json
    unseal -> read from json

    write root token to disk
    write unseal keys to disk

    create CA's for
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
