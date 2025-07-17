package handler

import (
	"go-rest-api-template/internal/domain/usecase"
	"go-rest-api-template/internal/model"
	"go-rest-api-template/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	userService    usecase.UserUsecase
	responseHelper *response.I18nResponseHelper
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userService usecase.UserUsecase, responseHelper *response.I18nResponseHelper) *AuthHandler {
	return &AuthHandler{
		userService:    userService,
		responseHelper: responseHelper,
	}
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
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

// Login handles user authentication
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return h.responseHelper.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_request", nil)
	}

	// Validate request
	loginReq := &model.UserLoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
	if err := loginReq.Validate(); err != nil {
		return h.responseHelper.ErrorWithI18n(c, fiber.StatusBadRequest, "validation_failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Authenticate user
	user, token, err := h.userService.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return h.responseHelper.ErrorWithI18n(c, fiber.StatusUnauthorized, "login_failed", map[string]interface{}{
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

	response := LoginResponse{
		User:  userResponse,
		Token: token,
	}

	return h.responseHelper.SuccessWithI18n(c, "login_success", response, nil)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return h.responseHelper.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_request", nil)
	}

	// Refresh token
	newToken, err := h.userService.RefreshToken(c.Context(), req.Token)
	if err != nil {
		return h.responseHelper.ErrorWithI18n(c, fiber.StatusUnauthorized, "token_refresh_failed", map[string]interface{}{
			"error": err.Error(),
		})
	}

	response := RefreshTokenResponse{
		Token: newToken,
	}

	return h.responseHelper.SuccessWithI18n(c, "token_refresh_success", response, nil)
}

// Logout handles user logout (optional - for token blacklist in the future)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// For JWT, logout is typically handled client-side by removing the token
	// In the future, you could implement token blacklisting here
	return h.responseHelper.SuccessWithI18n(c, "logout_success", nil, nil)
}
