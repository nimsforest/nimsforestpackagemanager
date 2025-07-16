package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nimsforest/nimsforestpackagemanager/internal/workspace"
)

func TestFullBuildInstallValidateCycle(t *testing.T) {
	// Create a temporary workspace
	tempDir := t.TempDir()
	
	// Step 1: Build the example tool binary
	t.Log("Step 1: Building example tool binary...")
	binaryPath := filepath.Join(tempDir, "bin", "example-tool")
	err := os.MkdirAll(filepath.Dir(binaryPath), 0755)
	if err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}
	
	cmd := exec.Command("go", "build", "-o", binaryPath, "./example-tool")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build example tool: %v", err)
	}
	
	// Verify binary was created
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatal("Example tool binary was not created")
	}
	t.Log("✓ Binary built successfully")

	// Step 2: Create workspace and install tool
	t.Log("Step 2: Creating workspace and installing tool...")
	workspaceFile := filepath.Join(tempDir, "nimsforest.workspace")
	ws := workspace.NewWorkspace()
	ws.Version = "1.0"
	ws.Organization = "./test-org"
	
	// Add tool to workspace (this simulates the install process)
	ws.AddTool(workspace.ToolEntry{
		Name:    "example-tool",
		Mode:    "binary",
		Path:    "bin/example-tool",
		Version: "latest",
	})

	if err := ws.Save(workspaceFile); err != nil {
		t.Fatalf("Failed to save workspace: %v", err)
	}
	t.Log("✓ Tool installed in workspace")

	// Step 3: Validate installation through package manager
	t.Log("Step 3: Validating installation through package manager...")
	
	// Build the main package manager binary
	pmBinaryPath := filepath.Join(tempDir, "nimsforestpm")
	cmd = exec.Command("go", "build", "-o", pmBinaryPath, "../cmd")
	cmd.Dir = "."
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build package manager: %v\nOutput: %s", err, output)
	}
	
	// Run status command from workspace directory
	cmd = exec.Command(pmBinaryPath, "status")
	cmd.Dir = tempDir
	output, err = cmd.Output()
	if err != nil {
		t.Fatalf("Status command failed: %v", err)
	}
	
	statusOutput := string(output)
	
	// Verify workspace was detected
	if !strings.Contains(statusOutput, "nimsforest.workspace") {
		t.Error("Workspace file not detected in status output")
	}
	
	// Verify tool is listed
	if !strings.Contains(statusOutput, "example-tool") {
		t.Error("Example tool not found in status output")
	}
	
	// Verify tool mode is binary
	if !strings.Contains(statusOutput, "binary") {
		t.Error("Tool mode not correctly shown as binary")
	}
	
	t.Log("✓ Package manager correctly detected installed tool")

	// Step 4: Validate tool commands work through workspace
	t.Log("Step 4: Validating tool commands...")
	
	// Test direct binary execution
	testCases := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "hello command",
			args:     []string{"hello"},
			expected: "Hello, World!",
		},
		{
			name:     "version command", 
			args:     []string{"version"},
			expected: "example version 1.0.0",
		},
		{
			name:     "status command",
			args:     []string{"status"},
			expected: "Tool: example",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tc.args...)
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}

			if !strings.Contains(string(output), tc.expected) {
				t.Errorf("Expected output to contain %q, got %q", tc.expected, string(output))
			}
		})
	}
	
	t.Log("✓ All tool commands validated successfully")
	
	// Step 5: Test workspace validation
	t.Log("Step 5: Testing workspace validation...")
	
	// Load and validate workspace
	loadedWs, err := workspace.LoadWorkspace(workspaceFile)
	if err != nil {
		t.Fatalf("Failed to load workspace: %v", err)
	}
	
	// Get tool from workspace
	toolEntry, err := loadedWs.GetTool("example-tool")
	if err != nil {
		t.Fatalf("Tool not found in workspace: %v", err)
	}
	
	// Validate tool properties
	if toolEntry.Name != "example-tool" {
		t.Errorf("Expected tool name 'example-tool', got %s", toolEntry.Name)
	}
	
	if toolEntry.Mode != "binary" {
		t.Errorf("Expected tool mode 'binary', got %s", toolEntry.Mode)
	}
	
	if toolEntry.Path != "bin/example-tool" {
		t.Errorf("Expected tool path 'bin/example-tool', got %s", toolEntry.Path)
	}
	
	t.Log("✓ Workspace validation completed successfully")
	
	t.Log("🎉 Full build → install → validate cycle completed successfully!")
}

func TestPackageManagerToolIntegration(t *testing.T) {
	// Test that the package manager can work with tools built using pkg/tool
	tempDir := t.TempDir()
	
	// Build example tool
	binaryPath := filepath.Join(tempDir, "bin", "example-tool")
	err := os.MkdirAll(filepath.Dir(binaryPath), 0755)
	if err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}
	
	cmd := exec.Command("go", "build", "-o", binaryPath, "./example-tool")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build example tool: %v", err)
	}
	
	// Create workspace
	workspaceFile := filepath.Join(tempDir, "nimsforest.workspace")
	ws := workspace.NewWorkspace()
	ws.Version = "1.0"
	ws.AddTool(workspace.ToolEntry{
		Name:    "example-tool",
		Mode:    "binary",
		Path:    "bin/example-tool", 
		Version: "1.0.0",
	})
	
	if err := ws.Save(workspaceFile); err != nil {
		t.Fatalf("Failed to save workspace: %v", err)
	}
	
	// Build package manager
	pmBinaryPath := filepath.Join(tempDir, "nimsforestpm")
	cmd = exec.Command("go", "build", "-o", pmBinaryPath, "../cmd")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build package manager: %v", err)
	}
	
	// Test that package manager can discover and work with the tool
	cmd = exec.Command(pmBinaryPath, "status")
	cmd.Dir = tempDir
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Status command failed: %v", err)
	}
	
	// Verify integration
	statusOutput := string(output)
	if !strings.Contains(statusOutput, "example-tool") {
		t.Error("Package manager did not detect the tool")
	}
	
	if !strings.Contains(statusOutput, "binary") {
		t.Error("Package manager did not correctly identify install mode")
	}
	
	t.Log("✓ Package manager successfully integrated with pkg/tool-based binary")
}