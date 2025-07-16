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
	userGroup.Post("/", userHandler.CreateUser)
	userGroup.Get("/", userHandler.GetAllUsers)
	userGroup.Get("/:id", userHandler.GetUserByID)
	userGroup.Put("/:id", userHandler.UpdateUser)
	userGroup.Delete("/:id", userHandler.DeleteUser)

	// Password management routes
	userGroup.Post("/forgot-password", userHandler.ForgotPassword)
	userGroup.Post("/reset-password", userHandler.ResetPassword)
	userGroup.Post("/:id/change-password", userHandler.ChangePassword)
}
