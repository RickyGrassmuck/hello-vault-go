package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	vault                *Vault
	secureServiceAddress string
}

// (POST /payments) : demonstrates fetching a static secret from Vault and using it to talk to another service
func (h *Handlers) CreatePayment(c *gin.Context) {
	// retrieve the secret from Vault
	secret, err := h.vault.GetSecretAPIKey(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// check that our expected key is in the returned secret
	apiKey, ok := secret["api_key"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "the secret retrieved from vault is missing 'api_key' field"})
		return
	}

	c.JSON(200, apiKey)
}
