module github.com/RickyGrassmuck/hello-vault-go

go 1.16

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/hashicorp/vault/api v1.3.0
	github.com/hashicorp/vault/api/auth/ldap v0.1.0
	github.com/jessevdk/go-flags v1.5.0
)

replace github.com/hashicorp/vault/api => /workspaces/vault/api
replace github.com/hashicorp/vault/api/auth/ldap => /workspaces/vault/api/auth/ldap
