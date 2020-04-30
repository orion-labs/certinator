package certinator

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
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
					Certificate:    secret.Data["certificate"].(string),
					Expiration:     exp,
				}
			}
		}
	}

	return certs, err
}

func (c *Certinator) ListCerts(caName string) (certs []string, err error) {
	// Have to query vault for serial numbers (vault won't give out CN's)
	path := fmt.Sprintf("%s/certs", caName)

	serials := make([]string, 0)

	secret, err := c.Client.Logical().List(path)
	if err != nil {
		err = errors.Wrapf(err, "failed listing certs on %s", caName)
		return certs, err
	}

	if secret != nil {
		if secret.Data != nil {
			s, ok := secret.Data["keys"].([]interface{})
			if ok {
				for _, n := range s {
					sn, ok := n.(string)
					if ok {
						serials = append(serials, sn)
					}
				}
			}
		}
	}

	certs = make([]string, 0)

	// Then we pull the cert from Vault by SN and read it to find it's CN.
	for _, sn := range serials {
		path := fmt.Sprintf("%s/cert/%s", caName, sn)
		secret, err := c.Client.Logical().Read(path)
		if err != nil {
			err = errors.Wrapf(err, "failed listing certs on %s", caName)
			return certs, err
		}

		if secret != nil {
			if secret.Data != nil {
				// pull out the data
				p, ok := secret.Data["certificate"].(string)
				if ok {
					// pull out the data from the PEM
					block, _ := pem.Decode([]byte(p))
					if block == nil {
						err = errors.New(fmt.Sprintf("failed decoding PEM for certificate serial %s", sn))
						return certs, err
					}
					// Parse the data into a certificate
					cert, err := x509.ParseCertificate(block.Bytes)
					if err != nil {
						err = errors.Wrapf(err, "failed parsing certificate for SN %s", sn)
						return certs, err
					}

					// pull out the CN on the cert, add it to the list.
					certs = append(certs, cert.Subject.CommonName)
				}
			}
		}
	}

	return certs, err
}

func (c *Certinator) RevokeCert(cn string, ca string) (err error) {
	// revoke can only be done via serial number, which we rarely know.
	// get all the serial numbers
	path := fmt.Sprintf("%s/certs", ca)
	serials := make([]string, 0)

	secret, err := c.Client.Logical().List(path)
	if err != nil {
		err = errors.Wrapf(err, "failed listing certs on %s", ca)
		return err
	}

	if secret != nil {
		if secret.Data != nil {
			s, ok := secret.Data["keys"].([]interface{})
			if ok {
				for _, n := range s {
					sn, ok := n.(string)
					if ok {
						serials = append(serials, sn)
					}
				}
			}
		}
	}

	// Then we pull the cert from Vault by SN and read it to find it's CN.
	for _, sn := range serials {
		path := fmt.Sprintf("%s/cert/%s", ca, sn)
		secret, err := c.Client.Logical().Read(path)
		if err != nil {
			err = errors.Wrapf(err, "failed reading cert %s on %s", sn, ca)
			return err
		}

		if secret != nil {
			if secret.Data != nil {
				// pull out the data
				p, ok := secret.Data["certificate"].(string)
				if ok {
					// pull out the data from the PEM
					block, _ := pem.Decode([]byte(p))
					if block == nil {
						err = errors.New(fmt.Sprintf("failed decoding PEM for certificate serial %s", sn))
						return err
					}
					// Parse the data into a certificate
					cert, err := x509.ParseCertificate(block.Bytes)
					if err != nil {
						err = errors.Wrapf(err, "failed parsing certificate for SN %s", sn)
						return err
					}

					// if the CN matches the one we're looking for
					if cert.Subject.CommonName == cn {
						path := fmt.Sprintf("%s/revoke", ca)
						data := map[string]interface{}{
							"serial_number": sn,
						}

						// tell Vault to revoke it
						_, err := c.Client.Logical().Write(path, data)
						if err != nil {
							err = errors.Wrapf(err, "failed revoking cert with CN %s Serial Number %s", cn, sn)
							return err
						}

						return err
					}
				}
			}
		}
	}

	return err
}
