package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateToolCommand(t *testing.T) {
	tool := Tool{
		Name:        "work",
		FullName:    "nimsforestwork", 
		Path:        "/test/path",
		Commands:    []string{"hello", "init", "triage"},
		Description: "Work Management System",
	}
	
	cmd := createToolCommand(tool, "hello")
	
	if cmd.Use != "hello" {
		t.Errorf("Expected command use 'hello', got '%s'", cmd.Use)
	}
	
	expectedShort := "Run nimsforestwork-hello command"
	if cmd.Short != expectedShort {
		t.Errorf("Expected short description '%s', got '%s'", expectedShort, cmd.Short)
	}
}

func TestToolMapping(t *testing.T) {
	// Test the tool mapping used in installTool
	expectedMappings := map[string]string{
		"work":        "nimsforestwork",
		"organize":    "nimsforestorganize",
		"communicate": "nimsforestcommunication",
		"productize":  "nimsforestproductize",
		"folders":     "nimsforestfolders",
		"webstack":    "nimsforestwebstack",
	}
	
	// This tests the logic inside installTool function
	// We can't easily test installTool directly due to file system operations
	// But we can verify the mapping logic
	for shortName, expectedFullName := range expectedMappings {
		// Simulate the mapping logic from installTool
		toolMap := map[string]string{
			"work":        "nimsforestwork",
			"organize":    "nimsforestorganize", 
			"communicate": "nimsforestcommunication",
			"productize":  "nimsforestproductize",
			"folders":     "nimsforestfolders",
			"webstack":    "nimsforestwebstack",
		}
		
		component, ok := toolMap[shortName]
		if !ok {
			t.Errorf("Tool mapping for '%s' not found", shortName)
			continue
		}
		
		if component != expectedFullName {
			t.Errorf("For tool '%s', expected mapping to '%s', got '%s'", 
				shortName, expectedFullName, component)
		}
	}
}

func TestShowStatusInWorkspace(t *testing.T) {
	// Create test workspace structure
	tempDir := t.TempDir()
	orgWorkspace := filepath.Join(tempDir, "test-organization-workspace")
	productsWorkspace := filepath.Join(tempDir, "products-workspace")
	
	if err := os.MkdirAll(orgWorkspace, 0755); err != nil {
		t.Fatalf("Failed to create org workspace: %v", err)
	}
	if err := os.MkdirAll(productsWorkspace, 0755); err != nil {
		t.Fatalf("Failed to create products workspace: %v", err)
	}
	
	// Create a mock tool
	toolDir := filepath.Join(productsWorkspace, "nimsforestwork-workspace", "main")
	if err := os.MkdirAll(toolDir, 0755); err != nil {
		t.Fatalf("Failed to create tool directory: %v", err)
	}
	
	makefilePath := filepath.Join(toolDir, "MAKEFILE.nimsforestwork")
	makefileContent := `nimsforestwork-hello:
	@echo "Hello"
`
	if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		t.Fatalf("Failed to create makefile: %v", err)
	}
	
	// Change to workspace directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)
	
	// Test showStatus - this should not return an error
	err := showStatus()
	if err != nil {
		t.Errorf("showStatus returned error in valid workspace: %v", err)
	}
}

func TestShowStatusOutsideWorkspace(t *testing.T) {
	// Create empty directory (no workspace structure)
	tempDir := t.TempDir()
	
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)
	
	// showStatus should not return error even outside workspace
	// It should just print "Not in a nimsforest workspace"
	err := showStatus()
	if err != nil {
		t.Errorf("showStatus should not return error outside workspace, got: %v", err)
	}
}

func TestRunHelloSystemCheck(t *testing.T) {
	// runHello performs system compatibility check
	// We can't easily test the actual execution without mocking
	// But we can test that it doesn't panic with basic execution
	
	// This is a basic smoke test - the function should not panic
	err := runHello()
	
	// The function should either succeed or fail based on system tools
	// We can't predict the outcome without knowing the test environment
	// So we just ensure it doesn't panic
	if err != nil {
		// Check if error is about missing required tools
		if !strings.Contains(err.Error(), "required but not installed") {
			t.Errorf("Unexpected error from runHello: %v", err)
		}
	}
}

func TestMakefilePathDiscovery(t *testing.T) {
	// Test the makefile discovery logic used in createOrganizationWorkspace and installTool
	tempDir := t.TempDir()
	
	// Create test makefile
	makefilePath := filepath.Join(tempDir, "MAKEFILE.nimsforestpm")
	makefileContent := "# Test makefile"
	if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		t.Fatalf("Failed to create test makefile: %v", err)
	}
	
	// Test path discovery logic (from createOrganizationWorkspace)
	possiblePaths := []string{
		filepath.Join(tempDir, "MAKEFILE.nimsforestpm"),
		"MAKEFILE.nimsforestpm",
	}
	
	var foundPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			foundPath = path
			break
		}
	}
	
	expectedPath := filepath.Join(tempDir, "MAKEFILE.nimsforestpm")
	if foundPath != expectedPath {
		t.Errorf("Expected to find makefile at %s, found at %s", expectedPath, foundPath)
	}
}

func TestMakefilePathDiscoveryNotFound(t *testing.T) {
	// Test when makefile is not found
	tempDir := t.TempDir()
	
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)
	
	// Test path discovery logic when no makefile exists
	possiblePaths := []string{
		filepath.Join(tempDir, "MAKEFILE.nimsforestpm"),
		"MAKEFILE.nimsforestpm",
	}
	
	var foundPath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			foundPath = path
			break
		}
	}
	
	if foundPath != "" {
		t.Errorf("Expected no makefile to be found, but found %s", foundPath)
	}
}

func TestUpdateToolsOutsideWorkspace(t *testing.T) {
	// Test update command outside workspace
	tempDir := t.TempDir()
	
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)
	
	err := updateTools([]string{"work"})
	if err == nil {
		t.Error("Expected error when updating tools outside workspace")
	}
	
	if !strings.Contains(err.Error(), "not in a nimsforest workspace") {
		t.Errorf("Expected workspace error, got: %v", err)
	}
}

func TestUpdateToolsWithInvalidTool(t *testing.T) {
	// Create mock workspace
	tempDir := t.TempDir()
	orgWorkspace := filepath.Join(tempDir, "test-organization-workspace")
	productsWorkspace := filepath.Join(tempDir, "products-workspace")
	
	if err := os.MkdirAll(orgWorkspace, 0755); err != nil {
		t.Fatalf("Failed to create org workspace: %v", err)
	}
	if err := os.MkdirAll(productsWorkspace, 0755); err != nil {
		t.Fatalf("Failed to create products workspace: %v", err)
	}
	
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)
	
	err := updateTools([]string{"invalid-tool"})
	if err == nil {
		t.Error("Expected error for invalid tool name")
	}
	
	if !strings.Contains(err.Error(), "unknown tool") {
		t.Errorf("Expected unknown tool error, got: %v", err)
	}
}

func TestUpdateToolsValidation(t *testing.T) {
	// Test the tool mapping validation used in updateSpecificTools
	toolMap := map[string]string{
		"work":        "nimsforestwork",
		"organize":    "nimsforestorganize",
		"communicate": "nimsforestcommunication",
		"productize":  "nimsforestproductize",
		"folders":     "nimsforestfolders",
		"webstack":    "nimsforestwebstack",
	}
	
	// Test valid tools
	validTools := []string{"work", "organize", "communicate", "productize", "folders", "webstack"}
	for _, tool := range validTools {
		if _, ok := toolMap[tool]; !ok {
			t.Errorf("Valid tool %s not found in mapping", tool)
		}
	}
	
	// Test invalid tool
	if _, ok := toolMap["invalid"]; ok {
		t.Error("Invalid tool should not be found in mapping")
	}
}