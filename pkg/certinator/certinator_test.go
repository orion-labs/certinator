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

	testAddress := fmt.Sprintf("127.0.0.1:%d", port)
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

	inputs := []struct {
		name string
	}{
		{
			"service",
		},
		{
			"client",
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

			err = c.DeleteCa(i.name)
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
