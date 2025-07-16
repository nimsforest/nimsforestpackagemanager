package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/nimsforest/nimsforestpackagemanager/internal/workspace"
	"github.com/nimsforest/nimsforestpackagemanager/pkg/tool"
)

func init() {
	rootCmd.AddCommand(createOrganizationWorkspaceCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(helloCmd)
	rootCmd.AddCommand(workspaceCmd)
	rootCmd.AddCommand(validateCmd)
	
	// Add flags for binary installation
	installCmd.Flags().StringP("name", "n", "", "Name of the tool to install")
	installCmd.Flags().StringP("path", "p", "", "Path to the binary to install")
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
Available tools: work, organize, communicate, productize, folders, webstack

Or install a custom binary tool:
  nimsforest install --name example-tool --path ./bin/example-tool`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if custom binary installation is requested
		toolName, _ := cmd.Flags().GetString("name")
		binaryPath, _ := cmd.Flags().GetString("path")
		
		if toolName != "" && binaryPath != "" {
			// Install custom binary
			if err := installBinaryTool(toolName, binaryPath); err != nil {
				fmt.Fprintf(os.Stderr, "Error installing binary tool %s: %v\n", toolName, err)
				os.Exit(1)
			}
			return
		}
		
		// Require at least one tool name for standard installation
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Error: must specify either tool names or --name and --path for binary installation\n")
			os.Exit(1)
		}
		
		// Standard tool installation
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

// createOrganizationWorkspace creates a new workspace using pure Go implementation
func createOrganizationWorkspace(orgName string) error {
	// Validate organization name
	if strings.TrimSpace(orgName) == "" {
		return fmt.Errorf("organization name cannot be empty")
	}

	orgName = strings.TrimSpace(orgName)

	// Validate organization name format
	if err := validateOrganizationName(orgName); err != nil {
		return fmt.Errorf("invalid organization name: %w", err)
	}

	fmt.Printf("=== Creating %s Organization Workspace ===\n", orgName)

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	// Create workspace directory structure
	workspaceDir := filepath.Join(currentDir, orgName+"-workspace")
	
	// Check if workspace already exists
	if _, err := os.Stat(workspaceDir); err == nil {
		return fmt.Errorf("workspace directory already exists: %s", workspaceDir)
	}

	fmt.Println("Creating organizational workspace structure...")
	
	// Create directory structure
	if err := createWorkspaceDirectories(workspaceDir, orgName); err != nil {
		return fmt.Errorf("failed to create workspace directories: %w", err)
	}

	// Generate nimsforest.workspace file
	if err := generateWorkspaceFile(workspaceDir, orgName); err != nil {
		return fmt.Errorf("failed to generate workspace file: %w", err)
	}

	// Initialize git repository
	if err := initializeGitRepository(workspaceDir, orgName); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	// Create template files
	if err := createTemplateFiles(workspaceDir, orgName); err != nil {
		return fmt.Errorf("failed to create template files: %w", err)
	}

	// Create initial commit
	if err := createInitialCommit(workspaceDir, orgName); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	fmt.Printf("✓ %s organization workspace created successfully!\n", orgName)
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s/%s-organization-workspace/main\n", orgName+"-workspace", orgName)
	fmt.Println("  nimsforestpm install work communicate organize webstack")
	fmt.Println("")
	fmt.Println("Or install all components:")
	fmt.Println("  nimsforestpm install work communicate organize webstack folders")

	return nil
}

// validateOrganizationName validates the organization name format
func validateOrganizationName(orgName string) error {
	if orgName == "" {
		return fmt.Errorf("organization name cannot be empty")
	}
	
	// Check for valid characters (alphanumeric, hyphens, underscores)
	for _, char := range orgName {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			return fmt.Errorf("organization name contains invalid character '%c'. Only letters, numbers, hyphens, and underscores allowed", char)
		}
	}
	
	// Check that it doesn't start or end with hyphen/underscore
	if strings.HasPrefix(orgName, "-") || strings.HasSuffix(orgName, "-") || strings.HasPrefix(orgName, "_") || strings.HasSuffix(orgName, "_") {
		return fmt.Errorf("organization name cannot start or end with hyphen or underscore")
	}
	
	return nil
}

// createWorkspaceDirectories creates the workspace directory structure
func createWorkspaceDirectories(workspaceDir, orgName string) error {
	// Create root workspace directory
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		return fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// Create organization workspace directory
	orgWorkspaceDir := filepath.Join(workspaceDir, orgName+"-organization-workspace")
	if err := os.MkdirAll(orgWorkspaceDir, 0755); err != nil {
		return fmt.Errorf("failed to create organization workspace directory: %w", err)
	}

	// Create main directory within organization workspace
	mainDir := filepath.Join(orgWorkspaceDir, "main")
	if err := os.MkdirAll(mainDir, 0755); err != nil {
		return fmt.Errorf("failed to create main directory: %w", err)
	}

	// Create products workspace directory
	productsDir := filepath.Join(workspaceDir, "products-workspace")
	if err := os.MkdirAll(productsDir, 0755); err != nil {
		return fmt.Errorf("failed to create products workspace directory: %w", err)
	}

	// Create basic organizational structure directories
	orgStructureDirs := []string{
		"actors/nims",
		"actors/humans",
		"actors/machines/mobile",
		"actors/machines/fixed",
		"assets/documentation",
		"assets/data",
		"assets/media",
		"assets/templates",
		"tools/shared",
		"tools/org-specific",
		"products",
	}

	for _, dir := range orgStructureDirs {
		fullPath := filepath.Join(mainDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateWorkspaceFile creates the nimsforest.workspace file at the workspace root
func generateWorkspaceFile(workspaceDir, orgName string) error {
	// Create workspace instance
	ws := workspace.NewWorkspace()
	ws.Organization = fmt.Sprintf("./%s-organization-workspace", orgName)

	// Save workspace file at workspace root
	workspaceFilePath := filepath.Join(workspaceDir, workspace.WorkspaceFileName)
	if err := ws.Save(workspaceFilePath); err != nil {
		return fmt.Errorf("failed to save workspace file: %w", err)
	}

	fmt.Printf("✓ Created workspace file: %s\n", workspaceFilePath)
	return nil
}

// initializeGitRepository initializes git repository in the organization workspace
func initializeGitRepository(workspaceDir, orgName string) error {
	fmt.Println("Initializing git repositories...")
	
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git command not found. Please install git")
	}

	// Initialize git repository in organization workspace
	orgWorkspaceDir := filepath.Join(workspaceDir, orgName+"-organization-workspace")
	cmd := exec.Command("git", "init")
	cmd.Dir = orgWorkspaceDir
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %w", err)
	}

	fmt.Printf("✓ Initialized git repository in %s\n", orgWorkspaceDir)
	return nil
}

// createTemplateFiles creates template files like README.md and Makefile
func createTemplateFiles(workspaceDir, orgName string) error {
	fmt.Println("Creating template files...")
	
	orgWorkspaceDir := filepath.Join(workspaceDir, orgName+"-organization-workspace")
	mainDir := filepath.Join(orgWorkspaceDir, "main")

	// Create README.md
	if err := createReadmeFile(mainDir, orgName); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	// Create Makefile
	if err := createMakefileTemplate(mainDir, orgName); err != nil {
		return fmt.Errorf("failed to create Makefile: %w", err)
	}

	return nil
}

// createReadmeFile creates the organization README.md from template
func createReadmeFile(mainDir, orgName string) error {
	readmeTemplate := `# {{ORG_NAME}} Organization

Organizational workspace powered by NimsForest components.

## Why This Structure?

This organizational structure is inspired by **game engine architecture** - treating your organization as a dynamic scene where different entities interact to create value.

We also drew inspiration from **Pixar's USD (Universal Scene Description)** for hierarchical organization and **git worktree** patterns for flexible development.

## Organizational Structure

### actors/
All entities that can perform actions in your organizational scene:

- **nims/**: Intelligent advisory actors that learn optimal patterns through reinforcement learning. These are the smart shadows that provide transparent, objective guidance to help work flow better.

- **humans/**: People in the organization - employees, stakeholders, customers. The decision makers and creative force.

- **machines/**: Physical systems that perform work:
  - **mobile/**: Drones, robots, vehicles - systems that can move around
  - **fixed/**: Servers, ASML machines, production equipment - stationary systems

### assets/
Resources and files that actors use:

- **documentation/**: Knowledge and processes
- **data/**: Information and datasets
- **media/**: Images, videos, presentations
- **templates/**: Reusable patterns and structures

### tools/
Capabilities and utilities that enable work:

- **shared/**: Tools shared across the organization (via tools-repository)
- **org-specific/**: Tools specific to this organization

### products/
What the organization builds and delivers:

- Each product has its own repository with the same actor/asset/tool structure
- Products are linked as git submodules for version control
- Can be software, hardware, or services

## Workspace Architecture

The workspace follows a **three-repository pattern**:

1. **Organization Repository** (` + "`{{ORG_NAME}}-repository`" + `): Core organizational structure
2. **Tools Repository** (` + "`tools-repository`" + `): Shared utilities and NimsForest components
3. **Product Repositories** (` + "`product-repositories/`" + `): Individual product development

This separation allows for:
- **Independent versioning**: Each component can evolve at its own pace
- **Flexible permissions**: Different access levels for different repositories
- **Git worktree support**: The ` + "`/main/`" + ` structure supports branching strategies

## Getting Started

### Install NimsForest Components

` + "```bash" + `
# Install core organizational intelligence cycle
nimsforestpm install work communicate organize webstack

# Add advanced folder management
nimsforestpm install folders

# Or install everything at once
nimsforestpm install work communicate organize webstack folders
` + "```" + `

### Add Your First Product

` + "```bash" + `
# Create a software product
nimsforestpm add-product my-app software

# Or hardware product
nimsforestpm add-product my-device hardware

# Or service product
nimsforestpm add-product my-service service
` + "```" + `

### Validate Your Setup

` + "```bash" + `
# Check organizational structure and installed components
nimsforestpm status
` + "```" + `

## Design Philosophy

This structure embodies several key principles:

1. **Game Engine Thinking**: Organizations are dynamic scenes with interacting entities
2. **Learning Systems**: Nims provide objective, transparent optimization
3. **Hierarchical Tools**: Each level (workspace, org, product) has its own tools
4. **Clean Separation**: Actors do things, assets are resources, tools enable work
5. **Git Worktree Ready**: Structure supports advanced branching workflows

## Next Steps for {{ORG_NAME}}

1. **Install components**: ` + "`nimsforestpm install work communicate organize webstack folders`" + `
2. **Add your first product**: ` + "`nimsforestpm add-product my-app software`" + `
3. **Validate setup**: ` + "`nimsforestpm status`" + `
4. **Initialize components**: Run the respective init commands for installed components

---

*Created with [NimsForest Package Manager](https://github.com/nimsforest/nimsforestpackagemanager)*
`

	// Replace template variables
	readmeContent := strings.ReplaceAll(readmeTemplate, "{{ORG_NAME}}", orgName)
	
	// Write README.md
	readmePath := filepath.Join(mainDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	fmt.Printf("✓ Created README.md from template\n")
	return nil
}

// createMakefileTemplate creates the organization Makefile
func createMakefileTemplate(mainDir, orgName string) error {
	makefileContent := fmt.Sprintf(`# %s Organization Makefile
# Include NimsForest Package Manager

# nimsforestpm available as global binary
# Products will be added here when installed

.PHONY: help
help:
	@echo "=== %s Organization ==="
	@echo "Use nimsforestpm commands to manage this organization:"
	@echo "  nimsforestpm install <component>  - Install NimsForest components"
	@echo "  nimsforestpm status              - Show organization status"
	@echo "  nimsforestpm update              - Update installed components"
	@echo ""
	@echo "Available components: work, communicate, organize, webstack, folders"
`, orgName, orgName)

	// Write Makefile
	makefilePath := filepath.Join(mainDir, "Makefile")
	if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		return fmt.Errorf("failed to write Makefile: %w", err)
	}

	fmt.Printf("✓ Created organization Makefile\n")
	return nil
}

// createInitialCommit creates the initial git commit
func createInitialCommit(workspaceDir, orgName string) error {
	orgWorkspaceDir := filepath.Join(workspaceDir, orgName+"-organization-workspace")
	
	// Add all files to git
	addCmd := exec.Command("git", "add", ".")
	addCmd.Dir = orgWorkspaceDir
	if err := addCmd.Run(); err != nil {
		return fmt.Errorf("failed to add files to git: %w", err)
	}

	// Create initial commit
	commitMsg := fmt.Sprintf("Initial %s organization setup", orgName)
	commitCmd := exec.Command("git", "commit", "-m", commitMsg)
	commitCmd.Dir = orgWorkspaceDir
	if err := commitCmd.Run(); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	fmt.Printf("✓ Created initial commit\n")
	return nil
}

// installTool installs a tool using pure Go implementation without makefile delegation
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
	
	// Validate component name
	if err := validateComponentName(component); err != nil {
		return fmt.Errorf("invalid component name: %w", err)
	}

	fmt.Printf("=== Installing %s ===\n", component)
	fmt.Printf("Adding %s as product workspace...\n", component)

	// Find workspace root
	workspaceRoot, err := findWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("not in a nimsforest workspace: %w", err)
	}

	fmt.Printf("📁 Found workspace at %s\n", workspaceRoot)

	// Ensure products-workspace directory exists
	productsDir := filepath.Join(workspaceRoot, "products-workspace")
	if err := os.MkdirAll(productsDir, 0755); err != nil {
		return fmt.Errorf("failed to create products-workspace directory: %w", err)
	}

	// Check if component is already installed
	componentDir := filepath.Join(productsDir, component+"-workspace")
	if _, err := os.Stat(componentDir); err == nil {
		fmt.Printf("✅ %s already installed\n", component)
		return nil
	}

	// Install component as binary (default mode)
	if err := installComponentAsBinary(component); err != nil {
		return err
	}

	// Update workspace file if present
	if err := updateWorkspaceAfterInstall(component); err != nil {
		// Don't fail the install if we can't update the workspace file
		fmt.Printf("Warning: Failed to update workspace file: %v\n", err)
	}

	// Register tool in the tool registry for dynamic command access
	if err := registerInstalledTool(workspaceRoot, component); err != nil {
		fmt.Printf("Warning: Failed to register tool in registry: %v\n", err)
	}

	fmt.Printf("✓ %s installed successfully!\n", component)
	fmt.Printf("Initialize with: nimsforestpm %s init\n", tool)

	return nil
}

// installBinaryTool installs a custom binary tool to the workspace
func installBinaryTool(toolName, binaryPath string) error {
	fmt.Printf("=== Installing Binary Tool: %s ===\n", toolName)
	
	// Validate binary path exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found: %s", binaryPath)
	}
	
	// Find workspace root
	workspaceRoot, err := findWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("not in a nimsforest workspace: %w", err)
	}
	
	fmt.Printf("📁 Found workspace at %s\n", workspaceRoot)
	
	// Create bin directory if it doesn't exist
	binDir := filepath.Join(workspaceRoot, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}
	
	// Copy binary to workspace bin directory
	destPath := filepath.Join(binDir, toolName)
	if err := copyFile(binaryPath, destPath); err != nil {
		return fmt.Errorf("failed to copy binary: %w", err)
	}
	
	// Make it executable
	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}
	
	fmt.Printf("📦 Binary installed at: %s\n", destPath)
	
	// Update workspace file
	if err := updateWorkspaceWithBinaryTool(workspaceRoot, toolName, destPath); err != nil {
		fmt.Printf("Warning: Failed to update workspace file: %v\n", err)
	}
	
	fmt.Printf("✓ %s installed successfully!\n", toolName)
	fmt.Printf("You can now use: %s <command>\n", destPath)
	
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	
	_, err = dstFile.ReadFrom(srcFile)
	return err
}

// updateWorkspaceWithBinaryTool updates the workspace file with the binary tool
func updateWorkspaceWithBinaryTool(workspaceRoot, toolName, binaryPath string) error {
	// Try to find workspace file
	workspaceFile, err := workspace.FindWorkspaceFile(workspaceRoot)
	if err != nil {
		// No workspace file found, create one
		if err := createWorkspaceFileWithBinaryTool(workspaceRoot, toolName, binaryPath); err != nil {
			return fmt.Errorf("failed to create workspace file: %w", err)
		}
		return nil
	}
	
	// Load workspace
	ws, err := workspace.LoadWorkspace(workspaceFile)
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}
	
	// Make path relative to workspace root
	relPath, err := filepath.Rel(workspaceRoot, binaryPath)
	if err != nil {
		relPath = binaryPath
	}
	
	// Create tool entry
	toolEntry := workspace.ToolEntry{
		Name:    toolName,
		Mode:    "binary",
		Path:    relPath,
		Version: "latest",
	}
	
	// Add tool to workspace
	ws.AddTool(toolEntry)
	
	// Save workspace
	if err := ws.Save(workspaceFile); err != nil {
		return fmt.Errorf("failed to save workspace: %w", err)
	}
	
	fmt.Printf("✓ Updated workspace file: added %s\n", toolName)
	return nil
}

// createWorkspaceFileWithBinaryTool creates a new workspace file with the binary tool
func createWorkspaceFileWithBinaryTool(workspaceRoot, toolName, binaryPath string) error {
	workspaceFilePath := filepath.Join(workspaceRoot, workspace.WorkspaceFileName)
	
	// Create new workspace
	ws := workspace.NewWorkspace()
	
	// Try to determine organization name from workspace structure
	if orgName := getOrganizationName(workspaceRoot); orgName != "" {
		ws.Organization = orgName
	}
	
	// Make path relative to workspace root
	relPath, err := filepath.Rel(workspaceRoot, binaryPath)
	if err != nil {
		relPath = binaryPath
	}
	
	// Create tool entry
	toolEntry := workspace.ToolEntry{
		Name:    toolName,
		Mode:    "binary",
		Path:    relPath,
		Version: "latest",
	}
	
	// Add tool to workspace
	ws.AddTool(toolEntry)
	
	// Save the workspace file
	if err := ws.Save(workspaceFilePath); err != nil {
		return fmt.Errorf("failed to save new workspace file: %w", err)
	}
	
	fmt.Printf("✓ Created workspace file: %s\n", workspaceFilePath)
	fmt.Printf("✓ Added %s to workspace\n", toolName)
	return nil
}


// showStatus shows workspace and tool status
func showStatus() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	// First try to use workspace file if available
	if _, err := workspace.FindWorkspaceFile(currentDir); err == nil {
		fmt.Println("Using nimsforest.workspace file for status")
		return showWorkspaceStatus()
	}

	// Fall back to makefile-based approach
	fmt.Println("Using makefile-based status (nimsforest.workspace not found)")
	
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
		fmt.Println("Not in a nimsforest workspace (neither nimsforest.workspace nor MAKEFILE.nimsforestpm found)")
		return nil
	}

	// Delegate to Makefile lint command which shows comprehensive status
	makeCmd := exec.Command("make", "-f", makefilePath, "nimsforestpm-lint")
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr

	return makeCmd.Run()
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

// Workspace management commands
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage workspace configuration",
	Long: `Manage workspace configuration using nimsforest.workspace files.
This provides a pure Go alternative to makefile-based workspace management.`,
}

var workspaceInitCmd = &cobra.Command{
	Use:   "init [org-name]",
	Short: "Initialize a new workspace with nimsforest.workspace file",
	Long: `Create a new nimsforest.workspace file in the current directory.
This file tracks workspace configuration and installed products.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		orgName := args[0]
		if err := initWorkspace(orgName); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing workspace: %v\n", err)
			os.Exit(1)
		}
	},
}

var workspaceStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show workspace status using workspace file",
	Long: `Display current workspace status including organization, products,
and validation results from the nimsforest.workspace file.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := showWorkspaceStatus(); err != nil {
			fmt.Fprintf(os.Stderr, "Error showing workspace status: %v\n", err)
			os.Exit(1)
		}
	},
}

var workspaceAddCmd = &cobra.Command{
	Use:   "add [product-path]",
	Short: "Add a product to the workspace file",
	Long: `Add a product directory to the nimsforest.workspace file.
The path can be relative to the workspace file or absolute.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		productPath := args[0]
		if err := addProductToWorkspace(productPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error adding product to workspace: %v\n", err)
			os.Exit(1)
		}
	},
}

var workspaceRemoveCmd = &cobra.Command{
	Use:   "remove [product-path]",
	Short: "Remove a product from the workspace file",
	Long: `Remove a product directory from the nimsforest.workspace file.
The path should match exactly what is stored in the file.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		productPath := args[0]
		if err := removeProductFromWorkspace(productPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error removing product from workspace: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Add subcommands to workspace command
	workspaceCmd.AddCommand(workspaceInitCmd)
	workspaceCmd.AddCommand(workspaceStatusCmd)
	workspaceCmd.AddCommand(workspaceAddCmd)
	workspaceCmd.AddCommand(workspaceRemoveCmd)
}

// initWorkspace creates a new workspace file with the given organization name
func initWorkspace(orgName string) error {
	// Validate organization name
	if strings.TrimSpace(orgName) == "" {
		return fmt.Errorf("organization name cannot be empty")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	workspaceFilePath := filepath.Join(currentDir, workspace.WorkspaceFileName)

	// Check if workspace file already exists
	if _, err := os.Stat(workspaceFilePath); err == nil {
		return fmt.Errorf("workspace file already exists: %s\nUse 'nimsforestpm workspace status' to view current configuration", workspaceFilePath)
	}

	// Create new workspace
	ws := workspace.NewWorkspace()
	ws.Organization = strings.TrimSpace(orgName)

	// Save the workspace file
	if err := ws.Save(workspaceFilePath); err != nil {
		return fmt.Errorf("failed to save workspace file: %w", err)
	}

	fmt.Printf("✓ Workspace initialized: %s\n", workspaceFilePath)
	fmt.Printf("  Organization: %s\n", ws.Organization)
	fmt.Println("  Version: " + ws.Version)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  - Use 'nimsforestpm workspace add <path>' to add products")
	fmt.Println("  - Use 'nimsforestpm workspace status' to view configuration")
	fmt.Println("  - Use 'nimsforestpm install <tool>' to install nimsforest tools")
	
	return nil
}

// showWorkspaceStatus displays the current workspace status
func showWorkspaceStatus() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	// Try to find workspace file
	workspaceFile, err := workspace.FindWorkspaceFile(currentDir)
	if err != nil {
		return fmt.Errorf("workspace file not found: %w", err)
	}

	// Load workspace
	ws, err := workspace.LoadWorkspace(workspaceFile)
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}

	// Display workspace information
	fmt.Printf("=== Workspace Status ===\n")
	fmt.Printf("File: %s\n", workspaceFile)
	fmt.Printf("Version: %s\n", ws.Version)
	
	if ws.Organization != "" {
		fmt.Printf("Organization: %s\n", ws.Organization)
	}

	fmt.Printf("Products (%d):\n", len(ws.Products))
	if len(ws.Products) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, product := range ws.Products {
			fmt.Printf("  - %s\n", product)
		}
	}

	fmt.Printf("Tools (%d):\n", len(ws.Tools))
	if len(ws.Tools) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, tool := range ws.Tools {
			fmt.Printf("  - %s (%s) at %s\n", tool.Name, tool.Mode, tool.Path)
		}
	}

	// Validate workspace
	fmt.Printf("\n=== Validation ===\n")
	if err := ws.Validate(); err != nil {
		fmt.Printf("⚠️  Validation warnings: %v\n", err)
	} else {
		fmt.Printf("✓ Workspace is valid\n")
	}

	// Show absolute paths
	orgPath, productPaths, err := ws.GetAbsolutePaths()
	if err != nil {
		fmt.Printf("⚠️  Error resolving paths: %v\n", err)
	} else {
		fmt.Printf("\n=== Resolved Paths ===\n")
		if orgPath != "" {
			fmt.Printf("Organization: %s\n", orgPath)
		}
		if len(productPaths) > 0 {
			fmt.Printf("Products:\n")
			for _, path := range productPaths {
				fmt.Printf("  - %s\n", path)
			}
		}
	}

	return nil
}

// addProductToWorkspace adds a product to the workspace file
func addProductToWorkspace(productPath string) error {
	// Validate product path
	if strings.TrimSpace(productPath) == "" {
		return fmt.Errorf("product path cannot be empty")
	}

	productPath = strings.TrimSpace(productPath)

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	// Find workspace file
	workspaceFile, err := workspace.FindWorkspaceFile(currentDir)
	if err != nil {
		return fmt.Errorf("workspace file not found: %w\nRun 'nimsforestpm workspace init <org-name>' to create one", err)
	}

	// Load workspace
	ws, err := workspace.LoadWorkspace(workspaceFile)
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}

	// Check if product already exists
	for _, existing := range ws.Products {
		if existing == productPath {
			fmt.Printf("Product already exists in workspace: %s\n", productPath)
			return nil
		}
	}

	// Add product
	ws.AddProduct(productPath)

	// Save workspace
	if err := ws.Save(workspaceFile); err != nil {
		return fmt.Errorf("failed to save workspace: %w", err)
	}

	fmt.Printf("✓ Added product to workspace: %s\n", productPath)
	fmt.Printf("  Total products: %d\n", len(ws.Products))
	return nil
}

// removeProductFromWorkspace removes a product from the workspace file
func removeProductFromWorkspace(productPath string) error {
	// Validate product path
	if strings.TrimSpace(productPath) == "" {
		return fmt.Errorf("product path cannot be empty")
	}

	productPath = strings.TrimSpace(productPath)

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	// Find workspace file
	workspaceFile, err := workspace.FindWorkspaceFile(currentDir)
	if err != nil {
		return fmt.Errorf("workspace file not found: %w\nRun 'nimsforestpm workspace init <org-name>' to create one", err)
	}

	// Load workspace
	ws, err := workspace.LoadWorkspace(workspaceFile)
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}

	// Check if product exists
	found := false
	for _, product := range ws.Products {
		if product == productPath {
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Available products:\n")
		if len(ws.Products) == 0 {
			fmt.Println("  (none)")
		} else {
			for _, product := range ws.Products {
				fmt.Printf("  - %s\n", product)
			}
		}
		return fmt.Errorf("product not found in workspace: %s", productPath)
	}

	// Remove product
	ws.RemoveProduct(productPath)

	// Save workspace
	if err := ws.Save(workspaceFile); err != nil {
		return fmt.Errorf("failed to save workspace: %w", err)
	}

	fmt.Printf("✓ Removed product from workspace: %s\n", productPath)
	fmt.Printf("  Remaining products: %d\n", len(ws.Products))
	return nil
}

// updateWorkspaceAfterInstall updates the workspace file after a successful installation
func updateWorkspaceAfterInstall(component string) error {
	// Find workspace root
	workspaceRoot, err := findWorkspaceRoot()
	if err != nil {
		return fmt.Errorf("workspace root not found: %w", err)
	}

	// Try to find workspace file
	workspaceFile, err := workspace.FindWorkspaceFile(workspaceRoot)
	if err != nil {
		// No workspace file found, create one
		if err := createWorkspaceFile(workspaceRoot, component); err != nil {
			return fmt.Errorf("failed to create workspace file: %w", err)
		}
		return nil
	}

	// Load workspace
	ws, err := workspace.LoadWorkspace(workspaceFile)
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}

	// Determine the product path based on typical workspace structure
	// Usually components are installed in products-workspace/{component}-workspace
	productPath := fmt.Sprintf("products-workspace/%s-workspace", component)

	// Add product to workspace if not already present
	ws.AddProduct(productPath)

	// Save workspace
	if err := ws.Save(workspaceFile); err != nil {
		return fmt.Errorf("failed to save workspace: %w", err)
	}

	fmt.Printf("✓ Updated workspace file: added %s\n", productPath)
	return nil
}

// createWorkspaceFile creates a new workspace file when one doesn't exist
func createWorkspaceFile(workspaceRoot, component string) error {
	workspaceFilePath := filepath.Join(workspaceRoot, workspace.WorkspaceFileName)
	
	// Create new workspace
	ws := workspace.NewWorkspace()
	
	// Try to determine organization name from workspace structure
	if orgName := getOrganizationName(workspaceRoot); orgName != "" {
		ws.Organization = orgName
	}
	
	// Add the component as first product
	productPath := fmt.Sprintf("products-workspace/%s-workspace", component)
	ws.AddProduct(productPath)
	
	// Save the workspace file
	if err := ws.Save(workspaceFilePath); err != nil {
		return fmt.Errorf("failed to save new workspace file: %w", err)
	}
	
	fmt.Printf("✓ Created workspace file: %s\n", workspaceFilePath)
	fmt.Printf("✓ Added %s to workspace\n", productPath)
	return nil
}

// getOrganizationName extracts organization name from workspace structure
func getOrganizationName(workspaceRoot string) string {
	entries, err := os.ReadDir(workspaceRoot)
	if err != nil {
		return ""
	}
	
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), "-organization-workspace") {
			// Extract organization name by removing the suffix
			return strings.TrimSuffix(entry.Name(), "-organization-workspace")
		}
	}
	return ""
}

// findWorkspaceFileFromDir is a helper function that wraps workspace.FindWorkspaceFile
// and provides better integration with existing code
func findWorkspaceFileFromDir(dir string) (string, bool) {
	workspaceFile, err := workspace.FindWorkspaceFile(dir)
	return workspaceFile, err == nil
}

// hasWorkspaceFile checks if there's a workspace file in the current directory tree
func hasWorkspaceFile() bool {
	currentDir, err := os.Getwd()
	if err != nil {
		return false
	}
	_, found := findWorkspaceFileFromDir(currentDir)
	return found
}

// findWorkspaceRoot finds the workspace root directory using the same logic as detectWorkspace
func findWorkspaceRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current directory: %w", err)
	}

	// Walk up directory tree looking for workspace structure
	for dir := cwd; dir != "/" && dir != "."; dir = filepath.Dir(dir) {
		if hasWorkspaceStructure(dir) {
			return dir, nil
		}
	}
	return "", fmt.Errorf("no NimsForest workspace found. Look for a directory with 'products-workspace/'")
}

// hasWorkspaceStructure checks if directory contains organizational workspace and products-workspace
func hasWorkspaceStructure(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	hasOrgWorkspace := false
	hasProductsWorkspace := false

	for _, entry := range entries {
		if entry.IsDir() {
			if strings.HasSuffix(entry.Name(), "-organization-workspace") {
				hasOrgWorkspace = true
			}
			if entry.Name() == "products-workspace" {
				hasProductsWorkspace = true
			}
		}
	}
	return hasOrgWorkspace && hasProductsWorkspace
}

// installComponentAsSubmodule installs a component as a git submodule
// TODO: This will be used for --mode submodule installation
// It should install the binary first, then ALSO clone the submodule for development
func installComponentAsSubmodule(productsDir, component string) error {
	fmt.Printf("📦 Installing %s as submodule...\n", component)
	
	// Validate that we're in a git repository
	if err := validateGitRepository(productsDir); err != nil {
		return fmt.Errorf("git repository validation failed: %w", err)
	}
	
	// Construct the GitHub URL
	repoURL := fmt.Sprintf("https://github.com/nimsforest/%s.git", component)
	submodulePath := component + "-workspace"
	
	// Validate that the remote repository exists
	if err := validateRemoteRepository(repoURL); err != nil {
		return fmt.Errorf("remote repository validation failed: %w", err)
	}
	
	// Add the submodule
	cmd := exec.Command("git", "submodule", "add", repoURL, submodulePath)
	cmd.Dir = productsDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add git submodule for %s: %w\n\nThis could be due to:\n- Network connectivity issues\n- Repository %s doesn't exist\n- Git authentication problems\n- Directory %s already exists", component, err, repoURL, submodulePath)
	}
	
	return nil
}

// registerInstalledTool registers the installed tool in the nimsforest.workspace file
func registerInstalledTool(workspaceRoot, component string) error {
	// Load the workspace file
	ws, err := workspace.LoadWorkspaceFromDir(workspaceRoot)
	if err != nil {
		return fmt.Errorf("failed to load workspace: %w", err)
	}
	
	// Create a tool entry for binary installation
	toolEntry := workspace.ToolEntry{
		Name:    component,
		Mode:    "binary",
		Path:    fmt.Sprintf("bin/%s", component),
		Version: "latest",
	}
	
	// Add the tool to the workspace
	ws.AddTool(toolEntry)
	
	// Save the workspace file
	if err := ws.Save(ws.FilePath); err != nil {
		return fmt.Errorf("failed to save workspace: %w", err)
	}
	
	fmt.Printf("🔗 Registered %s in workspace file\n", component)
	return nil
}

// installComponentAsBinary installs a component as a binary
func installComponentAsBinary(component string) error {
	fmt.Printf("📦 Installing %s as binary...\n", component)
	
	// TODO: Implement binary installation
	// This should:
	// 1. Download the binary from GitHub releases or build from source
	// 2. Install it to workspace/bin/ directory
	// 3. Make it executable
	// 4. Verify it works
	
	// For now, we'll simulate binary installation
	fmt.Printf("⚠️  Binary installation not implemented yet - tool will be available when implemented\n")
	return nil
}

// validateGitRepository checks if the directory is a valid git repository
func validateGitRepository(dir string) error {
	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git command not found. Please install git")
	}
	
	// Check if we're in a git repository
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not in a git repository. Please run 'git init' in the workspace root")
	}
	
	return nil
}

// validateRemoteRepository checks if the remote repository exists and is accessible
func validateRemoteRepository(repoURL string) error {
	fmt.Printf("🔍 Validating remote repository...\n")
	
	// Use git ls-remote to check if repository exists and is accessible
	cmd := exec.Command("git", "ls-remote", "--exit-code", repoURL)
	cmd.Stdout = nil // Suppress output
	cmd.Stderr = nil // Suppress error output
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("repository %s is not accessible. Please check:\n- Repository exists\n- Network connectivity\n- Git authentication if needed", repoURL)
	}
	
	fmt.Printf("✓ Remote repository validated\n")
	return nil
}

// validateComponentName validates the component name format
func validateComponentName(component string) error {
	if component == "" {
		return fmt.Errorf("component name cannot be empty")
	}
	
	// Check that it starts with "nimsforest"
	if !strings.HasPrefix(component, "nimsforest") {
		return fmt.Errorf("component name must start with 'nimsforest', got: %s", component)
	}
	
	// Check for valid characters (alphanumeric and no spaces)
	for _, char := range component {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
			return fmt.Errorf("component name contains invalid character '%c'. Only lowercase letters and numbers allowed", char)
		}
	}
	
	return nil
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a tool implementation",
	Long: `Validate a tool implementation to ensure it follows nimsforest tool interface standards.

This command validates:
- Tool interface implementation
- Metadata quality (name, version, description)
- Command structure and handlers
- Health check functionality
- Installation mode support`,
	Run: func(cmd *cobra.Command, args []string) {
		toolPath, _ := cmd.Flags().GetString("tool-path")
		verbose, _ := cmd.Flags().GetBool("verbose")
		testCommands, _ := cmd.Flags().GetBool("test-commands")
		timeout, _ := cmd.Flags().GetDuration("timeout")

		if toolPath == "" {
			fmt.Println("Error: --tool-path is required")
			os.Exit(1)
		}

		// Create validation options
		options := tool.ValidationOptions{
			ToolPath:         toolPath,
			InterfaceVersion: "1.0.0",
			TestCommands:     testCommands,
			Verbose:          verbose,
			Timeout:          timeout,
		}

		// Create validator and validate
		validator := tool.NewToolValidator(options)
		
		fmt.Printf("Validating tool at: %s\n", toolPath)
		if testCommands {
			fmt.Println("Running command tests...")
		}
		
		result, err := validator.ValidateTool(context.Background(), toolPath)
		if err != nil {
			fmt.Printf("Validation failed: %v\n", err)
			os.Exit(1)
		}

		// Format and display results
		output := tool.FormatValidationResult(result, verbose)
		fmt.Print(output)

		// Exit with appropriate code
		if !result.Valid {
			os.Exit(1)
		}
	},
}

func init() {
	// Add flags for validate command
	validateCmd.Flags().String("tool-path", "", "Path to the tool to validate (required)")
	validateCmd.Flags().Bool("verbose", false, "Show detailed validation information")
	validateCmd.Flags().Bool("test-commands", true, "Test command execution")
	validateCmd.Flags().Duration("timeout", 30*time.Second, "Timeout for command tests")
	
	// Mark tool-path as required
	validateCmd.MarkFlagRequired("tool-path")
}