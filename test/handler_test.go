package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
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
	"github.com/stretchr/testify/mock"
)

// MockUserUsecase is a mock implementation of UserUsecase
type MockUserUsecase struct {
	mock.Mock
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

func (m *MockUserUsecase) CreateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUsecase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUsecase) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUsecase) UpdateUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserUsecase) GetAllUsers(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserUsecase) GetUserCount(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockUserUsecase) ForgotPassword(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockUserUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}

func (m *MockUserUsecase) ChangePassword(ctx context.Context, userID int, currentPassword, newPassword string) error {
	args := m.Called(ctx, userID, currentPassword, newPassword)
	return args.Error(0)
}

func TestUserHandler_CreateUser(t *testing.T) {
	// Setup
	mockUsecase := new(MockUserUsecase)
	responseHelper := createTestResponseHelper()
	userHandler := handler.NewUserHandler(mockUsecase, responseHelper)

	app := fiber.New()
	app.Post("/users", userHandler.CreateUser)

	// Test data
	createReq := model.UserCreateRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Mock expectations
	mockUsecase.On("CreateUser", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

	// Prepare request
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}
