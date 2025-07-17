//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestPackageManagerBasics tests core package manager functionality
func TestPackageManagerBasics(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Build the package manager binary
	pmPath := buildPackageManagerBinary(t, tempDir)
	
	// Change to project root so tools.json can be found
	wd, _ := os.Getwd()
	projectRoot := filepath.Dir(wd)
	os.Chdir(projectRoot)

	// Test 1: Hello command (system check)
	t.Run("HelloCommand", func(t *testing.T) {
		cmd := exec.Command(pmPath, "hello")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Fatalf("Hello command failed: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "NimsForest Package Manager") {
			t.Error("Hello output should contain package manager information")
		}
	})

	// Test 2: Status command
	t.Run("StatusCommand", func(t *testing.T) {
		cmd := exec.Command(pmPath, "status")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Fatalf("Status command failed: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "NimsForest Tools Status") {
			t.Error("Status should show tools status")
		}
		
		// Should show available tools
		if !strings.Contains(outputStr, "Available tools:") {
			t.Error("Status should list available tools")
		}
	})

	// Test 3: Help command
	t.Run("HelpCommand", func(t *testing.T) {
		cmd := exec.Command(pmPath, "--help")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Fatalf("Help command failed: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		expectedCommands := []string{"install", "update", "status", "hello", "validate"}
		
		for _, cmd := range expectedCommands {
			if !strings.Contains(outputStr, cmd) {
				t.Errorf("Help output missing command: %s", cmd)
			}
		}
	})

	// Test 4: Invalid command handling
	t.Run("InvalidCommand", func(t *testing.T) {
		cmd := exec.Command(pmPath, "nonexistent-command")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error for invalid command")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "unknown command") {
			t.Error("Should show unknown command error")
		}
	})

	// Test 5: Install command validation
	t.Run("InstallValidation", func(t *testing.T) {
		// Test install with no arguments
		cmd := exec.Command(pmPath, "install")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error when no tools specified")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "requires at least 1 arg") {
			t.Error("Should indicate required arguments")
		}
	})

	// Test 6: Install unknown tool
	t.Run("InstallUnknownTool", func(t *testing.T) {
		cmd := exec.Command(pmPath, "install", "nonexistent-tool")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error for unknown tool")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "unknown tool") {
			t.Error("Should indicate unknown tool error")
		}
	})
}

// TestPackageManagerToolRegistry tests the tool registry functionality
func TestPackageManagerToolRegistry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	pmPath := buildPackageManagerBinary(t, tempDir)
	
	// Change to project root so tools.json can be found
	wd, _ := os.Getwd()
	projectRoot := filepath.Dir(wd)
	os.Chdir(projectRoot)

	t.Run("ToolRegistry", func(t *testing.T) {
		cmd := exec.Command(pmPath, "status")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Fatalf("Status command failed: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		
		// Should show expected tools from registry
		expectedTools := []string{"workspace", "organize", "work", "communicate"}
		for _, tool := range expectedTools {
			if !strings.Contains(outputStr, tool) {
				t.Errorf("Registry should contain tool: %s", tool)
			}
		}
	})
}

// TestPackageManagerValidation tests tool validation functionality
func TestPackageManagerValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	pmPath := buildPackageManagerBinary(t, tempDir)
	
	// Change to project root so tools.json can be found
	wd, _ := os.Getwd()
	projectRoot := filepath.Dir(wd)
	os.Chdir(projectRoot)

	t.Run("ValidateCommand", func(t *testing.T) {
		// Test validate with no arguments
		cmd := exec.Command(pmPath, "validate")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error when no tool specified for validation")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "accepts 1 arg") {
			t.Error("Should indicate required tool name")
		}
	})

	t.Run("ValidateNonexistentTool", func(t *testing.T) {
		cmd := exec.Command(pmPath, "validate", "nonexistent-tool")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error for nonexistent tool")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "not installed") {
			t.Error("Should indicate tool not installed")
		}
	})
}

// buildPackageManagerBinary builds the package manager binary for testing
func buildPackageManagerBinary(t *testing.T, tempDir string) string {
	binaryPath := filepath.Join(tempDir, "nimsforestpm")
	
	// Get the project root directory (parent of integration)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	
	// Navigate to project root (parent of integration directory)
	projectRoot := filepath.Dir(wd)
	
	// Build the binary
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd")
	cmd.Dir = projectRoot
	
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build package manager binary: %v\nOutput: %s", err, string(output))
	}
	
	return binaryPath
}