# JWT Authentication Flow Documentation

## 🔄 **Authentication Flow: Public First, Then Private**

### **Konsep Utama:**
1. **Public Token** - Untuk akses endpoint public (login, register, dll)
2. **Private Token** - Untuk akses endpoint private (user profile, dll)

### **Alur Penggunaan:**

```
1. Client                    2. Server
   |                            |
   |-- Get API Key ------------>|
   |<-- API Key + Public Token--|
   |                            |
   |-- Login (Public Token) --->|
   |<-- Private Token -----------|
   |                            |
   |-- Access Private (Private Token) ->|
   |<-- Protected Data ----------|
```

## 📋 **Step by Step Implementation**

### **Step 1: Generate Public Token**
```bash
# Client memerlukan API Key terlebih dahulu
go run jwt_tester.go
```

**Public Token Contains:**
```json
{
  "api_key_id": 1,
  "api_key_name": "client-app",
  "iss": "go-rest-api",
  "sub": "public-access",
  "exp": 1234567890
}
```

### **Step 2: Login Process**
```bash
curl -X POST http://localhost:8080/api/v1/public/auth/login \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -H "Authorization: Bearer PUBLIC_JWT_TOKEN" \
  -d '{"username":"user","password":"pass"}'
```

**Response (Private Token):**
```json
{
  "success": true,
  "data": {
    "user": {...},
    "token": "PRIVATE_JWT_TOKEN"
  }
}
```

**Private Token Contains:**
```json
{
  "api_key_id": 1,
  "api_key_name": "client-app",
  "user_id": 999,
  "username": "user",
  "email": "user@example.com",
  "iss": "go-rest-api",
  "sub": "private-access",
  "exp": 1234567890
}
```

### **Step 3: Access Private Endpoints**
```bash
curl -X GET http://localhost:8080/api/v1/private/users/profile \
  -H "X-API-Key: your-api-key" \
  -H "Authorization: Bearer PRIVATE_JWT_TOKEN"
```

## 🔒 **Security Features**

### **Token Validation:**
- ✅ Public token dapat digunakan untuk endpoint public
- ✅ Private token dapat digunakan untuk endpoint private
- ❌ Public token TIDAK dapat digunakan untuk endpoint private
- ❌ Private token TIDAK perlu digunakan untuk endpoint public

### **API Key Integration:**
- Semua request memerlukan API Key (`X-API-Key` header)
- JWT token harus sesuai dengan API Key yang digunakan
- Validasi IP whitelist pada API Key

### **Token Expiration:**
- Public Token: 2 jam
- Private Token: 24 jam

## 🧪 **Manual Testing**

### **1. Generate Test Tokens:**
```bash
go run jwt_tester.go
```

### **2. Test dengan cURL:**
```bash
./manual_test_guide.sh
```

### **3. Test dengan Postman:**
Import: `jwt_testing.postman_collection.json`

### **4. Debug Token Online:**
- URL: https://jwt.io/
- Secret: `test-secret-key-for-manual-testing`

## 🚀 **Quick Start Testing**

```bash
# 1. Generate tokens
go run jwt_tester.go

# 2. Start server
go run ./cmd/api/ --port=8080

# 3. Copy PUBLIC token and test login
curl -X POST http://localhost:8080/api/v1/public/auth/login \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -H "Authorization: Bearer PUBLIC_TOKEN_HERE" \
  -d '{"username":"testuser","password":"password"}'

# 4. Copy PRIVATE token from response and test private endpoint
curl -X GET http://localhost:8080/api/v1/private/users/profile \
  -H "X-API-Key: your-api-key" \
  -H "Authorization: Bearer PRIVATE_TOKEN_HERE"
```

## ✅ **Test Scenarios**

1. **Happy Path:**
   - Generate public token ✅
   - Login with public token → get private token ✅
   - Access private endpoint with private token ✅

2. **Security Tests:**
   - Use public token on private endpoint → should fail ❌
   - Use wrong API key → should fail ❌
   - Use expired token → should fail ❌
   - No API key → should fail ❌

## 🎯 **Integration Points**

- **Middleware:** Validates tokens and API keys
- **Container:** Manages service dependencies
- **Handler:** Implements login flow and token generation
- **Service:** Business logic for token validation
