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
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func (c *Certinator) FetchCRL(caName string) (crlPem []byte, err error) {
	path := fmt.Sprintf("%s/crl/pem", caName)

	url := fmt.Sprintf("%s/v1/%s", c.Client.Address(), path)

	fmt.Printf("Fetching crl from %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		err = errors.Wrapf(err, "failed to get crl from %s", url)
		return crlPem, err
	}

	defer resp.Body.Close()

	crlPem, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrapf(err, "failed to read response body")
	}

	return crlPem, err
}

func (c *Certinator) RotateCRL(caName string) (err error) {
	path := fmt.Sprintf("%s/crl/rotate", caName)

	_, err = c.Client.Logical().Read(path)
	if err != nil {
		err = errors.Wrapf(err, "failed rotating CRL for %s", caName)
		return err
	}

	return err
}
