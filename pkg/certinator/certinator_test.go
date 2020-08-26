// Copyright © Copyright 2020 Orion Labs, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func TestVaultStatus(t *testing.T) {
	c := Certinator{
		Client: testClient,
	}

	status, err := c.VaultStatus()
	if err != nil {
		t.Errorf("failed checking vault seal: %s", err)
	}

	assert.True(t, status != nil, "Nil Vault Status!")

}

func TestUsingRootToken(t *testing.T) {
	c := Certinator{
		Client: testClient,
	}

	root, err := c.UsingRootToken()
	if err != nil {
		t.Errorf("Error checking my own token: %s\n", err)
	}

	assert.True(t, root, "Not using root token")
}
