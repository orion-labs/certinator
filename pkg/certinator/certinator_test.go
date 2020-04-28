package certinator

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/phayes/freeport"
	"github.com/scribd/vaulttest/pkg/vaulttest"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

var tmpDir string
var testServer *vaulttest.VaultDevServer
var testClient *api.Client
var testAddress string

func TestMain(m *testing.M) {
	setUp()

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func setUp() {
	dir, err := ioutil.TempDir("", "certinator")
	if err != nil {
		fmt.Printf("Error creating temp dir %q: %s\n", tmpDir, err)
		os.Exit(1)
	}

	tmpDir = dir

	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatalf("unable to get a free port on which to run the test vault server: %s", err)
	}

	testAddress = fmt.Sprintf("127.0.0.1:%d", port)
	testServer = vaulttest.NewVaultDevServer(testAddress)

	if !testServer.Running {
		testServer.ServerStart()
		testClient = testServer.VaultTestClient()

	}

}

func tearDown() {
	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
		_ = os.Remove(tmpDir)
	}

	if testServer != nil {
		testServer.ServerShutDown()
	}
}

func TestVaultInitialized(t *testing.T) {
	c := Certinator{
		Client: testClient,
	}

	ok, err := c.VaultInitialized()
	if err != nil {
		t.Errorf("failed to check init status: %s", err)
	}

	assert.True(t, ok, "Vault is not initialized")
}

func TestVaultSealed(t *testing.T) {
	c := Certinator{
		Client: testClient,
	}

	ok, err := c.VaultSealed()
	if err != nil {
		t.Errorf("failed checking vault seal: %s", err)
	}

	assert.False(t, ok, "Vault is sealed")
}

func TestCaCrud(t *testing.T) {
	c := Certinator{
		Client: testClient,
	}

	cr := CertificateIssuingRole{
		DEFAULT_CERTIFICATE_ROLE,
		[]string{"test.com"},
		true,
		true,
		true,
		"8760h",
		"8760h",
	}

	inputs := []struct {
		name  string
		role  CertificateIssuingRole
		certs []CertificateRequest
	}{
		{
			"service",
			cr,
			[]CertificateRequest{
				{
					"foo.test.com",
					"",
					"",
					"8760",
				},
			},
		},
		{
			"client",
			cr,
			[]CertificateRequest{
				{
					"bar.test.com",
					"",
					"",
					"8760",
				},
			},
		},
	}

	for _, i := range inputs {
		t.Run(i.name, func(t *testing.T) {
			err := c.CreateCa(i.name)
			if err != nil {
				t.Errorf("failed to create %s", i.name)
			}

			exists, err := c.CaExists(i.name)
			if err != nil {
				t.Errorf("failed to check if %s CA exists", i.name)
			}

			assert.True(t, exists, "CA %s does not exist", i.name)

			err = c.TuneCA(i.name)
			if err != nil {
				t.Errorf("failed to tune CA %s", i.name)
			}

			err = c.ConfigureCRL(i.name, testAddress)
			if err != nil {
				t.Errorf("failed to configure CRL for CA %s", i.name)
			}

			info, err := c.GenerateCaCert(i.name, "test.orionlabs.io", true)
			if err != nil {
				t.Errorf("failed to create CA certificate for %s", i.name)
			}

			assert.True(t, info.Data != nil)

			if info.Data != nil {
				assert.True(t, info.Data["private_key"] != nil, "No Private Key")
				assert.True(t, info.Data["certificate"] != nil, "No CA Certificate")

			}

			err = c.CreateIssuingRole(i.name, i.role)
			if err != nil {
				t.Errorf("failed to create issuing role %s in CA %s", i.role.Name, i.name)
			}

			certs, err := c.CreateCerts(i.name, i.certs)
			if err != nil {
				t.Errorf("failed to create certs %s", err)
			}

			for cn, cert := range certs {
				assert.True(t, cert.PrivateKey != "", "Cert for %s has no private key!", cn)
				assert.True(t, cert.Certificate != "", "Cert for %s has no certificate", cn)
			}

			err = c.DeleteCA(i.name)
			if err != nil {
				t.Errorf("failed to delete CA %s", i.name)
			}

			exists, err = c.CaExists(i.name)
			if err != nil {
				t.Errorf("failed to check if %s CA exists", i.name)
			}
			assert.False(t, exists, "CA %s exists", i.name)

		})
	}
}
