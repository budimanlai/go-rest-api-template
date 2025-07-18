package application

import (
	"fmt"
	"go-rest-api-template/internal/middleware"
	"go-rest-api-template/internal/routes"

	"go-rest-api-template/pkg/database"

	gocli "github.com/budimanlai/go-cli"

	"github.com/gofiber/fiber/v2"
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

	// Add i18n middleware
	app.Use(middleware.I18nMiddleware(middleware.I18nConfig{
		DefaultLanguage: "en",
		SupportedLangs:  []string{"en", "id", "es"},
	}))

	// Initialize dependencies using dependency injection container
	container := NewContainer(db, c)

	// Setup all routes using route manager
	routeConfig := &routes.RouteConfig{
		UserHandler:   container.UserHandler,
		AuthHandler:   container.AuthHandler,
		JWTService:    container.JWTService,
		ApiKeyService: container.ApiKeyService,
		// Future: Add more handlers here
		// ProductHandler: container.ProductHandler,
		// OrderHandler:   container.OrderHandler,
	}
	routes.SetupAllRoutes(app, routeConfig)

	if err := app.Listen(":" + port); err != nil {
		c.Log(fmt.Sprintf("Failed to start server: %v", err))
	}
}
