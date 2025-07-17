package handler

import (
	"go-rest-api-template/internal/domain/entity"
	"go-rest-api-template/internal/domain/usecase"
	"go-rest-api-template/internal/model"
	"go-rest-api-template/internal/service"
	"go-rest-api-template/pkg/response"
	"go-rest-api-template/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	userService   usecase.UserUsecase
	jwtService    service.JWTService
	apiKeyService service.ApiKeyService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService usecase.UserUsecase, jwtService service.JWTService, apiKeyService service.ApiKeyService) *AuthHandler {
	return &AuthHandler{
		userService:   userService,
		jwtService:    jwtService,
		apiKeyService: apiKeyService,
	}
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// RegisterRequest represents the register request payload
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	FullName string `json:"full_name" validate:"required,min=2,max=100"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User  *model.UserResponse `json:"user"`
	Token string              `json:"token"`
}

// RefreshTokenRequest represents the refresh token request
type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

// RefreshTokenResponse represents the refresh token response
type RefreshTokenResponse struct {
	Token string `json:"token"`
}

// PublicTokenResponse represents the public token response
type PublicTokenResponse struct {
	PublicToken string `json:"public_token"`
	ExpiresIn   int    `json:"expires_in"` // seconds
	TokenType   string `json:"token_type"`
}

// GetPublicToken generates a public token for accessing public endpoints
func (h *AuthHandler) GetPublicToken(c *fiber.Ctx) error {
	// Get API key from context (should be set by middleware)
	apiKeyID, ok := c.Locals("api_key_id").(int)
	if !ok {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "api_key_required", nil)
	}

	apiKeyName, ok := c.Locals("api_key_name").(string)
	if !ok {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "api_key_invalid", nil)
	}

	// Create API key entity for token generation
	apiKey := &entity.ApiKey{
		ID:   apiKeyID,
		Name: apiKeyName,
	}

	// Generate public token
	publicToken, err := h.jwtService.GeneratePublicToken(apiKey)
	if err != nil {
		return response.ErrorWithI18n(c, fiber.StatusInternalServerError, "token_generation_failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tokenResponse := PublicTokenResponse{
		PublicToken: publicToken,
		ExpiresIn:   2 * 3600, // 2 hours in seconds
		TokenType:   "Bearer",
	}

	return response.SuccessWithI18n(c, "public_token_generated", tokenResponse, nil)
}

// Login handles user authentication and returns private token
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_request", nil)
	}

	// Validate request
	loginReq := &model.UserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	if err := loginReq.Validate(); err != nil {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "validation_failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Get API key from context (should be set by middleware)
	apiKeyID, ok := c.Locals("api_key_id").(int)
	if !ok {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "api_key_required", nil)
	}

	apiKeyName, ok := c.Locals("api_key_name").(string)
	if !ok {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "api_key_invalid", nil)
	}

	// Authenticate user (note: user service now returns empty token)
	user, _, err := h.userService.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return response.ErrorWithI18n(c, fiber.StatusUnauthorized, "login_failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Create entities for token generation
	apiKey := &entity.ApiKey{
		ID:   apiKeyID,
		Name: apiKeyName,
	}

	// Generate private token (contains API key + user info)
	privateToken, err := h.jwtService.GeneratePrivateToken(apiKey, user)
	if err != nil {
		return response.ErrorWithI18n(c, fiber.StatusInternalServerError, "token_generation_failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Convert to response format
	userResponse := &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	loginResponse := LoginResponse{
		User:  userResponse,
		Token: privateToken, // This is now a private token
	}

	return response.SuccessWithI18n(c, "login_success", loginResponse, nil)
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_request_body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Validate request
	if err := validator.ValidateStruct(&req); err != nil {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "validation_failed", map[string]interface{}{
			"errors": err,
		})
	}

	// Get API key info from context for user creation
	apiKeyID, ok := c.Locals("api_key_id").(int)
	if !ok {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "api_key_required", nil)
	}

	apiKeyName, ok := c.Locals("api_key_name").(string)
	if !ok {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "api_key_invalid", nil)
	}

	// Create user entity
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Status:   "active",
	}

	// Set password (will be hashed)
	if err := user.HashPassword(req.Password); err != nil {
		return response.ErrorWithI18n(c, fiber.StatusInternalServerError, "password_hash_failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Create user through usecase
	if err := h.userService.CreateUser(c.Context(), user); err != nil {
		return response.ErrorWithI18n(c, fiber.StatusInternalServerError, "registration_failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Create API key entity for token generation
	apiKey := &entity.ApiKey{
		ID:   apiKeyID,
		Name: apiKeyName,
	}

	// Generate private token for new user
	privateToken, err := h.jwtService.GeneratePrivateToken(apiKey, user)
	if err != nil {
		return response.ErrorWithI18n(c, fiber.StatusInternalServerError, "token_generation_failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Convert to response format
	userResponse := &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	loginResponse := LoginResponse{
		User:  userResponse,
		Token: privateToken,
	}

	return response.SuccessWithI18n(c, "registration_success", loginResponse, nil)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_request", nil)
	}

	// Refresh token
	newToken, err := h.userService.RefreshToken(c.Context(), req.Token)
	if err != nil {
		return response.ErrorWithI18n(c, fiber.StatusUnauthorized, "token_refresh_failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	refreshResponse := RefreshTokenResponse{
		Token: newToken,
	}

	return response.SuccessWithI18n(c, "token_refresh_success", refreshResponse, nil)
}

// Logout handles user logout (optional - for token blacklist in the future)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// For JWT, logout is typically handled client-side by removing the token
	// In the future, you could implement token blacklisting here
	return response.SuccessWithI18n(c, "logout_success", nil, nil)
}
