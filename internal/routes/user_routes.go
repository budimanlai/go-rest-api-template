package routes

import (
	"go-rest-api-template/internal/handler"
	"go-rest-api-template/internal/middleware"
	"go-rest-api-template/internal/service"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes sets up user-related routes with Private JWT middleware
func SetupUserRoutes(app *fiber.App, userHandler *handler.UserHandler, apiKeyService service.ApiKeyService, jwtService service.JWTService) {
	// API versioning
	v1 := app.Group("/api/v1")

	// Private middleware - requires API key + private JWT token
	privateMiddleware := middleware.PrivateMiddleware(apiKeyService, jwtService)

	// Create user routes group with private middleware
	userGroup := v1.Group("/users", privateMiddleware)

	// User CRUD routes
	userGroup.Get("/", userHandler.GetAllUsers)
	userGroup.Get("/:id", userHandler.GetUserByID)
	userGroup.Put("/:id", userHandler.UpdateUser)
	userGroup.Delete("/:id", userHandler.DeleteUser)

	// Password management routes
	userGroup.Post("/forgot-password", userHandler.ForgotPassword)
	userGroup.Post("/reset-password", userHandler.ResetPassword)
	userGroup.Post("/:id/change-password", userHandler.ChangePassword)
}
