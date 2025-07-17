package handler

import (
	"go-rest-api-template/internal/domain/entity"
	"go-rest-api-template/internal/domain/usecase"
	"go-rest-api-template/internal/model"
	"go-rest-api-template/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req model.UserCreateRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_request", nil)
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "validation_failed", nil)
	}

	// Convert to domain entity
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
	}

	// Hash password before saving
	if err := user.HashPassword(req.Password); err != nil {
		return response.ErrorWithI18n(c, fiber.StatusInternalServerError, "internal_server", nil)
	}

	// Create user
	if err := h.userUsecase.CreateUser(c.Context(), user); err != nil {
		if err.Error() == "username already exists" {
			return response.ErrorWithI18n(c, fiber.StatusConflict, "username_exists", nil)
		}
		if err.Error() == "email already exists" {
			return response.ErrorWithI18n(c, fiber.StatusConflict, "email_exists", nil)
		}
		return response.ErrorWithI18n(c, fiber.StatusInternalServerError, "internal_server", nil)
	}

	// Convert to response model
	userResponse := &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}

	return response.CreatedWithI18n(c, "user_created", userResponse, nil)
}

// GetUserByID handles GET /users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.ErrorWithI18n(c, 400, "invalid_user_id", nil)
	}

	user, err := h.userUsecase.GetUserByID(c.Context(), id)
	if err != nil {
		return response.ErrorWithI18n(c, 404, "user_not_found", nil)
	}

	// Convert entity to response
	userResponse := &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response.SuccessWithI18n(c, "user_retrieved", userResponse, nil)
}

// GetAllUsers handles GET /users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	// Parse query parameters
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	// Get users
	users, err := h.userUsecase.GetAllUsers(c.Context(), limit, offset)
	if err != nil {
		return response.ErrorWithI18n(c, 500, "failed_retrieve_users", nil)
	}

	// Get total count for pagination
	totalCount, err := h.userUsecase.GetUserCount(c.Context())
	if err != nil {
		return response.ErrorWithI18n(c, 500, "failed_get_user_count", nil)
	}

	// Convert entities to responses
	userResponses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	// Prepare pagination info
	totalPages := (totalCount + limit - 1) / limit
	paginationData := map[string]interface{}{
		"users": userResponses,
		"pagination": map[string]interface{}{
			"current_page": page,
			"total_pages":  totalPages,
			"total_count":  totalCount,
			"limit":        limit,
		},
	}

	return response.SuccessWithI18n(c, "users_retrieved", paginationData, nil)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.ErrorWithI18n(c, 400, "invalid_user_id", nil)
	}

	var req model.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorWithI18n(c, 400, "invalid_request_body", nil)
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return response.ErrorWithI18n(c, 400, "validation_failed", nil)
	}

	// Convert request to entity
	user := &entity.User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Status:   req.Status,
	}

	// Hash password if provided
	if req.Password != "" {
		if err := user.HashPassword(req.Password); err != nil {
			return response.ErrorWithI18n(c, 400, "password_hashing_failed", nil)
		}
	}

	// Update user
	if err := h.userUsecase.UpdateUser(c.Context(), user); err != nil {
		return response.ErrorWithI18n(c, 500, "failed_update_user", nil)
	}

	// Get updated user
	updatedUser, err := h.userUsecase.GetUserByID(c.Context(), id)
	if err != nil {
		return response.ErrorWithI18n(c, 500, "failed_retrieve_updated_user", nil)
	}

	// Convert entity to response
	userResponse := &model.UserResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Status:    updatedUser.Status,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	return response.SuccessWithI18n(c, "user_updated", userResponse, nil)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.ErrorWithI18n(c, 400, "invalid_user_id", nil)
	}

	if err := h.userUsecase.DeleteUser(c.Context(), id); err != nil {
		return response.ErrorWithI18n(c, 500, "failed_delete_user", nil)
	}

	return response.SuccessWithI18n(c, "user_deleted", nil, nil)
}

// ForgotPassword handles POST /users/forgot-password
func (h *UserHandler) ForgotPassword(c *fiber.Ctx) error {
	var req model.ForgotPasswordRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorWithI18n(c, 400, "invalid_request_body", nil)
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return response.ErrorWithI18n(c, 400, "validation_failed", nil)
	}

	// Process forgot password
	if err := h.userUsecase.ForgotPassword(c.Context(), req.Email); err != nil {
		return response.ErrorWithI18n(c, 500, "failed_forgot_password", nil)
	}

	return response.SuccessWithI18n(c, "reset_password_sent", nil, nil)
}

// ResetPassword handles POST /users/reset-password
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	var req model.ResetPasswordRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorWithI18n(c, 400, "invalid_request_body", nil)
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return response.ErrorWithI18n(c, 400, "validation_failed", nil)
	}

	// Process reset password
	if err := h.userUsecase.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
		return response.ErrorWithI18n(c, 400, "reset_password_failed", nil)
	}

	return response.Success(c, "Password reset successfully", nil)
}

// ChangePassword handles POST /users/:id/change-password
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.ErrorWithI18n(c, 400, "invalid_user_id", nil)
	}

	var req model.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return response.ErrorWithI18n(c, 400, "invalid_request_body", nil)
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return response.ErrorWithI18n(c, 400, "validation_failed", nil)
	}

	// Process change password
	if err := h.userUsecase.ChangePassword(c.Context(), id, req.CurrentPassword, req.NewPassword); err != nil {
		return response.BadRequest(c, "Change password failed", err.Error())
	}

	return response.Success(c, "Password changed successfully", nil)
}
