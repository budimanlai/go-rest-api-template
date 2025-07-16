package application

import (
	"fmt"
	"go-rest-api-template/internal/handler"
	"go-rest-api-template/internal/repository"
	"go-rest-api-template/internal/routes"
	"go-rest-api-template/internal/service"

	"go-rest-api-template/pkg/database"

	gocli "github.com/budimanlai/go-cli"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func RestApi(c *gocli.Cli) {
	c.Log("Starting Rest API Service...")
	port := c.Args.GetString("port")
	if len(port) == 0 {
		c.Log("Port is required. Example: --port=8080")
		return
	}

	c.Log(fmt.Sprintf("Run on port: %s", port))

	c.LoadConfig()

	// Setup database connection
	dbConfig := database.Config{
		Host:     c.Config.GetString("database.hostname"),
		Port:     c.Config.GetString("database.port"),
		Username: c.Config.GetString("database.username"),
		Password: c.Config.GetString("database.password"),
		Database: c.Config.GetString("database.database"),
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		c.Log(fmt.Sprintf("Failed to connect to database: %v", err))
		return
	}

	// Set database to app context
	AppContext.Db = db

	app := fiber.New()
	defer func() {
		if err := AppContext.Db.Close(); err != nil {
			c.Log(fmt.Sprintf("Failed to close DB: %v", err))
		}
		app.Shutdown()
	}()

	app.Use(logger.New(logger.Config{
		TimeZone:   "Asia/Jakarta",
		TimeFormat: "2006-Jan-02 15:04:05",
		Format:     "${time} | :" + port + " | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n${body}\n${resBody}\n\n",
	}))
	app.Use(keyauth.New(keyauth.Config{
		KeyLookup: "header:x-api-key",
		Validator: validateAPIKey,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized. Invalid api key",
			})
		},
	}))

	// Initialize dependencies using dependency injection
	// Repository layer
	userRepo := repository.NewUserRepository(db)

	// Service layer (Business logic)
	userService := service.NewUserService(userRepo)

	// Handler layer (HTTP controllers)
	userHandler := handler.NewUserHandler(userService)

	// Setup routes
	routes.SetupUserRoutes(app, userHandler)

	if err := app.Listen(":" + port); err != nil {
		c.Log(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func validateAPIKey(c *fiber.Ctx, key string) (bool, error) {
	// For now, let's use a simple validation
	// You can implement proper API key validation here
	// For example, check against database or hardcoded keys
	validKeys := []string{"test-api-key", "development-key"}

	for _, validKey := range validKeys {
		if key == validKey {
			return true, nil
		}
	}

	return false, nil
}
