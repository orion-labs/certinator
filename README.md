# certinator

Instrument Hashicorp Vault to create CA's and Certificates for Clients and Services

Connect to Vault via VAULT_ADDR

# Commands

For any command, run `certinator help <command>` for usage instructions.

* vault init (not yet implemented)
* vault unseal (not yet implemented)
* vault status - returns status information similar to `vault status`

* ca list - Lists CA's (Certificate Authorities) in Vault
* ca create <name> - Creates a new CA
* ca delete <name> - Deletes a CA

* cert create - Creates a single certificate, or multiple certificates via config file.
* cert list - Lists certs the CA has created
* cert revoke - Revokes a certificate by it's CommonName

* crl fetch - Fetches the current CRL (Certificate Revocation List) for a CA
* crl rotate - Rotates the CRL (CRL's in Vault are short lived documents by default)
