package routes

import (
	"go-rest-api-template/internal/handler"

	"github.com/gofiber/fiber/v2"
)

// SetupUserRoutes sets up user-related routes
func SetupUserRoutes(app *fiber.App, userHandler *handler.UserHandler) {
	// Create user routes group
	userGroup := app.Group("/api/v1/users")

	// User CRUD routes
	userGroup.Post("/", userHandler.CreateUser)      // POST /api/v1/users
	userGroup.Get("/", userHandler.GetAllUsers)      // GET /api/v1/users
	userGroup.Get("/:id", userHandler.GetUserByID)   // GET /api/v1/users/:id
	userGroup.Put("/:id", userHandler.UpdateUser)    // PUT /api/v1/users/:id
	userGroup.Delete("/:id", userHandler.DeleteUser) // DELETE /api/v1/users/:id
}
