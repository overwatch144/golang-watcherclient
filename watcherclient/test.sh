#!/bin/bash

# Golang Watcher Client - Test Script
# Bu script tüm testleri çalıştırır

echo "======================================"
echo "  Golang Watcher Client - Test Suite"
echo "======================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counter
TESTS_PASSED=0
TESTS_FAILED=0

# Function to run test
run_test() {
    local test_name=$1
    local test_command=$2
    
    echo -e "${YELLOW}Running: ${test_name}${NC}"
    if eval $test_command; then
        echo -e "${GREEN}✓ ${test_name} PASSED${NC}"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗ ${test_name} FAILED${NC}"
        ((TESTS_FAILED++))
    fi
    echo ""
}

# 1. Go Vet
run_test "Go Vet" "go vet ./..."

# 2. Go Format Check
run_test "Go Format" "test -z \$(gofmt -l .)"

# 3. Build
run_test "Build" "go build ./..."

# 4. Unit Tests
run_test "Unit Tests" "go test ./watcherclient -v"

# 5. Test Coverage
echo -e "${YELLOW}Running: Test Coverage${NC}"
if go test ./watcherclient -coverprofile=coverage.out > /dev/null 2>&1; then
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    echo -e "${GREEN}✓ Test Coverage: ${COVERAGE}${NC}"
    ((TESTS_PASSED++))
else
    echo -e "${RED}✗ Test Coverage FAILED${NC}"
    ((TESTS_FAILED++))
fi
echo ""

# Summary
echo "======================================"
echo "  Test Summary"
echo "======================================"
echo -e "${GREEN}Passed: ${TESTS_PASSED}${NC}"
echo -e "${RED}Failed: ${TESTS_FAILED}${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed!${NC}"
    exit 1
fi