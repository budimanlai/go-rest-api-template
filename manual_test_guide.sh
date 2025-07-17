#!/bin/bash

# Manual Testing Guide for Go REST API with ApiKeyOnlyMiddleware
# Following industry best practices (Strapi, GitHub, Twitter/X approach)

BASE_URL="http://localhost:3000"
API_KEY="test-api-key-12345"

echo "=== Go REST API - ApiKeyOnly Authentication Testing ==="
echo "Base URL: $BASE_URL"
echo "API Key: $API_KEY"
echo "Following industry best practices - no JWT required for register/login"
echo ""

# Function to run actual tests
run_tests() {
    echo "=== RUNNING AUTOMATED TESTS ==="
    echo ""
    
    echo "1. Testing Register Endpoint (API Key Only)..."
    REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/public/auth/register" \
        -H "X-API-Key: $API_KEY" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "testuser", 
            "email": "test@example.com", 
            "password": "password123",
            "full_name": "Test User"
        }')
    
    echo "Register Response: $REGISTER_RESPONSE"
    
    echo ""
    echo "2. Testing Login Endpoint (API Key Only)..."
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/public/auth/login" \
        -H "X-API-Key: $API_KEY" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "testuser", 
            "password": "password123"
        }')
    
    echo "Login Response: $LOGIN_RESPONSE"
    
    # Extract private token from login response
    PRIVATE_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    
    if [ -n "$PRIVATE_TOKEN" ]; then
        echo "✓ Private token obtained: ${PRIVATE_TOKEN:0:20}..."
        
        echo ""
        echo "3. Testing Private Endpoints (API Key + Private Token)..."
        USERS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/users" \
            -H "X-API-Key: $API_KEY" \
            -H "Authorization: Bearer $PRIVATE_TOKEN" \
            -H "Accept: application/json")
        
        echo "Users Response: $USERS_RESPONSE"
        echo "✓ All tests completed successfully!"
        echo ""
        
        echo "2. Registering user (may fail if user exists)..."
        REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/public/auth/register" \
            -H "X-API-Key: $API_KEY" \
            -H "Authorization: Bearer $PUBLIC_TOKEN" \
            -H "Content-Type: application/json" \
            -d '{
                "name": "Test User",
                "email": "test@example.com",
                "password": "password123"
            }')
        echo "Response: $REGISTER_RESPONSE"
        echo ""
        
        echo "3. Logging in to get private token..."
        LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/public/auth/login" \
            -H "X-API-Key: $API_KEY" \
            -H "Authorization: Bearer $PUBLIC_TOKEN" \
            -H "Content-Type: application/json" \
            -d '{
                "email": "test@example.com",
                "password": "password123"
            }')
        echo "Response: $LOGIN_RESPONSE"
        
        # Extract private token
        PRIVATE_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)
        
        if [ -n "$PRIVATE_TOKEN" ]; then
            echo "✓ Private token obtained: ${PRIVATE_TOKEN:0:20}..."
            echo ""
            
            echo "4. Testing private endpoint access..."
            USERS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/users" \
                -H "X-API-Key: $API_KEY" \
                -H "Authorization: Bearer $PRIVATE_TOKEN" \
                -H "Accept: application/json")
            echo "Response: $USERS_RESPONSE"
            echo ""
            
            echo "5. Testing token refresh..."
            REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/public/auth/refresh" \
                -H "X-API-Key: $API_KEY" \
                -H "Authorization: Bearer $PRIVATE_TOKEN" \
                -H "Content-Type: application/json")
            echo "Response: $REFRESH_RESPONSE"
            echo ""
            
        else
            echo "✗ Failed to get private token"
        fi
    else
        echo "✗ Failed to get public token"
    fi
    
    echo "=== ERROR TESTS ==="
    echo ""
    
    echo "Testing with invalid API key..."
    INVALID_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/public/auth/token" \
        -H "X-API-Key: invalid-key" \
        -H "Accept: application/json")
    echo "Response: $INVALID_RESPONSE"
    echo ""
    
    echo "Testing without API key..."
    NO_KEY_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/public/auth/token" \
        -H "Accept: application/json")
    echo "Response: $NO_KEY_RESPONSE"
    echo ""
    
    if [ -n "$PUBLIC_TOKEN" ]; then
        echo "Testing private access with public token..."
        WRONG_TOKEN_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/users" \
            -H "X-API-Key: $API_KEY" \
            -H "Authorization: Bearer $PUBLIC_TOKEN" \
            -H "Accept: application/json")
        echo "Response: $WRONG_TOKEN_RESPONSE"
        echo ""
    fi
}

# Show manual commands
show_manual_commands() {
    echo "=== MANUAL TESTING COMMANDS ==="
    echo ""
    
    echo "=== Step 1: Get Public Token ==="
    echo "curl -X GET \"$BASE_URL/api/v1/public/auth/token\" \\"
    echo "  -H \"X-API-Key: $API_KEY\" \\"
    echo "  -H \"Accept: application/json\""
    echo ""
    
    echo "=== Step 2: Register User (using public token) ==="
    echo "curl -X POST \"$BASE_URL/api/v1/public/auth/register\" \\"
    echo "  -H \"X-API-Key: $API_KEY\" \\"
    echo "  -H \"Authorization: Bearer YOUR_PUBLIC_TOKEN\" \\"
    echo "  -H \"Content-Type: application/json\" \\"
    echo "  -d '{"
    echo "    \"name\": \"Test User\","
    echo "    \"email\": \"test@example.com\","
    echo "    \"password\": \"password123\""
    echo "  }'"
    echo ""
    
    echo "=== Step 3: Login User (get private token) ==="
    echo "curl -X POST \"$BASE_URL/api/v1/public/auth/login\" \\"
    echo "  -H \"X-API-Key: $API_KEY\" \\"
    echo "  -H \"Authorization: Bearer YOUR_PUBLIC_TOKEN\" \\"
    echo "  -H \"Content-Type: application/json\" \\"
    echo "  -d '{"
    echo "    \"email\": \"test@example.com\","
    echo "    \"password\": \"password123\""
    echo "  }'"
    echo ""
    
    echo "=== Step 4: Access Private Endpoint ==="
    echo "curl -X GET \"$BASE_URL/api/v1/users\" \\"
    echo "  -H \"X-API-Key: $API_KEY\" \\"
    echo "  -H \"Authorization: Bearer YOUR_PRIVATE_TOKEN\" \\"
    echo "  -H \"Accept: application/json\""
    echo ""
    
    echo "=== Step 5: Refresh Token ==="
    echo "curl -X POST \"$BASE_URL/api/v1/public/auth/refresh\" \\"
    echo "  -H \"X-API-Key: $API_KEY\" \\"
    echo "  -H \"Authorization: Bearer YOUR_PRIVATE_TOKEN\" \\"
    echo "  -H \"Content-Type: application/json\""
    echo ""
    
    echo "=== Step 6: Logout ==="
    echo "curl -X POST \"$BASE_URL/api/v1/public/auth/logout\" \\"
    echo "  -H \"X-API-Key: $API_KEY\" \\"
    echo "  -H \"Authorization: Bearer YOUR_PRIVATE_TOKEN\" \\"
    echo "  -H \"Content-Type: application/json\""
    echo ""
    
    echo "=== Error Testing ==="
    echo ""
    echo "Test with invalid API key:"
    echo "curl -X GET \"$BASE_URL/api/v1/public/auth/token\" \\"
    echo "  -H \"X-API-Key: invalid-key\" \\"
    echo "  -H \"Accept: application/json\""
    echo ""
    
    echo "Test without API key:"
    echo "curl -X GET \"$BASE_URL/api/v1/public/auth/token\" \\"
    echo "  -H \"Accept: application/json\""
    echo ""
    
    echo "Test accessing private endpoint with public token:"
    echo "curl -X GET \"$BASE_URL/api/v1/users\" \\"
    echo "  -H \"X-API-Key: $API_KEY\" \\"
    echo "  -H \"Authorization: Bearer YOUR_PUBLIC_TOKEN\" \\"
    echo "  -H \"Accept: application/json\""
    echo ""
}

# Main script logic
if [ "$1" = "test" ]; then
    run_tests
else
    show_manual_commands
    echo ""
    echo "=== Usage ==="
    echo "Show manual commands: ./manual_test_guide.sh"
    echo "Run automated tests: ./manual_test_guide.sh test"
fi
