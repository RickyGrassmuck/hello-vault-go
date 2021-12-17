package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
)

type Environment struct {
	// The address of this service
	MyAddress string `               env:"MY_ADDRESS"                    default:":8080"                        description:"Listen to http traffic on this tcp address"             long:"my-address"`

	// Vault address, approle login credentials, and secret locations
	VaultAddress      string `env:"VAULT_ADDRESS"                 default:"localhost:8200"               description:"Vault address"                                          long:"vault-address"`
	VaultLDAPUsername string `env:"VAULT_LDAP_USERNAME"           required:"true"                        description:"AppRole RoleID to log in to Vault"                      long:"vault-ldap-username"`
	VaultLDAPassword  string `env:"VAULT_LDAP_PASSWORD"           default:"kv-v2/data/api-key"           description:"Path to the API key used by 'secure-sevice'"            long:"vault-ldap-password"`
	VaultAPIKeyPath   string `env:"VAULT_API_KEY_PATH"            default:"kv-v2/data/api-key"           description:"Path to the API key used by 'secure-sevice'"            long:"vault-api-key-path"`
}

func main() {
	log.Println("hello!")
	defer log.Println("goodbye!")

	var env Environment

	// parse & validate environment variables
	_, err := flags.Parse(&env)
	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		}
		log.Fatalf("unable to parse environment variables: %v", err)
	}

	if err := run(context.Background(), env); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run(ctx context.Context, env Environment) error {
	// WARNING: the goroutines in this function have simplified error handling
	// and could escape the scope of the function. Production applications
	// may want to add more complex error handling and leak protection logic.

	ctx, cancelContextFunc := context.WithCancel(ctx)
	defer cancelContextFunc()

	// vault
	vault, token, err := NewVaultLDAPClient(
		ctx,
		VaultParameters{
			env.VaultAddress,
			env.VaultLDAPUsername,
			env.VaultLDAPassword,
			env.VaultAPIKeyPath,
		},
	)
	if err != nil {
		return fmt.Errorf("unable to initialize vault connection @ %s: %w", env.VaultAddress, err)
	}
	go vault.RenewLoginPeriodically(ctx, token) // keep alive

	// handlers & routes
	h := Handlers{
		vault: vault,
	}

	r := gin.New()
	r.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/healthcheck"), // don't log healthcheck requests
	)

	// healthcheck
	r.GET("/healthcheck", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// demonstrates fetching a static secret from vault and using it to talk to another service
	r.POST("/payments", h.CreatePayment)

	r.Run(env.MyAddress)

	return nil
}
