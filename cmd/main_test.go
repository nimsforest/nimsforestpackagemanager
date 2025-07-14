package main

import (
	"os"
	"path/filepath"
	"strings"
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

func TestDiscoverInstalledTools(t *testing.T) {
	// Create test workspace structure
	tempDir := t.TempDir()
	productsDir := filepath.Join(tempDir, "products-workspace")
	
	// Create mock tool workspaces
	toolDirs := []string{
		"nimsforestwork-workspace",
		"nimsforestcommunication-workspace",
		"non-nimsforest-workspace", // Should be ignored
	}
	
	for _, dir := range toolDirs {
		toolPath := filepath.Join(productsDir, dir, "main")
		if err := os.MkdirAll(toolPath, 0755); err != nil {
			t.Fatalf("Failed to create tool directory: %v", err)
		}
		
		// Create mock MAKEFILE for nimsforest tools only
		if strings.HasPrefix(dir, "nimsforest") {
			toolName := strings.TrimSuffix(dir, "-workspace")
			makefilePath := filepath.Join(toolPath, "MAKEFILE."+toolName)
			makefileContent := `# Test makefile
` + toolName + `-hello:
	@echo "Hello"
` + toolName + `-init:
	@echo "Init"
` + toolName + `-test:
	@echo "Test"
`
			if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
				t.Fatalf("Failed to create makefile: %v", err)
			}
		}
	}
	
	tools, err := discoverInstalledTools(tempDir)
	if err != nil {
		t.Fatalf("Failed to discover tools: %v", err)
	}
	
	// Should find 2 nimsforest tools, ignore the non-nimsforest one
	if len(tools) != 2 {
		t.Errorf("Expected 2 tools, found %d", len(tools))
	}
	
	// Check tool properties
	for _, tool := range tools {
		if !strings.HasPrefix(tool.FullName, "nimsforest") {
			t.Errorf("Tool %s should start with 'nimsforest'", tool.FullName)
		}
		if len(tool.Commands) == 0 {
			t.Errorf("Tool %s should have commands", tool.Name)
		}
		// Check for expected commands
		expectedCommands := []string{"hello", "init", "test"}
		for _, expected := range expectedCommands {
			found := false
			for _, cmd := range tool.Commands {
				if cmd == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Tool %s missing expected command %s", tool.Name, expected)
			}
		}
	}
}

func TestExtractCommands(t *testing.T) {
	// Create temporary makefile
	tempFile := filepath.Join(t.TempDir(), "MAKEFILE.test")
	content := `# Test makefile
.PHONY: test-hello test-init test-lint

test-hello:
	@echo "Hello"

test-init:
	@echo "Init"

test-lint:
	@echo "Lint"

# This should be ignored (no tool prefix)
some-other-command:
	@echo "Other"

# This should also be ignored (different tool)
other-tool-command:
	@echo "Different tool"
`
	
	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test makefile: %v", err)
	}
	
	commands := extractCommands(tempFile, "test")
	
	expectedCommands := []string{"hello", "init", "lint"}
	if len(commands) != len(expectedCommands) {
		t.Errorf("Expected %d commands, got %d", len(expectedCommands), len(commands))
	}
	
	for _, expected := range expectedCommands {
		found := false
		for _, cmd := range commands {
			if cmd == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command %s not found in %v", expected, commands)
		}
	}
}

func TestExtractCommandsWithNonexistentFile(t *testing.T) {
	commands := extractCommands("/nonexistent/file", "test")
	
	// Should return fallback commands
	expectedFallback := []string{"hello", "init", "lint"}
	if len(commands) != len(expectedFallback) {
		t.Errorf("Expected fallback commands, got %v", commands)
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

func TestGetAvailableTools(t *testing.T) {
	tools := getAvailableTools()
	
	expectedTools := []string{"work", "organize", "communicate", "productize", "folders", "webstack"}
	if len(tools) != len(expectedTools) {
		t.Errorf("Expected %d tools, got %d", len(expectedTools), len(tools))
	}
	
	for _, expected := range expectedTools {
		found := false
		for _, tool := range tools {
			if tool == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected tool %s not found in available tools", expected)
		}
	}
}