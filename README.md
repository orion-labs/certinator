# certinator

Instrument Hashicorp Vault to create CA's and Certificates for Clients and Services

Connect to Vault via VAULT_ADDR

Init:

    init -> read json
    unseal -> read from json

    write root token to disk
    write unseal keys to disk

    create CA's for 
        * Services
        * Clients
    
Create Service Certs:

    Read service list from config file or CLI.

    Create Service Certs -> Write to:
        * K8S secrets
        * PEM Files
        
Create Client Certs:
    
    Read client list from config file or CLI.

    Create Client Certs -> Write to:
        * K8S Secerts
        * PEM FIles

# Commands
* vault init
* vault unseal
* vault status

* ca list
* ca create <name>
* ca delete <name>

* cert create
* cert list
* cert revoke

* crl fetch
* crl rotate
