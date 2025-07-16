package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/nimsforest/nimsforestpackagemanager/internal/workspace"
	"github.com/nimsforest/nimsforestpackagemanager/pkg/tool"
)

func TestExampleToolBinary(t *testing.T) {
	// Build the example tool binary
	binaryPath := filepath.Join(t.TempDir(), "example-tool")
	cmd := exec.Command("go", "build", "-o", binaryPath, "./example-tool")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build example tool: %v", err)
	}

	// Test binary exists and is executable
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatal("Example tool binary was not created")
	}

	// Test basic commands
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
			name:     "hello with name",
			args:     []string{"hello", "Test"},
			expected: "Hello, Test!",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(binaryPath, tc.args...)
			output, err := cmd.Output()
			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}

			if string(output) != tc.expected+"\n" {
				t.Errorf("Expected %q, got %q", tc.expected, string(output))
			}
		})
	}
}

func TestExampleToolInWorkspace(t *testing.T) {
	// Create a temporary workspace
	tempDir := t.TempDir()
	
	// Build the example tool binary
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

	// Create a workspace file
	workspaceFile := filepath.Join(tempDir, "nimsforest.workspace")
	ws := workspace.NewWorkspace()
	ws.Version = "1.0"
	ws.Organization = "./test-org"
	ws.AddTool(workspace.ToolEntry{
		Name:    "example-tool",
		Mode:    "binary",
		Path:    "bin/example-tool",
		Version: "latest",
	})

	if err := ws.Save(workspaceFile); err != nil {
		t.Fatalf("Failed to save workspace: %v", err)
	}

	// Load the workspace
	loadedWs, err := workspace.LoadWorkspace(workspaceFile)
	if err != nil {
		t.Fatalf("Failed to load workspace: %v", err)
	}

	// Verify the tool is in the workspace
	toolEntry, err := loadedWs.GetTool("example-tool")
	if err != nil {
		t.Fatalf("Tool not found in workspace: %v", err)
	}

	if toolEntry.Name != "example-tool" {
		t.Errorf("Expected tool name 'example-tool', got %s", toolEntry.Name)
	}

	if toolEntry.Mode != "binary" {
		t.Errorf("Expected tool mode 'binary', got %s", toolEntry.Mode)
	}

	// Test that the binary path resolves correctly
	expectedPath := filepath.Join(tempDir, "bin", "example-tool")
	actualPath := filepath.Join(filepath.Dir(workspaceFile), toolEntry.Path)
	if actualPath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, actualPath)
	}
}

func TestPkgToolInterface(t *testing.T) {
	// Test that the pkg/tool interface can be used to create a tool
	
	// Create a basic tool using the interface
	base := tool.NewBaseTool("test-tool", "1.0.0", "A test tool")
	
	// Add a simple command
	base.AddCommand(tool.Command{
		Name:        "test",
		Description: "Test command",
		Handler: func(ctx context.Context, args []string) error {
			return nil
		},
	})

	// Test basic properties
	if base.Name() != "test-tool" {
		t.Errorf("Expected name 'test-tool', got %s", base.Name())
	}

	if base.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %s", base.Version())
	}

	// Test command execution
	err := base.Execute(context.Background(), "test", []string{})
	if err != nil {
		t.Errorf("Command execution failed: %v", err)
	}

	// Test installation
	err = base.Install(context.Background(), tool.InstallOptions{
		Mode: tool.InstallModeBinary,
		Path: filepath.Join(t.TempDir(), "test-tool"),
	})
	if err != nil {
		t.Errorf("Installation failed: %v", err)
	}

	// Test status
	if base.Status() != tool.ToolStatusInstalled {
		t.Errorf("Expected status installed, got %s", base.Status())
	}

	// Test validation
	err = base.Validate(context.Background())
	if err != nil {
		t.Errorf("Validation failed: %v", err)
	}
}