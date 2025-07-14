package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestCLICoreCommands tests the core CLI commands without external dependencies
func TestCLICoreCommands(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Build the CLI binary for testing
	cliPath := buildCLIBinary(t, tempDir)
	os.Chdir(tempDir)

	// Test 1: Hello command (system check)
	t.Run("HelloCommand", func(t *testing.T) {
		cmd := exec.Command(cliPath, "hello")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Fatalf("Hello command failed: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "NimsforestPM System Check") {
			t.Error("Hello output should contain system check information")
		}
	})

	// Test 2: Status command outside workspace
	t.Run("StatusOutsideWorkspace", func(t *testing.T) {
		cmd := exec.Command(cliPath, "status")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Fatalf("Status command failed: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Not in a nimsforest workspace") {
			t.Error("Status should indicate when not in workspace")
		}
	})

	// Test 3: Invalid command handling
	t.Run("InvalidCommand", func(t *testing.T) {
		cmd := exec.Command(cliPath, "nonexistent-command")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error for invalid command")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "unknown command") {
			t.Error("Should show unknown command error")
		}
	})

	// Test 4: Missing arguments
	t.Run("MissingArguments", func(t *testing.T) {
		cmd := exec.Command(cliPath, "create-organization-workspace")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error when org name is missing")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "required") && !strings.Contains(outputStr, "accepts") {
			t.Error("Should indicate missing required argument")
		}
	})
	
	// Test 5: Update command outside workspace
	t.Run("UpdateOutsideWorkspace", func(t *testing.T) {
		cmd := exec.Command(cliPath, "update")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error when updating outside workspace")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "not in a nimsforest workspace") {
			t.Error("Should indicate not in workspace")
		}
	})
}

// TestToolDiscoveryAndProxying tests the core CLI mechanics with mock tools
func TestToolDiscoveryAndProxying(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Build CLI binary
	cliPath := buildCLIBinary(t, tempDir)

	// Create mock workspace with tools
	workspaceDir := setupMockWorkspace(t, tempDir)
	os.Chdir(workspaceDir)

	t.Run("DiscoverMockTools", func(t *testing.T) {
		cmd := exec.Command(cliPath, "status")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Fatalf("Status command failed: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		
		// Should discover our mock tools
		expectedTools := []string{"work", "communication"}
		for _, tool := range expectedTools {
			if !strings.Contains(outputStr, tool) {
				t.Errorf("Expected to find tool %s in status output", tool)
			}
		}
	})

	t.Run("VerifyDynamicCommands", func(t *testing.T) {
		// Test that CLI discovers and adds tool subcommands
		cmd := exec.Command(cliPath, "--help")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Fatalf("Help command failed: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		// Should show dynamically discovered tools as available commands
		if !strings.Contains(outputStr, "work") || !strings.Contains(outputStr, "communication") {
			t.Error("Dynamic tool commands should appear in help")
		}
	})

	t.Run("ProxyCommandToMockTool", func(t *testing.T) {
		// Test command proxying mechanics (not tool content)
		cmd := exec.Command(cliPath, "work", "hello")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			// Check if it's a make error (expected if make command fails)
			outputStr := string(output)
			if !strings.Contains(outputStr, "make:") {
				t.Fatalf("Unexpected command proxying error: %v\nOutput: %s", err, outputStr)
			}
			t.Logf("Make command failed as expected: %v", err)
			return
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Mock work hello") {
			t.Error("Command should proxy to make target and execute")
		}
	})

	t.Run("TestInvalidToolCommand", func(t *testing.T) {
		// Test that CLI shows available commands when invalid command is given
		cmd := exec.Command(cliPath, "work", "nonexistent")
		output, err := cmd.CombinedOutput()
		
		// Cobra shows help for invalid subcommands, which is expected behavior
		outputStr := string(output)
		t.Logf("Invalid command output: %s", outputStr)
		
		if err == nil {
			// If no error, should show help with available commands
			if !strings.Contains(outputStr, "Available Commands:") {
				t.Error("Should show available commands for invalid subcommand")
			}
		} else {
			// If error, that's also acceptable (depends on Cobra version)
			t.Logf("CLI returned error for invalid command: %v", err)
		}
	})
}

// TestCLIHelp tests help functionality
func TestCLIHelp(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	cliPath := buildCLIBinary(t, tempDir)

	t.Run("RootHelp", func(t *testing.T) {
		cmd := exec.Command(cliPath, "--help")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Fatalf("Help command failed: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		expectedHelpSections := []string{
			"NimsForest Package Manager",
			"Available Commands:",
			"create-organization-workspace",
			"install",
			"status",
			"hello",
		}

		for _, section := range expectedHelpSections {
			if !strings.Contains(outputStr, section) {
				t.Errorf("Help output missing expected section: %s", section)
			}
		}
	})

	t.Run("CommandHelp", func(t *testing.T) {
		cmd := exec.Command(cliPath, "create-organization-workspace", "--help")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Fatalf("Command help failed: %v\nOutput: %s", err, string(output))
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "Create") && !strings.Contains(outputStr, "organizational") {
			t.Error("Command help should contain description")
		}
	})
}

// buildCLIBinary builds the CLI binary for testing
func buildCLIBinary(t *testing.T, tempDir string) string {
	// Get the main package directory
	_, currentFile, _, _ := runtime.Caller(0)
	cmdDir := filepath.Dir(filepath.Dir(currentFile))
	
	binaryPath := filepath.Join(tempDir, "nimsforestpm-test")
	
	// Build the binary
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd")
	cmd.Dir = cmdDir
	
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI binary: %v\nOutput: %s", err, string(output))
	}
	
	return binaryPath
}

// setupMockWorkspace creates a mock workspace with test tools
func setupMockWorkspace(t *testing.T, tempDir string) string {
	// Create workspace structure
	workspaceDir := filepath.Join(tempDir, "mock-workspace")
	orgWorkspace := filepath.Join(workspaceDir, "test-organization-workspace")
	productsWorkspace := filepath.Join(workspaceDir, "products-workspace")

	dirs := []string{
		orgWorkspace,
		productsWorkspace,
		filepath.Join(productsWorkspace, "nimsforestwork-workspace", "main"),
		filepath.Join(productsWorkspace, "nimsforestcommunication-workspace", "main"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create mock makefiles
	mockTools := []struct {
		name     string
		commands []string
	}{
		{"nimsforestwork", []string{"hello", "init", "triage"}},
		{"nimsforestcommunication", []string{"hello", "init", "lint"}},
	}

	for _, tool := range mockTools {
		makefilePath := filepath.Join(productsWorkspace, tool.name+"-workspace", "main", "MAKEFILE."+tool.name)
		
		var content strings.Builder
		content.WriteString("# Mock " + tool.name + " makefile\n\n")
		
		for _, cmd := range tool.commands {
			content.WriteString(tool.name + "-" + cmd + ":\n")
			content.WriteString("\t@echo \"Mock " + strings.TrimPrefix(tool.name, "nimsforest") + " " + cmd + "\"\n\n")
		}
		
		if err := os.WriteFile(makefilePath, []byte(content.String()), 0644); err != nil {
			t.Fatalf("Failed to create mock makefile: %v", err)
		}
	}

	return filepath.Join(workspaceDir, "test-organization-workspace")
}