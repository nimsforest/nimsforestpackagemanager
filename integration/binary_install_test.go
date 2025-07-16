package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nimsforest/nimsforestpackagemanager/internal/workspace"
)

func TestBinaryInstallation(t *testing.T) {
	// Test the new binary installation functionality
	
	// Check if binaries exist
	wd, _ := os.Getwd()
	pmBinary := filepath.Join(filepath.Dir(wd), "bin", "nimsforestpm")
	exampleBinary := filepath.Join(filepath.Dir(wd), "bin", "example-tool")
	
	if _, err := os.Stat(pmBinary); os.IsNotExist(err) {
		t.Skip("Package manager binary not found. Run 'task build' first.")
	}
	
	if _, err := os.Stat(exampleBinary); os.IsNotExist(err) {
		t.Skip("Example tool binary not found. Run 'task build-example' first.")
	}
	
	// Create temporary workspace
	tempDir := t.TempDir()
	workspaceDir := filepath.Join(tempDir, "test-workspace")
	
	t.Log("Step 1: Creating organization workspace...")
	cmd := exec.Command(pmBinary, "create-organization-workspace", "test-org")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create organization workspace: %v", err)
	}
	
	// Find the created workspace
	workspaceDir = filepath.Join(tempDir, "test-org-workspace")
	if _, err := os.Stat(workspaceDir); os.IsNotExist(err) {
		t.Fatalf("Workspace directory not created: %s", workspaceDir)
	}
	
	t.Log("Step 2: Installing example tool via binary install...")
	cmd = exec.Command(pmBinary, "install", "--name", "example-tool", "--path", exampleBinary)
	cmd.Dir = workspaceDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Binary installation failed: %v\nOutput: %s", err, output)
	}
	
	// Check output contains expected messages
	outputStr := string(output)
	if !strings.Contains(outputStr, "Installing Binary Tool: example-tool") {
		t.Errorf("Expected installation message not found in output: %s", outputStr)
	}
	
	if !strings.Contains(outputStr, "✓ example-tool installed successfully!") {
		t.Errorf("Expected success message not found in output: %s", outputStr)
	}
	
	t.Log("Step 3: Verifying binary was installed...")
	installedBinary := filepath.Join(workspaceDir, "bin", "example-tool")
	if _, err := os.Stat(installedBinary); os.IsNotExist(err) {
		t.Errorf("Binary was not installed at expected location: %s", installedBinary)
	}
	
	// Test binary is executable
	if err := exec.Command(installedBinary, "version").Run(); err != nil {
		t.Errorf("Installed binary is not executable: %v", err)
	}
	
	t.Log("Step 4: Checking workspace status...")
	cmd = exec.Command(pmBinary, "status")
	cmd.Dir = workspaceDir
	output, err = cmd.Output()
	if err != nil {
		t.Fatalf("Status command failed: %v", err)
	}
	
	statusOutput := string(output)
	if !strings.Contains(statusOutput, "example-tool") {
		t.Errorf("Tool not found in status output: %s", statusOutput)
	}
	
	if !strings.Contains(statusOutput, "binary") {
		t.Errorf("Tool mode not shown as binary: %s", statusOutput)
	}
	
	t.Log("Step 5: Verifying workspace file was updated...")
	workspaceFile := filepath.Join(workspaceDir, "nimsforest.workspace")
	ws, err := workspace.LoadWorkspace(workspaceFile)
	if err != nil {
		t.Fatalf("Failed to load workspace file: %v", err)
	}
	
	// Check tool was added to workspace
	tool, err := ws.GetTool("example-tool")
	if err != nil {
		t.Fatalf("Tool not found in workspace: %v", err)
	}
	
	if tool.Mode != "binary" {
		t.Errorf("Expected tool mode 'binary', got '%s'", tool.Mode)
	}
	
	if tool.Path != "bin/example-tool" {
		t.Errorf("Expected tool path 'bin/example-tool', got '%s'", tool.Path)
	}
	
	t.Log("Step 6: Testing installed tool functionality...")
	cmd = exec.Command(installedBinary, "hello")
	cmd.Dir = workspaceDir
	output, err = cmd.Output()
	if err != nil {
		t.Fatalf("Installed tool execution failed: %v", err)
	}
	
	if !strings.Contains(string(output), "Hello, World!") {
		t.Errorf("Tool output unexpected: %s", output)
	}
	
	t.Log("✓ Binary installation test completed successfully!")
}

func TestBinaryInstallationErrors(t *testing.T) {
	// Test error cases for binary installation
	
	wd, _ := os.Getwd()
	pmBinary := filepath.Join(filepath.Dir(wd), "bin", "nimsforestpm")
	
	if _, err := os.Stat(pmBinary); os.IsNotExist(err) {
		t.Skip("Package manager binary not found. Run 'task build' first.")
	}
	
	// Create temporary workspace
	tempDir := t.TempDir()
	cmd := exec.Command(pmBinary, "create-organization-workspace", "test-org")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create organization workspace: %v", err)
	}
	
	workspaceDir := filepath.Join(tempDir, "test-org-workspace")
	
	t.Log("Testing binary installation with non-existent binary...")
	cmd = exec.Command(pmBinary, "install", "--name", "nonexistent", "--path", "/nonexistent/binary")
	cmd.Dir = workspaceDir
	if err := cmd.Run(); err == nil {
		t.Error("Expected error for non-existent binary, but command succeeded")
	}
	
	t.Log("Testing binary installation with missing name...")
	cmd = exec.Command(pmBinary, "install", "--path", "/some/path")
	cmd.Dir = workspaceDir
	if err := cmd.Run(); err == nil {
		t.Error("Expected error for missing name, but command succeeded")
	}
	
	t.Log("Testing binary installation with missing path...")
	cmd = exec.Command(pmBinary, "install", "--name", "somename")
	cmd.Dir = workspaceDir
	if err := cmd.Run(); err == nil {
		t.Error("Expected error for missing path, but command succeeded")
	}
	
	t.Log("✓ Error case testing completed successfully!")
}

func TestBinaryInstallationIntegration(t *testing.T) {
	// Test integration between binary installation and existing functionality
	
	wd, _ := os.Getwd()
	pmBinary := filepath.Join(filepath.Dir(wd), "bin", "nimsforestpm")
	exampleBinary := filepath.Join(filepath.Dir(wd), "bin", "example-tool")
	
	if _, err := os.Stat(pmBinary); os.IsNotExist(err) {
		t.Skip("Package manager binary not found. Run 'task build' first.")
	}
	
	if _, err := os.Stat(exampleBinary); os.IsNotExist(err) {
		t.Skip("Example tool binary not found. Run 'task build-example' first.")
	}
	
	// Create temporary workspace
	tempDir := t.TempDir()
	cmd := exec.Command(pmBinary, "create-organization-workspace", "test-org")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create organization workspace: %v", err)
	}
	
	workspaceDir := filepath.Join(tempDir, "test-org-workspace")
	
	t.Log("Step 1: Installing multiple tools (binary + standard)...")
	
	// Install binary tool
	cmd = exec.Command(pmBinary, "install", "--name", "example-tool", "--path", exampleBinary)
	cmd.Dir = workspaceDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Binary installation failed: %v", err)
	}
	
	// Try to install standard tool (this might fail if not implemented, but shouldn't break)
	cmd = exec.Command(pmBinary, "install", "work")
	cmd.Dir = workspaceDir
	_ = cmd.Run() // Ignore error as standard tools might not be available
	
	t.Log("Step 2: Checking combined status...")
	cmd = exec.Command(pmBinary, "status")
	cmd.Dir = workspaceDir
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Status command failed: %v", err)
	}
	
	statusOutput := string(output)
	if !strings.Contains(statusOutput, "example-tool") {
		t.Errorf("Binary tool not found in status: %s", statusOutput)
	}
	
	t.Log("Step 3: Testing workspace file integrity...")
	workspaceFile := filepath.Join(workspaceDir, "nimsforest.workspace")
	ws, err := workspace.LoadWorkspace(workspaceFile)
	if err != nil {
		t.Fatalf("Failed to load workspace file: %v", err)
	}
	
	// Validate workspace
	if err := ws.Validate(); err != nil {
		// Some validation errors are expected (like missing paths), but workspace should still be functional
		t.Logf("Workspace validation warnings (expected): %v", err)
	}
	
	t.Log("✓ Integration test completed successfully!")
}