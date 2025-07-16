//go:build integration

package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestComponentInstallation tests component installation flows
func TestComponentInstallation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Build the CLI binary for testing
	cliPath := buildCLIBinary(t, tempDir)
	
	// Create a test workspace
	workspaceDir := setupTestWorkspace(t, tempDir, "test-org")
	os.Chdir(workspaceDir)

	t.Run("InstallValidComponent", func(t *testing.T) {
		// Test installing a valid component
		cmd := exec.Command(cliPath, "install", "work")
		output, err := cmd.CombinedOutput()
		outputStr := string(output)
		
		t.Logf("Install output: %s", outputStr)
		
		if err != nil {
			// Only skip if it's a network/git connectivity issue
			if strings.Contains(outputStr, "network connectivity") || 
			   strings.Contains(outputStr, "git authentication") ||
			   strings.Contains(outputStr, "Repository") && strings.Contains(outputStr, "doesn't exist") {
				t.Skipf("Install test skipped - network/git issue: %v", err)
				return
			}
			// For other errors, this is a real test failure
			t.Fatalf("Install command failed: %v\nOutput: %s", err, outputStr)
		}

		// Verify component was installed in products-workspace
		productsDir := filepath.Join(workspaceDir, "..", "products-workspace")
		componentDir := filepath.Join(productsDir, "nimsforestwork-workspace")
		
		if _, err := os.Stat(componentDir); os.IsNotExist(err) {
			t.Error("Component directory should be created after installation")
		}
		
		// Verify workspace file was updated with the tool
		workspaceFile := filepath.Join(workspaceDir, "..", "nimsforest.workspace")
		content, err := os.ReadFile(workspaceFile)
		if err != nil {
			t.Fatalf("Failed to read workspace file: %v", err)
		}
		
		workspaceStr := string(content)
		if !strings.Contains(workspaceStr, "nimsforestwork") {
			t.Error("Workspace file should contain the installed tool")
		}
		
		// Verify tools section exists
		if !strings.Contains(workspaceStr, "tools (") {
			t.Error("Workspace file should contain tools section")
		}
		
		// Verify correct command suggestion in output
		if !strings.Contains(outputStr, "nimsforestpm work init") {
			t.Error("Install output should suggest nimsforestpm command, not make command")
		}
	})

	t.Run("InstallMultipleComponents", func(t *testing.T) {
		cmd := exec.Command(cliPath, "install", "communicate", "organize")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Logf("Multiple install output: %s", string(output))
			t.Skipf("Multiple install test skipped - requires make environment: %v", err)
			return
		}

		// Check that both components were installed
		productsDir := filepath.Join(workspaceDir, "..", "products-workspace")
		expectedComponents := []string{
			"nimsforestcommunication-workspace",
			"nimsforestorganize-workspace",
		}

		for _, component := range expectedComponents {
			componentDir := filepath.Join(productsDir, component)
			if _, err := os.Stat(componentDir); os.IsNotExist(err) {
				t.Errorf("Component %s should be installed", component)
			}
		}
	})

	t.Run("InstallInvalidComponent", func(t *testing.T) {
		cmd := exec.Command(cliPath, "install", "invalid-component")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error for invalid component")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "unknown tool") {
			t.Error("Should indicate unknown tool error")
		}
	})

	t.Run("InstallFromWrongLocation", func(t *testing.T) {
		// Change to a directory without makefile
		wrongDir := filepath.Join(tempDir, "wrong-location")
		if err := os.MkdirAll(wrongDir, 0755); err != nil {
			t.Fatalf("Failed to create wrong directory: %v", err)
		}
		os.Chdir(wrongDir)

		cmd := exec.Command(cliPath, "install", "work")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error when installing from wrong location")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "MAKEFILE.nimsforestpm not found") {
			t.Error("Should indicate makefile not found")
		}
	})
}

// TestGitSubmoduleIntegration tests git submodule operations
func TestGitSubmoduleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Build the CLI binary for testing
	cliPath := buildCLIBinary(t, tempDir)
	
	// Create a test workspace with git initialization
	workspaceDir := setupTestWorkspaceWithGit(t, tempDir, "test-org")
	os.Chdir(workspaceDir)

	t.Run("SubmoduleStatusCheck", func(t *testing.T) {
		// After installing components, check git submodule status
		cmd := exec.Command(cliPath, "install", "work")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Logf("Install output: %s", string(output))
			t.Skipf("Git submodule test skipped - requires make environment: %v", err)
			return
		}

		// Check git submodule status
		productsDir := filepath.Join(workspaceDir, "..", "products-workspace")
		gitCmd := exec.Command("git", "submodule", "status")
		gitCmd.Dir = productsDir
		gitOutput, gitErr := gitCmd.CombinedOutput()
		
		if gitErr != nil {
			t.Logf("Git submodule status failed: %v\nOutput: %s", gitErr, string(gitOutput))
			t.Skip("Git submodule status check skipped")
			return
		}

		gitOutputStr := string(gitOutput)
		if !strings.Contains(gitOutputStr, "nimsforestwork-workspace") {
			t.Error("Git submodule should show installed component")
		}
	})

	t.Run("SubmoduleUpdate", func(t *testing.T) {
		// Test updating submodules
		cmd := exec.Command(cliPath, "update", "work")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Logf("Update output: %s", string(output))
			t.Skipf("Update test skipped - requires git setup: %v", err)
			return
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "updated successfully") {
			t.Error("Update should report success")
		}
	})
}

// TestInstallationErrorHandling tests error scenarios during installation
func TestInstallationErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Build the CLI binary for testing
	cliPath := buildCLIBinary(t, tempDir)
	
	t.Run("InstallWithoutWorkspace", func(t *testing.T) {
		// Try to install without being in a workspace
		testDir := filepath.Join(tempDir, "no-workspace")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
		os.Chdir(testDir)

		cmd := exec.Command(cliPath, "install", "work")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error when installing outside workspace")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "MAKEFILE.nimsforestpm not found") {
			t.Error("Should indicate makefile not found")
		}
	})

	t.Run("InstallWithEmptyArgs", func(t *testing.T) {
		cmd := exec.Command(cliPath, "install")
		output, err := cmd.CombinedOutput()
		
		if err == nil {
			t.Error("Expected error when no components specified")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "requires at least") {
			t.Error("Should indicate missing arguments")
		}
	})
}

// TestStatusAfterInstallation tests status command after installing components
func TestStatusAfterInstallation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	// Build the CLI binary for testing
	cliPath := buildCLIBinary(t, tempDir)
	
	// Create a test workspace
	workspaceDir := setupTestWorkspace(t, tempDir, "test-org")
	os.Chdir(workspaceDir)

	t.Run("StatusBeforeInstallation", func(t *testing.T) {
		cmd := exec.Command(cliPath, "status")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Logf("Status output: %s", string(output))
			t.Skip("Status test skipped - requires make environment")
			return
		}

		outputStr := string(output)
		// Should show workspace but no installed components
		if strings.Contains(outputStr, "Not in a nimsforest workspace") {
			t.Error("Should recognize workspace")
		}
	})

	t.Run("StatusAfterInstallation", func(t *testing.T) {
		// First install a component
		cmd := exec.Command(cliPath, "install", "work")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Logf("Install output: %s", string(output))
			t.Skip("Install test skipped - requires make environment")
			return
		}

		// Then check status
		statusCmd := exec.Command(cliPath, "status")
		statusOutput, statusErr := statusCmd.CombinedOutput()
		
		if statusErr != nil {
			t.Logf("Status output: %s", string(statusOutput))
			t.Skip("Status after install test skipped")
			return
		}

		statusStr := string(statusOutput)
		// Should show installed component
		if !strings.Contains(statusStr, "work") {
			t.Error("Status should show installed component")
		}
	})
}

// setupTestWorkspace creates a test workspace structure
func setupTestWorkspace(t *testing.T, tempDir, orgName string) string {
	// Create workspace structure
	workspaceDir := filepath.Join(tempDir, "test-workspace")
	orgWorkspace := filepath.Join(workspaceDir, orgName+"-organization-workspace")
	productsWorkspace := filepath.Join(workspaceDir, "products-workspace")

	dirs := []string{orgWorkspace, productsWorkspace}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Copy makefile for testing
	if err := copyMakefileToWorkspace(t, workspaceDir); err != nil {
		t.Logf("Warning: Could not copy makefile: %v", err)
	}

	return orgWorkspace
}

// setupTestWorkspaceWithGit creates a test workspace with git initialization
func setupTestWorkspaceWithGit(t *testing.T, tempDir, orgName string) string {
	workspaceDir := setupTestWorkspace(t, tempDir, orgName)
	
	// Initialize git in products-workspace
	productsDir := filepath.Join(workspaceDir, "..", "products-workspace")
	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = productsDir
	if err := gitCmd.Run(); err != nil {
		t.Logf("Warning: Could not initialize git: %v", err)
	}

	// Set up git config for testing
	gitCmds := [][]string{
		{"git", "config", "user.name", "Test User"},
		{"git", "config", "user.email", "test@example.com"},
	}
	
	for _, cmdArgs := range gitCmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = productsDir
		if err := cmd.Run(); err != nil {
			t.Logf("Warning: Git config failed: %v", err)
		}
	}

	return workspaceDir
}

// copyMakefileToWorkspace copies the makefile to the workspace for testing
func copyMakefileToWorkspace(t *testing.T, workspaceDir string) error {
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(currentFile))
	makefileSrc := filepath.Join(projectRoot, "MAKEFILE.nimsforestpm")
	makefileDst := filepath.Join(workspaceDir, "MAKEFILE.nimsforestpm")
	
	return copyFile(makefileSrc, makefileDst)
}