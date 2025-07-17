# 🌍 Multilingual Support (i18n) Guide

## 📁 i18n Structure

```
pkg/i18n/
├── manager.go          # i18n manager implementation
locales/
├── en.json            # English translations
├── id.json            # Indonesian translations
internal/middleware/
├── i18n.go            # Language detection middleware
```

## 🎯 Features

### ✅ **Language Detection**
1. **Query Parameter**: `?lang=id`
2. **Accept-Language Header**: `Accept-Language: id,en;q=0.9`
3. **Default Fallback**: English (`en`)

### ✅ **Supported Languages**
- 🇺🇸 **English** (`en`) - Default
- 🇮🇩 **Indonesian** (`id`) - Bahasa Indonesia

### ✅ **Translation Categories**
- **Errors**: `error.user_not_found`, `error.validation_failed`
- **Success**: `success.user_created`, `success.user_updated`
- **Validation**: `validation.required`, `validation.min_length`
- **Fields**: `field.username`, `field.email`

## 🚀 Usage Examples

### **1. API Request with Language**

```bash
# Using query parameter
curl -H "x-api-key: test-api-key" \
     "http://localhost:8080/api/v1/users?lang=id"

# Using Accept-Language header
curl -H "x-api-key: test-api-key" \
     -H "Accept-Language: id" \
     "http://localhost:8080/api/v1/users"
```

### **2. Response Examples**

**English Response:**
```json
{
  "success": false,
  "message": "User not found",
  "error": "user_not_found"
}
```

**Indonesian Response:**
```json
{
  "success": false,
  "message": "Pengguna tidak ditemukan", 
  "error": "user_not_found"
}
```

### **3. In Handler Code**

```go
// Success response with i18n
return h.responseHelper.SuccessWithI18n(c, "user_created", userData, nil)

// Error response with i18n
return h.responseHelper.ErrorWithI18n(c, 404, "user_not_found", nil)

// With template data
return h.responseHelper.ErrorWithI18n(c, 400, "username_min_length", 
    map[string]interface{}{
        "MinLength": 3,
    })
```

## 📝 Adding New Languages

### **Step 1: Create Language File**
```bash
# Create new language file
touch locales/es.json  # Spanish
```

### **Step 2: Add Translations**
```json
{
  "error.user_not_found": "Usuario no encontrado",
  "success.user_created": "Usuario creado exitosamente",
  // ... more translations
}
```

### **Step 3: Update Configuration**
```go
// In container.go
i18nConfig := i18n.Config{
    DefaultLanguage: "en",
    LocalesPath:     "./locales",
    SupportedLangs:  []string{"en", "id", "es"}, // Add 'es'
}

// In rest_api.go middleware
app.Use(middleware.I18nMiddleware(middleware.I18nConfig{
    DefaultLanguage: "en",
    SupportedLangs:  []string{"en", "id", "es"}, // Add 'es'
}))
```

## 🔧 Translation Keys Naming Convention

### **Error Messages**
- `error.{category}_{specific}`: `error.user_not_found`
- `error.validation_{rule}`: `error.validation_required`

### **Success Messages**
- `success.{entity}_{action}`: `success.user_created`

### **Validation Messages**
- `validation.{rule}`: `validation.required`, `validation.min_length`

### **Field Names**
- `field.{field_name}`: `field.username`, `field.email`

## 📊 Template Data Usage

```go
// Using template data for dynamic content
templateData := map[string]interface{}{
    "MinLength": 6,
    "MaxLength": 50,
    "Field": "password",
}

return h.responseHelper.ErrorWithI18n(c, 400, "password_min_length", templateData)
```

**Translation with template:**
```json
{
  "error.password_min_length": "Password must be at least {{.MinLength}} characters"
}
```

**Result:**
- English: "Password must be at least 6 characters"
- Indonesian: "Password harus minimal 6 karakter"

## 🧪 Testing i18n

### **Test Different Languages**
```bash
# Test English (default)
curl -H "x-api-key: test-api-key" \
     "http://localhost:8080/api/v1/users"

# Test Indonesian
curl -H "x-api-key: test-api-key" \
     "http://localhost:8080/api/v1/users?lang=id"

# Test with Accept-Language header
curl -H "x-api-key: test-api-key" \
     -H "Accept-Language: id,en;q=0.9" \
     "http://localhost:8080/api/v1/users"
```

## 🎯 Best Practices

1. **Always provide fallback**: English as default language
2. **Consistent key naming**: Use clear, hierarchical naming
3. **Template data**: Use for dynamic content like field names, lengths
4. **Translation completeness**: Ensure all keys exist in all language files
5. **Context-aware**: Consider cultural context in translations

## 🚀 Production Ready

- ✅ **Performance**: Translations loaded once at startup
- ✅ **Memory efficient**: Localizers cached per language
- ✅ **Fallback safe**: Always returns meaningful text
- ✅ **Extensible**: Easy to add new languages
- ✅ **Standards compliant**: Follows i18n best practices

**Your API now speaks multiple languages! 🌍**
