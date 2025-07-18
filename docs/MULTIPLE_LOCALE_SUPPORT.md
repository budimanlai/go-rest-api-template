# Go-i18n Multiple Locale Support Guide

## Current Implementation Status ✅

**go-i18n SUDAH SUPPORT multiple locale JSON files!** Current implementation kita sudah configured dengan baik:

### 📁 Current Locale Files
```
locales/
├── en.json    # English translations (110 keys)
└── id.json    # Indonesian translations (110 keys)
```

### 🔧 Current Configuration
```go
// pkg/i18n/manager.go
func NewManager(config Config) (*Manager, error) {
    bundle := i18n.NewBundle(language.English)
    bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

    // Load message files for each supported language
    for _, lang := range config.SupportedLangs {
        filename := fmt.Sprintf("%s.json", lang)
        filePath := filepath.Join(config.LocalesPath, filename)
        
        // Loads each locale file automatically!
        if _, err := bundle.LoadMessageFile(filePath); err != nil {
            fmt.Printf("Warning: Could not load language file %s: %v\n", filePath, err)
        }
    }
    
    // Create localizers for each supported language
    localizers := make(map[string]*i18n.Localizer)
    for _, lang := range config.SupportedLangs {
        localizers[lang] = i18n.NewLocalizer(bundle, lang)
    }
}
```

## 🌍 Adding New Locales

### Step 1: Create New Locale File

Let's add Japanese (ja) support:

```bash
# Create Japanese locale file
touch locales/ja.json
```

### Step 2: Add Translations

Copy structure from existing files:

```json
// locales/ja.json
{
  "error.user_not_found": "ユーザーが見つかりません",
  "error.username_required": "ユーザー名は必須です",
  "error.username_min_length": "ユーザー名は{{.MinLength}}文字以上である必要があります",
  "error.email_required": "メールアドレスは必須です",
  "error.email_invalid": "無効なメール形式",
  "error.password_required": "パスワードは必須です",
  "error.password_min_length": "パスワードは{{.MinLength}}文字以上である必要があります",
  "success.user_retrieved": "ユーザーが正常に取得されました",
  "success.users_retrieved": "{{.Count}}人のユーザーが正常に取得されました"
}
```

### Step 3: Update Configuration

Update your app configuration to include new language:

```go
// In your main.go or config
config := i18n.Config{
    DefaultLanguage: "en",
    LocalesPath:     "locales",
    SupportedLangs:  []string{"en", "id", "ja"}, // Add "ja" here
}
```

### Step 4: Test New Locale

```go
// Usage remains exactly the same!
manager, _ := i18n.NewManager(config)

// Japanese translation
message := manager.TranslateError("ja", "error.user_not_found", nil)
fmt.Println(message) // Output: "ユーザーが見つかりません"
```

## 🚀 Advanced Multiple Locale Features

### 1. **Fallback Language Support**
```go
// If translation missing in 'ja', falls back to 'en'
localizer := i18n.NewLocalizer(bundle, "ja", "en")
```

### 2. **Nested JSON Structure**
```json
// locales/fr.json - French with nested structure
{
  "user": {
    "error": {
      "not_found": "Utilisateur introuvable",
      "validation": {
        "username_required": "Le nom d'utilisateur est requis",
        "email_invalid": "Format d'email invalide"
      }
    },
    "success": {
      "created": "Utilisateur {{.Username}} créé avec succès",
      "retrieved": "Utilisateur récupéré avec succès"
    }
  }
}
```

### 3. **Regional Variants**
```
locales/
├── en.json       # English (default)
├── en-US.json    # American English
├── en-GB.json    # British English
├── id.json       # Indonesian
├── ja.json       # Japanese
├── zh.json       # Chinese (Simplified)
├── zh-TW.json    # Chinese (Traditional)
└── fr.json       # French
```

### 4. **Dynamic Language Loading**
```go
// Load additional languages at runtime
func (m *Manager) AddLanguage(lang, filePath string) error {
    if _, err := m.bundle.LoadMessageFile(filePath); err != nil {
        return err
    }
    
    m.localizers[lang] = i18n.NewLocalizer(m.bundle, lang)
    return nil
}

// Usage
manager.AddLanguage("es", "locales/es.json") // Add Spanish
```

## 📝 Example: Adding Spanish Support

Let's create a complete example:

### 1. Create Spanish Translations
```json
// locales/es.json
{
  "error.user_not_found": "Usuario no encontrado",
  "error.username_required": "El nombre de usuario es requerido",
  "error.username_min_length": "El nombre de usuario debe tener al menos {{.MinLength}} caracteres",
  "error.email_required": "El email es requerido",
  "error.email_invalid": "Formato de email inválido",
  "error.password_required": "La contraseña es requerida",
  "error.password_min_length": "La contraseña debe tener al menos {{.MinLength}} caracteres",
  "error.username_exists": "El nombre de usuario ya existe",
  "error.email_exists": "El email ya existe",
  "error.invalid_credentials": "Nombre de usuario o contraseña inválidos",
  
  "success.user_created": "Usuario {{.Username}} creado exitosamente",
  "success.user_retrieved": "Usuario recuperado exitosamente",
  "success.users_retrieved": "{{.Count}} usuarios recuperados exitosamente",
  "success.user_updated": "Usuario actualizado exitosamente",
  "success.user_deleted": "Usuario eliminado exitosamente",
  
  "field.username": "nombre de usuario",
  "field.email": "email",
  "field.password": "contraseña",
  "field.fullname": "nombre completo",
  
  "validation.required": "{{.Field}} es requerido",
  "validation.email": "{{.Field}} debe ser un email válido",
  "validation.min": "{{.Field}} debe tener al menos {{.MinLength}} caracteres",
  "validation.max": "{{.Field}} no puede tener más de {{.MaxLength}} caracteres"
}
```

### 2. Update Main Configuration
```go
// In internal/application/app.go or wherever you initialize i18n
config := i18n.Config{
    DefaultLanguage: "en",
    LocalesPath:     "locales",
    SupportedLangs:  []string{"en", "id", "es"}, // Add Spanish
}

i18nManager, err := i18n.NewManager(config)
```

### 3. Test Spanish Translations
```go
// In your handlers, Spanish will work automatically!
lang := c.Get("lang", "en").(string) // Could be "es"

return response.Success(
    c,
    i18nManager.TranslateSuccess(lang, "success.user_retrieved", nil),
    user,
)

// If lang = "es", returns: "Usuario recuperado exitosamente"
```

## 🔧 Current Implementation Benefits

### ✅ **Already Working Features:**
1. **Multiple JSON files** - ✅ Supports en.json, id.json, and can add more
2. **Automatic loading** - ✅ Loads all files in SupportedLangs array
3. **Template variables** - ✅ Supports {{.Variable}} interpolation
4. **Fallback support** - ✅ Falls back to default language if key missing
5. **Dynamic access** - ✅ Runtime language switching works perfectly

### ✅ **File Structure Already Optimal:**
```
locales/
├── en.json    # 110 translation keys
├── id.json    # 110 translation keys
└── [new].json # Just add more files here!
```

### ✅ **Zero Code Changes Needed:**
- Current API already supports any language
- Just add new JSON files and update config
- Handlers automatically work with new languages

## 🎯 Recommendations

### **For Adding New Languages:**

1. **Copy existing structure** from en.json or id.json
2. **Translate all keys** to maintain consistency
3. **Add language code** to SupportedLangs config
4. **Test with API calls** using new language

### **Best Practices:**

1. **Keep key structure consistent** across all locale files
2. **Use template variables** for dynamic content: `{{.Variable}}`
3. **Test fallback behavior** when translations are missing
4. **Organize keys logically** (error.*, success.*, field.*, validation.*)

## 🚀 Ready to Add More Languages!

Current go-i18n implementation is **perfectly set up** for multiple locales. You can add as many languages as needed:

```bash
# Add any new language files
touch locales/fr.json    # French
touch locales/ja.json    # Japanese  
touch locales/de.json    # German
touch locales/es.json    # Spanish
touch locales/pt.json    # Portuguese
touch locales/ar.json    # Arabic
touch locales/zh.json    # Chinese
```

**The system will automatically load and use them!** 🎉
