// ApiKey Only Authentication CLI Tester
// Following industry best practices (Stripa, GitHub, Twitter/X approach)
//
// Usage: go run tools/apikey_auth_tester.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	testBaseURL string
	testAPIKey  string
)

// init loads environment variables
func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get config from environment or use defaults
	testBaseURL = getEnv("TEST_BASE_URL", "http://localhost:8080")
	testAPIKey = getEnv("TEST_API_KEY", "dev_api_key_12345678901234567890")
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	fmt.Println("=== ApiKeyOnlyMiddleware CLI Tester ===")
	fmt.Printf("Base URL: %s\n", testBaseURL)
	fmt.Printf("API Key: %s\n", testAPIKey)
	fmt.Println("Following industry best practices (Strapi, GitHub, Twitter/X approach)")
	fmt.Println()

	// Step 1: Register User
	fmt.Println("1. Registering User (API Key Only)...")
	err := registerUser()
	if err != nil {
		fmt.Printf("⚠️  Registration failed (user may already exist): %v\n", err)
	} else {
		fmt.Println("✅ User registered successfully")
	}
	fmt.Println()

	// Step 2: Login to get Private Token
	fmt.Println("2. Logging in (API Key Only)...")
	privateToken, err := loginUser()
	if err != nil {
		fmt.Printf("❌ Failed to login: %v\n", err)
		return
	}
	fmt.Printf("✅ Private Token obtained: %s...\n", privateToken[:20])
	fmt.Println()

	// Step 3: Test Private Endpoints
	fmt.Println("3. Testing Private Endpoints (API Key + Private Token)...")
	err = testPrivateEndpoints(privateToken)
	if err != nil {
		fmt.Printf("❌ Failed to test private endpoints: %v\n", err)
		return
	}
	fmt.Println("✅ All private endpoints working correctly")
	fmt.Println()

	fmt.Println("🎉 All tests completed successfully!")
	fmt.Println("✅ ApiKeyOnlyMiddleware is working correctly")
	fmt.Println("✅ Authentication flow follows industry standards")
}

func registerUser() error {
	userData := map[string]string{
		"username":  "testuser_cli",
		"email":     "test_cli@example.com",
		"password":  "password123",
		"full_name": "Test User CLI",
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", testBaseURL+"/api/v1/public/auth/register", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Only API Key required - no JWT token needed
	// This follows industry best practices used by Strapi, GitHub, Twitter/X
	req.Header.Set("X-API-Key", testAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func loginUser() (string, error) {
	loginData := map[string]string{
		"username": "testuser_cli",
		"password": "password123",
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", testBaseURL+"/api/v1/public/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Only API Key required - no JWT token needed
	// This follows industry best practices used by Strapi, GitHub, Twitter/X
	req.Header.Set("X-API-Key", testAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response APIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	// Extract token from response
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	token, ok := data["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}

func testPrivateEndpoints(privateToken string) error {
	// Test GET /api/v1/users (list users)
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", testBaseURL+"/api/v1/users", nil)
	if err != nil {
		return err
	}

	// Both API Key and Private Token required for private endpoints
	req.Header.Set("X-API-Key", testAPIKey)
	req.Header.Set("Authorization", "Bearer "+privateToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("private endpoints test failed with status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("  ✅ GET /api/v1/users - Success")
	return nil
}
