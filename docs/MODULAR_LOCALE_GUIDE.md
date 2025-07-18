# Modular Locale Structure Guide

## ✅ **IMPLEMENTED**: Modular Translation Files by Module/Feature

Sekarang project sudah support **multiple locale JSON files yang diorganisir berdasarkan module/feature**!

## 📁 New Locale Structure

```
locales/
├── en/                    # English translations
│   ├── common.json        # Common errors, validations, fields, API messages
│   ├── user.json          # User management translations
│   └── auth.json          # Authentication translations
├── id/                    # Indonesian translations  
│   ├── common.json        # Error umum, validasi, field, pesan API
│   ├── user.json          # Terjemahan manajemen pengguna
│   └── auth.json          # Terjemahan autentikasi
├── es/                    # Spanish translations
│   ├── common.json        # Errores comunes, validaciones, campos, mensajes API
│   ├── user.json          # Traducciones de gestión de usuarios
│   └── auth.json          # Traducciones de autenticación
└── [old files]            # Legacy single files (for backward compatibility)
    ├── en.json            # Will be used if modular files don't exist
    ├── id.json            # Fallback to single file
    └── es.json            # Backward compatibility
```

## 🔧 Updated Configuration

### I18n Manager Configuration
```go
// internal/application/container.go
i18nConfig := i18n.Config{
    DefaultLanguage: "en",
    LocalesPath:     "./locales",
    SupportedLangs:  []string{"en", "id", "es"}, // Multiple languages
    Modules:         []string{"common", "user", "auth"}, // Modular files
}
```

### Smart Loading Strategy
1. **Try legacy single file first** (en.json, id.json, es.json)
2. **If not found, load modular files** (en/common.json, en/user.json, etc.)
3. **Backward compatibility** maintained
4. **Graceful degradation** if files missing

## 📄 Module Breakdown

### 1. **common.json** - Shared Translations
- `error.*` - General API errors (validation, database, auth, etc.)
- `validation.*` - Validation rules and messages
- `field.*` - Form field labels
- `message.*` - General system messages
- `api.*` - API operation messages

**Keys Examples:**
```json
{
  "error.validation_failed": "Validation failed",
  "error.internal_server": "Internal server error", 
  "error.unauthorized": "Unauthorized. Invalid API key",
  "validation.required": "{{.Field}} is required",
  "field.username": "Username",
  "message.welcome": "Welcome to our API"
}
```

### 2. **user.json** - User Management
- `error.user_*` - User-specific errors
- `success.user_*` - User operation success messages
- User CRUD operations, profile management

**Keys Examples:**
```json
{
  "error.user_not_found": "User not found",
  "error.username_exists": "Username already exists",
  "success.user_created": "User created successfully",
  "success.users_retrieved": "Users retrieved successfully"
}
```

### 3. **auth.json** - Authentication & Authorization
- `error.reset_*` - Password reset errors
- `success.login` - Authentication success messages
- `auth.*` - Detailed auth messages
- Registration, login, logout, token management

**Keys Examples:**
```json
{
  "error.reset_token_invalid": "Invalid or expired reset token",
  "success.login": "Login successful",
  "auth.login_success": "Welcome back! You have successfully logged in"
}
```

## 🚀 Usage (No Changes Required!)

**API handlers work exactly the same:**

```go
// All these still work without any code changes!

// Common errors
return response.ErrorWithI18n(c, fiber.StatusBadRequest, "validation_failed", nil)

// User-specific messages  
return response.SuccessWithI18n(c, "user_retrieved", user, nil)

// Auth messages
return response.SuccessWithI18n(c, "login", tokenData, nil)
```

**Translation keys automatically found across all module files!**

## ✨ Benefits of Modular Structure

### 1. **Better Organization**
- Related translations grouped together
- Easy to find specific module translations
- Clear separation of concerns

### 2. **Team Collaboration**
- Different teams can work on different modules
- Less merge conflicts in translation files
- Easier to assign translation tasks

### 3. **Maintainability**
- Smaller, focused files are easier to manage
- Quick identification of missing translations per module
- Easier to add new modules

### 4. **Scalability**
- Easy to add new modules (product.json, order.json, etc.)
- Can add new languages by copying folder structure
- Supports unlimited modules and languages

### 5. **Backward Compatibility**
- Legacy single files still supported
- Gradual migration possible
- No breaking changes

## 📋 File Organization Guide

### When to Create New Module Files:

1. **Large Feature Sets** (20+ translation keys)
   ```
   locales/en/product.json    # Product management
   locales/en/order.json      # Order processing  
   locales/en/payment.json    # Payment handling
   ```

2. **Domain-Specific Translations**
   ```
   locales/en/admin.json      # Admin panel specific
   locales/en/public.json     # Public website
   locales/en/mobile.json     # Mobile app specific
   ```

3. **Business Logic Modules**
   ```
   locales/en/inventory.json  # Inventory management
   locales/en/reports.json    # Reporting system
   locales/en/notifications.json # Notification templates
   ```

### File Naming Convention:
- **Use lowercase** - `user.json`, not `User.json`
- **Use singular nouns** - `user.json`, not `users.json` 
- **Be descriptive** - `auth.json` better than `a.json`
- **Group related features** - `user-profile.json` for complex modules

## 🔄 Migration from Single Files

### Automatic Migration Strategy:
1. **Current single files kept** for backward compatibility
2. **New modular files take precedence** when available
3. **Gradual migration possible** - can migrate one language at a time

### Manual Migration Steps:
```bash
# 1. Keep old files as backup
mv locales/en.json locales/en-backup.json

# 2. New structure automatically works
# Files already created in locales/en/, locales/id/, locales/es/

# 3. Test application
go run cmd/api/main.go

# 4. Remove backup after confirming everything works
rm locales/en-backup.json
```

## 🌍 Adding New Languages

### To add French (fr):
```bash
# 1. Create directory
mkdir locales/fr

# 2. Copy and translate module files
cp locales/en/common.json locales/fr/common.json
cp locales/en/user.json locales/fr/user.json  
cp locales/en/auth.json locales/fr/auth.json

# 3. Translate content in each file
# Edit locales/fr/*.json files

# 4. Update configuration
SupportedLangs: []string{"en", "id", "es", "fr"}
```

### To add new modules:
```bash
# 1. Create module file for all languages
touch locales/en/product.json
touch locales/id/product.json
touch locales/es/product.json

# 2. Add module to configuration
Modules: []string{"common", "user", "auth", "product"}

# 3. Add translations to each file
```

## 🎯 Current Implementation Status

### ✅ **Completed:**
- [x] Modular file structure created
- [x] English translations (common, user, auth)
- [x] Indonesian translations (common, user, auth)  
- [x] Spanish translations (common, user, auth)
- [x] Updated i18n manager for modular loading
- [x] Backward compatibility maintained
- [x] Application configuration updated
- [x] Build and functionality tested

### 🔄 **Available for Extension:**
- [ ] Additional modules (product, order, admin, etc.)
- [ ] More languages (French, German, Japanese, etc.)
- [ ] Advanced features (pluralization, gender, etc.)

## 🎉 Summary

**Perfect! Modular locale structure is now implemented:**

1. **Organized by modules** - common, user, auth
2. **Multiple languages support** - en, id, es  
3. **Backward compatible** - legacy files still work
4. **Easy to extend** - add modules or languages easily
5. **Zero code changes** - existing API handlers work unchanged
6. **Better maintainability** - smaller, focused files

**The system now supports exactly what you requested: per-module translation files organized in language folders!** 🚀
