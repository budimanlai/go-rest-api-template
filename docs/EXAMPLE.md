# Example Implementation: User CRUD API

This example shows how to implement a complete CRUD API using the template structure.

## 1. Database Schema

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    status ENUM('active', 'inactive') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE api_keys (
    id INT AUTO_INCREMENT PRIMARY KEY,
    key VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    status ENUM('active', 'inactive') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample API keys
INSERT INTO api_keys (key, name) VALUES 
('test-api-key', 'Development Key'),
('production-key', 'Production Key');
```

## 2. Complete User Implementation

### Entity (Business Logic)
```go
// internal/domain/entity/user.go
package entity

import (
    "errors"
    "time"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Business rules
func (u *User) ValidateForCreate() error {
    if len(u.Username) < 3 {
        return errors.New("username must be at least 3 characters")
    }
    if len(u.Password) < 6 {
        return errors.New("password must be at least 6 characters")
    }
    return nil
}

func (u *User) HashPassword() error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}
```

### Repository Interface
```go
// internal/domain/repository/user_repository.go
package repository

import (
    "context"
    "go-rest-api-template/internal/domain/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id int) (*entity.User, error)
    GetByEmail(ctx context.Context, email string) (*entity.User, error)
    GetByUsername(ctx context.Context, username string) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id int) error
    GetAll(ctx context.Context, limit, offset int) ([]*entity.User, error)
    GetCount(ctx context.Context) (int, error)
}
```

### Use Case (Business Logic)
```go
// internal/domain/usecase/user_usecase.go
package usecase

import (
    "context"
    "errors"
    "go-rest-api-template/internal/domain/entity"
    "go-rest-api-template/internal/domain/repository"
)

var (
    ErrUserNotFound      = errors.New("user not found")
    ErrUserAlreadyExists = errors.New("user already exists")
    ErrInvalidPassword   = errors.New("invalid password")
)

type UserUsecase interface {
    CreateUser(ctx context.Context, user *entity.User) error
    GetUser(ctx context.Context, id int) (*entity.User, error)
    UpdateUser(ctx context.Context, id int, updates map[string]interface{}) (*entity.User, error)
    DeleteUser(ctx context.Context, id int) error
    GetUsers(ctx context.Context, limit, offset int) ([]*entity.User, int, error)
    AuthenticateUser(ctx context.Context, email, password string) (*entity.User, error)
}

type userUsecase struct {
    userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
    return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *entity.User) error {
    // Validate business rules
    if err := user.ValidateForCreate(); err != nil {
        return err
    }

    // Check if user already exists
    existingUser, _ := u.userRepo.GetByEmail(ctx, user.Email)
    if existingUser != nil {
        return ErrUserAlreadyExists
    }

    existingUser, _ = u.userRepo.GetByUsername(ctx, user.Username)
    if existingUser != nil {
        return ErrUserAlreadyExists
    }

    // Hash password
    if err := user.HashPassword(); err != nil {
        return err
    }

    // Set default status
    user.Status = "active"

    return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) GetUser(ctx context.Context, id int) (*entity.User, error) {
    user, err := u.userRepo.GetByID(ctx, id)
    if err != nil {
        return nil, ErrUserNotFound
    }
    return user, nil
}

func (u *userUsecase) AuthenticateUser(ctx context.Context, email, password string) (*entity.User, error) {
    user, err := u.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return nil, ErrInvalidPassword
    }

    if !user.CheckPassword(password) {
        return nil, ErrInvalidPassword
    }

    return user, nil
}

// ... implement other methods
```

## 3. HTTP Layer Implementation

### Request/Response Models
```go
// internal/model/user_dto.go
package model

type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
    Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
    Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

type UserResponse struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Status    string `json:"status"`
    CreatedAt string `json:"created_at"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}
```

### HTTP Handler
```go
// internal/handler/user_handler.go
package handler

import (
    "strconv"
    "go-rest-api-template/internal/domain/entity"
    "go-rest-api-template/internal/domain/usecase"
    "go-rest-api-template/internal/model"
    "go-rest-api-template/pkg/response"
    "go-rest-api-template/pkg/validator"

    "github.com/gofiber/fiber/v2"
)

type UserHandler struct {
    userUsecase usecase.UserUsecase
    validator   *validator.Validator
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
    return &UserHandler{
        userUsecase: userUsecase,
        validator:   validator.New(),
    }
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req model.CreateUserRequest
    
    if err := c.BodyParser(&req); err != nil {
        return response.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
    }

    if err := h.validator.Validate(req); err != nil {
        return response.SendError(c, fiber.StatusBadRequest, validator.FormatValidationErrorString(err), err)
    }

    user := &entity.User{
        Username: req.Username,
        Email:    req.Email,
        Password: req.Password,
    }

    if err := h.userUsecase.CreateUser(c.Context(), user); err != nil {
        if err == usecase.ErrUserAlreadyExists {
            return response.SendError(c, fiber.StatusConflict, err.Error(), err)
        }
        return response.SendError(c, fiber.StatusInternalServerError, "Failed to create user", err)
    }

    userResponse := model.UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Email:     user.Email,
        Status:    user.Status,
        CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
    }

    return response.SendSuccess(c, "User created successfully", userResponse)
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        return response.SendError(c, fiber.StatusBadRequest, "Invalid user ID", err)
    }

    user, err := h.userUsecase.GetUser(c.Context(), id)
    if err != nil {
        if err == usecase.ErrUserNotFound {
            return response.SendError(c, fiber.StatusNotFound, err.Error(), err)
        }
        return response.SendError(c, fiber.StatusInternalServerError, "Failed to get user", err)
    }

    userResponse := model.UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Email:     user.Email,
        Status:    user.Status,
        CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
    }

    return response.SendSuccess(c, "User retrieved successfully", userResponse)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
    var req model.LoginRequest
    
    if err := c.BodyParser(&req); err != nil {
        return response.SendError(c, fiber.StatusBadRequest, "Invalid request body", err)
    }

    if err := h.validator.Validate(req); err != nil {
        return response.SendError(c, fiber.StatusBadRequest, validator.FormatValidationErrorString(err), err)
    }

    user, err := h.userUsecase.AuthenticateUser(c.Context(), req.Email, req.Password)
    if err != nil {
        return response.SendError(c, fiber.StatusUnauthorized, "Invalid credentials", err)
    }

    userResponse := model.UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Email:     user.Email,
        Status:    user.Status,
        CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
    }

    return response.SendSuccess(c, "Login successful", userResponse)
}

// ... implement UpdateUser, DeleteUser, GetUsers
```

## 4. Routes Setup
```go
// internal/routes/user_routes.go
package routes

import (
    "go-rest-api-template/internal/handler"
    "github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, userHandler *handler.UserHandler) {
    api := app.Group("/api/v1")
    
    // Public routes
    api.Post("/auth/login", userHandler.Login)
    
    // Protected routes
    users := api.Group("/users")
    users.Post("/", userHandler.CreateUser)
    users.Get("/:id", userHandler.GetUser)
    users.Put("/:id", userHandler.UpdateUser)
    users.Delete("/:id", userHandler.DeleteUser)
    users.Get("/", userHandler.GetUsers)
}
```

## 5. Wire Everything Together
```go
// internal/application/rest_api.go (additions)

func RestApi(c *gocli.Cli) {
    // ... existing setup code ...

    // Dependency injection
    userRepo := repository.NewUserRepository(AppContext.Db)
    userUsecase := usecase.NewUserUsecase(userRepo)
    userHandler := handler.NewUserHandler(userUsecase)

    // Setup routes
    routes.SetupUserRoutes(app, userHandler)

    // Health check route
    app.Get("/health", func(c *fiber.Ctx) error {
        return response.SendSuccess(c, "Server is healthy", fiber.Map{
            "status": "OK",
            "timestamp": time.Now(),
        })
    })

    // ... rest of the code ...
}
```

## 6. API Usage Examples

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "x-api-key: test-api-key" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get User
```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "x-api-key: test-api-key"
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "x-api-key: test-api-key" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

This example demonstrates a complete implementation following Clean Architecture principles with proper separation of concerns.
