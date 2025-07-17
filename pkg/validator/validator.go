package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// CustomValidator wraps the validator instance for production use
type CustomValidator struct {
	validator *validator.Validate
}

// New creates a new production-ready validator instance
func New() *CustomValidator {
	validate := validator.New()

	// Register custom tag name function to use JSON tags
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{
		validator: validate,
	}
}

// Validate validates a struct and returns detailed error information
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

// ValidationError represents a structured validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Param   string `json:"param"`
	Value   string `json:"value,omitempty"`
}

// FormatValidationErrors formats validation errors into structured format
func FormatValidationErrors(err error) []ValidationError {
	var errors []ValidationError

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			ve := ValidationError{
				Field: err.Field(),
				Tag:   err.Tag(),
				Param: err.Param(),
				Value: fmt.Sprintf("%v", err.Value()),
			}

			// Generate human-readable messages
			switch err.Tag() {
			case "required":
				ve.Message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				ve.Message = fmt.Sprintf("%s must be a valid email address", err.Field())
			case "min":
				if err.Kind() == reflect.String {
					ve.Message = fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
				} else {
					ve.Message = fmt.Sprintf("%s must be at least %s", err.Field(), err.Param())
				}
			case "max":
				if err.Kind() == reflect.String {
					ve.Message = fmt.Sprintf("%s must be at most %s characters long", err.Field(), err.Param())
				} else {
					ve.Message = fmt.Sprintf("%s must be at most %s", err.Field(), err.Param())
				}
			case "len":
				ve.Message = fmt.Sprintf("%s must be exactly %s characters long", err.Field(), err.Param())
			case "oneof":
				ve.Message = fmt.Sprintf("%s must be one of: %s", err.Field(), err.Param())
			case "alpha":
				ve.Message = fmt.Sprintf("%s must contain only alphabetic characters", err.Field())
			case "alphanum":
				ve.Message = fmt.Sprintf("%s must contain only alphanumeric characters", err.Field())
			case "numeric":
				ve.Message = fmt.Sprintf("%s must be a valid number", err.Field())
			case "url":
				ve.Message = fmt.Sprintf("%s must be a valid URL", err.Field())
			case "uri":
				ve.Message = fmt.Sprintf("%s must be a valid URI", err.Field())
			case "gte":
				ve.Message = fmt.Sprintf("%s must be greater than or equal to %s", err.Field(), err.Param())
			case "lte":
				ve.Message = fmt.Sprintf("%s must be less than or equal to %s", err.Field(), err.Param())
			case "gt":
				ve.Message = fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param())
			case "lt":
				ve.Message = fmt.Sprintf("%s must be less than %s", err.Field(), err.Param())
			case "ne":
				ve.Message = fmt.Sprintf("%s must not be equal to %s", err.Field(), err.Param())
			case "eq":
				ve.Message = fmt.Sprintf("%s must be equal to %s", err.Field(), err.Param())
			default:
				ve.Message = fmt.Sprintf("%s is invalid", err.Field())
			}

			errors = append(errors, ve)
		}
	}

	return errors
}

// FormatValidationErrorsAsString formats validation errors as a single string
func FormatValidationErrorsAsString(err error) string {
	errors := FormatValidationErrors(err)
	var messages []string
	for _, e := range errors {
		messages = append(messages, e.Message)
	}
	return strings.Join(messages, "; ")
}

// Global production-ready validator instance
var globalValidator = New()

// ValidateStruct validates a struct using the global validator instance
func ValidateStruct(i interface{}) error {
	return globalValidator.Validate(i)
}

// GetValidationErrors returns structured validation errors
func GetValidationErrors(err error) []ValidationError {
	return FormatValidationErrors(err)
}

// FiberValidationErrorHandler creates a Fiber-compatible validation error handler
func FiberValidationErrorHandler(c *fiber.Ctx, err error) error {
	if validationErr := GetValidationErrors(err); len(validationErr) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  validationErr,
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"success": false,
		"message": "Invalid request data",
		"error":   err.Error(),
	})
}
