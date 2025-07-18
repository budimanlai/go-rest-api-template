package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Manager handles internationalization
type Manager struct {
	bundle     *i18n.Bundle
	localizers map[string]*i18n.Localizer
}

// Config for i18n manager
type Config struct {
	DefaultLanguage string
	LocalesPath     string
	SupportedLangs  []string
	Modules         []string // Optional: specific modules to load (e.g., ["common", "user", "auth"])
}

// NewManager creates a new i18n manager
func NewManager(config Config) (*Manager, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Default modules if not specified
	modules := config.Modules
	if len(modules) == 0 {
		modules = []string{"common", "user", "auth"} // Default modules
	}

	// Load message files for each supported language and module
	for _, lang := range config.SupportedLangs {
		// Try to load legacy single file first (for backward compatibility)
		legacyFile := fmt.Sprintf("%s.json", lang)
		legacyPath := filepath.Join(config.LocalesPath, legacyFile)
		if _, err := os.Stat(legacyPath); err == nil {
			if _, err := bundle.LoadMessageFile(legacyPath); err != nil {
				fmt.Printf("Warning: Could not load legacy language file %s: %v\n", legacyPath, err)
			}
			continue // Skip module loading if legacy file exists
		}

		// Load modular files
		for _, module := range modules {
			filename := fmt.Sprintf("%s.%s.json", module, lang)
			filePath := filepath.Join(config.LocalesPath, lang, filename)

			if _, err := bundle.LoadMessageFile(filePath); err != nil {
				// Log warning but don't fail if file doesn't exist
				fmt.Printf("Warning: Could not load language file %s: %v\n", filePath, err)
			}
		}
	}

	// Create localizers for each supported language
	localizers := make(map[string]*i18n.Localizer)
	for _, lang := range config.SupportedLangs {
		localizers[lang] = i18n.NewLocalizer(bundle, lang)
	}

	return &Manager{
		bundle:     bundle,
		localizers: localizers,
	}, nil
}

// GetLocalizer returns localizer for given language
func (m *Manager) GetLocalizer(lang string) *i18n.Localizer {
	if localizer, exists := m.localizers[lang]; exists {
		return localizer
	}
	// Return default (English) if language not found
	return m.localizers["en"]
}

// Translate translates a message key to the specified language
func (m *Manager) Translate(lang, messageID string, templateData map[string]interface{}) string {
	localizer := m.GetLocalizer(lang)

	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
	})

	if err != nil {
		// Return message ID if translation fails
		return messageID
	}

	return message
}

// TranslateError translates error messages
func (m *Manager) TranslateError(lang, errorKey string, data map[string]interface{}) string {
	return m.Translate(lang, "error."+errorKey, data)
}

// TranslateSuccess translates success messages
func (m *Manager) TranslateSuccess(lang, successKey string, data map[string]interface{}) string {
	return m.Translate(lang, "success."+successKey, data)
}

// GetSupportedLanguages returns list of supported languages
func (m *Manager) GetSupportedLanguages() []string {
	langs := make([]string, 0, len(m.localizers))
	for lang := range m.localizers {
		langs = append(langs, lang)
	}
	return langs
}
