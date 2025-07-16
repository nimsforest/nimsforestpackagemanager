package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nimsforest/nimsforestpackagemanager/internal/workspace"
)

func TestPreBuiltBinaryCycle(t *testing.T) {
	// This test assumes binaries are already built via task system
	// Usage: task build-and-test
	
	// Check if binaries exist (using absolute paths)
	wd, _ := os.Getwd()
	pmBinary := filepath.Join(filepath.Dir(wd), "bin", "nimsforestpm")
	exampleBinary := filepath.Join(filepath.Dir(wd), "bin", "example-tool")
	
	if _, err := os.Stat(pmBinary); os.IsNotExist(err) {
		t.Skip("Package manager binary not found. Run 'task build' first.")
	}
	
	if _, err := os.Stat(exampleBinary); os.IsNotExist(err) {
		t.Skip("Example tool binary not found. Run 'task build-example' first.")
	}
	
	// Step 1: Test example tool works
	t.Log("Step 1: Testing example tool binary...")
	cmd := exec.Command(exampleBinary, "hello")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Example tool failed: %v", err)
	}
	if !strings.Contains(string(output), "Hello, World!") {
		t.Errorf("Expected 'Hello, World!' but got: %s", output)
	}
	t.Log("✓ Example tool binary works")
	
	// Step 2: Create workspace and add tool
	t.Log("Step 2: Creating workspace...")
	tempDir := t.TempDir()
	
	// Copy example binary to temp workspace
	tempBinary := filepath.Join(tempDir, "bin", "example-tool")
	err = os.MkdirAll(filepath.Dir(tempBinary), 0755)
	if err != nil {
		t.Fatalf("Failed to create temp bin dir: %v", err)
	}
	
	// Copy the binary
	srcFile, err := os.Open(exampleBinary)
	if err != nil {
		t.Fatalf("Failed to open source binary: %v", err)
	}
	
	dstFile, err := os.Create(tempBinary)
	if err != nil {
		srcFile.Close()
		t.Fatalf("Failed to create dest binary: %v", err)
	}
	
	_, err = dstFile.ReadFrom(srcFile)
	srcFile.Close()
	dstFile.Close()
	if err != nil {
		t.Fatalf("Failed to copy binary: %v", err)
	}
	
	// Make it executable
	err = os.Chmod(tempBinary, 0755)
	if err != nil {
		t.Fatalf("Failed to make binary executable: %v", err)
	}
	
	// Create workspace file
	workspaceFile := filepath.Join(tempDir, "nimsforest.workspace")
	ws := workspace.NewWorkspace()
	ws.Version = "1.0"
	ws.Organization = "./test-org"
	ws.AddTool(workspace.ToolEntry{
		Name:    "example-tool",
		Mode:    "binary",
		Path:    "bin/example-tool",
		Version: "1.0.0",
	})
	
	if err := ws.Save(workspaceFile); err != nil {
		t.Fatalf("Failed to save workspace: %v", err)
	}
	t.Log("✓ Workspace created with tool installed")
	
	// Step 3: Test package manager can detect the tool
	t.Log("Step 3: Testing package manager status...")
	cmd = exec.Command(pmBinary, "status")
	cmd.Dir = tempDir
	output, err = cmd.Output()
	if err != nil {
		t.Fatalf("Package manager status failed: %v", err)
	}
	
	statusOutput := string(output)
	if !strings.Contains(statusOutput, "example-tool") {
		t.Errorf("Package manager did not detect example-tool in status output: %s", statusOutput)
	}
	
	if !strings.Contains(statusOutput, "binary") {
		t.Errorf("Package manager did not show binary mode: %s", statusOutput)
	}
	
	t.Log("✓ Package manager detected the tool")
	
	// Step 4: Test tool still works in workspace context
	t.Log("Step 4: Testing tool in workspace context...")
	cmd = exec.Command(tempBinary, "version")
	cmd.Dir = tempDir
	output, err = cmd.Output()
	if err != nil {
		t.Fatalf("Tool execution in workspace failed: %v", err)
	}
	
	if !strings.Contains(string(output), "example version 1.0.0") {
		t.Errorf("Expected version output but got: %s", output)
	}
	
	t.Log("✓ Tool works in workspace context")
	
	t.Log("🎉 Complete build → install → validate cycle successful!")
}

func TestTaskSystemIntegration(t *testing.T) {
	// Test that the task system properly builds and integrates everything
	
	// Check if task command is available
	if _, err := exec.LookPath("task"); err != nil {
		t.Skip("Task command not found. Please install task first.")
	}
	
	// Run just the build tasks to avoid recursion
	t.Log("Running task build...")
	cmd := exec.Command("task", "build")
	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Task build failed: %v", err)
	}
	
	t.Log("Running task build-example...")
	cmd = exec.Command("task", "build-example")
	_, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Task build-example failed: %v", err)
	}
	
	t.Log("Running task example-tool-test...")
	cmd = exec.Command("task", "example-tool-test")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Task example-tool-test failed: %v\nOutput: %s", err, output)
	}
	
	// Check that the output contains success messages
	outputStr := string(output)
	if !strings.Contains(outputStr, "Example tool binary tests passed!") {
		t.Errorf("Expected success message not found in output: %s", outputStr)
	}
	
	t.Log("✓ Task system integration successful")
}