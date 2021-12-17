package main

import (
	"context"
	"fmt"
	"log"

	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/ldap"
)

type VaultParameters struct {
	// connection parameters
	address      string
	ldapUsername string
	ldapPassword string

	// the locations of our two secrets
	apiKeyPath string
}

type Vault struct {
	client     *vault.Client
	parameters VaultParameters
}

// NewVaultLDAPClient logs in to Vault using the LDAP authentication
// method, returning an authenticated client and the auth token itself, which
// can be periodically renewed.
func NewVaultLDAPClient(ctx context.Context, parameters VaultParameters) (*Vault, *vault.Secret, error) {
	log.Printf("connecting to vault @ %s", parameters.address)

	config := vault.DefaultConfig() // modify for more granular configuration
	config.Address = parameters.address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize vault client: %w", err)
	}

	vault := &Vault{
		client:     client,
		parameters: parameters,
	}

	token, err := vault.login(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("vault login error: %w", err)
	}

	log.Println("connecting to vault: success!")

	return vault, token, nil
}

// A combination of a RoleID and a SecretID is required to log into Vault
// with AppRole authentication method. The SecretID is a value that needs
// to be protected, so instead of the app having knowledge of the SecretID
// directly, we have a trusted orchestrator (simulated with a script here)
// give the app access to a short-lived response-wrapping token.
//
// ref: https://www.vaultproject.io/docs/concepts/response-wrapping
// ref: https://learn.hashicorp.com/tutorials/vault/secure-introduction?in=vault/app-integration#trusted-orchestrator
// ref: https://learn.hashicorp.com/tutorials/vault/approle-best-practices?in=vault/auth-methods#secretid-delivery-best-practices
func (v *Vault) login(ctx context.Context) (*vault.Secret, error) {
	log.Printf("logging in to vault with ldap auth; username id: %s", v.parameters.ldapUsername)

	ldapPassword := &ldap.Password{
		FromString: v.parameters.ldapPassword,
	}

	ldapAuth, err := ldap.NewLDAPAuth(
		v.parameters.ldapUsername,
		ldapPassword,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize LDAP authentication method: %w", err)
	}

	authInfo, err := v.client.Auth().Login(ctx, ldapAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login using LDAP auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no LDAP info was returned after login")
	}

	log.Println("logging in to vault with LDAP auth: success!")

	return authInfo, nil
}

// GetSecretAPIKey fetches the latest version of secret api key from kv-v2
func (v *Vault) GetSecretAPIKey(ctx context.Context) (map[string]interface{}, error) {
	log.Println("getting secret api key from vault")

	secret, err := v.client.Logical().Read(v.parameters.apiKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read secret: %w", err)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("malformed secret returned: %v", data)
	}

	log.Println("getting secret api key from vault: success!")

	return data, nil
}
