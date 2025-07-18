package application

import (
	"go-rest-api-template/internal/domain/repository"
	"go-rest-api-template/internal/domain/usecase"
	"go-rest-api-template/internal/handler"
	"go-rest-api-template/internal/middleware"
	repositoryImpl "go-rest-api-template/internal/repository"
	"go-rest-api-template/internal/service"
	"go-rest-api-template/pkg/i18n"
	"go-rest-api-template/pkg/response"

	gocli "github.com/budimanlai/go-cli"
	"github.com/jmoiron/sqlx"
)

// Container holds all application dependencies
type Container struct {
	// Database
	DB *sqlx.DB

	// JWT Configuration
	jwtSecret          string
	publicTokenExpiry  int
	privateTokenExpiry int

	// I18n
	I18nManager *i18n.Manager

	// Repositories
	UserRepo   repository.UserRepository
	ApiKeyRepo repository.ApiKeyRepository

	// Services (Business Logic)
	JWTService    service.JWTService
	UserService   usecase.UserUsecase
	ApiKeyService service.ApiKeyService

	// Handlers (HTTP Controllers)
	UserHandler *handler.UserHandler
	AuthHandler *handler.AuthHandler
}

// NewContainer creates and initializes all dependencies
func NewContainer(db *sqlx.DB, config *gocli.Cli) *Container {
	container := &Container{
		DB: db,
		// Initialize JWT configuration from config file
		jwtSecret:          config.Config.GetString("jwt.secret"),
		publicTokenExpiry:  config.Config.GetInt("jwt.public_token_expiry_hours"),
		privateTokenExpiry: config.Config.GetInt("jwt.private_token_expiry_hours"),
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

	// Create response helper and set it globally
	responseHelper := response.NewI18nResponseHelper(manager)

	// Set global middleware helper to avoid import cycle
	middleware.SetGlobalI18nResponseHelper(responseHelper)

	// Set global response helper for direct usage in handlers
	response.GlobalI18nResponseHelper = responseHelper
}

// initRepositories initializes all repository implementations
func (c *Container) initRepositories() {
	c.UserRepo = repositoryImpl.NewUserRepository(c.DB)
	c.ApiKeyRepo = repositoryImpl.NewApiKeyRepository(c.DB)
}

// initServices initializes all service implementations
func (c *Container) initServices() {
	c.ApiKeyService = service.NewApiKeyService(c.ApiKeyRepo)
	c.JWTService = service.NewJWTService(c.jwtSecret, c.publicTokenExpiry, c.privateTokenExpiry, c.ApiKeyService)
	c.UserService = service.NewUserService(c.UserRepo, c.JWTService)
}

// initHandlers initializes all HTTP handlers
func (c *Container) initHandlers() {
	c.UserHandler = handler.NewUserHandler()
	c.AuthHandler = handler.NewAuthHandler(c.UserService, c.JWTService, c.ApiKeyService)
}

// Future: Add more dependencies here
// Example when adding Product module:
// func (c *Container) initProductDependencies() {
//     c.ProductRepo = repositoryImpl.NewProductRepository(c.DB)
//     c.ProductService = service.NewProductService(c.ProductRepo)
//     c.ProductHandler = handler.NewProductHandler(c.ProductService)
// }
