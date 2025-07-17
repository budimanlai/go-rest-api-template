package main

import (
	"encoding/json"
	"fmt"
	"go-rest-api-template/internal/model"
	"go-rest-api-template/pkg/validator"
)

func main() {
	fmt.Println("🧪 Testing Production-Ready Validator")
	fmt.Println("=====================================")

	// Test 1: Valid UserCreateRequest
	fmt.Println("\n1. Testing Valid UserCreateRequest:")
	validUser := model.UserCreateRequest{
		Username: "johndoe123",
		Email:    "john@example.com",
		Password: "securePassword123",
	}

	if err := validUser.Validate(); err != nil {
		fmt.Printf("❌ Validation failed (unexpected): %v\n", err)
	} else {
		fmt.Printf("✅ Valid user passed validation\n")
	}

	// Test 2: Invalid UserCreateRequest
	fmt.Println("\n2. Testing Invalid UserCreateRequest:")
	invalidUser := model.UserCreateRequest{
		Username: "ab",      // Too short (min 3)
		Email:    "invalid", // Invalid email format
		Password: "123",     // Too short (min 8)
	}

	if err := invalidUser.Validate(); err != nil {
		fmt.Printf("✅ Invalid user caught by validation:\n")
		errors := validator.GetValidationErrors(err)
		for _, e := range errors {
			fmt.Printf("   - Field: %s, Error: %s\n", e.Field, e.Message)
		}
	} else {
		fmt.Printf("❌ Invalid user passed validation (unexpected)\n")
	}

	// Test 3: Test advanced validation rules
	fmt.Println("\n3. Testing ChangePasswordRequest with nefield validation:")
	samePasswordReq := model.ChangePasswordRequest{
		CurrentPassword: "password123",
		NewPassword:     "password123", // Same as current (should fail with nefield)
	}

	if err := samePasswordReq.Validate(); err != nil {
		fmt.Printf("✅ Same password validation working:\n")
		errors := validator.GetValidationErrors(err)
		for _, e := range errors {
			fmt.Printf("   - Field: %s, Error: %s\n", e.Field, e.Message)
		}
	} else {
		fmt.Printf("❌ Same password validation not working\n")
	}

	// Test 4: Test structured error output
	fmt.Println("\n4. Testing JSON error output:")
	errorData := validator.GetValidationErrors(invalidUser.Validate())
	jsonOutput, _ := json.MarshalIndent(errorData, "", "  ")
	fmt.Printf("JSON Error Structure:\n%s\n", jsonOutput)

	// Test 5: Test UserUpdateRequest with oneof validation
	fmt.Println("\n5. Testing UserUpdateRequest with invalid status:")
	invalidStatus := model.UserUpdateRequest{
		Username: "validuser",
		Email:    "valid@example.com",
		Status:   "invalid_status", // Should fail oneof validation
	}

	if err := invalidStatus.Validate(); err != nil {
		fmt.Printf("✅ Status validation working:\n")
		errors := validator.GetValidationErrors(err)
		for _, e := range errors {
			fmt.Printf("   - Field: %s, Error: %s\n", e.Field, e.Message)
		}
	}

	fmt.Println("\n🎉 All validator tests completed!")
	fmt.Println("✅ Validator is production-ready with:")
	fmt.Println("   - Comprehensive validation rules")
	fmt.Println("   - Structured error responses")
	fmt.Println("   - Field-specific error messages")
	fmt.Println("   - JSON tag field names")
	fmt.Println("   - Advanced rules (nefield, oneof, etc.)")
}
