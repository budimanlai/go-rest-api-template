package middleware

import (
	"context"
	"go-rest-api-template/internal/service"
	"go-rest-api-template/pkg/response"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// PublicMiddleware validates API keys for public endpoints and optionally validates public JWT tokens
func PublicMiddleware(apiKeyService service.ApiKeyService, jwtService service.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Step 1: Get and validate API key from X-API-Key header
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			// Try from Authorization header as API key (not Bearer token)
			auth := c.Get("Authorization")
			if strings.HasPrefix(auth, "ApiKey ") {
				apiKey = strings.TrimPrefix(auth, "ApiKey ")
			}
		}

		if apiKey == "" {
			return response.BadRequest(c, "API key is required for public endpoints", "")
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

		// Store API key info in context
		c.Locals("api_key_id", apiKeyEntity.ID)
		c.Locals("api_key_name", apiKeyEntity.Name)
		c.Locals("api_key_h2h", apiKeyEntity.IsH2HEnabled())

		// Step 2: Check for public JWT token (optional)
		auth := c.Get("Authorization")
		if auth != "" && strings.HasPrefix(auth, "Bearer ") {
			tokenString := strings.TrimPrefix(auth, "Bearer ")
			if tokenString != "" {
				// Validate public JWT token (contains API key info)
				claims, tokenApiKey, err := jwtService.ValidatePublicToken(tokenString)
				if err == nil && tokenApiKey != nil {
					// Verify token API key matches request API key
					if tokenApiKey.ID == apiKeyEntity.ID {
						// Store additional token info if needed
						c.Locals("token_valid", true)
						c.Locals("token_api_key_id", claims.ApiKeyID)
						c.Locals("token_api_key_name", claims.ApiKeyName)
					} else {
						// Token API key doesn't match request API key
						c.Locals("token_valid", false)
					}
				} else {
					// Token exists but invalid
					c.Locals("token_valid", false)
				}
			}
		} else {
			// No token provided
			c.Locals("token_valid", false)
		}

		// Log API key access (async with timeout)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = apiKeyService.LogApiKeyAccess(ctx, apiKeyEntity.ID)
		}()

		return c.Next()
	}
}

// PrivateMiddleware validates both API keys and private JWT tokens for private endpoints
func PrivateMiddleware(apiKeyService service.ApiKeyService, jwtService service.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Step 1: Validate API Key first
		apiKey := c.Get("X-API-Key")
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

		// Step 2: Validate Private JWT Token
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return response.BadRequest(c, "Bearer token is required for private endpoints", "")
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == "" {
			return response.BadRequest(c, "Invalid token format", "")
		}

		// Validate private JWT token (contains API key + user info)
		claims, tokenApiKey, user, err := jwtService.ValidatePrivateToken(tokenString)
		if err != nil {
			return response.BadRequest(c, "Invalid or expired token", err.Error())
		}

		// Verify token API key matches request API key
		if tokenApiKey.ID != apiKeyEntity.ID {
			return response.BadRequest(c, "Token API key doesn't match request API key", "")
		}

		// Store API key info in context
		c.Locals("api_key_id", apiKeyEntity.ID)
		c.Locals("api_key_name", apiKeyEntity.Name)
		c.Locals("api_key_h2h", apiKeyEntity.IsH2HEnabled())

		// Store user info from JWT token
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("user_email", claims.Email)
		c.Locals("user", user)

		// Log API key access (async with timeout)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = apiKeyService.LogApiKeyAccess(ctx, apiKeyEntity.ID)
		}()

		return c.Next()
	}
}

// OptionalPrivateJWTMiddleware validates API key and optionally validates private JWT token
// Useful for endpoints that can work with or without user authentication
func OptionalPrivateJWTMiddleware(apiKeyService service.ApiKeyService, jwtService service.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Step 1: Validate API Key (required)
		apiKey := c.Get("X-API-Key")
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

		// Store API key info in context
		c.Locals("api_key_id", apiKeyEntity.ID)
		c.Locals("api_key_name", apiKeyEntity.Name)
		c.Locals("api_key_h2h", apiKeyEntity.IsH2HEnabled())

		// Step 2: Try to validate private JWT Token (optional)
		auth := c.Get("Authorization")
		if auth != "" && strings.HasPrefix(auth, "Bearer ") {
			tokenString := strings.TrimPrefix(auth, "Bearer ")
			if tokenString != "" {
				// Validate private JWT token
				claims, tokenApiKey, user, err := jwtService.ValidatePrivateToken(tokenString)
				if err == nil && tokenApiKey != nil {
					// Verify token API key matches request API key
					if tokenApiKey.ID == apiKeyEntity.ID {
						// Store user info if token is valid
						c.Locals("user_id", claims.UserID)
						c.Locals("username", claims.Username)
						c.Locals("user_email", claims.Email)
						c.Locals("user", user)
						c.Locals("authenticated", true)
					} else {
						// Token API key doesn't match
						c.Locals("authenticated", false)
					}
				} else {
					// Token exists but invalid
					c.Locals("authenticated", false)
				}
			}
		} else {
			// No token provided
			c.Locals("authenticated", false)
		}

		// Log API key access (async with timeout)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = apiKeyService.LogApiKeyAccess(ctx, apiKeyEntity.ID)
		}()

		return c.Next()
	}
}
