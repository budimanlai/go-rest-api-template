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
	resp, err := client.Get(baseURL + "/api/v1/health")
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
	req.Header.Set("X-API-Key", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

// makeRequestWithAuth makes authenticated request with JWT token
func makeRequestWithAuth(method, endpoint string, token string, body interface{}) (*http.Response, error) {
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
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Authorization", "Bearer "+token)

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

	var jwtToken string
	var createdUserID int

	// Step 1: Register user using public endpoint
	t.Run("Register User", func(t *testing.T) {
		resp, err := makeRequest("POST", "/api/v1/public/auth/register", testUser)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should succeed or user already exists
		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			var response APIResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			assert.True(t, response.Success)
		}
	})

	// Step 2: Login to get JWT token
	t.Run("Login User", func(t *testing.T) {
		loginData := map[string]interface{}{
			"username": testUser["username"],
			"password": testUser["password"],
		}

		resp, err := makeRequest("POST", "/api/v1/public/auth/login", loginData)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)

		// Extract JWT token
		if tokenData, ok := response.Data.(map[string]interface{}); ok {
			if token, ok := tokenData["token"].(string); ok {
				jwtToken = token
			}
		}
		require.NotEmpty(t, jwtToken, "JWT token should not be empty")
	})

	// Step 3: Create User using private endpoint with JWT
	t.Run("Create User", func(t *testing.T) {
		newUser := map[string]interface{}{
			"username": "newuser123",
			"email":    "newuser@example.com",
			"password": "password123",
		}

		resp, err := makeRequestWithAuth("POST", "/api/v1/users", jwtToken, newUser)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "created")

		// Extract user ID from response
		if userData, ok := response.Data.(map[string]interface{}); ok {
			if id, ok := userData["id"].(float64); ok {
				createdUserID = int(id)
			}
		}
	})

	t.Run("Get All Users", func(t *testing.T) {
		resp, err := makeRequestWithAuth("GET", "/api/v1/users", jwtToken, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
	})

	t.Run("Get User By ID", func(t *testing.T) {
		if createdUserID == 0 {
			t.Skip("Skipping because no user was created")
		}

		url := fmt.Sprintf("/api/v1/users/%d", createdUserID)
		resp, err := makeRequestWithAuth("GET", url, jwtToken, nil)
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
			t.Skip("Skipping because no user was created")
		}

		updateData := map[string]interface{}{
			"username": "updateduser123",
			"email":    "updated@example.com",
		}

		url := fmt.Sprintf("/api/v1/users/%d", createdUserID)
		resp, err := makeRequestWithAuth("PUT", url, jwtToken, updateData)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "updated")
	})

	t.Run("Delete User", func(t *testing.T) {
		if createdUserID == 0 {
			t.Skip("Skipping because no user was created")
		}

		url := fmt.Sprintf("/api/v1/users/%d", createdUserID)
		resp, err := makeRequestWithAuth("DELETE", url, jwtToken, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "deleted")
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("Invalid API Key", func(t *testing.T) {
		req, err := http.NewRequest("GET", baseURL+"/api/v1/users", nil)
		require.NoError(t, err)

		req.Header.Set("X-API-Key", "invalid-key")

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Missing API Key", func(t *testing.T) {
		req, err := http.NewRequest("GET", baseURL+"/api/v1/users", nil)
		require.NoError(t, err)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Invalid JWT Token", func(t *testing.T) {
		resp, err := makeRequestWithAuth("GET", "/api/v1/users", "invalid-token", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Access Private Endpoint Without JWT", func(t *testing.T) {
		resp, err := makeRequest("GET", "/api/v1/users", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("Health Check", func(t *testing.T) {
		resp, err := makeRequest("GET", "/api/v1/health", nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response APIResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.True(t, response.Success)
		assert.Contains(t, response.Message, "healthy")
	})
}
