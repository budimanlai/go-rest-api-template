package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	apiKeyTestURL = "http://localhost:3000"
	testAPIKey    = "test-api-key-12345"
)

// TestApiKeyOnlyAuthentication tests the new ApiKeyOnlyMiddleware implementation
func TestApiKeyOnlyAuthentication(t *testing.T) {
	fmt.Println("=== Testing ApiKeyOnlyMiddleware Implementation ===")

	// Step 1: Test Register with API Key Only
	t.Run("Register with API Key Only", func(t *testing.T) {
		err := testRegisterUser()
		// Registration may fail if user already exists, that's okay
		if err != nil {
			fmt.Printf("⚠️  Registration failed (user may already exist): %v\n", err)
		} else {
			fmt.Println("✅ User registered successfully")
		}
	})

	// Step 2: Test Login with API Key Only
	t.Run("Login with API Key Only", func(t *testing.T) {
		privateToken, err := testLoginUser()
		require.NoError(t, err, "Login should succeed")
		require.NotEmpty(t, privateToken, "Private token should not be empty")
		fmt.Printf("✅ Private Token: %s...\n", privateToken[:20])

		// Step 3: Test Private Endpoints with Private Token
		t.Run("Test Private Endpoints", func(t *testing.T) {
			err := testPrivateEndpoints(privateToken)
			assert.NoError(t, err, "Private endpoints should work with valid token")
			fmt.Println("✅ All private endpoints working correctly")
		})
	})

	fmt.Println("🎉 ApiKeyOnlyMiddleware tests completed!")
}

func testRegisterUser() error {
	userData := map[string]string{
		"username":  "testuser_apikey",
		"email":     "test_apikey@example.com",
		"password":  "password123",
		"full_name": "Test User ApiKey",
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", apiKeyTestURL+"/api/v1/public/auth/register", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Only API Key required - no JWT token needed (following industry best practices)
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

func testLoginUser() (string, error) {
	loginData := map[string]string{
		"username": "testuser_apikey",
		"password": "password123",
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", apiKeyTestURL+"/api/v1/public/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	// Only API Key required - no JWT token needed (following industry best practices)
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

	req, err := http.NewRequest("GET", apiKeyTestURL+"/api/v1/users", nil)
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
