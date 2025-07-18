package handler_test

import (
	"bytes"
	"encoding/json"
	"go-rest-api-template/internal/handler"
	"go-rest-api-template/internal/model"
	"go-rest-api-template/pkg/i18n"
	"go-rest-api-template/pkg/response"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// createTestResponseHelper creates a response helper for testing with minimal i18n setup
func createTestResponseHelper() *response.I18nResponseHelper {
	// Create a simple i18n manager for testing
	config := i18n.Config{
		DefaultLanguage: "en",
		LocalesPath:     "../locales", // Relative to test directory
		SupportedLangs:  []string{"en"},
	}

	// Create manager (if fails, create nil manager for basic testing)
	manager, err := i18n.NewManager(config)
	if err != nil {
		// For testing, create a simple manager that just returns keys
		return createSimpleResponseHelper()
	}

	return response.NewI18nResponseHelper(manager)
}

// setupTestGlobalHelpers sets up global helpers for testing
func setupTestGlobalHelpers() {
	responseHelper := createTestResponseHelper()
	response.GlobalI18nResponseHelper = responseHelper
}

// createSimpleResponseHelper creates a response helper with minimal setup
func createSimpleResponseHelper() *response.I18nResponseHelper {
	// For testing, we'll create a basic setup
	// This is a workaround since we can't easily mock the concrete types
	config := i18n.Config{
		DefaultLanguage: "en",
		LocalesPath:     ".",
		SupportedLangs:  []string{"en"},
	}

	manager, _ := i18n.NewManager(config)
	return response.NewI18nResponseHelper(manager)
}

func TestUserHandler_CreateUser(t *testing.T) {
	// Setup global helpers first
	setupTestGlobalHelpers()

	// Setup
	userHandler := handler.NewUserHandler()

	app := fiber.New()
	app.Post("/users", userHandler.CreateUser)

	// Test data
	createReq := model.UserCreateRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Prepare request
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode) // Handler returns 200, not 201
}

func TestUserHandler_GetUserByID(t *testing.T) {
	// Setup global helpers first
	setupTestGlobalHelpers()

	// Setup
	userHandler := handler.NewUserHandler()

	app := fiber.New()
	app.Get("/users/:id", userHandler.GetUserByID)

	// Prepare request
	req := httptest.NewRequest("GET", "/users/1", nil)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	// Setup global helpers first
	setupTestGlobalHelpers()

	// Setup
	userHandler := handler.NewUserHandler()

	app := fiber.New()
	app.Get("/users", userHandler.GetAllUsers)

	// Prepare request
	req := httptest.NewRequest("GET", "/users", nil)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
