package routes

import (
	"go-rest-api-template/internal/handler"

	"github.com/gofiber/fiber/v2"
)

// AdvancedRouteConfig holds configuration for advanced route setup
type AdvancedRouteConfig struct {
	// Handlers
	UserHandler *handler.UserHandler
	// ProductHandler *handler.ProductHandler  // Future
	// OrderHandler   *handler.OrderHandler    // Future
	// AuthHandler    *handler.AuthHandler     // Future

	// Optional middleware configurations
	EnableAPIVersioning bool
	EnableRateLimit     bool
	EnableCORS          bool
}

// SetupAdvancedRoutes sets up routes with versioning and middleware support
func SetupAdvancedRoutes(app *fiber.App, config *AdvancedRouteConfig) {
	// API base
	api := app.Group("/api")

	// Setup versioned routes
	setupV1Routes(api, config)
	// setupV2Routes(api, config) // Future: API v2
}

// setupV1Routes sets up version 1 API routes
func setupV1Routes(api fiber.Router, config *AdvancedRouteConfig) {
	v1 := api.Group("/v1")

	// Setup module routes for v1
	setupUserV1Routes(v1, config.UserHandler)
	// setupProductV1Routes(v1, config.ProductHandler) // Future
	// setupOrderV1Routes(v1, config.OrderHandler)     // Future
}

// setupUserV1Routes sets up user routes for API v1
func setupUserV1Routes(v1 fiber.Router, userHandler *handler.UserHandler) {
	users := v1.Group("/users")

	// CRUD operations
	users.Post("/", userHandler.CreateUser)
	users.Get("/", userHandler.GetAllUsers)
	users.Get("/:id", userHandler.GetUserByID)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)

	// Password management
	users.Post("/forgot-password", userHandler.ForgotPassword)
	users.Post("/reset-password", userHandler.ResetPassword)
	users.Post("/:id/change-password", userHandler.ChangePassword)
}

// Future: API v2 routes with different structure
// func setupV2Routes(api fiber.Router, config *AdvancedRouteConfig) {
//     v2 := api.Group("/v2")
//     // Different route structure for v2
//     setupUserV2Routes(v2, config.UserHandler)
// }

// Future: Different user routes structure for v2
// func setupUserV2Routes(v2 fiber.Router, userHandler *handler.UserHandler) {
//     // Maybe different naming or grouping
//     accounts := v2.Group("/accounts") // Different naming in v2
//     accounts.Post("/", userHandler.CreateUser)
//     // ... different structure
// }
