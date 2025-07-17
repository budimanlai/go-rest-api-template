package main

import (
	"context"
	"fmt"
	"go-rest-api-template/internal/domain/entity"
	"go-rest-api-template/internal/service"
	"os"
	"strings"
)

// MockApiKeyService for manual testing
type MockApiKeyService struct{}

func (m *MockApiKeyService) ValidateApiKey(ctx context.Context, key string) (*entity.ApiKey, error) {
	return &entity.ApiKey{ID: 1, Name: "test-api-key"}, nil
}

func (m *MockApiKeyService) ValidateAuthKey(ctx context.Context, authKey string) (*entity.ApiKey, error) {
	return &entity.ApiKey{ID: 1, Name: "test-api-key"}, nil
}

func (m *MockApiKeyService) GetApiKeyByID(ctx context.Context, id int) (*entity.ApiKey, error) {
	return &entity.ApiKey{ID: id, Name: "test-api-key"}, nil
}

func (m *MockApiKeyService) GetAllApiKeys(ctx context.Context, limit, offset int) ([]*entity.ApiKey, error) {
	return []*entity.ApiKey{{ID: 1, Name: "test-api-key"}}, nil
}

func (m *MockApiKeyService) GetActiveApiKeys(ctx context.Context, limit, offset int) ([]*entity.ApiKey, error) {
	return []*entity.ApiKey{{ID: 1, Name: "test-api-key"}}, nil
}

func (m *MockApiKeyService) LogApiKeyAccess(ctx context.Context, apiKeyID int) error {
	return nil
}

func printSeparator() {
	fmt.Println("\n" + strings.Repeat("=", 60))
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "help" {
		fmt.Println("Manual JWT Testing Tool")
		fmt.Println("Usage:")
		fmt.Println("  go run manual_test.go          - Run all tests")
		fmt.Println("  go run manual_test.go help     - Show this help")
		return
	}

	fmt.Println("🧪 MANUAL JWT TESTING TOOL")
	printSeparator()

	// Initialize services
	mockApiKeyService := &MockApiKeyService{}
	jwtService := service.NewJWTService("test-secret-key-for-manual-testing", 2, 24, mockApiKeyService)

	// Test data
	apiKey := &entity.ApiKey{
		ID:   1,
		Name: "manual-test-api-key",
	}

	user := &entity.User{
		ID:       999,
		Username: "manualuser",
		Email:    "manual@test.com",
	}

	// Test 1: Generate and Display Public Token (for accessing public endpoints)
	fmt.Println("1️⃣  GENERATING PUBLIC TOKEN")
	fmt.Printf("   📝 Purpose: Access public endpoints (login, register, etc.)\n")
	fmt.Printf("   API Key: ID=%d, Name=%s\n", apiKey.ID, apiKey.Name)

	publicToken, err := jwtService.GeneratePublicToken(apiKey)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Public Token Generated:\n")
	fmt.Printf("   %s\n", publicToken)
	fmt.Printf("   💡 Use this token for login endpoint\n")

	// Test 2: Validate Public Token
	printSeparator()
	fmt.Println("2️⃣  VALIDATING PUBLIC TOKEN")

	claims, returnedApiKey, err := jwtService.ValidatePublicToken(publicToken)
	if err != nil {
		fmt.Printf("❌ Validation Error: %v\n", err)
	} else {
		fmt.Printf("✅ Public Token Valid!\n")
		fmt.Printf("   Claims API Key ID: %d\n", claims.ApiKeyID)
		fmt.Printf("   Claims API Key Name: %s\n", claims.ApiKeyName)
		fmt.Printf("   Retrieved API Key ID: %d\n", returnedApiKey.ID)
		fmt.Printf("   Expires At: %v\n", claims.ExpiresAt.Time)
		fmt.Printf("   Subject: %s\n", claims.Subject)
	}

	// Test 3: Simulate Login Process - Generate Private Token
	printSeparator()
	fmt.Println("3️⃣  SIMULATE LOGIN PROCESS")
	fmt.Printf("   📝 Purpose: After successful login, get private token for private endpoints\n")
	fmt.Printf("   🔐 Login with Public Token → Get Private Token\n")
	fmt.Printf("   API Key: ID=%d, Name=%s\n", apiKey.ID, apiKey.Name)
	fmt.Printf("   User: ID=%d, Username=%s, Email=%s\n", user.ID, user.Username, user.Email)

	privateToken, err := jwtService.GeneratePrivateToken(apiKey, user)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	fmt.Printf("✅ Private Token Generated (after login):\n")
	fmt.Printf("   %s\n", privateToken)
	fmt.Printf("   💡 Use this token for private endpoints (user profile, etc.)\n")

	// Test 4: Validate Private Token
	printSeparator()
	fmt.Println("4️⃣  VALIDATING PRIVATE TOKEN")

	privateClaims, returnedApiKey2, returnedUser, err := jwtService.ValidatePrivateToken(privateToken)
	if err != nil {
		fmt.Printf("❌ Validation Error: %v\n", err)
	} else {
		fmt.Printf("✅ Private Token Valid!\n")
		fmt.Printf("   Claims API Key ID: %d\n", privateClaims.ApiKeyID)
		fmt.Printf("   Claims API Key Name: %s\n", privateClaims.ApiKeyName)
		fmt.Printf("   Claims User ID: %d\n", privateClaims.UserID)
		fmt.Printf("   Claims Username: %s\n", privateClaims.Username)
		fmt.Printf("   Claims Email: %s\n", privateClaims.Email)
		fmt.Printf("   Retrieved API Key ID: %d\n", returnedApiKey2.ID)
		fmt.Printf("   Retrieved User ID: %d\n", returnedUser.ID)
		fmt.Printf("   Expires At: %v\n", privateClaims.ExpiresAt.Time)
		fmt.Printf("   Subject: %s\n", privateClaims.Subject)
	}

	// Test 5: Security Test - Try to use Public Token for Private Validation
	printSeparator()
	fmt.Println("5️⃣  SECURITY TEST - Token Type Validation")
	fmt.Printf("   📝 Testing: Public token should NOT work for private endpoints\n")
	fmt.Println("   Trying to validate Public Token as Private Token...")

	_, _, _, err = jwtService.ValidatePrivateToken(publicToken)
	if err != nil {
		fmt.Printf("✅ Security OK: Public token correctly rejected for private endpoint\n")
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("❌ Security FAILED: Public token incorrectly accepted for private endpoint!\n")
	}

	// Test 6: Manual Token Inspection
	printSeparator()
	fmt.Println("6️⃣  USAGE FLOW SUMMARY")
	fmt.Println("   📱 Typical API Usage Flow:")
	fmt.Println("")
	fmt.Println("   1. Client gets PUBLIC TOKEN first (with API key)")
	fmt.Println("   2. Client calls LOGIN endpoint using PUBLIC TOKEN")
	fmt.Println("   3. Login successful → Server returns PRIVATE TOKEN")
	fmt.Println("   4. Client uses PRIVATE TOKEN for protected endpoints")
	fmt.Println("")
	fmt.Println("   🔧 Manual Testing Tools:")
	fmt.Println("   🌐 Online JWT Debugger: https://jwt.io/")
	fmt.Println("   🔑 Secret Key: test-secret-key-for-manual-testing")
	fmt.Println("")
	fmt.Println("   📋 Copy tokens above for manual testing")

	printSeparator()
	fmt.Println("✅ Manual testing completed!")
	fmt.Println("💡 Next: Use these tokens to test your API endpoints")
	fmt.Println("🚀 Run: ./manual_test_guide.sh for cURL examples")
}
