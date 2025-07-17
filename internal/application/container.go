package application

import (
	"go-rest-api-template/internal/domain/repository"
	"go-rest-api-template/internal/domain/usecase"
	"go-rest-api-template/internal/handler"
	repositoryImpl "go-rest-api-template/internal/repository"
	"go-rest-api-template/internal/service"

	"github.com/jmoiron/sqlx"
)

// Container holds all application dependencies
type Container struct {
	// Database
	DB *sqlx.DB

	// Repositories
	UserRepo repository.UserRepository

	// Services (Business Logic)
	UserService usecase.UserUsecase

	// Handlers (HTTP Controllers)
	UserHandler *handler.UserHandler
}

// NewContainer creates and initializes all dependencies
func NewContainer(db *sqlx.DB) *Container {
	container := &Container{
		DB: db,
	}

	// Initialize dependencies in order
	container.initRepositories()
	container.initServices()
	container.initHandlers()

	return container
}

// initRepositories initializes all repository implementations
func (c *Container) initRepositories() {
	c.UserRepo = repositoryImpl.NewUserRepository(c.DB)
}

// initServices initializes all service implementations
func (c *Container) initServices() {
	c.UserService = service.NewUserService(c.UserRepo)
}

// initHandlers initializes all HTTP handlers
func (c *Container) initHandlers() {
	c.UserHandler = handler.NewUserHandler(c.UserService)
}

// Future: Add more dependencies here
// Example when adding Product module:
// func (c *Container) initProductDependencies() {
//     c.ProductRepo = repositoryImpl.NewProductRepository(c.DB)
//     c.ProductService = service.NewProductService(c.ProductRepo)
//     c.ProductHandler = handler.NewProductHandler(c.ProductService)
// }
