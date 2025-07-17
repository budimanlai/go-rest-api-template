package middleware

import (
	"context"
	"go-rest-api-template/internal/service"
	"go-rest-api-template/pkg/response"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ApiKeyMiddleware validates API keys for JWT middleware integration
func ApiKeyMiddleware(apiKeyService service.ApiKeyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get API key from header
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			// Try from Authorization header as Bearer token
			auth := c.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				apiKey = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if apiKey == "" {
			return response.BadRequest(c, "API key is required", "")
		}

		// Validate API key
		ctx := context.Background()
		apiKeyEntity, err := apiKeyService.ValidateApiKey(ctx, apiKey)
		if err != nil {
			return response.InternalServerError(c, "Internal server error", err.Error())
		}

		if apiKeyEntity == nil {
			return response.BadRequest(c, "Invalid or inactive API key", "")
		}

		// Check IP whitelist if configured
		clientIP := c.IP()
		if !apiKeyEntity.IsIPWhitelisted(clientIP) {
			return response.BadRequest(c, "IP address not whitelisted", "")
		}

		// Store API key info in context for later use
		c.Locals("api_key_id", apiKeyEntity.ID)
		c.Locals("api_key_name", apiKeyEntity.Name)
		c.Locals("api_key_h2h", apiKeyEntity.IsH2HEnabled())

		// Log API key access (async)
		go func() {
			_ = apiKeyService.LogApiKeyAccess(context.Background(), apiKeyEntity.ID)
		}()

		return c.Next()
	}
}

// AuthKeyMiddleware validates auth keys (alternative authentication method)
func AuthKeyMiddleware(apiKeyService service.ApiKeyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get auth key from header
		authKey := c.Get("X-Auth-Key")
		if authKey == "" {
			return response.BadRequest(c, "Auth key is required", "")
		}

		// Validate auth key
		ctx := context.Background()
		apiKeyEntity, err := apiKeyService.ValidateAuthKey(ctx, authKey)
		if err != nil {
			return response.InternalServerError(c, "Internal server error", err.Error())
		}

		if apiKeyEntity == nil {
			return response.BadRequest(c, "Invalid or inactive auth key", "")
		}

		// Check IP whitelist if configured
		clientIP := c.IP()
		if !apiKeyEntity.IsIPWhitelisted(clientIP) {
			return response.BadRequest(c, "IP address not whitelisted", "")
		}

		// Store API key info in context
		c.Locals("api_key_id", apiKeyEntity.ID)
		c.Locals("api_key_name", apiKeyEntity.Name)
		c.Locals("api_key_h2h", apiKeyEntity.IsH2HEnabled())

		// Log access
		go func() {
			_ = apiKeyService.LogApiKeyAccess(context.Background(), apiKeyEntity.ID)
		}()

		return c.Next()
	}
}
