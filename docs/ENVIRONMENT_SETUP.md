# 🚀 Environment Setup Guide

## 📋 Quick Setup for New Developers

### 1. **Clone Repository**
```bash
git clone <repository-url>
cd go-rest_api
go mod tidy
```

### 2. **Environment Configuration**
```bash
# Copy environment template
cp .env.example .env

# Edit .env with your actual values
nano .env  # or use your preferred editor
```

### 3. **Required .env Values**

#### **For Testing:**
- `TEST_API_KEY`: Get from database `api_key` table
- `TEST_BASE_URL`: Usually `http://localhost:8080`

#### **Optional Overrides:**
- Database settings (if not using config.json)
- JWT settings (if not using config.json)
- Server settings

### 4. **Get API Key from Database**

```sql
-- Run this query to see available API keys
SELECT api_key, name, status FROM api_key WHERE status = 'active';

-- Example result:
-- dev_api_key_12345678901234567890 | Development API Key | active
```

### 5. **Example .env File**
```bash
# Testing Configuration
TEST_API_KEY=dev_api_key_12345678901234567890
TEST_BASE_URL=http://localhost:8080

# Optional: Database override
# DATABASE_HOST=127.0.0.1
# DATABASE_PASSWORD=your_password
```

## 🔒 Security Best Practices

### ✅ **Do:**
- Keep `.env` file local only
- Use different API keys for different environments
- Rotate API keys regularly
- Use strong JWT secrets in production

### ❌ **Don't:**
- Commit `.env` file to git
- Share API keys in chat/email
- Use production keys in development
- Hardcode secrets in source code

## 🧪 Testing Setup

### **Run Tests:**
```bash
# Test with environment configuration
make test-env

# CLI authentication test
make test-auth

# All tests
make test
```

### **Expected Output:**
```
✅ Server is running, starting tests...
✅ API Key validation working
⚠️  Registration/Login tests (may need database setup)
```

## 🆘 Troubleshooting

### **Common Issues:**

1. **"Invalid or inactive API key"**
   - Check TEST_API_KEY in .env
   - Verify API key exists in database
   - Ensure API key status is 'active'

2. **".env file not found"**
   - Copy from .env.example
   - Check file location (root directory)

3. **"Server not running"**
   - Start server: `make run`
   - Check port 8080 is available

## 🎯 Ready to Go!

After setup, you should have:
- ✅ .env file with your configuration
- ✅ API key working for tests
- ✅ Database connection
- ✅ Tests passing

**Happy coding!** 🚀
