# User API Endpoints Documentation

## Base URL
```
http://localhost:8080/api/v1/users
```

## Authentication
All endpoints require API key in header:
```
x-api-key: test-api-key
```

## Endpoints

### 1. Create User
**POST** `/api/v1/users`

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "status": "active",
    "created_at": "2025-07-16T10:00:00Z",
    "updated_at": "2025-07-16T10:00:00Z"
  }
}
```

### 2. Get User by ID
**GET** `/api/v1/users/:id`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "status": "active",
    "created_at": "2025-07-16T10:00:00Z",
    "updated_at": "2025-07-16T10:00:00Z"
  }
}
```

### 3. Get All Users (with Pagination)
**GET** `/api/v1/users?page=1&limit=10`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "id": 1,
        "username": "john_doe",
        "email": "john@example.com",
        "status": "active",
        "created_at": "2025-07-16T10:00:00Z",
        "updated_at": "2025-07-16T10:00:00Z"
      }
    ],
    "pagination": {
      "current_page": 1,
      "total_pages": 1,
      "total_count": 1,
      "limit": 10
    }
  }
}
```

### 4. Update User
**PUT** `/api/v1/users/:id`

**Request Body (all fields optional):**
```json
{
  "username": "john_updated",
  "email": "john_updated@example.com",
  "password": "newpassword123",
  "status": "active"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "username": "john_updated",
    "email": "john_updated@example.com",
    "status": "active",
    "created_at": "2025-07-16T10:00:00Z",
    "updated_at": "2025-07-16T10:01:00Z"
  }
}
```

### 5. Delete User (Soft Delete)
**DELETE** `/api/v1/users/:id`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User deleted successfully",
  "data": null
}
```

## Error Responses

### Validation Error (400 Bad Request)
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "username must be at least 3 characters"
}
```

### Not Found (404 Not Found)
```json
{
  "success": false,
  "message": "User not found",
  "error": "user not found"
}
```

### Server Error (500 Internal Server Error)
```json
{
  "success": false,
  "message": "Failed to create user",
  "error": "database connection failed"
}
```

## Business Rules

1. **Username**: Minimum 3 characters, maximum 50 characters, must be unique
2. **Email**: Must be valid email format, must be unique  
3. **Password**: Minimum 6 characters, automatically hashed using bcrypt
4. **Status**: Only accepts "active" or "inactive" values
5. **Soft Delete**: Users are not physically deleted, only marked as deleted

## Testing with cURL

### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "x-api-key: test-api-key" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Get All Users
```bash
curl -X GET "http://localhost:8080/api/v1/users?page=1&limit=5" \
  -H "x-api-key: test-api-key"
```

### Get User by ID
```bash
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "x-api-key: test-api-key"
```

### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -H "x-api-key: test-api-key" \
  -d '{
    "username": "updated_user",
    "status": "inactive"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "x-api-key: test-api-key"
```
