package certinator

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

func (c *Certinator) CreateCerts(caName string, requests []CertificateRequest) (certs CertificateBundle, err error) {

	certs = make(map[string]CertInfo)

	for _, r := range requests {
		data := map[string]interface{}{
			"common_name": r.CommonName,
			"ttl":         r.Ttl,
		}

		if r.IpSans != "" {
			data["ip_sans"] = r.IpSans
		}

		if r.Sans != "" {
			data["alt_names"] = r.Sans
		}

		path := fmt.Sprintf("%s/issue/%s", caName, DEFAULT_CERTIFICATE_ROLE)

		secret, err := c.Client.Logical().Write(path, data)
		if err != nil {
			err = errors.Wrapf(err, "failed creating certificate for %s on service CA", r.CommonName)
			return certs, err
		}

		if secret != nil {
			if secret.Data != nil {
				exp, err := secret.Data["expiration"].(json.Number).Int64()
				if err != nil {
					err = errors.Wrapf(err, "failed converting %v to int64", secret.Data["expiration"])
				}
				certs[r.CommonName] = CertInfo{
					IssuingCA:      secret.Data["issuing_ca"].(string),
					PrivateKey:     secret.Data["private_key"].(string),
					PrivateKeyType: secret.Data["private_key_type"].(string),
					SerialNumber:   secret.Data["serial_number"].(string),
					Certificate:    secret.Data["serial_number"].(string),
					Expiration:     exp,
				}
			}
		}
	}

	return certs, err
}

func (c *Certinator) ListCerts() (err error) {
	return err
}

func (c *Certinator) RevokeCert(cn string, ca string) (err error) {
	return err
}
