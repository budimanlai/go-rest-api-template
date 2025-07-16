package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

// Validator wraps the validator instance
type Validator struct {
	validator *validator.Validate
}

// New creates a new validator instance
func New() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

// Validate validates a struct
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

// FormatValidationError formats validation errors into readable messages
func FormatValidationError(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			switch err.Tag() {
			case "required":
				errors = append(errors, fmt.Sprintf("%s is required", err.Field()))
			case "email":
				errors = append(errors, fmt.Sprintf("%s must be a valid email", err.Field()))
			case "min":
				errors = append(errors, fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param()))
			case "max":
				errors = append(errors, fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param()))
			default:
				errors = append(errors, fmt.Sprintf("%s is invalid", err.Field()))
			}
		}
	}

	return errors
}

// FormatValidationErrorString formats validation errors into a single string
func FormatValidationErrorString(err error) string {
	errors := FormatValidationError(err)
	return strings.Join(errors, ", ")
}

// Global validator instance
var globalValidator = New()

// ValidateStruct validates a struct using global validator
func ValidateStruct(i interface{}) error {
	return globalValidator.Validate(i)
}
