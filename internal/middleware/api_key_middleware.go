package middleware

import (
	"context"
	"go-rest-api-template/internal/service"
	"go-rest-api-template/pkg/response"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Global I18n response helper for middleware usage
var GlobalI18nResponseHelper *response.I18nResponseHelper

// SetGlobalI18nResponseHelper sets the global response helper
func SetGlobalI18nResponseHelper(helper *response.I18nResponseHelper) {
	GlobalI18nResponseHelper = helper
}

// ApiKeyOnlyMiddleware creates a simple API key validation middleware
func ApiKeyOnlyMiddleware(apiKeyService service.ApiKeyService) fiber.Handler {
	return ApiKeyMiddleware(apiKeyService, GlobalI18nResponseHelper)
}

// ApiKeyMiddleware validates API keys for JWT middleware integration
func ApiKeyMiddleware(apiKeyService service.ApiKeyService, responseHelper *response.I18nResponseHelper) fiber.Handler {
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
			if responseHelper != nil {
				return responseHelper.ErrorWithI18n(c, fiber.StatusBadRequest, "api_key_required", nil)
			}
			return response.BadRequest(c, "API key is required", "")
		}

		// Validate API key
		ctx := context.Background()
		apiKeyEntity, err := apiKeyService.ValidateApiKey(ctx, apiKey)
		if err != nil {
			if responseHelper != nil {
				return responseHelper.ErrorWithI18n(c, fiber.StatusInternalServerError, "internal_server", nil)
			}
			return response.InternalServerError(c, "Internal server error", err.Error())
		}

		if apiKeyEntity == nil {
			if responseHelper != nil {
				return responseHelper.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_api_key", nil)
			}
			return response.BadRequest(c, "Invalid or inactive API key", "")
		}

		// Check IP whitelist if configured
		clientIP := c.IP()
		if !apiKeyEntity.IsIPWhitelisted(clientIP) {
			if responseHelper != nil {
				return responseHelper.ErrorWithI18n(c, fiber.StatusBadRequest, "ip_not_whitelisted", nil)
			}
			return response.BadRequest(c, "IP address not whitelisted", "")
		}

		// Store API key info in context for later use
		c.Locals("api_key_id", apiKeyEntity.ID)
		c.Locals("api_key_name", apiKeyEntity.Name)
		c.Locals("api_key_h2h", apiKeyEntity.IsH2HEnabled())

		// Log API key access (async with timeout)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = apiKeyService.LogApiKeyAccess(ctx, apiKeyEntity.ID)
		}()

		return c.Next()
	}
}

// AuthKeyMiddleware validates auth keys (alternative authentication method)
func AuthKeyMiddleware(apiKeyService service.ApiKeyService, responseHelper *response.I18nResponseHelper) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get auth key from header
		authKey := c.Get("X-Auth-Key")
		if authKey == "" {
			if responseHelper != nil {
				return responseHelper.ErrorWithI18n(c, fiber.StatusBadRequest, "auth_key_required", nil)
			}
			return response.BadRequest(c, "Auth key is required", "")
		}

		// Validate auth key
		ctx := context.Background()
		apiKeyEntity, err := apiKeyService.ValidateAuthKey(ctx, authKey)
		if err != nil {
			if responseHelper != nil {
				return responseHelper.ErrorWithI18n(c, fiber.StatusInternalServerError, "internal_server", nil)
			}
			return response.InternalServerError(c, "Internal server error", err.Error())
		}

		if apiKeyEntity == nil {
			if responseHelper != nil {
				return responseHelper.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_auth_key", nil)
			}
			return response.BadRequest(c, "Invalid or inactive auth key", "")
		}

		// Check IP whitelist if configured
		clientIP := c.IP()
		if !apiKeyEntity.IsIPWhitelisted(clientIP) {
			if responseHelper != nil {
				return responseHelper.ErrorWithI18n(c, fiber.StatusBadRequest, "ip_not_whitelisted", nil)
			}
			return response.BadRequest(c, "IP address not whitelisted", "")
		}

		// Store API key info in context
		c.Locals("api_key_id", apiKeyEntity.ID)
		c.Locals("api_key_name", apiKeyEntity.Name)
		c.Locals("api_key_h2h", apiKeyEntity.IsH2HEnabled())

		// Log access (async with timeout)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = apiKeyService.LogApiKeyAccess(ctx, apiKeyEntity.ID)
		}()

		return c.Next()
	}
}
