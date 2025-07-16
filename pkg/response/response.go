package response

import (
	"github.com/gofiber/fiber/v2"
)

// StandardResponse represents standard API response format
type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse represents paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// SendSuccess sends successful response
func SendSuccess(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError sends error response
func SendError(c *fiber.Ctx, statusCode int, message string, err error) error {
	response := StandardResponse{
		Success: false,
		Message: message,
	}

	if err != nil {
		response.Error = err.Error()
	}

	return c.Status(statusCode).JSON(response)
}

// SendPaginated sends paginated response
func SendPaginated(c *fiber.Ctx, message string, data interface{}, pagination Pagination) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// Helper functions for common HTTP responses

// Success sends 200 OK response
func Success(c *fiber.Ctx, message string, data interface{}) error {
	return SendSuccess(c, message, data)
}

// Created sends 201 Created response
func Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends 400 Bad Request response
func BadRequest(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusBadRequest).JSON(StandardResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// NotFound sends 404 Not Found response
func NotFound(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusNotFound).JSON(StandardResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// InternalServerError sends 500 Internal Server Error response
func InternalServerError(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(StandardResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}
