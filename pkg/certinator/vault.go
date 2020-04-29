package certinator

import (
	tls2 "crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const VAULT_TOKEN_ENV_VAR = "VAULT_TOKEN"
const DEFAULT_VAULT_TOKEN_FILE = ".vault-token"

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

func (c *Certinator) InitVault() (err error) {
	// TODO implement and test InitVault
	/*
		PUT /v1/sys/init HTTP/1.1
		Host: localhost:8080
		User-Agent: Go-http-client/1.1
		Content-Length: 166
		X-Vault-Request: true
		X-Vault-Token: s.GJKSSV5Pv5e9jqx2Eoz5dtoD
		Accept-Encoding: gzip

		{"secret_shares":5,"secret_threshold":3,"stored_shares":0,"pgp_keys":null,"recovery_shares":5,"recovery_threshold":3,"recovery_pgp_keys":null,"root_token_pgp_key":""}


	*/

	return err
}

func (c *Certinator) UnsealVault() (err error) {
	// TODO implement and test UnsealVault
	/*
		PUT /v1/sys/unseal HTTP/1.1
		Host: localhost:8080
		User-Agent: Go-http-client/1.1
		Content-Length: 43
		X-Vault-Request: true
		X-Vault-Token: s.GJKSSV5Pv5e9jqx2Eoz5dtoD
		Accept-Encoding: gzip

		{"key":"foo","reset":false,"migrate":false}
	*/
	return err
}

func (c *Certinator) VaultStatus() (status *api.SealStatusResponse, err error) {
	// GET /v1/sys/seal-status HTTP/1.1
	status, err = c.Client.Sys().SealStatus()
	if err != nil {
		err = errors.Wrapf(err, "failed to get status")
		return status, err
	}

	return status, err
}

func VaultClient(address string, cacert string, verbose bool) (client *api.Client, err error) {
	apiConfig, err := ApiConfig(address, cacert)

	if verbose {
		fmt.Printf("Vault Address: %s\n", apiConfig.Address)
		if cacert != "" {
			fmt.Printf("Private CA Cert in use.\n")
		}
	}

	client, err = api.NewClient(apiConfig)
	if err != nil {
		err = errors.Wrapf(err, "failed to create vault api client")
		return client, err
	}

	// Straight up take the token from the environment if provided
	if os.Getenv(VAULT_TOKEN_ENV_VAR) != "" {
		client.SetToken(os.Getenv(VAULT_TOKEN_ENV_VAR))
		return client, err
	}

	// Attempt to use a token on the filesystem if it exists
	ok, err := UseFSToken(client, verbose)
	if err != nil {
		err = errors.Wrapf(err, "failed to make use of filesystem token")
		return client, err
	}

	if ok {
		return client, err
	}

	return client, err
}

// ApiConfig creates a vault api config in a standard fashion.  Stolen from vault-authenticator.
func ApiConfig(address string, cacert string) (config *api.Config, err error) {
	// read the environment and use that over anything
	config = api.DefaultConfig()

	err = config.ReadEnvironment()
	if err != nil {
		err = errors.Wrapf(err, "failed to inject environment into client config")
		return config, err
	}

	if config.Address == "https://127.0.0.1:8200" {
		if address != "" {
			config.Address = address
		}
	}

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		err = errors.Wrapf(err, "failed to get system cert pool")
		return config, err
	}

	if cacert != "" {
		ok := rootCAs.AppendCertsFromPEM([]byte(cacert))
		if !ok {
			err = errors.New("Failed to add scribd root cert to system CA bundle")
			return config, err
		}
	}

	clientConfig := &tls2.Config{
		RootCAs: rootCAs,
	}

	config.HttpClient.Transport = &http.Transport{TLSClientConfig: clientConfig}

	return config, err
}

// UseFSToken Attempts to use a Vault Token found on the filesystem.
func UseFSToken(client *api.Client, verbose bool) (ok bool, err error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		err = errors.Wrap(err, "failed to get user's homedir")
		return ok, err
	}

	tokenFilePath := fmt.Sprintf("%s/%s", homeDir, DEFAULT_VAULT_TOKEN_FILE)
	VerboseOutput(verbose, "Looking for a potential vault token at %s", tokenFilePath)

	if _, existErr := os.Stat(tokenFilePath); !os.IsNotExist(existErr) {
		VerboseOutput(verbose, "  It exists.")

		b, err := ioutil.ReadFile(tokenFilePath)
		if err != nil {
			err = errors.Wrapf(err, "failed to read token out of %s", tokenFilePath)
			return ok, err
		}

		token := string(b)

		// chomp
		token = strings.TrimRight(token, "\n")

		client.SetToken(token) // set token

		if token == "" {
			VerboseOutput(verbose, "  token file has no content.")
			return ok, err
		}

		_, tokOkErr := client.Auth().Token().LookupSelf()
		if tokOkErr != nil {
			VerboseOutput(verbose, "  token is not valid.")
			// don't blow up, just return false, and let auth proceed
			return ok, err
		}

		VerboseOutput(verbose, "  token set.")

		err = RenewTokenIfStale(client, verbose)
		if err != nil {
			return ok, err
		}

		ok = true
		return ok, err
	}

	VerboseOutput(verbose, "  no token found.  Moving on to other auth methods.")

	return ok, err
}

// RenewTokenIfStale renews a Vault token if it happens to be near expiration.
func RenewTokenIfStale(client *api.Client, verbose bool) (err error) {
	// at this point, we have a token, either from the env or the filesystem.
	// renew the token, since it may be near expiration
	// don't really care about the error result of this call, as some tokens are not refreshable, and this is mostly a convenience feature so the user doesn't have to login.
	_, _ = client.Auth().Token().RenewSelf(0)

	return err
}
