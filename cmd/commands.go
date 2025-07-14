package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createOrganizationWorkspaceCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(helloCmd)
}

var createOrganizationWorkspaceCmd = &cobra.Command{
	Use:   "create-organization-workspace [org-name]",
	Short: "Create a new organizational workspace",
	Long: `Create a complete organizational workspace structure with:
- {org-name}-organization-workspace for coordination
- products-workspace for product development`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		orgName := args[0]
		if err := createOrganizationWorkspace(orgName); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating workspace: %v\n", err)
			os.Exit(1)
		}
	},
}

var installCmd = &cobra.Command{
	Use:   "install [tool1] [tool2] ...",
	Short: "Install nimsforest tools as product workspaces",
	Long: `Install tools by adding them as git submodules in products-workspace.
Available tools: work, organize, communicate, productize, folders, webstack`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, tool := range args {
			if err := installTool(tool); err != nil {
				fmt.Fprintf(os.Stderr, "Error installing %s: %v\n", tool, err)
				os.Exit(1)
			}
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show workspace status and installed tools",
	Run: func(cmd *cobra.Command, args []string) {
		if err := showStatus(); err != nil {
			fmt.Fprintf(os.Stderr, "Error showing status: %v\n", err)
			os.Exit(1)
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [tool1] [tool2] ...",
	Short: "Update installed nimsforest tools",
	Long: `Update tools by pulling latest changes from their git submodules.
If no tools are specified, updates all installed tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := updateTools(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating tools: %v\n", err)
			os.Exit(1)
		}
	},
}

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "System compatibility check",
	Run: func(cmd *cobra.Command, args []string) {
		if err := runHello(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// createOrganizationWorkspace creates a new workspace using the make command
func createOrganizationWorkspace(orgName string) error {
	// Find the makefile in the current binary's directory or use current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	// Try to find MAKEFILE.nimsforestpm in various locations
	possiblePaths := []string{
		filepath.Join(currentDir, "MAKEFILE.nimsforestpm"),
		"MAKEFILE.nimsforestpm",
	}

	var makefilePath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			makefilePath = path
			break
		}
	}

	if makefilePath == "" {
		return fmt.Errorf("MAKEFILE.nimsforestpm not found. Run from nimsforestpm directory")
	}

	// Run make command
	makeCmd := exec.Command("make", "-f", makefilePath, "nimsforestpm-create-organisation", fmt.Sprintf("ORG_NAME=%s", orgName))
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr
	makeCmd.Stdin = os.Stdin

	return makeCmd.Run()
}

// installTool installs a tool using the make command
func installTool(tool string) error {
	// Map short names to full component names
	toolMap := map[string]string{
		"work":        "nimsforestwork",
		"organize":    "nimsforestorganize", 
		"communicate": "nimsforestcommunication",
		"productize":  "nimsforestproductize",
		"folders":     "nimsforestfolders",
		"webstack":    "nimsforestwebstack",
	}

	component, ok := toolMap[tool]
	if !ok {
		return fmt.Errorf("unknown tool: %s. Available: %s", tool, strings.Join(getAvailableTools(), ", "))
	}

	// Try to find MAKEFILE.nimsforestpm
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	possiblePaths := []string{
		filepath.Join(currentDir, "MAKEFILE.nimsforestpm"),
		"MAKEFILE.nimsforestpm",
	}

	var makefilePath string
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			makefilePath = path
			break
		}
	}

	if makefilePath == "" {
		return fmt.Errorf("MAKEFILE.nimsforestpm not found. Run from organization workspace")
	}

	// Run make command
	makeCmd := exec.Command("make", "-f", makefilePath, "nimsforestpm-install-component", fmt.Sprintf("COMPONENT=%s", component))
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr
	makeCmd.Stdin = os.Stdin

	return makeCmd.Run()
}

// showStatus shows workspace and tool status
func showStatus() error {
	workspaceRoot, err := detectWorkspace()
	if err != nil {
		fmt.Println("Not in a nimsforest workspace")
		return nil
	}

	fmt.Printf("Workspace: %s\n", workspaceRoot)

	tools, err := discoverInstalledTools(workspaceRoot)
	if err != nil {
		return fmt.Errorf("error discovering tools: %w", err)
	}

	fmt.Printf("Installed tools (%d):\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("  • %s - %s\n", tool.Name, tool.Description)
	}

	return nil
}

// runHello performs system compatibility check
func runHello() error {
	fmt.Println("=== NimsforestPM System Check ===")
	fmt.Println("Checking system compatibility...")

	// Check for required tools
	requiredTools := []string{"git", "make"}
	for _, tool := range requiredTools {
		if _, err := exec.LookPath(tool); err != nil {
			return fmt.Errorf("%s is required but not installed", tool)
		}
		fmt.Printf("✓ %s available\n", tool)
	}

	fmt.Println()
	fmt.Println("NimsforestPM - Package Manager for Organizational Components")
	fmt.Println("Bootstrap and manage organizational workspaces")
	fmt.Println()
	fmt.Println("Next: Run 'nimsforestpm create-organization-workspace my-org' to create workspace")
	
	return nil
}

// updateTools updates installed tools via git submodule update
func updateTools(toolNames []string) error {
	workspaceRoot, err := detectWorkspace()
	if err != nil {
		return fmt.Errorf("not in a nimsforest workspace: %w", err)
	}

	// If no specific tools specified, update all installed tools
	if len(toolNames) == 0 {
		fmt.Println("Updating all installed tools...")
		return updateAllTools(workspaceRoot)
	}

	// Update specific tools
	fmt.Printf("Updating tools: %s\n", strings.Join(toolNames, ", "))
	return updateSpecificTools(workspaceRoot, toolNames)
}

// updateAllTools updates all installed tools
func updateAllTools(workspaceRoot string) error {
	// Change to products-workspace directory
	productsDir := filepath.Join(workspaceRoot, "products-workspace")
	if _, err := os.Stat(productsDir); os.IsNotExist(err) {
		return fmt.Errorf("products-workspace not found")
	}

	// Run git submodule update --remote for all submodules
	cmd := exec.Command("git", "submodule", "update", "--remote")
	cmd.Dir = productsDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update submodules: %w", err)
	}

	fmt.Println("✓ All tools updated successfully!")
	return nil
}

// updateSpecificTools updates only the specified tools
func updateSpecificTools(workspaceRoot string, toolNames []string) error {
	toolMap := map[string]string{
		"work":        "nimsforestwork",
		"organize":    "nimsforestorganize",
		"communicate": "nimsforestcommunication",
		"productize":  "nimsforestproductize",
		"folders":     "nimsforestfolders",
		"webstack":    "nimsforestwebstack",
	}

	productsDir := filepath.Join(workspaceRoot, "products-workspace")
	if _, err := os.Stat(productsDir); os.IsNotExist(err) {
		return fmt.Errorf("products-workspace not found")
	}

	for _, toolName := range toolNames {
		fullName, ok := toolMap[toolName]
		if !ok {
			return fmt.Errorf("unknown tool: %s. Available: %s", toolName, strings.Join(getAvailableTools(), ", "))
		}

		toolWorkspaceDir := filepath.Join(productsDir, fullName+"-workspace")
		if _, err := os.Stat(toolWorkspaceDir); os.IsNotExist(err) {
			fmt.Printf("⚠️  Tool %s is not installed, skipping...\n", toolName)
			continue
		}

		fmt.Printf("Updating %s...\n", toolName)
		
		// Run git submodule update --remote for specific submodule
		cmd := exec.Command("git", "submodule", "update", "--remote", fullName+"-workspace")
		cmd.Dir = productsDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to update %s: %w", toolName, err)
		}

		fmt.Printf("✓ %s updated successfully!\n", toolName)
	}

	return nil
}

// getAvailableTools returns list of available tool names
func getAvailableTools() []string {
	return []string{"work", "organize", "communicate", "productize", "folders", "webstack"}
}