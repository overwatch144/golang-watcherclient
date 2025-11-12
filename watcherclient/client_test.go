package watcherclient

import (
	"testing"
	"time"
)

// Test 1: Client oluşturma (token ile)
func TestNewClientWithToken(t *testing.T) {
	endpoint := "http://localhost:9322"
	token := "test-token-123"

	client := NewClientWithToken(endpoint, token)

	if client == nil {
		t.Fatal("Client is nil")
	}

	expectedEndpoint := endpoint + "/v1"
	if client.GetEndpoint() != expectedEndpoint {
		t.Errorf("Expected endpoint %s, got %s", expectedEndpoint, client.GetEndpoint())
	}
}

// Test 2: API Version
func TestAPIVersion(t *testing.T) {
	client := NewClientWithToken("http://localhost:9322", "test-token")

	if client.GetAPIVersion() != "v1" {
		t.Errorf("Expected v1, got %s", client.GetAPIVersion())
	}

	client.SetAPIVersion("v2")
	if client.GetAPIVersion() != "v2" {
		t.Errorf("Expected v2 after set, got %s", client.GetAPIVersion())
	}
}

// Test 3: Timeout
func TestClientTimeout(t *testing.T) {
	client := NewClientWithToken("http://localhost:9322", "test-token")

	defaultTimeout := client.GetTimeout()
	if defaultTimeout != DefaultTimeout {
		t.Errorf("Expected default timeout %v, got %v", DefaultTimeout, defaultTimeout)
	}

	newTimeout := 60 * time.Second
	client.SetTimeout(newTimeout)
	if client.GetTimeout() != newTimeout {
		t.Errorf("Expected timeout %v, got %v", newTimeout, client.GetTimeout())
	}
}

// Test 4: TokenAuthenticator
func TestTokenAuthenticator(t *testing.T) {
	token := "my-test-token"
	endpoint := "http://test:9322"

	auth := NewTokenAuthenticator(endpoint, token)

	gotToken, err := auth.GetToken()
	if err != nil {
		t.Fatalf("GetToken failed: %v", err)
	}

	if gotToken != token {
		t.Errorf("Expected token %s, got %s", token, gotToken)
	}

	if auth.GetEndpoint() != endpoint {
		t.Errorf("Expected endpoint %s, got %s", endpoint, auth.GetEndpoint())
	}
}

// Test 5: Empty token error
func TestEmptyToken(t *testing.T) {
	auth := NewTokenAuthenticator("http://test:9322", "")

	_, err := auth.GetToken()
	if err == nil {
		t.Error("Expected error for empty token, got nil")
	}
}

// Test 6: Audit struct
func TestAuditStruct(t *testing.T) {
	audit := &Audit{
		Name:        "Test Audit",
		AuditType:   "ONESHOT",
		Goal:        "server_consolidation",
		AutoTrigger: true,
	}

	if audit.Name != "Test Audit" {
		t.Errorf("Expected name 'Test Audit', got %s", audit.Name)
	}

	if audit.AuditType != "ONESHOT" {
		t.Errorf("Expected type ONESHOT, got %s", audit.AuditType)
	}

	if audit.Goal != "server_consolidation" {
		t.Errorf("Expected goal server_consolidation, got %s", audit.Goal)
	}

	if !audit.AutoTrigger {
		t.Error("Expected AutoTrigger to be true")
	}
}

// Test 7: AuditTemplate struct
func TestAuditTemplateStruct(t *testing.T) {
	template := &AuditTemplate{
		Name:        "Test Template",
		Description: "Test Description",
		Goal:        "energy_efficiency",
		Strategy:    "basic",
	}

	if template.Name != "Test Template" {
		t.Errorf("Expected name 'Test Template', got %s", template.Name)
	}

	if template.Goal != "energy_efficiency" {
		t.Errorf("Expected goal energy_efficiency, got %s", template.Goal)
	}
}

// Test 8: Error helpers - IsNotFound
func TestIsNotFound(t *testing.T) {
	notFoundErr := &APIError{StatusCode: 404}
	if !IsNotFound(notFoundErr) {
		t.Error("IsNotFound should return true for 404")
	}

	otherErr := &APIError{StatusCode: 500}
	if IsNotFound(otherErr) {
		t.Error("IsNotFound should return false for non-404")
	}
}

// Test 9: Error helpers - IsUnauthorized
func TestIsUnauthorized(t *testing.T) {
	unauthorizedErr := &APIError{StatusCode: 401}
	if !IsUnauthorized(unauthorizedErr) {
		t.Error("IsUnauthorized should return true for 401")
	}

	otherErr := &APIError{StatusCode: 200}
	if IsUnauthorized(otherErr) {
		t.Error("IsUnauthorized should return false for non-401")
	}
}

// Test 10: Error helpers - IsForbidden
func TestIsForbidden(t *testing.T) {
	forbiddenErr := &APIError{StatusCode: 403}
	if !IsForbidden(forbiddenErr) {
		t.Error("IsForbidden should return true for 403")
	}

	otherErr := &APIError{StatusCode: 200}
	if IsForbidden(otherErr) {
		t.Error("IsForbidden should return false for non-403")
	}
}

// Test 11: Error helpers - IsConflict
func TestIsConflict(t *testing.T) {
	conflictErr := &APIError{StatusCode: 409}
	if !IsConflict(conflictErr) {
		t.Error("IsConflict should return true for 409")
	}

	otherErr := &APIError{StatusCode: 200}
	if IsConflict(otherErr) {
		t.Error("IsConflict should return false for non-409")
	}
}

// Test 12: APIError message
func TestAPIErrorMessage(t *testing.T) {
	err := &APIError{
		StatusCode: 404,
		Message:    "Not Found",
		URL:        "http://test/api/resource",
		Method:     "GET",
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error message should not be empty")
	}

	// Error message should contain status code
	if errMsg == "" {
		t.Error("Error message should contain information")
	}
}

// Test 13: ListOptions
func TestListOptions(t *testing.T) {
	opts := &ListOptions{
		Limit:   10,
		Marker:  "marker-123",
		SortKey: "created_at",
		SortDir: "desc",
	}

	if opts.Limit != 10 {
		t.Errorf("Expected limit 10, got %d", opts.Limit)
	}

	if opts.SortDir != "desc" {
		t.Errorf("Expected sort_dir desc, got %s", opts.SortDir)
	}
}

// Test 14: Default constants
func TestConstants(t *testing.T) {
	if DefaultAPIVersion != "v1" {
		t.Errorf("Expected DefaultAPIVersion v1, got %s", DefaultAPIVersion)
	}

	if DefaultTimeout != 30*time.Second {
		t.Errorf("Expected DefaultTimeout 30s, got %v", DefaultTimeout)
	}
}

// Test 15: SessionManager
func TestSessionManager(t *testing.T) {
	sm := NewSessionManager()

	if sm == nil {
		t.Fatal("SessionManager is nil")
	}

	// List should be empty initially
	sessions := sm.ListSessions()
	if len(sessions) != 0 {
		t.Errorf("Expected 0 sessions, got %d", len(sessions))
	}

	// Test GetSession for non-existent session
	_, err := sm.GetSession("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent session")
	}
}

// Test 16: Utility functions
func TestStringPtr(t *testing.T) {
	str := "test"
	ptr := StringPtr(str)

	if ptr == nil {
		t.Fatal("StringPtr returned nil")
	}

	if *ptr != str {
		t.Errorf("Expected %s, got %s", str, *ptr)
	}
}

func TestBoolPtr(t *testing.T) {
	val := true
	ptr := BoolPtr(val)

	if ptr == nil {
		t.Fatal("BoolPtr returned nil")
	}

	if *ptr != val {
		t.Errorf("Expected %v, got %v", val, *ptr)
	}
}

func TestIntPtr(t *testing.T) {
	val := 42
	ptr := IntPtr(val)

	if ptr == nil {
		t.Fatal("IntPtr returned nil")
	}

	if *ptr != val {
		t.Errorf("Expected %d, got %d", val, *ptr)
	}
}

// Test 17: Goal struct
func TestGoalStruct(t *testing.T) {
	goal := &Goal{
		Name:        "test_goal",
		DisplayName: "Test Goal",
	}

	if goal.Name != "test_goal" {
		t.Errorf("Expected name test_goal, got %s", goal.Name)
	}

	if goal.DisplayName != "Test Goal" {
		t.Errorf("Expected display name 'Test Goal', got %s", goal.DisplayName)
	}
}

// Test 18: Strategy struct
func TestStrategyStruct(t *testing.T) {
	strategy := &Strategy{
		Name:        "test_strategy",
		DisplayName: "Test Strategy",
		GoalUUID:    "goal-123",
	}

	if strategy.Name != "test_strategy" {
		t.Errorf("Expected name test_strategy, got %s", strategy.Name)
	}

	if strategy.GoalUUID != "goal-123" {
		t.Errorf("Expected goal UUID goal-123, got %s", strategy.GoalUUID)
	}
}

// Test 19: Action struct
func TestActionStruct(t *testing.T) {
	action := &Action{
		ActionType:     "migrate",
		State:          "PENDING",
		ActionPlanUUID: "plan-123",
	}

	if action.ActionType != "migrate" {
		t.Errorf("Expected action type migrate, got %s", action.ActionType)
	}

	if action.State != "PENDING" {
		t.Errorf("Expected state PENDING, got %s", action.State)
	}
}

// Test 20: ActionPlan struct
func TestActionPlanStruct(t *testing.T) {
	plan := &ActionPlan{
		AuditUUID: "audit-123",
		State:     "RECOMMENDED",
		Strategy:  "basic",
	}

	if plan.AuditUUID != "audit-123" {
		t.Errorf("Expected audit UUID audit-123, got %s", plan.AuditUUID)
	}

	if plan.State != "RECOMMENDED" {
		t.Errorf("Expected state RECOMMENDED, got %s", plan.State)
	}
}
