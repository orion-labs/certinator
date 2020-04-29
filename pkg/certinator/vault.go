package certinator

import "github.com/pkg/errors"

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
	return err
}

func (c *Certinator) UnsealVault() (err error) {
	// TODO implement and test UnsealVault
	return err
}
