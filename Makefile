ROOT_DIR := $(shell pwd)
BINARY_NAME := rest-api
BUILD_DIR := build
CMD_DIR := cmd/api

# Go related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/$(BUILD_DIR)
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Build information
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse HEAD)

# Ldflags for build information
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

.PHONY: all build clean test coverage help run docker-build docker-run

# Default target
all: clean build

## Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

## Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(CMD_DIR)
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./$(CMD_DIR)
	@echo "Multi-platform build completed"

## Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME) run --port=8080

## Run in development mode with live reload (requires air)
dev:
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install it with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without live reload..."; \
		$(MAKE) run; \
	fi

## Clean build artifacts
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean completed"

## Download dependencies
deps:
	@echo "Downloading dependencies..."
	@$(GOMOD) download
	@$(GOMOD) tidy
	@echo "Dependencies updated"

## Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

## Run API authentication tests (requires server running)
test-auth:
	@echo "Running API authentication tests..."
	@go run tools/apikey_auth_tester.go

## Run tests with .env configuration
test-env:
	@echo "Running tests with .env configuration..."
	@if [ ! -f .env ]; then echo "❌ .env file not found. Please create one with TEST_API_KEY and TEST_BASE_URL"; exit 1; fi
	@$(GOTEST) -v ./test -run TestApiKeyOnlyAuthentication

## Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(GOTEST) -v -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	@$(GOTEST) -v -race ./...

## Run benchmarks
bench:
	@echo "Running benchmarks..."
	@$(GOTEST) -bench=. -benchmem ./...

## Lint code
lint:
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/"; \
	fi

## Format code
fmt:
	@echo "Formatting code..."
	@$(GOCMD) fmt ./...

## Vet code
vet:
	@echo "Vetting code..."
	@$(GOCMD) vet ./...

## Security check
security:
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

## Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):$(VERSION) .
	@docker build -t $(BINARY_NAME):latest .

## Run with Docker Compose
docker-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

## Stop Docker Compose services
docker-down:
	@echo "Stopping Docker Compose services..."
	@docker-compose down

## View Docker Compose logs
docker-logs:
	@docker-compose logs -f

## Generate API documentation (if using swag)
docs:
	@if command -v swag > /dev/null; then \
		swag init -g cmd/api/main.go; \
	else \
		echo "swag not found. Install it with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

## Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "Development tools installed"

## Show help
help:
	@echo "Available commands:"
	@echo "  build          Build the application"
	@echo "  build-all      Build for multiple platforms"
	@echo "  run            Build and run the application"
	@echo "  dev            Run in development mode with live reload"
	@echo "  clean          Clean build artifacts"
	@echo "  deps           Download and tidy dependencies"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  test-race      Run tests with race detection"
	@echo "  bench          Run benchmarks"
	@echo "  lint           Lint code"
	@echo "  fmt            Format code"
	@echo "  vet            Vet code"
	@echo "  security       Run security checks"
	@echo "  docker-build   Build Docker image"
	@echo "  docker-up      Start services with Docker Compose"
	@echo "  docker-down    Stop Docker Compose services"
	@echo "  docker-logs    View Docker Compose logs"
	@echo "  docs           Generate API documentation"
	@echo "  install-tools  Install development tools"
	@echo "  help           Show this help message"
