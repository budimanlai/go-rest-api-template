package routes

import (
	"go-rest-api-template/internal/handler"
	"go-rest-api-template/internal/middleware"
	"go-rest-api-template/internal/service"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes sets up authentication routes
func SetupAuthRoutes(app *fiber.App, authHandler *handler.AuthHandler, apiKeyService service.ApiKeyService, jwtService service.JWTService) {
	// API versioning
	v1 := app.Group("/api/v1")

	// API Key Only Middleware - following industry best practices for register/login endpoints
	// This approach is used by major platforms like Strapi, GitHub, Twitter/X
	apiKeyOnlyMiddleware := middleware.ApiKeyOnlyMiddleware(apiKeyService)

	// Public endpoints - only require API key (no JWT tokens needed)
	public := v1.Group("/public", apiKeyOnlyMiddleware)

	// Auth endpoints - simplified authentication following industry standards
	auth := public.Group("/auth")
	auth.Get("/token", authHandler.GetPublicToken)  // GET /api/v1/public/auth/token - Get public token
	auth.Post("/register", authHandler.Register)    // POST /api/v1/public/auth/register - Register new user
	auth.Post("/login", authHandler.Login)          // POST /api/v1/public/auth/login - Login (get private token)
	auth.Post("/refresh", authHandler.RefreshToken) // POST /api/v1/public/auth/refresh - Refresh token
	auth.Post("/logout", authHandler.Logout)        // POST /api/v1/public/auth/logout - Logout

	// Private endpoints - require API key + private JWT token
	// privateMiddleware := middleware.PrivateMiddleware(apiKeyService, jwtService)
	// private := v1.Group("/private", privateMiddleware)

	// Private auth endpoints (if needed) - commented out until method is implemented
	// privateAuth := private.Group("/auth")
	// privateAuth.Get("/me", authHandler.GetCurrentUser) // GET /api/v1/private/auth/me - Get current user info
}
