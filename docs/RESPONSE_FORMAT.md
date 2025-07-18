# API Response Format Documentation

## Overview
This document describes the standardized response format used across all API endpoints in the application. All responses follow a consistent structure with `data` and `meta` fields to ensure predictable client integration.

## Standard Response Structure

All API responses use the following format:

```json
{
    "data": <object|array|null>,
    "meta": {
        "success": <boolean>,
        "message": <string>,
        // ... additional metadata based on response type
    }
}
```

### Key Principles
- **Consistent Structure**: All responses use the same `data` and `meta` format
- **Clear Success Indication**: The `meta.success` field always indicates request success/failure
- **Meaningful Messages**: The `meta.message` field provides human-readable feedback
- **Structured Errors**: Validation errors include detailed field-level information

## Response Types

### 1. Success Response (Single Object)

**Status Code**: 200 OK or 201 Created

```json
{
    "data": {
        "id": 1,
        "username": "john_doe",
        "email": "john@example.com",
        "status": "active"
    },
    "meta": {
        "success": true,
        "message": "User retrieved successfully"
    }
}
```

### 2. Success Response (Array)

**Status Code**: 200 OK

```json
{
    "data": [
        {
            "id": 1,
            "username": "john_doe",
            "email": "john@example.com",
            "status": "active"
        },
        {
            "id": 2,
            "username": "jane_doe", 
            "email": "jane@example.com",
            "status": "active"
        }
    ],
    "meta": {
        "success": true,
        "message": "Users retrieved successfully"
    }
}
```

### 3. Success Response (No Data)

**Status Code**: 200 OK

```json
{
    "data": null,
    "meta": {
        "success": true,
        "message": "User deleted successfully"
    }
}
```

### 4. Paginated Response

**Status Code**: 200 OK

```json
{
    "data": [
        // ... array of objects
    ],
    "meta": {
        "success": true,
        "message": "Users retrieved successfully",
        "pagination": {
            "page": 1,
            "limit": 10,
            "total": 100,
            "total_pages": 10
        }
    }
}
```

### 5. Error Response (General)

**Status Code**: 400, 404, 500, etc.

```json
{
    "data": null,
    "meta": {
        "success": false,
        "message": "User not found"
    }
}
```

### 6. Validation Error Response

**Status Code**: 400 Bad Request

```json
{
    "data": null,
    "meta": {
        "success": false,
        "message": "Validation failed. Please check the following fields",
        "errors": {
            "total_errors": 2,
            "validation_errors": [
                {
                    "field": "email",
                    "message": "email must be a valid email address"
                },
                {
                    "field": "password",
                    "message": "password must be at least 8 characters"
                }
            ]
        }
    }
}
```

## Implementation

### Response Helper Functions

The application provides several helper functions in `pkg/response/response.go`:

#### Success Responses
```go
// Send success response with data
response.Success(c, "User retrieved successfully", userData)

// Send created response
response.Created(c, "User created successfully", userData)

// Send paginated response
response.SendPaginated(c, "Users retrieved successfully", usersData, pagination)
```

#### Error Responses
```go
// Send general error
response.BadRequest(c, "Invalid input", "")
response.NotFound(c, "User not found", "")
response.InternalServerError(c, "Database error", "")

// Send validation error (automatically processes validator errors)
response.ValidationErrorResponse(c, "Validation failed", validationErr)
```

#### I18n Responses
```go
// Success with internationalization
response.SuccessWithI18n(c, "user_retrieved", userData, nil)

// Error with internationalization
response.ErrorWithI18n(c, fiber.StatusBadRequest, "invalid_user_id", nil)
```

### Migration from Old Format

If you're migrating from the old response format, here's the comparison:

#### Old Format (Deprecated)
```json
{
    "success": true,
    "message": "User retrieved successfully",
    "data": {
        "id": 1,
        "username": "john_doe"
    }
}
```

#### New Format (Current)
```json
{
    "data": {
        "id": 1,
        "username": "john_doe"
    },
    "meta": {
        "success": true,
        "message": "User retrieved successfully"
    }
}
```

## Benefits

1. **Consistency**: All responses follow the same structure
2. **Extensibility**: Easy to add metadata without breaking existing clients
3. **Client-Friendly**: Clear separation between data and metadata
4. **Error Handling**: Structured validation errors for better UX
5. **Future-Proof**: Can easily add features like caching info, rate limiting, etc.

## Client Integration Examples

### JavaScript/TypeScript
```javascript
// Handle response
fetch('/api/users/1')
  .then(response => response.json())
  .then(result => {
    if (result.meta.success) {
      console.log('User data:', result.data);
    } else {
      console.error('Error:', result.meta.message);
      
      // Handle validation errors
      if (result.meta.errors) {
        result.meta.errors.validation_errors.forEach(error => {
          console.log(`${error.field}: ${error.message}`);
        });
      }
    }
  });
```

### Go Client
```go
type APIResponse struct {
    Data interface{} `json:"data"`
    Meta struct {
        Success bool   `json:"success"`
        Message string `json:"message"`
        Errors  *struct {
            TotalErrors      int `json:"total_errors"`
            ValidationErrors []struct {
                Field   string `json:"field"`
                Message string `json:"message"`
            } `json:"validation_errors"`
        } `json:"errors,omitempty"`
    } `json:"meta"`
}
```

This standardized format ensures all clients can handle responses predictably and provides a solid foundation for API evolution.
