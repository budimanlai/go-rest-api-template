// Simple GoEasyI18n Demo
// Run this with: go run examples/goeasyi18n_demo.go

package main

import (
	"fmt"
	"log"
)

// Simulate goeasyi18n usage (since we don't have it installed)
type mockI18n struct {
	translations map[string]map[string]string
}

func newMockI18n() *mockI18n {
	return &mockI18n{
		translations: map[string]map[string]string{
			"en": {
				"user_created":        "User {{.Username}} created successfully",
				"user_retrieved":      "User retrieved successfully",
				"users_retrieved":     "{{.Count}} users retrieved successfully",
				"user_not_found":      "User not found",
				"validation_failed":   "Validation failed. Please check the following fields",
				"field.email":         "email",
				"field.username":      "username",
				"validation.required": "{{.Field}} is required",
			},
			"id": {
				"user_created":        "Pengguna {{.Username}} berhasil dibuat",
				"user_retrieved":      "Pengguna berhasil diambil",
				"users_retrieved":     "{{.Count}} pengguna berhasil diambil",
				"user_not_found":      "Pengguna tidak ditemukan",
				"validation_failed":   "Validasi gagal. Silakan periksa field berikut",
				"field.email":         "email",
				"field.username":      "nama pengguna",
				"validation.required": "{{.Field}} wajib diisi",
			},
		},
	}
}

func (m *mockI18n) T(lang, key string) string {
	if translations, exists := m.translations[lang]; exists {
		if text, exists := translations[key]; exists {
			return text
		}
	}
	return key // fallback
}

func main() {
	fmt.Println("=== GoEasyI18n Demo ===")
	fmt.Println("Menunjukkan betapa mudahnya menggunakan goeasyi18n")
	fmt.Println()

	// Create mock instance (would be: goeasyi18n.NewI18n())
	i18n := newMockI18n()

	fmt.Println("1. Terjemahan sederhana:")
	fmt.Printf("   EN: %s\n", i18n.T("en", "user_retrieved"))
	fmt.Printf("   ID: %s\n", i18n.T("id", "user_retrieved"))
	fmt.Println()

	fmt.Println("2. Terjemahan dengan template (would support variables):")
	fmt.Printf("   EN: %s\n", i18n.T("en", "user_created"))
	fmt.Printf("   ID: %s\n", i18n.T("id", "user_created"))
	fmt.Println()

	fmt.Println("3. Error messages:")
	fmt.Printf("   EN: %s\n", i18n.T("en", "user_not_found"))
	fmt.Printf("   ID: %s\n", i18n.T("id", "user_not_found"))
	fmt.Println()

	fmt.Println("4. Validation messages:")
	fmt.Printf("   EN field: %s\n", i18n.T("en", "field.email"))
	fmt.Printf("   ID field: %s\n", i18n.T("id", "field.username"))
	fmt.Printf("   EN validation: %s\n", i18n.T("en", "validation.required"))
	fmt.Printf("   ID validation: %s\n", i18n.T("id", "validation.required"))
	fmt.Println()

	fmt.Println("=== Perbandingan API ===")
	fmt.Println()

	fmt.Println("Current Implementation (go-i18n):")
	fmt.Println("```go")
	fmt.Println("bundle := i18n.NewBundle(language.English)")
	fmt.Println("bundle.RegisterUnmarshalFunc(\"json\", json.Unmarshal)")
	fmt.Println("bundle.LoadMessageFile(\"locales/en.json\")")
	fmt.Println("localizer := i18n.NewLocalizer(bundle, \"en\")")
	fmt.Println("message := localizer.Localize(&i18n.LocalizeConfig{")
	fmt.Println("    MessageID: \"user_created\",")
	fmt.Println("    TemplateData: map[string]interface{}{\"Username\": \"john\"},")
	fmt.Println("})")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("GoEasyI18n Implementation:")
	fmt.Println("```go")
	fmt.Println("i18n := goeasyi18n.NewI18n()")
	fmt.Println("i18n.LoadFromJSONFile(\"en\", \"locales/en.json\")")
	fmt.Println("message := i18n.T(\"en\", \"user_created\", goeasyi18n.Options{")
	fmt.Println("    Data: map[string]any{\"Username\": \"john\"},")
	fmt.Println("})")
	fmt.Println("```")
	fmt.Println()

	fmt.Println("Keunggulan GoEasyI18n:")
	fmt.Println("✅ API yang sangat sederhana - hanya i18n.T(lang, key)")
	fmt.Println("✅ Tidak butuh dependencies eksternal")
	fmt.Println("✅ Built-in template support untuk variables")
	fmt.Println("✅ Automatic consistency checking")
	fmt.Println("✅ Bisa load dari JSON, YAML, string, file, embed.FS")
	fmt.Println("✅ Language-specific functions (translateEn, translateId)")
	fmt.Println("✅ Zero configuration untuk basic usage")
	fmt.Println()

	fmt.Println("Kekurangan GoEasyI18n:")
	fmt.Println("❌ Community lebih kecil dari go-i18n")
	fmt.Println("❌ Kurang battle-tested")
	fmt.Println("❌ Perlu migration effort")
	fmt.Println()

	fmt.Println("=== Kesimpulan ===")
	fmt.Println("GoEasyI18n memang terlihat JAUH lebih mudah digunakan!")
	fmt.Println("Namun, current implementation dengan go-i18n sudah bekerja dengan baik.")
	fmt.Println("Rekomendasi: tetap gunakan go-i18n untuk stabilitas,")
	fmt.Println("pertimbangkan GoEasyI18n untuk project baru.")

	log.Println("Demo selesai! GoEasyI18n sangat menarik untuk dicoba.")
}
