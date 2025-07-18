package response

import (
	"fmt"
	"go-rest-api-template/pkg/i18n"
	"go-rest-api-template/pkg/logger"
	"go-rest-api-template/pkg/validator"
	"runtime"
	"strings"

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

// ValidationError represents validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Param   string `json:"param"`
}

// I18nResponseHelper helps create i18n responses
type I18nResponseHelper struct {
	i18nManager *i18n.Manager
}

// Global I18n response helper instance
var GlobalI18nResponseHelper *I18nResponseHelper

// NewI18nResponseHelper creates new i18n response helper
func NewI18nResponseHelper(manager *i18n.Manager) *I18nResponseHelper {
	return &I18nResponseHelper{
		i18nManager: manager,
	}
}

// Global helper functions that use the global instance

// SuccessWithI18n creates success response with i18n message using global helper
func SuccessWithI18n(c *fiber.Ctx, messageKey string, data interface{}, templateData map[string]interface{}) error {
	if GlobalI18nResponseHelper != nil {
		return GlobalI18nResponseHelper.SuccessWithI18n(c, messageKey, data, templateData)
	}
	// Fallback to regular response
	return Success(c, messageKey, data)
}

// ErrorWithI18n creates error response with i18n message using global helper
func ErrorWithI18n(c *fiber.Ctx, status int, errorKey string, templateData map[string]interface{}) error {
	// Log the error with caller information
	logErrorWithCaller(status, errorKey, templateData)

	if GlobalI18nResponseHelper != nil {
		return GlobalI18nResponseHelper.ErrorWithI18n(c, status, errorKey, templateData)
	}
	// Fallback to regular response
	return c.Status(status).JSON(StandardResponse{
		Success: false,
		Message: errorKey,
		Error:   errorKey,
	})
}

// logErrorWithCaller logs error with file and line information
func logErrorWithCaller(status int, errorKey string, templateData map[string]interface{}) {
	// Get caller information - we need to skip more levels to get to the actual handler
	var file string
	var line int
	var funcName string

	// Try different skip levels to find the handler function
	for skip := 2; skip <= 6; skip++ {
		pc, f, l, ok := runtime.Caller(skip)
		if !ok {
			continue
		}

		// Get function name
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			funcName = fn.Name()
			// If we find a handler function, use this location
			if strings.Contains(funcName, "handler") || strings.Contains(funcName, "Handler") {
				file = f
				line = l
				break
			}
		}

		// Fallback to first valid caller
		if file == "" {
			file = f
			line = l
		}
	}

	if file == "" {
		file = "unknown"
		line = 0
	}

	// Get just the filename, not the full path
	parts := strings.Split(file, "/")
	filename := parts[len(parts)-1]

	// Create concise error message
	errorMsg := fmt.Sprintf("HTTP %d: %s at %s:%d", status, errorKey, filename, line)

	if funcName != "" {
		// Extract just the function name from full path
		funcParts := strings.Split(funcName, ".")
		shortFuncName := funcParts[len(funcParts)-1]
		errorMsg += fmt.Sprintf(" in %s()", shortFuncName)
	}

	// Add template data if available
	if len(templateData) > 0 {
		errorMsg += fmt.Sprintf(" | Data: %+v", templateData)
	}

	// Log without stack trace for cleaner output
	logger.Error("%s", errorMsg)
}

// CreatedWithI18n creates 201 response with i18n message using global helper
func CreatedWithI18n(c *fiber.Ctx, messageKey string, data interface{}, templateData map[string]interface{}) error {
	if GlobalI18nResponseHelper != nil {
		return GlobalI18nResponseHelper.CreatedWithI18n(c, messageKey, data, templateData)
	}
	// Fallback to regular response
	return Created(c, messageKey, data)
}

// SuccessWithI18n creates success response with i18n message
func (h *I18nResponseHelper) SuccessWithI18n(c *fiber.Ctx, messageKey string, data interface{}, templateData map[string]interface{}) error {
	lang := getLanguageFromContext(c)
	message := h.i18nManager.TranslateSuccess(lang, messageKey, templateData)

	return c.Status(fiber.StatusOK).JSON(StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorWithI18n creates error response with i18n message
func (h *I18nResponseHelper) ErrorWithI18n(c *fiber.Ctx, status int, errorKey string, templateData map[string]interface{}) error {
	// Auto-log this error before processing
	logErrorWithCaller(status, errorKey, templateData)

	lang := getLanguageFromContext(c)
	message := h.i18nManager.TranslateError(lang, errorKey, templateData)

	return c.Status(status).JSON(StandardResponse{
		Success: false,
		Message: message,
		Error:   errorKey,
	})
}

// CreatedWithI18n creates 201 response with i18n message
func (h *I18nResponseHelper) CreatedWithI18n(c *fiber.Ctx, messageKey string, data interface{}, templateData map[string]interface{}) error {
	lang := getLanguageFromContext(c)
	message := h.i18nManager.TranslateSuccess(lang, messageKey, templateData)

	return c.Status(fiber.StatusCreated).JSON(StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ValidationErrorWithI18n creates validation error response with i18n
func (h *I18nResponseHelper) ValidationErrorWithI18n(c *fiber.Ctx, errors []ValidationError) error {
	// Auto-log validation errors
	templateData := map[string]interface{}{
		"validation_errors": errors,
	}
	logErrorWithCaller(fiber.StatusBadRequest, "validation_failed", templateData)

	lang := getLanguageFromContext(c)

	// Translate each validation error
	translatedErrors := make([]ValidationError, len(errors))
	for i, err := range errors {
		translatedErrors[i] = ValidationError{
			Field: h.i18nManager.Translate(lang, "field."+err.Field, nil),
			Message: h.i18nManager.Translate(lang, "validation."+err.Tag, map[string]interface{}{
				"Field":     h.i18nManager.Translate(lang, "field."+err.Field, nil),
				"MinLength": err.Param,
				"MaxLength": err.Param,
			}),
			Tag:   err.Tag,
			Param: err.Param,
		}
	}

	message := h.i18nManager.TranslateError(lang, "validation_failed", nil)

	return c.Status(fiber.StatusBadRequest).JSON(StandardResponse{
		Success: false,
		Message: message,
		Data:    translatedErrors,
	})
}

// getLanguageFromContext extracts language from fiber context
func getLanguageFromContext(c *fiber.Ctx) string {
	if lang, ok := c.Locals("language").(string); ok {
		return lang
	}
	return "en" // fallback to English
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
	// Auto-log this error
	templateData := map[string]interface{}{}
	if err != nil {
		templateData["error"] = err.Error()
	}
	logErrorWithCaller(statusCode, message, templateData)

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
	// Auto-log this error
	templateData := map[string]interface{}{}
	if err != "" {
		templateData["error"] = err
	}
	logErrorWithCaller(fiber.StatusBadRequest, message, templateData)

	return c.Status(fiber.StatusBadRequest).JSON(StandardResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// NotFound sends 404 Not Found response
func NotFound(c *fiber.Ctx, message string, err string) error {
	// Auto-log this error
	templateData := map[string]interface{}{}
	if err != "" {
		templateData["error"] = err
	}
	logErrorWithCaller(fiber.StatusNotFound, message, templateData)

	return c.Status(fiber.StatusNotFound).JSON(StandardResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// InternalServerError sends 500 Internal Server Error response
func InternalServerError(c *fiber.Ctx, message string, err string) error {
	// Auto-log this error with stack trace
	templateData := map[string]interface{}{}
	if err != "" {
		templateData["error"] = err
	}
	logErrorWithCaller(fiber.StatusInternalServerError, message, templateData)

	return c.Status(fiber.StatusInternalServerError).JSON(StandardResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// ValidationErrorResponse sends structured validation error response with new format
func ValidationErrorResponse(c *fiber.Ctx, message string, err error) error {
	// Get validation errors from the error
	validationErrors := validator.GetValidationErrors(err)

	// Simplify validation errors to only include field and message
	simplifiedErrors := make([]map[string]string, len(validationErrors))
	for i, ve := range validationErrors {
		simplifiedErrors[i] = map[string]string{
			"field":   ve.Field,
			"message": ve.Message,
		}
	}

	// Auto-log validation errors
	templateData := map[string]interface{}{
		"validation_errors": simplifiedErrors,
		"total_errors":      len(simplifiedErrors),
	}

	logErrorWithCaller(fiber.StatusBadRequest, "validation_failed", templateData)

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"data": nil,
		"meta": fiber.Map{
			"success": false,
			"message": message,
			"errors": fiber.Map{
				"total_errors":      len(simplifiedErrors),
				"validation_errors": simplifiedErrors,
			},
		},
	})
}
