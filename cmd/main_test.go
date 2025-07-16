package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectWorkspace(t *testing.T) {
	// Create temporary directory structure for testing
	tempDir := t.TempDir()
	
	// Test case 1: Valid workspace structure
	orgWorkspace := filepath.Join(tempDir, "test-organization-workspace")
	productsWorkspace := filepath.Join(tempDir, "products-workspace")
	
	if err := os.MkdirAll(orgWorkspace, 0755); err != nil {
		t.Fatalf("Failed to create org workspace: %v", err)
	}
	if err := os.MkdirAll(productsWorkspace, 0755); err != nil {
		t.Fatalf("Failed to create products workspace: %v", err)
	}
	
	// Change to temp directory and test detection
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tempDir)
	workspace, err := detectWorkspace()
	if err != nil {
		t.Errorf("Expected to detect workspace, got error: %v", err)
	}
	if workspace != tempDir {
		t.Errorf("Expected workspace %s, got %s", tempDir, workspace)
	}
	
	// Test case 2: No workspace structure
	emptyDir := t.TempDir()
	os.Chdir(emptyDir)
	_, err = detectWorkspace()
	if err == nil {
		t.Error("Expected error when no workspace structure exists")
	}
}

func TestHasWorkspaceStructure(t *testing.T) {
	// Test case 1: Valid structure
	tempDir := t.TempDir()
	orgWorkspace := filepath.Join(tempDir, "test-organization-workspace")
	productsWorkspace := filepath.Join(tempDir, "products-workspace")
	
	os.MkdirAll(orgWorkspace, 0755)
	os.MkdirAll(productsWorkspace, 0755)
	
	if !hasWorkspaceStructure(tempDir) {
		t.Error("Expected valid workspace structure to be detected")
	}
	
	// Test case 2: Missing products-workspace
	tempDir2 := t.TempDir()
	orgWorkspace2 := filepath.Join(tempDir2, "test-organization-workspace")
	os.MkdirAll(orgWorkspace2, 0755)
	
	if hasWorkspaceStructure(tempDir2) {
		t.Error("Expected invalid structure (missing products-workspace) to be rejected")
	}
	
	// Test case 3: Missing organization workspace
	tempDir3 := t.TempDir()
	productsWorkspace3 := filepath.Join(tempDir3, "products-workspace")
	os.MkdirAll(productsWorkspace3, 0755)
	
	if hasWorkspaceStructure(tempDir3) {
		t.Error("Expected invalid structure (missing org workspace) to be rejected")
	}
}

func TestAddDynamicCommands(t *testing.T) {
	// Test that addDynamicCommands doesn't error
	err := addDynamicCommands()
	if err != nil {
		t.Errorf("addDynamicCommands should not error: %v", err)
	}
}

func TestToolStruct(t *testing.T) {
	// Test Tool struct creation and properties
	tool := Tool{
		Name:        "work",
		FullName:    "nimsforestwork",
		Path:        "/test/path",
		Commands:    []string{"hello", "init"},
		Description: "Test tool",
	}
	
	if tool.Name != "work" {
		t.Errorf("Expected tool name 'work', got '%s'", tool.Name)
	}
	if tool.FullName != "nimsforestwork" {
		t.Errorf("Expected tool full name 'nimsforestwork', got '%s'", tool.FullName)
	}
	if len(tool.Commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(tool.Commands))
	}
}

func TestGetToolDescription(t *testing.T) {
	testCases := []struct {
		toolName string
		expected string
	}{
		{"nimsforestwork", "Work Management and Tracking System"},
		{"nimsforestorganize", "Organizational Structure Component"},
		{"nimsforestcommunication", "Communication Systems Component"},
		{"unknown-tool", "unknown-tool"}, // Should return tool name if not found
	}
	
	for _, tc := range testCases {
		result := getToolDescription(tc.toolName)
		if result != tc.expected {
			t.Errorf("For tool %s, expected %s, got %s", tc.toolName, tc.expected, result)
		}
	}
}

func TestGetAvailableToolsDescriptions(t *testing.T) {
	// Test that tool descriptions can be retrieved
	expectedTools := []string{"nimsforestwork", "nimsforestorganize", "nimsforestcommunication", "nimsforestproductize", "nimsforestfolders", "nimsforestwebstack"}
	
	for _, tool := range expectedTools {
		desc := getToolDescription(tool)
		if desc == "" {
			t.Errorf("Tool %s should have a non-empty description", tool)
		}
		// Description should not be the tool name itself for known tools
		if desc == tool {
			t.Errorf("Tool %s returned tool name as description, expected a proper description", tool)
		}
	}
	
	// Test unknown tool
	unknownDesc := getToolDescription("unknown-tool")
	if unknownDesc != "unknown-tool" {
		t.Errorf("Unknown tool should return tool name as description, got %s", unknownDesc)
	}
}