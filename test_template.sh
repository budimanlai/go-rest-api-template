#!/bin/bash

# 🧪 REST API Template Testing Script
# This script tests all aspects of the project template

echo "🚀 Starting comprehensive testing for REST API Template..."
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results tracking
TESTS_PASSED=0
TESTS_FAILED=0

function test_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ PASS${NC}: $2"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}❌ FAIL${NC}: $2"
        ((TESTS_FAILED++))
    fi
}

echo -e "${BLUE}1. Testing Go Build...${NC}"
go build -v ./cmd/api/ > /dev/null 2>&1
test_result $? "Go build compilation"

echo -e "${BLUE}2. Testing Go Modules...${NC}"
go mod tidy > /dev/null 2>&1
test_result $? "Go modules tidy"

echo -e "${BLUE}3. Testing Unit Tests...${NC}"
go test ./test/handler_test.go > /dev/null 2>&1
test_result $? "Handler unit tests"

echo -e "${BLUE}4. Testing Project Structure...${NC}"
# Check key files exist
if [ -f "./internal/application/container.go" ] && 
   [ -f "./internal/routes/route_manager.go" ] && 
   [ -f "./internal/service/user_service.go" ] &&
   [ -f "./internal/handler/user_handler.go" ]; then
    test_result 0 "Project structure"
else
    test_result 1 "Project structure"
fi

echo -e "${BLUE}5. Testing Clean Architecture Layers...${NC}"
# Check if layers are properly separated
if [ -d "./internal/domain/entity" ] && 
   [ -d "./internal/domain/repository" ] && 
   [ -d "./internal/domain/usecase" ] &&
   [ -d "./internal/repository" ] &&
   [ -d "./internal/service" ] &&
   [ -d "./internal/handler" ]; then
    test_result 0 "Clean Architecture layers"
else
    test_result 1 "Clean Architecture layers"
fi

echo -e "${BLUE}6. Testing Configuration Files...${NC}"
if [ -f "./configs/config.json" ] && 
   [ -f "./scripts/init.sql" ] && 
   [ -f "./go.mod" ] &&
   [ -f "./Dockerfile" ] &&
   [ -f "./docker-compose.yml" ]; then
    test_result 0 "Configuration files"
else
    test_result 1 "Configuration files"
fi

echo -e "${BLUE}7. Testing Docker Setup...${NC}"
if docker --version > /dev/null 2>&1; then
    # Test if Docker build works (without actually building)
    if [ -f "./Dockerfile" ]; then
        test_result 0 "Docker configuration"
    else
        test_result 1 "Docker configuration"
    fi
else
    echo -e "${YELLOW}⚠️  SKIP${NC}: Docker not available"
fi

echo -e "${BLUE}8. Testing Documentation...${NC}"
if [ -f "./README.md" ] && 
   [ -f "./docs/CONTAINER_GUIDE.md" ]; then
    test_result 0 "Documentation files"
else
    test_result 1 "Documentation files"
fi

echo -e "${BLUE}9. Testing Database Schema...${NC}"
if [ -f "./scripts/init.sql" ]; then
    # Check if SQL contains required tables and fields
    if grep -q "password_hash" ./scripts/init.sql && 
       grep -q "reset_password_token" ./scripts/init.sql &&
       grep -q "created_by" ./scripts/init.sql; then
        test_result 0 "Database schema"
    else
        test_result 1 "Database schema"
    fi
else
    test_result 1 "Database schema"
fi

echo -e "${BLUE}10. Testing API Security...${NC}"
# Check if API key validation exists
if grep -q "x-api-key" ./internal/application/rest_api.go; then
    test_result 0 "API security implementation"
else
    test_result 1 "API security implementation"
fi

echo ""
echo "=================================================="
echo -e "${GREEN}📊 Test Summary:${NC}"
echo -e "   ✅ Passed: ${GREEN}$TESTS_PASSED${NC}"
echo -e "   ❌ Failed: ${RED}$TESTS_FAILED${NC}"
echo -e "   📈 Success Rate: $(( TESTS_PASSED * 100 / (TESTS_PASSED + TESTS_FAILED) ))%"

if [ $TESTS_FAILED -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 All tests passed! Your REST API template is ready for development!${NC}"
    echo ""
    echo -e "${YELLOW}📋 Next Steps:${NC}"
    echo "   1. Update database configuration in configs/config.json"
    echo "   2. Run database migration: mysql < scripts/init.sql"
    echo "   3. Start the server: go run ./cmd/api/ --port=8080"
    echo "   4. Test endpoints with: curl -H 'x-api-key: test-api-key' http://localhost:8080/api/v1/users"
    exit 0
else
    echo ""
    echo -e "${RED}⚠️  Some tests failed. Please check the issues above.${NC}"
    exit 1
fi
