package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	baseURL string
	apiKey  string
)

// init loads environment variables
func init() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get config from environment or use defaults
	baseURL = getEnv("TEST_BASE_URL", "http://localhost:8080")
	apiKey = getEnv("TEST_API_KEY", "dev_api_key_12345678901234567890")
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// TestUser represents test user data
type TestUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Status   string `json:"status"`
}

// APIResponse represents standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func TestMain(m *testing.M) {
	// Check if server is running
	if !isServerRunning() {
		fmt.Println("❌ Server is not running on localhost:8080")
		fmt.Println("Please start the server first: go run ./cmd/api/ --port=8080")
		os.Exit(1)
	}

	fmt.Println("✅ Server is running, starting tests...")
	code := m.Run()
	os.Exit(code)
}

func isServerRunning() bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(baseURL + "/api/v1/users")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

func makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, baseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func TestAPIEndpoints(t *testing.T) {
	// Test data
	testUser := map[string]interface{}{
		"username": "testuser123",
		"email":    "test@example.com",
		"password": "password123",
	}

	var createdUserID int

	t.Run("Create User", func(t *testing.T) {
		resp, err := makeRequest("POST", "/api/v1/users", testUser)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "created successfully")

		// Extract user ID from response
		if userData, ok := response.Data.(map[string]interface{}); ok {
			if id, ok := userData["id"].(float64); ok {
				createdUserID = int(id)
			}
		}
	})

	t.Run("Get All Users", func(t *testing.T) {
		resp, err := makeRequest("GET", "/api/v1/users", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
	})

	t.Run("Get User by ID", func(t *testing.T) {
		if createdUserID == 0 {
			t.Skip("No user ID available from create test")
		}

		endpoint := fmt.Sprintf("/api/v1/users/%d", createdUserID)
		resp, err := makeRequest("GET", endpoint, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
	})

	t.Run("Update User", func(t *testing.T) {
		if createdUserID == 0 {
			t.Skip("No user ID available from create test")
		}

		updateData := map[string]interface{}{
			"username": "updateduser123",
			"email":    "updated@example.com",
		}

		endpoint := fmt.Sprintf("/api/v1/users/%d", createdUserID)
		resp, err := makeRequest("PUT", endpoint, updateData)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
	})

	t.Run("Test Password Management", func(t *testing.T) {
		// Test forgot password
		forgotData := map[string]interface{}{
			"email": "updated@example.com",
		}

		resp, err := makeRequest("POST", "/api/v1/users/forgot-password", forgotData)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 200 even if user doesn't exist (security)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete User", func(t *testing.T) {
		if createdUserID == 0 {
			t.Skip("No user ID available from create test")
		}

		endpoint := fmt.Sprintf("/api/v1/users/%d", createdUserID)
		resp, err := makeRequest("DELETE", endpoint, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("Invalid API Key", func(t *testing.T) {
		req, err := http.NewRequest("GET", baseURL+"/api/v1/users", nil)
		require.NoError(t, err)

		req.Header.Set("x-api-key", "invalid-key")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Missing API Key", func(t *testing.T) {
		req, err := http.NewRequest("GET", baseURL+"/api/v1/users", nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		resp, err := makeRequest("GET", "/api/v1/users/99999", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
