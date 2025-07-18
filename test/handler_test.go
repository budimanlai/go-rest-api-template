package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go-rest-api-template/internal/domain/entity"
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

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	users map[int]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[int]*entity.User),
	}
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	user.ID = len(m.users) + 1
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	if _, exists := m.users[user.ID]; !exists {
		return errors.New("user not found")
	}
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	if _, exists := m.users[id]; !exists {
		return errors.New("user not found")
	}
	delete(m.users, id)
	return nil
}

func (m *MockUserRepository) GetAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	users := make([]*entity.User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

func (m *MockUserRepository) GetCount(ctx context.Context) (int, error) {
	return len(m.users), nil
}

func (m *MockUserRepository) GetByVerificationToken(ctx context.Context, token string) (*entity.User, error) {
	return nil, errors.New("not implemented")
}

func (m *MockUserRepository) UpdateVerificationToken(ctx context.Context, user *entity.User) error {
	return errors.New("not implemented")
}

// Add test user to mock repository
func (m *MockUserRepository) AddTestUser(id int, username, email, status string) {
	m.users[id] = &entity.User{
		ID:       id,
		Username: username,
		Email:    email,
		Status:   status,
	}
}

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

	// Setup mock repository
	mockRepo := NewMockUserRepository()
	userHandler := handler.NewUserHandler(mockRepo)

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

	// Setup mock repository with test data
	mockRepo := NewMockUserRepository()
	mockRepo.AddTestUser(1, "testuser", "test@example.com", "active")
	userHandler := handler.NewUserHandler(mockRepo)

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

	// Setup mock repository
	mockRepo := NewMockUserRepository()
	userHandler := handler.NewUserHandler(mockRepo)

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
