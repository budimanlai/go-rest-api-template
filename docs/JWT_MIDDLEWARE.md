# JWT Authentication and API Key Middleware

This document explains how to use the JWT authentication and API key middleware system.

## Overview

The system provides three types of middleware:

1. **PublicMiddleware** - API key required, no user authentication
2. **PrivateMiddleware** - API key + JWT token required  
3. **OptionalJWTMiddleware** - API key required, JWT token optional

## Middleware Types

### 1. PublicMiddleware (API Key Only)
For endpoints that don't require user authentication, only API key validation. No user information is stored in context.

**Usage:**
```go
publicMiddleware := middleware.PublicMiddleware(container.ApiKeyService)
app.Use("/api/v1/public", publicMiddleware)
```

**Headers Required:**
- `X-API-Key: your-api-key` OR
- `Authorization: ApiKey your-api-key`

**Use Cases:**
- Login endpoint
- Register endpoint  
- Forgot/Reset password
- Public data endpoints

### 2. PrivateMiddleware (API Key + JWT)
For endpoints that require both API key and user authentication.

**Usage:**
```go
privateMiddleware := middleware.PrivateMiddleware(container.ApiKeyService, container.JWTService)
app.Use("/api/v1/private", privateMiddleware)
```

**Headers Required:**
- `X-API-Key: your-api-key`
- `Authorization: Bearer jwt-token`

**Use Cases:**
- User profile management
- Change password
- User-specific data
- Admin functions

### 3. OptionalJWTMiddleware (API Key + Optional JWT)
For endpoints that work with or without user authentication.

**Usage:**
```go
optionalMiddleware := middleware.OptionalJWTMiddleware(container.ApiKeyService, container.JWTService)
app.Use("/api/v1/optional", optionalMiddleware)
```

**Headers Required:**
- `X-API-Key: your-api-key`
- `Authorization: Bearer jwt-token` (optional)

## Example Route Setup

```go
func SetupRoutes(app *fiber.App, container *application.Container) {
    // Initialize middleware (note: PublicMiddleware only requires ApiKeyService)
    publicMW := middleware.PublicMiddleware(container.ApiKeyService)
    privateMW := middleware.PrivateMiddleware(container.ApiKeyService, container.JWTService)
    
    v1 := app.Group("/api/v1")
    
    // Public routes (API key required only, no user info)
    public := v1.Group("/public", publicMW)
    public.Post("/auth/login", authHandler.Login)
    public.Post("/auth/register", userHandler.CreateUser)
    public.Get("/location/countries", locationHandler.GetCountries)
    
    // Private routes (API key + JWT required)
    private := v1.Group("/private", privateMW)
    private.Get("/users/profile", userHandler.GetProfile)
    private.Put("/users/profile", userHandler.UpdateProfile)
    private.Get("/users", userHandler.GetAllUsers)
}
```

## Context Helper

Use the `ContextHelper` to extract data from the request context:

```go
contextHelper := middleware.NewContextHelper()

// In your handlers
func (h *Handler) SomeHandler(c *fiber.Ctx) error {
    // Get API key info (available in all middleware)
    apiKeyID := contextHelper.MustGetAPIKeyID(c)
    apiKeyName, _ := contextHelper.GetAPIKeyName(c)
    isH2H := contextHelper.IsH2HEnabled(c)
    
    // Get user info (available in private middleware)
    userID := contextHelper.MustGetUserID(c)
    username, _ := contextHelper.GetUsername(c)
    email, _ := contextHelper.GetUserEmail(c)
    
    // Check authentication status (for optional middleware)
    isAuth := contextHelper.IsAuthenticated(c)
    
    return c.JSON(fiber.Map{
        "api_key_id": apiKeyID,
        "user_id": userID,
        "authenticated": isAuth,
    })
}
```

## JWT Service Configuration

Configure JWT service in your container:

```go
func (c *Container) initServices() {
    // JWT configuration
    jwtSecret := os.Getenv("JWT_SECRET") // Use environment variable
    issuer := "go-rest-api-template"
    tokenExpiry := time.Hour * 24        // 24 hours
    refreshExpiry := time.Hour * 24 * 7  // 7 days for refresh
    
    c.JWTService = service.NewJWTService(jwtSecret, issuer, tokenExpiry, refreshExpiry)
    c.UserService = service.NewUserService(c.UserRepo, c.JWTService)
    c.ApiKeyService = service.NewApiKeyService(c.ApiKeyRepo)
}
```

## API Key Database Schema

```sql
CREATE TABLE `api_key` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(256) NOT NULL,
  `description` text DEFAULT NULL,
  `api_key` varchar(32) NOT NULL,
  `auth_key` varchar(64) NOT NULL,
  `status` varchar(15) NOT NULL DEFAULT 'active',
  `h2h` char(1) NOT NULL DEFAULT 'N',
  `last_access` datetime DEFAULT NULL,
  `ip_whitelist` varchar(256) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `created_by` int(11) unsigned NOT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` int(11) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_api_key` (`api_key`)
);
```

## Security Features

1. **API Key Validation**: All requests must have valid API key
2. **IP Whitelisting**: Optional IP restriction per API key
3. **H2H Support**: Host-to-Host capability flags
4. **Access Logging**: Last access time tracking
5. **JWT Expiry**: Configurable token expiration
6. **Token Refresh**: Secure token refresh mechanism

## Error Responses

All middleware returns standardized error responses:

```json
{
  "success": false,
  "message": "API key is required",
  "error": "authentication_error"
}
```

Common error scenarios:
- Missing API key
- Invalid/inactive API key  
- IP not whitelisted
- Missing JWT token (private endpoints)
- Invalid/expired JWT token
- Internal server errors
