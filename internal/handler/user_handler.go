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
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return response.BadRequest(c, "Validation failed", err.Error())
	}

	// Convert request to entity
	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// Create user
	if err := h.userUsecase.CreateUser(c.Context(), user); err != nil {
		return response.InternalServerError(c, "Failed to create user", err.Error())
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

	return response.Created(c, "User created successfully", userResponse)
}

// GetUserByID handles GET /users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err.Error())
	}

	user, err := h.userUsecase.GetUserByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "User not found", err.Error())
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

	return response.Success(c, "User retrieved successfully", userResponse)
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
		return response.InternalServerError(c, "Failed to retrieve users", err.Error())
	}

	// Get total count for pagination
	totalCount, err := h.userUsecase.GetUserCount(c.Context())
	if err != nil {
		return response.InternalServerError(c, "Failed to get user count", err.Error())
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

	return response.Success(c, "Users retrieved successfully", paginationData)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err.Error())
	}

	var req model.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", err.Error())
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return response.BadRequest(c, "Validation failed", err.Error())
	}

	// Convert request to entity
	user := &entity.User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Status:   req.Status,
	}

	// Update user
	if err := h.userUsecase.UpdateUser(c.Context(), user); err != nil {
		return response.InternalServerError(c, "Failed to update user", err.Error())
	}

	// Get updated user
	updatedUser, err := h.userUsecase.GetUserByID(c.Context(), id)
	if err != nil {
		return response.InternalServerError(c, "Failed to retrieve updated user", err.Error())
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

	return response.Success(c, "User updated successfully", userResponse)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid user ID", err.Error())
	}

	if err := h.userUsecase.DeleteUser(c.Context(), id); err != nil {
		return response.InternalServerError(c, "Failed to delete user", err.Error())
	}

	return response.Success(c, "User deleted successfully", nil)
}
