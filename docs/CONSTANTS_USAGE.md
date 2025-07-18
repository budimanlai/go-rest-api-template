# Constants Documentation

## Overview
This package contains all constants used throughout the application to ensure consistency, maintainability, and avoid magic strings/numbers.

## File Structure

```
internal/constant/
├── constant.go    # General application constants
├── database.go    # Database-specific constants
└── http.go        # HTTP-related constants
```

## Usage Examples

### 1. Using Database Constants

**Before (with magic strings):**
```go
// ❌ Bad - magic strings
query := `INSERT INTO user (username, auth_key, email, password_hash, status, created_by, created_at, updated_at) 
          VALUES (:username, :auth_key, :email, :password_hash, :status, :created_by, NOW(), NOW())`

authKey := common.GenerateRandomString(32)
status := "active"
```

**After (with constants):**
```go
// ✅ Good - using constants
import "go-rest-api-template/internal/constant"

result, err := r.db.NamedExecContext(ctx, constant.QueryInsertUser, userModel)

authKey := common.GenerateRandomString(constant.AuthKeyLength)
status := constant.UserStatusActive
```

### 2. Using HTTP Constants

**Before:**
```go
// ❌ Bad - magic numbers and strings
return response.ErrorWithI18n(c, 400, "invalid_request", nil)
return response.ErrorWithI18n(c, 409, "username_exists", nil)
return response.CreatedWithI18n(c, "user_created", userResponse, nil)
```

**After:**
```go
// ✅ Good - using constants
return response.ErrorWithI18n(c, constant.StatusBadRequest, constant.ErrInvalidRequest, nil)
return response.ErrorWithI18n(c, constant.StatusConflict, constant.ErrUsernameExists, nil)
return response.CreatedWithI18n(c, constant.MsgUserCreated, userResponse, nil)
```

### 3. Using Validation Constants

**Before:**
```go
// ❌ Bad - magic numbers in validation tags
type UserCreateRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
    Email    string `json:"email" validate:"required,email,max=100"`
    Password string `json:"password" validate:"required,min=8,max=100"`
}
```

**After:**
```go
// ✅ Good - but validation tags still need to use literals
// However, you can use constants in code validation:
func (r *UserCreateRequest) Validate() error {
    if len(r.Username) < constant.MinUsernameLength || len(r.Username) > constant.MaxUsernameLength {
        return fmt.Errorf("username must be between %d and %d characters", 
            constant.MinUsernameLength, constant.MaxUsernameLength)
    }
    return validator.ValidateStruct(r)
}
```

### 4. Using Default Values

**Before:**
```go
// ❌ Bad - magic values
user.SetCreatedBy(0)
limit := 10
authKeyLength := 32
```

**After:**
```go
// ✅ Good - using constants
user.SetCreatedBy(constant.DefaultCreatedBy)
limit := constant.DefaultPageLimit
authKeyLength := constant.AuthKeyLength
```

## Constant Categories

### 1. Database Constants (`database.go`)

#### Query Constants
- `QueryInsertUser` - User insertion query
- `QuerySelectUserByID` - Get user by ID
- `QuerySelectUserByEmail` - Get user by email
- `QueryUpdateUser` - Update user information
- And more...

#### Field Constants
- `FieldUserID` - "id"
- `FieldUserUsername` - "username"
- `FieldUserEmail` - "email"
- And more...

### 2. HTTP Constants (`http.go`)

#### Status Codes
- `StatusOK` - 200
- `StatusCreated` - 201
- `StatusBadRequest` - 400
- `StatusUnauthorized` - 401
- And more...

#### Route Patterns
- `RouteAuthLogin` - "/auth/login"
- `RouteUsers` - "/users"
- `RouteUserByID` - "/users/:id"
- And more...

### 3. General Constants (`constant.go`)

#### User Status
- `UserStatusActive` - "active"
- `UserStatusInactive` - "inactive"
- `UserStatusSuspended` - "suspended"

#### Messages
- `MsgUserCreated` - "user_created"
- `MsgUserRetrieved` - "user_retrieved"
- `ErrUserNotFound` - "user_not_found"
- And more...

## Benefits

### 1. **Consistency**
```go
// All user creation uses same default
user.SetCreatedBy(constant.DefaultCreatedBy)

// All queries use same pattern
r.db.GetContext(ctx, &userModel, constant.QuerySelectUserByID, id)
```

### 2. **Maintainability**
```go
// Change once, applies everywhere
const DefaultPageLimit = 20 // Changed from 10 to 20
```

### 3. **Type Safety**
```go
// Prevents typos
status := constant.UserStatusActive  // ✅ Compile-time checked
status := "activ"                   // ❌ Runtime error
```

### 4. **Documentation**
```go
// Self-documenting code
authKey := common.GenerateRandomString(constant.AuthKeyLength) // Clear intent
authKey := common.GenerateRandomString(32)                    // Magic number
```

## Best Practices

### 1. **Grouping Related Constants**
```go
// User Status Constants
const (
    UserStatusActive    = "active"
    UserStatusInactive  = "inactive"
    UserStatusSuspended = "suspended"
)
```

### 2. **Descriptive Names**
```go
// ✅ Good
const DefaultPageLimit = 10
const MaxPasswordLength = 100

// ❌ Bad
const DPL = 10
const MPL = 100
```

### 3. **Logical Organization**
- `constant.go` - General application constants
- `database.go` - Database-specific constants
- `http.go` - HTTP/API-specific constants

### 4. **Import Alias for Clarity**
```go
import "go-rest-api-template/internal/constant"

// Use with clear namespace
status := constant.UserStatusActive
query := constant.QuerySelectUserByID
```

## Migration Strategy

### Step 1: Replace Magic Strings in Handlers
```go
// Before
return response.ErrorWithI18n(c, 400, "invalid_request", nil)

// After
return response.ErrorWithI18n(c, constant.StatusBadRequest, constant.ErrInvalidRequest, nil)
```

### Step 2: Replace Magic Strings in Repository
```go
// Before
query := `SELECT * FROM user WHERE id = ? AND deleted_at IS NULL`

// After
err := r.db.GetContext(ctx, &userModel, constant.QuerySelectUserByID, id)
```

### Step 3: Replace Magic Values
```go
// Before
user.SetCreatedBy(0)

// After
user.SetCreatedBy(constant.DefaultCreatedBy)
```

## IDE Benefits

### 1. **Auto-completion**
```go
constant.User... // IDE shows: UserStatusActive, UserStatusInactive, etc.
constant.Query... // IDE shows: QueryInsertUser, QuerySelectUserByID, etc.
```

### 2. **Find All References**
- Easy to find where constants are used
- Safe refactoring with IDE support

### 3. **Go to Definition**
- Jump to constant definition quickly
- Understand value and context

## Testing with Constants

```go
func TestUserCreation(t *testing.T) {
    user := &entity.User{
        Username: "testuser",
        Email:    "test@example.com",
        Status:   constant.UserStatusActive, // ✅ Clear intent
    }
    
    assert.Equal(t, constant.UserStatusActive, user.Status)
}
```

## Production Benefits

1. **Reduced Bugs**: No more typos in status values or route patterns
2. **Easier Maintenance**: Change constants in one place
3. **Better Code Reviews**: Reviewers can easily understand intent
4. **Consistent API**: All endpoints use same error messages and status codes
5. **Database Safety**: All queries are predefined and consistent

## Future Enhancements

1. **Configuration Integration**: Load some constants from config files
2. **Validation Integration**: Use constants in validation rules
3. **Documentation Generation**: Auto-generate API docs from constants
4. **Internationalization**: Map message constants to multiple languages
