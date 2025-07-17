package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// ContextHelper provides utility functions to extract data from Fiber context
type ContextHelper struct{}

// NewContextHelper creates a new context helper
func NewContextHelper() *ContextHelper {
	return &ContextHelper{}
}

// GetAPIKeyID extracts API key ID from context
func (h *ContextHelper) GetAPIKeyID(c *fiber.Ctx) (int, error) {
	if id, ok := c.Locals("api_key_id").(int); ok {
		return id, nil
	}
	return 0, errors.New("api_key_id not found in context")
}

// GetAPIKeyName extracts API key name from context
func (h *ContextHelper) GetAPIKeyName(c *fiber.Ctx) (string, error) {
	if name, ok := c.Locals("api_key_name").(string); ok {
		return name, nil
	}
	return "", errors.New("api_key_name not found in context")
}

// IsH2HEnabled checks if H2H is enabled from context
func (h *ContextHelper) IsH2HEnabled(c *fiber.Ctx) bool {
	if h2h, ok := c.Locals("api_key_h2h").(bool); ok {
		return h2h
	}
	return false
}

// GetUserID extracts user ID from context (from JWT)
func (h *ContextHelper) GetUserID(c *fiber.Ctx) (int, error) {
	if id, ok := c.Locals("user_id").(int); ok {
		return id, nil
	}
	return 0, errors.New("user_id not found in context")
}

// GetUsername extracts username from context (from JWT)
func (h *ContextHelper) GetUsername(c *fiber.Ctx) (string, error) {
	if username, ok := c.Locals("username").(string); ok {
		return username, nil
	}
	return "", errors.New("username not found in context")
}

// GetUserEmail extracts user email from context (from JWT)
func (h *ContextHelper) GetUserEmail(c *fiber.Ctx) (string, error) {
	if email, ok := c.Locals("user_email").(string); ok {
		return email, nil
	}
	return "", errors.New("user_email not found in context")
}

// IsAuthenticated checks if user is authenticated (has valid JWT token)
func (h *ContextHelper) IsAuthenticated(c *fiber.Ctx) bool {
	if auth, ok := c.Locals("authenticated").(bool); ok {
		return auth
	}
	// If not explicitly set, check if user_id exists
	_, err := h.GetUserID(c)
	return err == nil
}

// MustGetUserID extracts user ID from context or panics (use in private endpoints only)
func (h *ContextHelper) MustGetUserID(c *fiber.Ctx) int {
	id, err := h.GetUserID(c)
	if err != nil {
		panic("user_id not found in context - middleware not properly configured")
	}
	return id
}

// MustGetAPIKeyID extracts API key ID from context or panics
func (h *ContextHelper) MustGetAPIKeyID(c *fiber.Ctx) int {
	id, err := h.GetAPIKeyID(c)
	if err != nil {
		panic("api_key_id not found in context - middleware not properly configured")
	}
	return id
}
