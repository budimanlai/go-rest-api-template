package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"go-rest-api-template/internal/domain/entity"
	"go-rest-api-template/internal/handler"
	"go-rest-api-template/internal/model"
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

func TestUserHandler_CreateUser(t *testing.T) {
	// Setup
	mockUsecase := new(MockUserUsecase)
	userHandler := handler.NewUserHandler(mockUsecase)

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
