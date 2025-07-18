# GoEasyI18n Implementation Example

This document shows how easy it would be to implement `goeasyi18n` as an alternative to our current i18n solution.

## Key Benefits of GoEasyI18n

1. **No external dependencies** - Pure Go implementation
2. **Simple API** - Easy to understand and use
3. **Template support** - Built-in variable interpolation
4. **Consistency checking** - Automatic validation of translation keys
5. **Multiple data sources** - JSON, YAML, strings, files, embed.FS
6. **Language-specific functions** - Create dedicated translate functions per language

## Basic Implementation

```go
package main

import (
    "fmt"
    "github.com/eduardolat/goeasyi18n"
)

func main() {
    // 1. Create i18n instance with configuration
    i18n := goeasyi18n.NewI18n(goeasyi18n.Config{
        FallbackLanguageName:    "en",
        DisableConsistencyCheck: false,
    })

    // 2. Define translations for our current structure
    enTranslations := goeasyi18n.TranslateStrings{
        // Success messages
        {
            Key:     "user_created",
            Default: "User {{.Username}} created successfully",
        },
        {
            Key:     "user_retrieved", 
            Default: "User retrieved successfully",
        },
        {
            Key:     "users_retrieved",
            Default: "{{.Count}} users retrieved successfully",
        },
        // Error messages
        {
            Key:     "user_not_found",
            Default: "User not found",
        },
        {
            Key:     "validation_failed",
            Default: "Validation failed. Please check the following fields",
        },
        // Field names
        {
            Key:     "field.email",
            Default: "email",
        },
        {
            Key:     "field.username", 
            Default: "username",
        },
        // Validation messages
        {
            Key:     "validation.required",
            Default: "{{.Field}} is required",
        },
        {
            Key:     "validation.email",
            Default: "{{.Field}} must be a valid email address",
        },
        {
            Key:     "validation.min",
            Default: "{{.Field}} must be at least {{.MinLength}} characters",
        },
    }

    idTranslations := goeasyi18n.TranslateStrings{
        // Success messages
        {
            Key:     "user_created",
            Default: "Pengguna {{.Username}} berhasil dibuat",
        },
        {
            Key:     "user_retrieved",
            Default: "Pengguna berhasil diambil",
        },
        {
            Key:     "users_retrieved", 
            Default: "{{.Count}} pengguna berhasil diambil",
        },
        // Error messages
        {
            Key:     "user_not_found",
            Default: "Pengguna tidak ditemukan",
        },
        {
            Key:     "validation_failed",
            Default: "Validasi gagal. Silakan periksa field berikut",
        },
        // Field names
        {
            Key:     "field.email",
            Default: "email",
        },
        {
            Key:     "field.username",
            Default: "nama pengguna",
        },
        // Validation messages
        {
            Key:     "validation.required",
            Default: "{{.Field}} wajib diisi",
        },
        {
            Key:     "validation.email",
            Default: "{{.Field}} harus berupa alamat email yang valid",
        },
        {
            Key:     "validation.min",
            Default: "{{.Field}} harus minimal {{.MinLength}} karakter",
        },
    }

    // 3. Add languages to i18n instance
    i18n.AddLanguage("en", enTranslations)
    i18n.AddLanguage("id", idTranslations)

    // 4. Create language-specific translate functions (very convenient!)
    translateEn := i18n.NewLangTranslateFunc("en")
    translateId := i18n.NewLangTranslateFunc("id")

    // 5. Usage examples
    // Simple translation
    fmt.Println("EN:", translateEn("user_retrieved"))
    fmt.Println("ID:", translateId("user_retrieved"))
    
    // With template data
    fmt.Println("EN:", translateEn("user_created", goeasyi18n.Options{
        Data: map[string]any{"Username": "john_doe"},
    }))
    
    // With validation message
    emailField := translateEn("field.email")
    fmt.Println("EN:", translateEn("validation.required", goeasyi18n.Options{
        Data: map[string]any{"Field": emailField},
    }))
    
    // General translate method (similar to current usage)
    message := i18n.T("id", "user_not_found")
    fmt.Println("ID:", message)
}
```

## Loading from JSON Files

```go
// Load from JSON files (similar to our current structure)
func loadTranslationsFromFiles() *goeasyi18n.I18n {
    i18n := goeasyi18n.NewI18n()
    
    // Load English translations
    enJSON := `{
        "success": {
            "user_created": "User {{.Username}} created successfully",
            "user_retrieved": "User retrieved successfully"
        },
        "error": {
            "user_not_found": "User not found",
            "validation_failed": "Validation failed"
        },
        "field": {
            "email": "email",
            "username": "username"
        }
    }`
    
    // Load from JSON string
    i18n.LoadFromJSONString("en", enJSON)
    
    // Or load from file
    // i18n.LoadFromJSONFile("en", "locales/en.json")
    
    return i18n
}
```

## Wrapper for Our Current API

```go
// EasyI18nManager wraps goeasyi18n to match our current interface
type EasyI18nManager struct {
    i18n *goeasyi18n.I18n
}

func NewEasyI18nManager() *EasyI18nManager {
    i18n := goeasyi18n.NewI18n(goeasyi18n.Config{
        FallbackLanguageName: "en",
    })
    
    // Load translations here...
    
    return &EasyI18nManager{i18n: i18n}
}

func (m *EasyI18nManager) TranslateSuccess(lang, key string, data map[string]interface{}) string {
    return m.i18n.T(lang, key, goeasyi18n.Options{Data: data})
}

func (m *EasyI18nManager) TranslateError(lang, key string, data map[string]interface{}) string {
    return m.i18n.T(lang, key, goeasyi18n.Options{Data: data})
}

func (m *EasyI18nManager) Translate(lang, key string, data map[string]interface{}) string {
    return m.i18n.T(lang, key, goeasyi18n.Options{Data: data})
}
```

## Migration Strategy

### 1. Install GoEasyI18n
```bash
go get github.com/eduardolat/goeasyi18n
```

### 2. Replace Current I18n Manager

**Current Implementation:**
```go
// pkg/i18n/manager.go
func NewManager() *Manager {
    // Current complex implementation
}
```

**New Implementation:**
```go
// pkg/i18n/manager.go
func NewManager() *EasyI18nManager {
    i18n := goeasyi18n.NewI18n(goeasyi18n.Config{
        FallbackLanguageName: "en",
    })
    
    // Load existing translation files
    i18n.LoadFromJSONFile("en", "locales/en.json")
    i18n.LoadFromJSONFile("id", "locales/id.json")
    
    return &EasyI18nManager{i18n: i18n}
}
```

### 3. Update Response Helpers

**Current Usage:**
```go
return response.Success(
    c,
    i18nManager.TranslateSuccess(lang, "user.retrieved", nil),
    user,
)
```

**New Usage (exactly the same!):**
```go
return response.Success(
    c,
    i18nManager.TranslateSuccess(lang, "user.retrieved", nil),
    user,
)
```

## Comparison with Current Implementation

| Feature | Current go-i18n | GoEasyI18n |
|---------|------------------|------------|
| **Dependencies** | go-i18n + text/template | None (pure Go) |
| **Setup Complexity** | Medium (requires Bundle, Localizer) | Simple (one struct) |
| **Template Variables** | Manual template parsing | Built-in support |
| **Consistency Check** | Manual | Automatic |
| **File Loading** | Manual file reading + parsing | Built-in methods |
| **Language Functions** | Manual creation | Automatic generation |
| **API Simplicity** | Complex | Very simple |
| **Performance** | Good | Good |
| **Community** | Large | Smaller but active |

## Advantages of Migration

1. **Simplified Code**: Much less boilerplate code required
2. **Better DX**: More intuitive API with automatic consistency checking
3. **Reduced Dependencies**: One less external dependency
4. **Built-in Features**: Template support, multiple loading methods
5. **Easy Testing**: Simple to mock and test

## Disadvantages

1. **Smaller Community**: Less mature ecosystem
2. **Migration Effort**: Need to update existing translation files
3. **Learning Curve**: Team needs to learn new API (though it's simpler)

## Recommendation

**Keep Current Implementation** for now because:
1. Current system is working well
2. Migration effort may not justify the benefits
3. go-i18n is more mature and battle-tested
4. Our current wrapper already simplifies the API

**Consider GoEasyI18n** for:
- New projects
- When current i18n becomes maintenance burden
- If team specifically needs simpler i18n solution

## Example Translation Files

### locales/en.json
```json
{
  "success": {
    "user_created": "User {{.Username}} created successfully",
    "user_retrieved": "User retrieved successfully",
    "users_retrieved": "{{.Count}} users retrieved successfully"
  },
  "error": {
    "user_not_found": "User not found",
    "validation_failed": "Validation failed"
  },
  "field": {
    "email": "email",
    "username": "username",
    "password": "password"
  },
  "validation": {
    "required": "{{.Field}} is required",
    "email": "{{.Field}} must be a valid email",
    "min": "{{.Field}} must be at least {{.MinLength}} characters"
  }
}
```

### locales/id.json
```json
{
  "success": {
    "user_created": "Pengguna {{.Username}} berhasil dibuat",
    "user_retrieved": "Pengguna berhasil diambil", 
    "users_retrieved": "{{.Count}} pengguna berhasil diambil"
  },
  "error": {
    "user_not_found": "Pengguna tidak ditemukan",
    "validation_failed": "Validasi gagal"
  },
  "field": {
    "email": "email",
    "username": "nama pengguna", 
    "password": "kata sandi"
  },
  "validation": {
    "required": "{{.Field}} wajib diisi",
    "email": "{{.Field}} harus berupa email yang valid",
    "min": "{{.Field}} harus minimal {{.MinLength}} karakter"
  }
}
```
