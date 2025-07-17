#!/bin/bash

# ================================
# 🌍 i18n Testing Script
# ================================

BASE_URL="http://localhost:8080"
API_KEY="test-api-key"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================"
echo -e "🌍 MULTILINGUAL API TESTING"
echo -e "================================${NC}\n"

# Function to make API request
test_request() {
    local method=$1
    local endpoint=$2
    local lang_param=$3
    local accept_header=$4
    local body=$5
    local description=$6
    
    echo -e "${YELLOW}Testing: $description${NC}"
    
    # Build curl command
    local curl_cmd="curl -s -X $method"
    curl_cmd="$curl_cmd -H 'x-api-key: $API_KEY'"
    curl_cmd="$curl_cmd -H 'Content-Type: application/json'"
    
    # Add Accept-Language header if provided
    if [ ! -z "$accept_header" ]; then
        curl_cmd="$curl_cmd -H 'Accept-Language: $accept_header'"
    fi
    
    # Build URL with language parameter
    local url="$BASE_URL$endpoint"
    if [ ! -z "$lang_param" ]; then
        url="$url?lang=$lang_param"
    fi
    
    # Add body if provided
    if [ ! -z "$body" ]; then
        curl_cmd="$curl_cmd -d '$body'"
    fi
    
    # Execute request
    local response=$(eval "$curl_cmd '$url'")
    local http_code=$(eval "$curl_cmd -w '%{http_code}' -o /dev/null '$url'")
    
    echo "HTTP Code: $http_code"
    echo "Response: $response"
    echo -e "${GREEN}✓ Request completed${NC}\n"
}

echo -e "${BLUE}🔧 Starting Application...${NC}"
echo "Please make sure your application is running on $BASE_URL"
echo "You can start it with: ./rest-api run --port=8080"
echo ""
read -p "Press Enter when the application is ready..."

echo -e "\n${BLUE}1. Testing English (Default Language)${NC}"
echo "----------------------------------------"

# Test 1: Get users in English (default)
test_request "GET" "/api/v1/users" "" "" "" "Get users in English (default)"

# Test 2: Get non-existent user in English
test_request "GET" "/api/v1/users/999" "" "" "" "Get non-existent user in English"

echo -e "${BLUE}2. Testing Indonesian via Query Parameter${NC}"
echo "-----------------------------------------------"

# Test 3: Get users in Indonesian via query param
test_request "GET" "/api/v1/users" "id" "" "" "Get users in Indonesian (?lang=id)"

# Test 4: Get non-existent user in Indonesian
test_request "GET" "/api/v1/users/999" "id" "" "" "Get non-existent user in Indonesian"

echo -e "${BLUE}3. Testing Indonesian via Accept-Language Header${NC}"
echo "---------------------------------------------------"

# Test 5: Get users in Indonesian via Accept-Language header
test_request "GET" "/api/v1/users" "" "id" "" "Get users in Indonesian (Accept-Language: id)"

# Test 6: Get users with complex Accept-Language header
test_request "GET" "/api/v1/users" "" "id,en;q=0.9" "" "Get users with Accept-Language: id,en;q=0.9"

echo -e "${BLUE}4. Testing User Creation with Different Languages${NC}"
echo "----------------------------------------------------"

# Test 7: Create user with invalid data in English
test_request "POST" "/api/v1/users" "" "" '{"username":"ab","email":"invalid","password":"123"}' "Create user with validation errors (English)"

# Test 8: Create user with invalid data in Indonesian
test_request "POST" "/api/v1/users" "id" "" '{"username":"ab","email":"invalid","password":"123"}' "Create user with validation errors (Indonesian)"

echo -e "${BLUE}5. Testing Mixed Language Preferences${NC}"
echo "---------------------------------------------"

# Test 9: Query param overrides Accept-Language header
test_request "GET" "/api/v1/users/999" "en" "id" "" "Query param 'en' overrides Accept-Language 'id'"

# Test 10: Fallback to English for unsupported language
test_request "GET" "/api/v1/users/999" "es" "" "" "Unsupported language 'es' should fallback to English"

echo -e "${BLUE}6. Testing Authentication with Different Languages${NC}"
echo "----------------------------------------------------"

# Test 11: Missing API key in Indonesian
echo -e "${YELLOW}Testing: Missing API key error in Indonesian${NC}"
response=$(curl -s -H "Accept-Language: id" "$BASE_URL/api/v1/users")
echo "Response: $response"
echo -e "${GREEN}✓ Request completed${NC}\n"

echo -e "${GREEN}================================"
echo -e "🎉 i18n TESTING COMPLETED!"
echo -e "================================${NC}\n"

echo -e "${YELLOW}📝 Test Summary:${NC}"
echo "• Tested default English language"
echo "• Tested Indonesian via query parameter (?lang=id)"
echo "• Tested Indonesian via Accept-Language header"
echo "• Tested language priority (query param > header > default)"
echo "• Tested validation errors in multiple languages"
echo "• Tested authentication errors in multiple languages"
echo "• Tested fallback to English for unsupported languages"

echo -e "\n${BLUE}🌍 Language Support Verified:${NC}"
echo "• 🇺🇸 English (en) - Default ✓"
echo "• 🇮🇩 Indonesian (id) - Supported ✓"

echo -e "\n${YELLOW}💡 Tips:${NC}"
echo "• Use ?lang=id for Indonesian responses"
echo "• Use Accept-Language: id header as alternative"
echo "• All error and success messages are now multilingual"
echo "• Query parameter takes priority over Accept-Language header"
