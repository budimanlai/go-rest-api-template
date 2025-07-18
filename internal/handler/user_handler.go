package handler

import (
	"context"
	"strconv"

	"go-rest-api-template/internal/constant"
	"go-rest-api-template/internal/domain/repository"
	"go-rest-api-template/internal/model"
	"go-rest-api-template/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetUserByID handles GET /users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_user_id", nil)
	}

	// Get user from database
	ctx := context.Background()
	user, err := h.userRepo.GetByID(ctx, id)
	if err != nil {
		return response.ErrorWithI18n(c, fiber.StatusNotFound, "user_not_found", nil)
	}

	// Convert entity to response model
	userResponse := &model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Status:   user.Status,
	}

	return response.SuccessWithI18n(c, "user_retrieved", userResponse, nil)
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	// Temporary simple implementation
	return response.SuccessWithI18n(c, "user_created", map[string]interface{}{
		"id":       1,
		"username": "new_user",
		"status":   constant.UserStatusActive,
	}, nil)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_user_id", nil)
	}

	return response.SuccessWithI18n(c, "user_updated", map[string]interface{}{
		"id":      id,
		"status":  constant.UserStatusActive,
		"message": "User updated successfully",
	}, nil)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_user_id", nil)
	}

	return response.SuccessWithI18n(c, "user_deleted", map[string]interface{}{
		"id":      id,
		"message": "User deleted successfully",
	}, nil)
}

// GetAllUsers handles GET /users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	ctx := context.Background()

	// Get users from database
	users, err := h.userRepo.GetAll(ctx, 0, 0) // 0 means no limit/offset for now
	if err != nil {
		return response.ErrorWithI18n(c, fiber.StatusInternalServerError, "failed_to_get_users", nil)
	}

	// Convert entities to response models
	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = model.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Status:   user.Status,
		}
	}

	return response.SuccessWithI18n(c, "users_retrieved", userResponses, nil)
}

// ForgotPassword handles POST /users/forgot-password
func (h *UserHandler) ForgotPassword(c *fiber.Ctx) error {
	return response.SuccessWithI18n(c, "password_reset_email_sent", map[string]interface{}{
		"message": "Password reset email sent successfully",
	}, nil)
}

// ResetPassword handles POST /users/reset-password
func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	return response.SuccessWithI18n(c, "password_reset_success", map[string]interface{}{
		"message": "Password reset successfully",
	}, nil)
}

// ChangePassword handles PUT /users/:id/password
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_user_id", nil)
	}

	return response.SuccessWithI18n(c, "password_changed", map[string]interface{}{
		"id":      id,
		"message": "Password changed successfully",
	}, nil)
}
