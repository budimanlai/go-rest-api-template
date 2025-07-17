package routes

import (
	"go-rest-api-template/internal/handler"
	"go-rest-api-template/internal/service"

	"github.com/gofiber/fiber/v2"
)

// RouteConfig holds all handlers needed for route setup
type RouteConfig struct {
	UserHandler   *handler.UserHandler
	AuthHandler   *handler.AuthHandler
	JWTService    service.JWTService
	ApiKeyService service.ApiKeyService
	// ProductHandler *handler.ProductHandler  // Future
	// OrderHandler   *handler.OrderHandler    // Future
}

// SetupAllRoutes automatically sets up all application routes
func SetupAllRoutes(app *fiber.App, config *RouteConfig) {
	// Setup authentication routes
	SetupAuthRoutes(app, config.AuthHandler, config.ApiKeyService, config.JWTService)

	// Setup user routes
	setupUserRoutes(app, config.UserHandler, config.ApiKeyService, config.JWTService)
	// setupProductRoutes(app, config.ProductHandler)  // Future
	// setupOrderRoutes(app, config.OrderHandler)      // Future
}

// setupUserRoutes sets up user-related routes
func setupUserRoutes(app *fiber.App, userHandler *handler.UserHandler, apiKeyService service.ApiKeyService, jwtService service.JWTService) {
	SetupUserRoutes(app, userHandler, apiKeyService, jwtService)
}

// Future route setups (examples)
// func setupProductRoutes(app *fiber.App, productHandler *handler.ProductHandler) {
//     SetupProductRoutes(app, productHandler)
// }

// func setupOrderRoutes(app *fiber.App, orderHandler *handler.OrderHandler) {
//     SetupOrderRoutes(app, orderHandler)
// }

// func setupAuthRoutes(app *fiber.App, authHandler *handler.AuthHandler) {
//     SetupAuthRoutes(app, authHandler)
// }
