# 🛡️ Production-Ready Validator Guide

## 📋 Overview

This project uses **`github.com/go-playground/validator/v10`** - a production-ready, comprehensive validation library that provides:

- ✅ **80+ built-in validation rules**
- ✅ **Custom error messages**
- ✅ **Structured error responses**
- ✅ **Fiber framework integration**
- ✅ **JSON field name mapping**
- ✅ **Advanced validation rules**

## 🚀 Features

### **Built-in Validation Rules (No Manual Implementation Needed)**

#### **String Validations**
```go
type Example struct {
    Username string `validate:"required,min=3,max=50,alphanum"`
    Email    string `validate:"required,email"`
    Password string `validate:"required,min=8,max=100"`
    Status   string `validate:"oneof=active inactive suspended"`
    URL      string `validate:"url"`
    Phone    string `validate:"e164"`  // International phone format
}
```

#### **Numeric Validations**
```go
type Example struct {
    Age    int     `validate:"min=0,max=120"`
    Price  float64 `validate:"gt=0"`
    Rating int     `validate:"gte=1,lte=5"`
}
```

#### **Advanced Validations**
```go
type ChangePassword struct {
    CurrentPassword string `validate:"required"`
    NewPassword     string `validate:"required,min=8,nefield=CurrentPassword"`
    ConfirmPassword string `validate:"required,eqfield=NewPassword"`
}
```

### **Supported Validation Tags**

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field is mandatory | `validate:"required"` |
| `email` | Valid email format | `validate:"email"` |
| `min=n` | Minimum length/value | `validate:"min=3"` |
| `max=n` | Maximum length/value | `validate:"max=50"` |
| `len=n` | Exact length | `validate:"len=10"` |
| `alpha` | Only alphabetic chars | `validate:"alpha"` |
| `alphanum` | Alphanumeric only | `validate:"alphanum"` |
| `numeric` | Numbers only | `validate:"numeric"` |
| `oneof` | One of specific values | `validate:"oneof=red blue green"` |
| `url` | Valid URL format | `validate:"url"` |
| `uri` | Valid URI format | `validate:"uri"` |
| `base64` | Valid base64 encoding | `validate:"base64"` |
| `jwt` | Valid JWT token | `validate:"jwt"` |
| `uuid` | Valid UUID format | `validate:"uuid"` |
| `ip` | Valid IP address | `validate:"ip"` |
| `mac` | Valid MAC address | `validate:"mac"` |
| `latitude` | Valid latitude | `validate:"latitude"` |
| `longitude` | Valid longitude | `validate:"longitude"` |

### **Comparison Tags**
| Tag | Description | Example |
|-----|-------------|---------|
| `eqfield` | Equal to another field | `validate:"eqfield=Password"` |
| `nefield` | Not equal to another field | `validate:"nefield=OldPassword"` |
| `gtfield` | Greater than another field | `validate:"gtfield=StartDate"` |
| `ltefield` | Less than or equal to field | `validate:"ltefield=EndDate"` |

## 🎯 Usage Examples

### **1. Basic Model Validation**

```go
// internal/model/user_model.go
type UserCreateRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
    Email    string `json:"email" validate:"required,email,max=100"`
    Password string `json:"password" validate:"required,min=8,max=100"`
}

func (r *UserCreateRequest) Validate() error {
    return validator.ValidateStruct(r)
}
```

### **2. Handler Usage**

```go
// internal/handler/user_handler.go
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var req model.UserCreateRequest
    
    if err := c.BodyParser(&req); err != nil {
        return h.responseHelper.ErrorWithI18n(c, 400, "invalid_request_body", nil)
    }
    
    // Validation with structured errors
    if err := req.Validate(); err != nil {
        validationErrors := validator.GetValidationErrors(err)
        return c.Status(400).JSON(fiber.Map{
            "success": false,
            "message": "Validation failed",
            "errors":  validationErrors,
        })
    }
    
    // Continue with business logic...
}
```

### **3. Structured Error Response**

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": [
    {
      "field": "username",
      "message": "username must be at least 3 characters long",
      "tag": "min",
      "param": "3",
      "value": "ab"
    },
    {
      "field": "email", 
      "message": "email must be a valid email address",
      "tag": "email",
      "param": "",
      "value": "invalid-email"
    }
  ]
}
```

## 🔧 Production Features

### **1. JSON Field Name Mapping**
Automatically uses JSON tag names in error messages:
```go
type User struct {
    Username string `json:"username" validate:"required"`  // Error shows "username", not "Username"
}
```

### **2. Fiber Integration**
```go
// Use the built-in Fiber error handler
return validator.FiberValidationErrorHandler(c, err)
```

### **3. Multiple Error Formats**
```go
// Get structured errors
errors := validator.GetValidationErrors(err)

// Get single string
message := validator.FormatValidationErrorsAsString(err)
```

## 📋 Common Validation Patterns

### **User Registration**
```go
type UserCreateRequest struct {
    Username        string `json:"username" validate:"required,min=3,max=30,alphanum"`
    Email           string `json:"email" validate:"required,email,max=100"`
    Password        string `json:"password" validate:"required,min=8,max=100"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
    Age             int    `json:"age" validate:"gte=13,lte=120"`
    Website         string `json:"website" validate:"omitempty,url"`
}
```

### **Business Data**
```go
type ProductRequest struct {
    Name        string  `json:"name" validate:"required,min=2,max=100"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Category    string  `json:"category" validate:"required,oneof=electronics clothing books"`
    SKU         string  `json:"sku" validate:"required,alphanum,len=8"`
    Description string  `json:"description" validate:"max=500"`
    Weight      float64 `json:"weight" validate:"gt=0"`
}
```

### **API Configuration**
```go
type APIConfig struct {
    BaseURL     string `json:"base_url" validate:"required,url"`
    APIKey      string `json:"api_key" validate:"required,min=32,max=64,alphanum"`
    Timeout     int    `json:"timeout" validate:"min=1,max=300"`
    RetryCount  int    `json:"retry_count" validate:"min=0,max=10"`
    Environment string `json:"environment" validate:"oneof=development staging production"`
}
```

## 🎉 Benefits

### ✅ **No Manual Validation Code**
- No need to write custom validation logic
- 80+ built-in rules cover most use cases
- Comprehensive error messages

### ✅ **Production-Ready**
- Battle-tested library used by thousands of projects
- High performance with minimal overhead
- Thread-safe and concurrent-safe

### ✅ **Fiber Compatible**
- Built-in error handlers for Fiber
- JSON response formatting
- Middleware integration support

### ✅ **Maintainable**
- Declarative validation with struct tags
- Clear error messages for debugging
- Consistent validation across the application

### ✅ **Extensible**
- Custom validation functions if needed
- Custom error messages
- Conditional validation rules

## 🚀 Testing Your Validator

Run the validator test:
```bash
go run ./cmd/validator-test/
```

This will test all validation rules and show structured error responses.

**Your application now has enterprise-grade validation! 🛡️**
