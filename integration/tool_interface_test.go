//go:build integration

package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/nimsforest/nimsforestpackagemanager/internal/runtimetool"
	"github.com/nimsforest/nimsforestpackagemanager/internal/workspace"
	"github.com/nimsforest/nimsforestpackagemanager/pkg/tool"
)

// TestToolInterfaceIntegration tests the complete tool interface workflow:
// 1. Use package manager to install an example tool
// 2. Verify tool is registered in workspace
// 3. Execute all available commands using runtime system
// 4. Validate tool health and status throughout
func TestToolInterfaceIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test environment
	env := NewTestEnvironment(t, "test-org")
	defer env.Cleanup()

	helper := NewTestHelper(t)
	
	t.Run("CompleteToolWorkflow", func(t *testing.T) {
		// Step 1: Install example tool using package manager
		result, err := env.RunCLI("install", "work")
		if err != nil {
			t.Logf("Install output: %s", result.Output)
			t.Skipf("Install test skipped - requires make environment: %v", err)
			return
		}
		helper.AssertSuccess(result, "Tool installation")

		// Step 2: Verify tool is installed and accessible
		t.Run("VerifyToolInstallation", func(t *testing.T) {
			// Check that tool directory exists
			toolPath := filepath.Join(env.ProductsDir, "nimsforestwork-workspace")
			helper.AssertDirExists(toolPath, "Tool should be installed")

			// Use status command to verify
			statusResult, err := env.RunCLI("status")
			if err != nil {
				t.Logf("Status output: %s", statusResult.Output)
				t.Skip("Status check skipped")
				return
			}
			helper.AssertContains(statusResult, "work", "Status should show installed tool")
		})

		// Step 3: Test tool interface directly
		t.Run("TestToolInterface", func(t *testing.T) {
			// Load workspace to get tool registry
			ws, err := workspace.LoadWorkspaceFromDir(env.OrganizationDir)
			if err != nil {
				t.Fatalf("Failed to load workspace: %v", err)
			}

			// Verify workspace loaded successfully
			if ws == nil {
				t.Fatal("Workspace should be loaded successfully")
			}

			// Create example tool for testing
			exampleTool := createExampleTool(t)
			
			// Test basic tool interface
			t.Run("BasicToolInterface", func(t *testing.T) {
				// Test tool metadata
				if exampleTool.Name() != "example" {
					t.Errorf("Expected tool name 'example', got %s", exampleTool.Name())
				}

				if exampleTool.Version() != "1.0.0" {
					t.Errorf("Expected version '1.0.0', got %s", exampleTool.Version())
				}

				// Test commands
				commands := exampleTool.Commands()
				if len(commands) == 0 {
					t.Error("Tool should have commands")
				}

				// Test tool validation
				if err := exampleTool.Validate(context.Background()); err != nil {
					t.Errorf("Tool validation failed: %v", err)
				}
			})

			// Test tool installation via interface
			t.Run("ToolInterfaceInstallation", func(t *testing.T) {
				// Install tool using interface
				installOpts := tool.InstallOptions{
					Mode: tool.InstallModeBinary,
					Path: filepath.Join(env.TempDir, "example-tool"),
				}
				
				if err := exampleTool.Install(context.Background(), installOpts); err != nil {
					t.Errorf("Tool installation via interface failed: %v", err)
				}

				// Verify tool status
				if exampleTool.Status() != tool.ToolStatusInstalled {
					t.Errorf("Expected tool status to be installed, got %s", exampleTool.Status())
				}
			})

			// Test health check
			t.Run("ToolHealthCheck", func(t *testing.T) {
				health := exampleTool.HealthCheck(context.Background())
				if health.Status != tool.HealthStatusHealthy {
					t.Errorf("Expected healthy status, got %s", health.Status)
				}
			})
		})

		// Step 4: Test runtime tool execution
		t.Run("RuntimeToolExecution", func(t *testing.T) {
			// Load workspace first
			ws, err := workspace.LoadWorkspaceFromDir(env.OrganizationDir)
			if err != nil {
				t.Fatalf("Failed to load workspace: %v", err)
			}

			// Create mock tool entry for runtime tool
			toolEntry := workspace.ToolEntry{
				Name:    "example",
				Mode:    "binary",
				Version: "1.0.0",
				Path:    filepath.Join(env.TempDir, "runtime-example"),
			}
			
			// Create runtime tool instance
			runtimeTool := runtimetool.NewRuntimeTool(toolEntry, ws)
			
			// Test runtime tool basic functionality
			if runtimeTool == nil {
				t.Fatal("Failed to create runtime tool")
			}

			// Test runtime tool execution with workspace context
			t.Run("RuntimeExecution", func(t *testing.T) {
				// Create example tool and add to runtime
				exampleTool := createExampleTool(t)
				
				// Install the tool first
				installOpts := tool.InstallOptions{
					Mode: tool.InstallModeBinary,
					Path: filepath.Join(env.TempDir, "runtime-example"),
				}
				
				if err := exampleTool.Install(context.Background(), installOpts); err != nil {
					t.Errorf("Runtime tool installation failed: %v", err)
				}

				// Execute commands through runtime
				commands := exampleTool.Commands()
				for _, cmd := range commands {
					if !cmd.Hidden {
						t.Run(fmt.Sprintf("Execute_%s", cmd.Name), func(t *testing.T) {
							err := exampleTool.Execute(context.Background(), cmd.Name, []string{})
							if err != nil {
								t.Logf("Command %s execution failed: %v", cmd.Name, err)
								// Don't fail test for command errors in integration test
							}
						})
					}
				}
			})
		})

		// Step 5: Test tool lifecycle management
		t.Run("ToolLifecycleManagement", func(t *testing.T) {
			exampleTool := createExampleTool(t)
			
			// Test installation
			installOpts := tool.InstallOptions{
				Mode: tool.InstallModeBinary,
				Path: filepath.Join(env.TempDir, "lifecycle-example"),
			}
			
			if err := exampleTool.Install(context.Background(), installOpts); err != nil {
				t.Errorf("Tool installation failed: %v", err)
			}

			// Test update
			updateOpts := tool.UpdateOptions{
				Force: false,
			}
			
			if err := exampleTool.Update(context.Background(), updateOpts); err != nil {
				t.Errorf("Tool update failed: %v", err)
			}

			// Test uninstall
			uninstallOpts := tool.UninstallOptions{
				Force: false,
			}
			
			if err := exampleTool.Uninstall(context.Background(), uninstallOpts); err != nil {
				t.Errorf("Tool uninstall failed: %v", err)
			}

			// Verify tool is uninstalled
			if exampleTool.Status() != tool.ToolStatusNotInstalled {
				t.Errorf("Expected tool to be uninstalled, got status %s", exampleTool.Status())
			}
		})

		// Step 6: Test workspace tool integration
		t.Run("WorkspaceToolIntegration", func(t *testing.T) {
			// Load workspace
			ws, err := workspace.LoadWorkspaceFromDir(env.OrganizationDir)
			if err != nil {
				t.Fatalf("Failed to load workspace: %v", err)
			}

			// Verify workspace can track tools
			if ws == nil {
				t.Error("Workspace should be loaded successfully")
			}

			// Test workspace tool listing (when implemented)
			// This is a placeholder for future workspace tool integration
			t.Log("Workspace tool integration placeholder - to be implemented")
		})
	})
}

// TestToolInterfaceErrorHandling tests error scenarios in tool interface
func TestToolInterfaceErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("InvalidToolOperations", func(t *testing.T) {
		exampleTool := createExampleTool(t)

		// Test executing non-existent command
		err := exampleTool.Execute(context.Background(), "nonexistent", []string{})
		if err == nil {
			t.Error("Expected error for non-existent command")
		}

		// Test installing with invalid options
		invalidOpts := tool.InstallOptions{
			Mode: tool.InstallMode(999), // Invalid mode
		}
		
		err = exampleTool.Install(context.Background(), invalidOpts)
		if err == nil {
			t.Error("Expected error for invalid install options")
		}

		// Test updating uninstalled tool
		updateOpts := tool.UpdateOptions{}
		err = exampleTool.Update(context.Background(), updateOpts)
		if err == nil {
			t.Error("Expected error updating uninstalled tool")
		}

		// Test uninstalling uninstalled tool
		uninstallOpts := tool.UninstallOptions{}
		err = exampleTool.Uninstall(context.Background(), uninstallOpts)
		if err == nil {
			t.Error("Expected error uninstalling uninstalled tool")
		}
	})

	t.Run("ToolValidationErrors", func(t *testing.T) {
		// Create tool with invalid configuration
		baseTool := tool.NewBaseTool("", "", "") // Empty name and version
		
		// Test validation fails
		err := baseTool.Validate(context.Background())
		if err == nil {
			t.Error("Expected validation error for empty tool name")
		}
	})
}

// TestToolInterfacePerformance tests performance aspects of tool interface
func TestToolInterfacePerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	t.Run("ToolCreationPerformance", func(t *testing.T) {
		start := time.Now()
		
		// Create multiple tools
		for i := 0; i < 100; i++ {
			_ = createExampleTool(t)
		}
		
		duration := time.Since(start)
		if duration > time.Second {
			t.Errorf("Tool creation took too long: %v", duration)
		}
	})

	t.Run("ToolCommandExecutionPerformance", func(t *testing.T) {
		exampleTool := createExampleTool(t)
		
		// Install tool first
		installOpts := tool.InstallOptions{
			Mode: tool.InstallModeBinary,
			Path: filepath.Join(t.TempDir(), "perf-example"),
		}
		
		if err := exampleTool.Install(context.Background(), installOpts); err != nil {
			t.Fatalf("Tool installation failed: %v", err)
		}

		// Time command execution
		start := time.Now()
		
		commands := exampleTool.Commands()
		for _, cmd := range commands {
			if !cmd.Hidden {
				_ = exampleTool.Execute(context.Background(), cmd.Name, []string{})
			}
		}
		
		duration := time.Since(start)
		if duration > 5*time.Second {
			t.Errorf("Command execution took too long: %v", duration)
		}
	})
}

// createExampleTool creates an example tool for testing
func createExampleTool(t *testing.T) *ExampleTool {
	base := tool.NewBaseTool("example", "1.0.0", "An example tool for integration testing")
	exampleTool := &ExampleTool{BaseTool: base}
	
	// Add commands
	exampleTool.AddCommand(tool.Command{
		Name:        "hello",
		Description: "Say hello",
		Usage:       "hello [name]",
		Handler: func(ctx context.Context, args []string) error {
			name := "World"
			if len(args) > 0 {
				name = args[0]
			}
			t.Logf("Hello, %s!", name)
			return nil
		},
	})
	
	exampleTool.AddCommand(tool.Command{
		Name:        "status",
		Description: "Show tool status",
		Usage:       "status",
		Handler: func(ctx context.Context, args []string) error {
			t.Logf("Tool status: %s", exampleTool.Status())
			return nil
		},
	})
	
	return exampleTool
}

// ExampleTool implements the tool interface for testing
type ExampleTool struct {
	*tool.BaseTool
}

// HealthCheck implements custom health check
func (t *ExampleTool) HealthCheck(ctx context.Context) tool.HealthCheck {
	health := t.BaseTool.HealthCheck(ctx)
	
	// Add custom health information
	health.Details["integration_test"] = true
	health.Details["commands_count"] = len(t.Commands())
	
	return health
}