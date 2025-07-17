# Dependency Injection & Route Management Guide

## 📁 Current Structure

```
internal/
├── application/
│   ├── container.go      # 🆕 Dependency injection container
│   └── rest_api.go       # 🔧 Minimal main application setup
└── routes/
    ├── route_manager.go  # 🆕 Centralized route management
    └── user_routes.go    # Existing user routes
```

## 🎯 Benefits

### ✅ Before (Tightly Coupled)
```go
// All dependencies and routes mixed in rest_api.go
userRepo := repository.NewUserRepository(db)
userService := service.NewUserService(userRepo)
userHandler := handler.NewUserHandler(userService)
routes.SetupUserRoutes(app, userHandler)
```

### ✅ After (Clean Separation)
```go
// Clean separation with minimal changes needed
container := NewContainer(db)
routeConfig := &routes.RouteConfig{
    UserHandler: container.UserHandler,
}
routes.SetupAllRoutes(app, routeConfig)
```

## 🚀 Adding New Modules

### Step 1: Add to Container struct
```go
type Container struct {
    // Existing
    UserHandler *handler.UserHandler
    
    // 🆕 New Product Module
    ProductRepo    repository.ProductRepository
    ProductService usecase.ProductUsecase
    ProductHandler *handler.ProductHandler
}
```

### Step 2: Add to RouteConfig
```go
type RouteConfig struct {
    UserHandler    *handler.UserHandler
    ProductHandler *handler.ProductHandler  // 🆕 Add this
}
```

### Step 3: Update Container initialization
```go
func (c *Container) initProductModule() {
    c.ProductRepo = repositoryImpl.NewProductRepository(c.DB)
    c.ProductService = service.NewProductService(c.ProductRepo)
    c.ProductHandler = handler.NewProductHandler(c.ProductService)
}
```

### Step 4: Add route setup
```go
// In route_manager.go
func SetupAllRoutes(app *fiber.App, config *RouteConfig) {
    setupUserRoutes(app, config.UserHandler)
    setupProductRoutes(app, config.ProductHandler)  // 🆕 Add this
}

func setupProductRoutes(app *fiber.App, handler *handler.ProductHandler) {
    SetupProductRoutes(app, handler)  // Call existing product routes
}
```

### Step 5: Update rest_api.go (MINIMAL CHANGE)
```go
routeConfig := &routes.RouteConfig{
    UserHandler:    container.UserHandler,
    ProductHandler: container.ProductHandler,  // 🆕 Just add this line
}
```

## 📈 Scalability Benefits

### ✅ **rest_api.go stays minimal**
- Only 1-2 lines added per new module
- No complex dependency setup
- Clear separation of concerns

### ✅ **Centralized management**
- All dependencies in `container.go`
- All routes in `route_manager.go`
- Easy to maintain and debug

### ✅ **Future-ready structure**
```go
// Easy to add advanced features later:
// - API versioning (/api/v1, /api/v2)
// - Route-specific middleware
// - Module-based rate limiting
// - Conditional route enabling
```

## 🔗 Cross-Module Dependencies

```go
// Example: OrderService needs UserService and ProductService
func (c *Container) initOrderModule() {
    c.OrderRepo = repositoryImpl.NewOrderRepository(c.DB)
    c.OrderService = service.NewOrderService(
        c.OrderRepo,
        c.UserService,    // Cross-dependency
        c.ProductService, // Cross-dependency
    )
    c.OrderHandler = handler.NewOrderHandler(c.OrderService)
}
```

## 🧪 Testing Benefits

```go
// Easy to test individual modules
func TestUserModule() {
    mockDB := &MockDB{}
    container := NewContainer(mockDB)
    
    routeConfig := &routes.RouteConfig{
        UserHandler: container.UserHandler,
    }
    
    app := fiber.New()
    routes.SetupAllRoutes(app, routeConfig)
    
    // Test specific routes
}
```

## 🎯 Final Architecture

```
Request Flow:
HTTP Request → rest_api.go → route_manager.go → specific_routes.go → handler → service → repository → database

Dependency Flow:
container.go → Initialize all dependencies → Pass to route_manager.go → Setup all routes
```

**Result**: Adding new modules requires minimal changes to main files! 🎉
