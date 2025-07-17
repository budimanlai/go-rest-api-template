package routes

import (
	"go-rest-api-template/internal/handler"

	"github.com/gofiber/fiber/v2"
)

// RouteConfig holds all handlers needed for route setup
type RouteConfig struct {
	UserHandler *handler.UserHandler
	// ProductHandler *handler.ProductHandler  // Future
	// OrderHandler   *handler.OrderHandler    // Future
	// AuthHandler    *handler.AuthHandler     // Future
}

// SetupAllRoutes automatically sets up all application routes
func SetupAllRoutes(app *fiber.App, config *RouteConfig) {
	// Setup all module routes
	setupUserRoutes(app, config.UserHandler)
	// setupProductRoutes(app, config.ProductHandler)  // Future
	// setupOrderRoutes(app, config.OrderHandler)      // Future
	// setupAuthRoutes(app, config.AuthHandler)        // Future
}

// setupUserRoutes sets up user-related routes
func setupUserRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	SetupUserRoutes(app, userHandler)
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
