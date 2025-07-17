package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nimsforest/nimsforestpackagemanager/internal/registry"
	"github.com/nimsforest/nimsforesttool/tool"
	"github.com/spf13/cobra"
)

// Command initialization
func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(helloCmd)
	rootCmd.AddCommand(validateCmd)

	// Initialize command flags
	helloCmd.Flags().BoolP("dev", "d", false, "Enable developer mode (checks for additional development tools)")
}

// ============================================================================
// COMMAND DEFINITIONS
// ============================================================================

// Workspace functionality has been moved to nimsforestworkspace tool
// Install with: nimsforestpm install workspace
// Use with: nimsforestworkspace create <org-name>

var installCmd = &cobra.Command{
	Use:   "install [tool1] [tool2] ...",
	Short: "Install nimsforest tools via go get",
	Long: fmt.Sprintf(`Install nimsforest tools using go get and go install.

Short names (recommended): %s
Full repository paths also supported.

Examples:
  nimsforestpm install organize
  nimsforestpm install work communicate
  nimsforestpm install all
  nimsforestpm install github.com/nimsforest/nimsforestorganize
  nimsforestpm install github.com/otherperson/customtool`, strings.Join(registry.AvailableTools(), ", ")),
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Handle 'all' argument
		if len(args) == 1 && args[0] == "all" {
			args = registry.AvailableTools()
		}

		// Install each tool
		for _, toolName := range args {
			if err := registry.InstallTool(toolName); err != nil {
				fmt.Fprintf(os.Stderr, "Error installing %s: %v\n", toolName, err)
				os.Exit(1)
			}
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show installed nimsforest tools",
	Run: func(cmd *cobra.Command, args []string) {
		showSimpleStatus()
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [tool1] [tool2] ...",
	Short: "Update installed nimsforest tools",
	Long: `Update tools using go get -u and go install.
If no tools are specified, all installed tools will be updated.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// Update all installed tools
			args = registry.InstalledTools()
			if len(args) == 0 {
				fmt.Println("No tools installed to update.")
				return
			}
		}

		// Update specific tools
		for _, toolName := range args {
			if err := registry.UpdateTool(toolName); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating %s: %v\n", toolName, err)
				os.Exit(1)
			}
		}
	},
}

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "System compatibility check",
	Run: func(cmd *cobra.Command, args []string) {
		devMode, _ := cmd.Flags().GetBool("dev")
		if err := runHello(devMode); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate <tool-name>",
	Short: "Validate a nimsforest tool",
	Long: `Validate that a tool conforms to the nimsforest package manager interface.
This checks if the tool supports the required commands and interface.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		toolName := args[0]
		if err := validateTool(toolName); err != nil {
			fmt.Fprintf(os.Stderr, "Error validating %s: %v\n", toolName, err)
			os.Exit(1)
		}
	},
}

// ============================================================================
// COMMAND IMPLEMENTATIONS
// ============================================================================

// showSimpleStatus displays the current status of installed tools
func showSimpleStatus() {
	fmt.Println("=== NimsForest Tools Status ===")

	available := registry.AvailableTools()
	installed := registry.InstalledTools()

	fmt.Printf("Available tools: %s\n", strings.Join(available, ", "))
	fmt.Printf("Installed tools: %s\n", strings.Join(installed, ", "))

	if len(installed) == 0 {
		fmt.Println("\nNo tools installed. Use 'nimsforestpm install <tool>' to install tools.")
		return
	}

	fmt.Println("\nTool Details:")
	for _, toolName := range available {
		status := "❌ Not installed"
		if registry.IsToolInstalled(toolName) {
			status = "✅ Installed"
		}

		// Get tool info for description
		if info, err := registry.GetToolInfo(toolName); err == nil {
			fmt.Printf("  %s: %s - %s\n", toolName, status, info.Description)
		} else {
			fmt.Printf("  %s: %s\n", toolName, status)
		}
	}
}

// runHello performs basic system compatibility checks
func runHello(devMode bool) error {
	fmt.Println("=== NimsForest Package Manager ===")
	fmt.Println("System Compatibility Check")
	fmt.Println("")

	// Check Go installation
	if _, err := exec.LookPath("go"); err != nil {
		fmt.Printf("❌ Go not found\n")
		fmt.Printf("Please install Go to use nimsforest tools:\n")
		fmt.Printf("  • Download: https://golang.org/dl/\n")
		fmt.Printf("  • Linux: sudo apt install golang-go\n")
		fmt.Printf("  • macOS: brew install go\n")
		fmt.Printf("  • Windows: winget install GoLang.Go\n")
		fmt.Printf("\nAfter installing Go, run 'nimsforestpm hello' again to verify.\n")
		return fmt.Errorf("Go installation required")
	}

	// Get Go version
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get Go version: %w", err)
	}
	fmt.Printf("✓ %s", output)

	// Check Git installation
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Printf("❌ Git not found\n")
		fmt.Printf("Please install Git for workspace management:\n")
		fmt.Printf("  • Download: https://git-scm.com/downloads\n")
		fmt.Printf("  • Linux: sudo apt install git\n")
		fmt.Printf("  • macOS: brew install git\n")
		fmt.Printf("  • Windows: winget install Git.Git\n")
		fmt.Printf("\nAfter installing Git, run 'nimsforestpm hello' again to verify.\n")
		return fmt.Errorf("Git installation required")
	}

	// Get Git version
	cmd = exec.Command("git", "--version")
	output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get Git version: %w", err)
	}
	fmt.Printf("✓ %s", output)

	// Developer mode checks
	if devMode {
		fmt.Println("=== Developer Mode Checks ===")

		// Check for Task (task runner)
		if _, err := exec.LookPath("task"); err != nil {
			fmt.Printf("❌ Task not found\n")
			fmt.Printf("Task is recommended for development:\n")
			fmt.Printf("  • Download: https://taskfile.dev/installation/\n")
			fmt.Printf("  • Linux: sudo snap install task --classic\n")
			fmt.Printf("  • macOS: brew install go-task/tap/go-task\n")
			fmt.Printf("  • Windows: winget install Task.Task\n")
			fmt.Printf("  • Go: go install github.com/go-task/task/v3/cmd/task@latest\n")
		} else {
			// Get Task version
			cmd = exec.Command("task", "--version")
			output, err = cmd.Output()
			if err != nil {
				fmt.Printf("✓ Task installed (version check failed)\n")
			} else {
				fmt.Printf("✓ Task %s", output)
			}
		}

		fmt.Println("")
	}

	fmt.Println("✓ System is ready for NimsForest!")
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Println("  nimsforestpm install workspace")
	fmt.Println("  nimsforestworkspace create <org-name>")
	fmt.Println("  nimsforestpm install <tool-name>")
	fmt.Println("  nimsforestpm status")

	if devMode {
		fmt.Println("")
		fmt.Println("Development commands:")
		fmt.Println("  task build         # Build the project")
		fmt.Println("  task test          # Run tests")
		fmt.Println("  task build-release # Build release binaries")
	}

	return nil
}

// validateTool validates that a tool conforms to the package manager interface
func validateTool(toolName string) error {
	var toolPath string

	// Check if it's a direct path to a tool binary
	if strings.Contains(toolName, "/") {
		// Direct path provided
		toolPath = toolName
	} else {
		// Check if tool is installed in registry
		if !registry.IsToolInstalled(toolName) {
			return fmt.Errorf("tool %s is not installed. Run 'nimsforestpm install %s' first", toolName, toolName)
		}

		// Get tool path from GOPATH
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %v", err)
			}
			gopath = filepath.Join(home, "go")
		}

		toolPath = filepath.Join(gopath, "bin", toolName)
	}

	// Validate tool using the package manager interface
	if err := tool.ValidateTool(toolPath); err != nil {
		return fmt.Errorf("tool validation failed: %v", err)
	}

	// Get tool info
	info, err := tool.QueryTool(toolPath)
	if err != nil {
		return fmt.Errorf("failed to query tool info: %v", err)
	}

	fmt.Printf("✓ Tool %s is valid\n", toolName)
	fmt.Printf("  Name: %s\n", info.Name)
	fmt.Printf("  Version: %s\n", info.Version)
	fmt.Printf("  Description: %s\n", info.Description)
	fmt.Printf("  Commands: %s\n", strings.Join(info.Commands, ", "))

	return nil
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================
