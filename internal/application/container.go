package application

import (
	"go-rest-api-template/internal/domain/repository"
	"go-rest-api-template/internal/domain/usecase"
	"go-rest-api-template/internal/handler"
	repositoryImpl "go-rest-api-template/internal/repository"
	"go-rest-api-template/internal/service"
	"go-rest-api-template/pkg/i18n"
	"go-rest-api-template/pkg/response"

	"github.com/jmoiron/sqlx"
)

// Container holds all application dependencies
type Container struct {
	// Database
	DB *sqlx.DB

	// I18n
	I18nManager    *i18n.Manager
	ResponseHelper *response.I18nResponseHelper

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
	container.initI18n()
	container.initRepositories()
	container.initServices()
	container.initHandlers()

	return container
}

// initI18n initializes internationalization
func (c *Container) initI18n() {
	i18nConfig := i18n.Config{
		DefaultLanguage: "en",
		LocalesPath:     "./locales",
		SupportedLangs:  []string{"en", "id"},
	}

	manager, err := i18n.NewManager(i18nConfig)
	if err != nil {
		// Log error but don't fail startup
		panic("Failed to initialize i18n: " + err.Error())
	}

	c.I18nManager = manager
	c.ResponseHelper = response.NewI18nResponseHelper(manager)
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
	c.UserHandler = handler.NewUserHandler(c.UserService, c.ResponseHelper)
}

// Future: Add more dependencies here
// Example when adding Product module:
// func (c *Container) initProductDependencies() {
//     c.ProductRepo = repositoryImpl.NewProductRepository(c.DB)
//     c.ProductService = service.NewProductService(c.ProductRepo)
//     c.ProductHandler = handler.NewProductHandler(c.ProductService)
// }
