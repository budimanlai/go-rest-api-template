package handler

import (
	"strconv"

	"go-rest-api-template/internal/constant"
	"go-rest-api-template/internal/model"
	"go-rest-api-template/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	// Add dependencies here when needed
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUserByID handles GET /users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.ErrorWithI18n(c, constant.StatusBadRequest, "invalid_user_id", nil)
	}

	// Temporary response
	userResponse := &model.UserResponse{
		ID:       id,
		Username: "example_user",
		Email:    "example@email.com",
		Status:   constant.UserStatusActive,
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
		return response.ErrorWithI18n(c, constant.StatusBadRequest, "invalid_user_id", nil)
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
		return response.ErrorWithI18n(c, constant.StatusBadRequest, "invalid_user_id", nil)
	}

	return response.SuccessWithI18n(c, "user_deleted", map[string]interface{}{
		"id":      id,
		"message": "User deleted successfully",
	}, nil)
}

// GetAllUsers handles GET /users
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users := []model.UserResponse{
		{
			ID:       1,
			Username: "user1",
			Email:    "user1@email.com",
			Status:   constant.UserStatusActive,
		},
		{
			ID:       2,
			Username: "user2",
			Email:    "user2@email.com",
			Status:   constant.UserStatusActive,
		},
	}

	return response.SuccessWithI18n(c, "users_retrieved", users, nil)
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
		return response.ErrorWithI18n(c, constant.StatusBadRequest, "invalid_user_id", nil)
	}

	return response.SuccessWithI18n(c, "password_changed", map[string]interface{}{
		"id":      id,
		"message": "Password changed successfully",
	}, nil)
}
