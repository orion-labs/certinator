# certinator

Instrument Hashicorp Vault to create CA's and Certificates for Clients and Services

Connect to Vault via `VAULT_ADDR`, using token in `~/.vault-token`.

# Commands

For any command, run `certinator help <command>` for usage instructions.

* ca list - Lists CA's (Certificate Authorities) in Vault
* ca create <name> - Creates a new CA
* ca delete <name> - Deletes a CA

* cert create - Creates a single certificate, or multiple certificates via config file.
* cert list - Lists certs the CA has created
* cert revoke - Revokes a certificate by it's CommonName

* crl fetch - Fetches the current CRL (Certificate Revocation List) for a CA
* crl rotate - Rotates the CRL (CRL's in Vault are short lived documents by default)

# Usage

Certinator is just bouncing commands off of Vault.  If you can run the commands in your Vault, it will work.  If you cannot, then you're out of luck.

Like most Vault tools, it depends on `VAULT_ADDR` and your token stored in `~/.vault-token`.

# Using Certinator with a Dev Mode Vault

Vault is pretty cool in that the vault binary distributed by Hashicorp is both a client and a server all in one.

The server even has a ‘dev mode’ flag that you can use to demo vault’s features without enabling the really high security features that to be honest, require significant understanding to work under.

Dev mode is great for one offs, demos, and tests.  Here’s how you use it.


## Get Vault

Download the appropriate version of vault from here:

https://www.vaultproject.io/downloads

Install it, and make sure you can run vault from the command line.

If you’re on a Mac, you can run brew install vault.


## Start Vault Server in Dev Mode

From the command line, run:

    vault server -dev

Much text will scroll by.   Look for something like this:

WARNING! dev mode is enabled! In this mode, Vault runs entirely in-memory
and starts unsealed with a single unseal key. The root token is already
authenticated to the CLI, so you can immediately begin using Vault.

You may need to set the following environment variable:

    $ export VAULT_ADDR='http://127.0.0.1:8200'

    The unseal key and root token are displayed below in case you want to
    seal/unseal the Vault or re-authenticate.

    Unseal Key: t26Rv5VAEUGpP1tXkcANuVKELeozPfrw3ztQKGiUiCg=
    jRoot Token: s.e4vegH8GHtNOhe1ooVcqzYBx

    Development mode should NOT be used in production installations!
    
You will need to export VAULT_ADDR='http://127.0.0.1:8200' into your environnment.

You will also need to take the value for Root Token `s.e4vegH8GHtNOhe1ooVcqzYBx` in the above example, and write it into a file named `~/.vault-token`.  

That file must contain a single line, with only the value supplied by the vault server command.  (It’ll be different every time you run vault server -dev.


## Test Vault Connection to Server
If the above is correct, you should be able to run the following commands with success:

#### Command: vault status

#### Output: 

    Key             Value
    ---             -----
    Seal Type       shamir
    Initialized     true
    Sealed          false
    Total Shares    1
    Threshold       1
    Version         1.4.0
    Cluster Name    vault-cluster-711027a9
    Cluster ID      39026b8b-15f7-2711-1cb3-6bfc7cc841ed
    HA Enabled      false

#### Command: vault secrets list

#### Output:

    Path          Type         Accessor              Description
    ----          ----         --------              -----------
    cubbyhole/    cubbyhole    cubbyhole_22dd5360    per-token private secret storage
    identity/     identity     identity_2bb2bcbe     identity store
    secret/       kv           kv_a1eb8ba8           key/value secret storage
    sys/          system       system_28b36865       system endpoints used for control, policy and debugging

## Create a CA

#### Command: certinator ca create -n onprem -e -d example.com,*.example.com onprem

#### Output:

    CA onprem created.
      CA Certificate written to onprem-ca.crt
      CA Private Key writen to onprem-ca.key

#### Additional Outputs:
You should see the following files in the directory in which you ran the command:

    onprem-ca.crt

    onprem-ca.key


## Create a Wildcard Certificate 

#### Command: certinator cert create -c onprem -n *.example.net -t 8760h

#### Output: 

    You are currently using the root token.  You should not be doing this unless it's really necessary.
    
(That’s just a warning.  You can ignore it when you’re working in dev mode)


#### Additional Outputs:

You should see the following files in the directory in which you ran the command:

    *.example.com.crt

    *.example.com.key