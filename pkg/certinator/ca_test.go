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
	"github.com/stretchr/testify/assert"
	"testing"
)

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
			err := c.CreateCA(i.name)
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

			crlPem, err := c.FetchCRL(i.name)
			if err != nil {
				t.Errorf("error: %s", err)
			}

			assert.True(t, len(crlPem) != 0, "No crl fetched!")

			err = c.RotateCRL(i.name)
			if err != nil {
				t.Errorf("error rotating crl: %s", err)
			}

			cas, err := c.ListCAs()
			if err != nil {
				t.Errorf("error listing CA's: %s", err)
			}

			assert.True(t, len(cas) > 0, "Cannot list CA's!")

			fetchedCerts, err := c.ListCerts(i.name)
			if err != nil {
				t.Errorf("error listing certs: %s", err)
			}

			assert.True(t, len(fetchedCerts) > 0, "Fetched no certs!")

			for _, cert := range i.certs {
				err = c.RevokeCert(cert.CommonName, i.name)
				if err != nil {
					t.Errorf("failed revoking certificate %s: %s", cert.CommonName, err)
				}
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
