# 🔧 Testing Tools untuk ApiKeyOnlyMiddleware

## 📁 Struktur Testing Tools

```
├── test/
│   ├── apikey_middleware_test.go  ✅ Unit test untuk ApiKeyOnlyMiddleware
│   ├── integration_test.go        ✅ Integration tests  
│   └── handler_test.go            ✅ Handler tests
├── tools/
│   └── apikey_auth_tester.go      ✅ CLI testing tool
├── postman_collection.json       ✅ Postman collection
└── manual_test_guide.sh          ✅ Shell script testing
```

## 🧪 Testing Options

### 1. **Unit Tests (Go Testing)**
```bash
# Run specific ApiKeyOnlyMiddleware tests
go test ./test -run TestApiKeyOnlyAuthentication -v

# Run all tests
go test ./test -v
```

### 2. **CLI Testing Tool**
```bash
# Run standalone CLI tester
go run tools/apikey_auth_tester.go
```

### 3. **Postman Collection**
```bash
# Import postman_collection.json into Postman
# Collection includes:
# - Register (API Key only)
# - Login (API Key only)  
# - Private endpoints (API Key + Token)
```

### 4. **Shell Script Testing**
```bash
# Run automated shell script tests
chmod +x manual_test_guide.sh
./manual_test_guide.sh
```

## 🔍 Test Coverage

### ✅ **ApiKeyOnlyMiddleware Tests**
- ✅ Register endpoint dengan API Key only
- ✅ Login endpoint dengan API Key only
- ✅ Private endpoints dengan API Key + Private Token
- ✅ Error handling dan validation

### ✅ **Security Tests**
- ✅ API Key validation
- ✅ IP whitelisting check
- ✅ Access logging verification
- ✅ Token format validation

### ✅ **Integration Tests**
- ✅ End-to-end authentication flow
- ✅ Cross-endpoint compatibility
- ✅ Error response formats

## 📊 Test Results Expected

### **Successful Flow:**
```
1. Register (API Key only) → 200 OK + Private Token
2. Login (API Key only) → 200 OK + Private Token  
3. Private Endpoints (API Key + Token) → 200 OK + Data
```

### **Error Cases:**
```
- Missing API Key → 400 Bad Request
- Invalid API Key → 400 Bad Request
- Missing Token (private endpoints) → 400 Bad Request
- Invalid Token → 400 Bad Request
```

## 🚀 Quick Start Testing

1. **Start the server:**
   ```bash
   make run
   ```

2. **Run CLI tester:**
   ```bash
   go run tools/apikey_auth_tester.go
   ```

3. **Expected output:**
   ```
   === ApiKeyOnlyMiddleware CLI Tester ===
   1. Registering User (API Key Only)...
   ✅ User registered successfully
   
   2. Logging in (API Key Only)...
   ✅ Private Token obtained: eyJhbGciOiJIUzI1NiIs...
   
   3. Testing Private Endpoints (API Key + Private Token)...
   ✅ GET /api/v1/users - Success
   
   🎉 All tests completed successfully!
   ✅ ApiKeyOnlyMiddleware is working correctly
   ```

## 🎯 Troubleshooting

### **Common Issues:**

1. **Server not running:**
   ```bash
   make run  # Start server on port 3000
   ```

2. **Invalid API Key:**
   ```bash
   # Check configs/config.json for valid API keys
   # Default: "test-api-key-12345"
   ```

3. **Database connection:**
   ```bash
   # Ensure PostgreSQL is running
   # Check database configuration
   ```

## 📋 Manual Testing Checklist

- [ ] Register new user dengan API Key only
- [ ] Login existing user dengan API Key only  
- [ ] Access GET /api/v1/users dengan API Key + Token
- [ ] Test invalid API Key responses
- [ ] Test missing token responses
- [ ] Verify IP whitelisting (if configured)
- [ ] Check access logs

## 🏆 Conclusion

Semua testing tools telah diperbaiki dan error redeclaration/unused sudah teratasi:

- ✅ **No compilation errors**
- ✅ **Clean package structure**  
- ✅ **Comprehensive test coverage**
- ✅ **Multiple testing options**
- ✅ **Industry standard compliance**

**Ready for production testing!** 🚀
