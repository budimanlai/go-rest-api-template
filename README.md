# Go REST API Template

A production-ready REST API template built with Go, following Clean Architecture principles and best practices.

## 🏗️ Architecture

This template follows **Clean Architecture** with clear separation of concerns:

```
cmd/api/                    # Application entry point
├── main.go                 # Main application file

internal/                   # Private application code
├── application/            # Application layer
│   ├── rest_api.go        # REST API setup and configuration
│   └── app.go             # Application context
├── domain/                 # Business logic layer (no dependencies)
│   ├── entity/            # Business entities
│   ├── repository/        # Repository interfaces
│   └── usecase/           # Business use cases
├── handler/               # HTTP handlers (interface adapters)
├── middleware/            # Custom middleware
├── model/                 # Data models and DTOs
├── repository/            # Repository implementations
├── routes/                # Route definitions
└── service/               # External service integrations

pkg/                       # Public packages (reusable)
├── database/              # Database utilities
├── response/              # HTTP response helpers
└── validator/             # Validation utilities

configs/                   # Configuration files
├── config.json            # Application configuration

docs/                      # Documentation
test/                      # Test files
```

## 🚀 Features

- ✅ **Clean Architecture** with proper layering
- ✅ **Go Fiber** v2 web framework
- ✅ **MySQL** database with SQLX
- ✅ **API Key Authentication** middleware
- ✅ **Request/Response** logging
- ✅ **Graceful Shutdown**
- ✅ **Configuration Management**
- ✅ **Input Validation**
- ✅ **Standardized Responses**
- ✅ **Database Connection Pooling**
- ✅ **Error Handling**

## 🛠️ Tech Stack

- **Framework**: Go Fiber v2
- **Database**: MySQL with SQLX
- **Validation**: go-playground/validator
- **CLI**: Custom CLI framework
- **Authentication**: API Key based
- **Password**: bcrypt hashing

## 📦 Dependencies

```go
require (
    github.com/gofiber/fiber/v2
    github.com/jmoiron/sqlx
    github.com/go-sql-driver/mysql
    github.com/go-playground/validator
    golang.org/x/crypto
)
```

## 🏃‍♂️ Quick Start

### 1. Clone & Install
```bash
git clone <repository-url>
cd go-rest_api
go mod tidy
```

### 2. Configuration
Edit `configs/config.json`:
```json
{
    "database": {
        "hostname": "127.0.0.1",
        "port": 3306,
        "username": "root",
        "password": "your_password",
        "database": "your_database"
    },
    "server": {
        "host": "127.0.0.1",
        "port": 8080,
        "debug": true
    }
}
```

### 3. Build & Run
```bash
# Build
go build -o rest-api ./cmd/api

# Run
./rest-api run --port=8080
```

### 4. Test API
```bash
# Valid API key
curl -X GET http://127.0.0.1:8080/api/health \
  -H "x-api-key: test-api-key"

# Invalid API key
curl -X GET http://127.0.0.1:8080/api/health \
  -H "x-api-key: invalid-key"
```

## 📋 Implementation Guide

### 1. Create Business Entity
```go
// internal/domain/entity/user.go
type User struct {
    ID       int
    Username string
    Email    string
    Password string
}

func (u *User) Validate() error {
    // Business validation rules
}
```

### 2. Define Repository Interface
```go
// internal/domain/repository/user_repository.go
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id int) (*entity.User, error)
    // ... other methods
}
```

### 3. Implement Repository
```go
// internal/repository/user_repository_impl.go
type userRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repository.UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
    // Database implementation
}
```

### 4. Create Use Case
```go
// internal/domain/usecase/user_usecase.go
type UserUsecase interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
}

type userUsecase struct {
    userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
    return &userUsecase{userRepo: userRepo}
}
```

### 5. Implement HTTP Handler
```go
// internal/handler/user_handler.go
type UserHandler struct {
    userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
    return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    // HTTP handler implementation
}
```

### 6. Setup Routes
```go
// internal/routes/user_routes.go
func SetupUserRoutes(app *fiber.App, userHandler *handler.UserHandler) {
    api := app.Group("/api/v1")
    
    users := api.Group("/users")
    users.Post("/", userHandler.CreateUser)
    users.Get("/:id", userHandler.GetUser)
    users.Put("/:id", userHandler.UpdateUser)
    users.Delete("/:id", userHandler.DeleteUser)
}
```

## 🔧 Development

### Project Structure Guidelines

1. **Domain Layer** (`internal/domain/`)
   - Contains pure business logic
   - No external dependencies
   - Defines interfaces for data access

2. **Infrastructure Layer** (`internal/repository/`, `internal/service/`)
   - Implements domain interfaces
   - Handles external dependencies (database, APIs)

3. **Interface Layer** (`internal/handler/`, `internal/routes/`)
   - HTTP request/response handling
   - Input validation and transformation

4. **Application Layer** (`internal/application/`)
   - Application configuration
   - Dependency injection
   - Service orchestration

### Best Practices

1. **Dependency Direction**: Always point inward to domain
2. **Interface Segregation**: Keep interfaces small and focused
3. **Single Responsibility**: Each layer has one reason to change
4. **Error Handling**: Use custom error types
5. **Testing**: Mock interfaces for unit tests

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/domain/usecase/...
```

## 📝 API Documentation

### Standard Response Format
```json
{
    "success": true,
    "message": "Operation successful",
    "data": {...}
}
```

### Error Response Format
```json
{
    "success": false,
    "message": "Error message",
    "error": "Detailed error information"
}
```

### Authentication
All endpoints require `x-api-key` header:
```
x-api-key: test-api-key
```

## 🚀 Deployment

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o rest-api ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/rest-api .
COPY --from=builder /app/configs ./configs
CMD ["./rest-api", "run", "--port=8080"]
```

### Environment Variables
```bash
export DATABASE_HOST=localhost
export DATABASE_PORT=3306
export DATABASE_USER=root
export DATABASE_PASS=password
export DATABASE_NAME=myapp
```

## 📚 Learning Resources

- [Clean Architecture by Robert Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Fiber Documentation](https://docs.gofiber.io/)
- [SQLX Documentation](http://jmoiron.github.io/sqlx/)

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
