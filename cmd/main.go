package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nimsforestpm",
	Short: "NimsForest Package Manager - Bootstrap and manage organizational workspaces",
	Long: `NimsForest Package Manager creates and manages organizational workspaces
where organizations can explicitly optimize their coordination (organize) and 
value creation (productize) in an endless improvement cycle.`,
}

func main() {
	// Discover and add dynamic subcommands for installed tools
	if err := addDynamicCommands(); err != nil {
		// Don't fail if we can't discover commands (might not be in workspace)
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// detectWorkspace finds the workspace root by walking up the directory tree
func detectWorkspace() (string, error) {
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
	return "", fmt.Errorf("not in a nimsforest workspace")
}

// hasWorkspaceStructure checks if directory contains {orgname}-organization-workspace and products-workspace
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

// Tool represents an installed nimsforest tool
type Tool struct {
	Name        string   // e.g., "work"
	FullName    string   // e.g., "nimsforestwork"
	Path        string   // Path to the tool workspace
	Commands    []string // Available commands for this tool
	Description string   // Tool description
}

// discoverInstalledTools scans products-workspace for installed tools
func discoverInstalledTools(workspaceRoot string) ([]Tool, error) {
	var tools []Tool
	
	productsDir := filepath.Join(workspaceRoot, "products-workspace")
	entries, err := os.ReadDir(productsDir)
	if err != nil {
		return nil, fmt.Errorf("cannot read products-workspace: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), "-workspace") {
			toolName := strings.TrimSuffix(entry.Name(), "-workspace")
			
			// Skip if not a nimsforest tool
			if !strings.HasPrefix(toolName, "nimsforest") {
				continue
			}

			makefilePath := filepath.Join(productsDir, entry.Name(), "main", "MAKEFILE."+toolName)
			if _, err := os.Stat(makefilePath); err == nil {
				tool := Tool{
					Name:        strings.TrimPrefix(toolName, "nimsforest"),
					FullName:    toolName,
					Path:        filepath.Join(productsDir, entry.Name(), "main"),
					Commands:    extractCommands(makefilePath, toolName),
					Description: getToolDescription(toolName),
				}
				tools = append(tools, tool)
			}
		}
	}
	return tools, nil
}

// extractCommands parses a MAKEFILE to extract available commands
func extractCommands(makefilePath, toolName string) []string {
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		return []string{"hello", "init", "lint"} // fallback
	}

	var commands []string
	lines := strings.Split(string(content), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Look for lines like "toolname-command:"
		if strings.HasPrefix(line, toolName+"-") && strings.HasSuffix(line, ":") {
			// Extract command name
			full := strings.TrimSuffix(line, ":")
			if cmd := strings.TrimPrefix(full, toolName+"-"); cmd != full {
				commands = append(commands, cmd)
			}
		}
	}
	
	if len(commands) == 0 {
		return []string{"hello", "init", "lint"} // fallback
	}
	
	return commands
}

// getToolDescription returns description for a tool
func getToolDescription(toolName string) string {
	descriptions := map[string]string{
		"nimsforestwork":          "Work Management and Tracking System",
		"nimsforestorganize":      "Organizational Structure Component", 
		"nimsforestcommunication": "Communication Systems Component",
		"nimsforestproductize":    "Product Development Orchestrator",
		"nimsforestfolders":       "Advanced Folder Management System",
		"nimsforestwebstack":      "Web Development and Deployment Tools",
	}
	if desc, ok := descriptions[toolName]; ok {
		return desc
	}
	return toolName
}

// addDynamicCommands discovers installed tools and adds them as subcommands
func addDynamicCommands() error {
	workspaceRoot, err := detectWorkspace()
	if err != nil {
		return err
	}

	tools, err := discoverInstalledTools(workspaceRoot)
	if err != nil {
		return err
	}

	for _, tool := range tools {
		toolCmd := &cobra.Command{
			Use:   tool.Name,
			Short: tool.Description,
			Long:  fmt.Sprintf("%s - %s", tool.FullName, tool.Description),
		}

		// Add subcommands for each tool command
		for _, cmdName := range tool.Commands {
			cmd := createToolCommand(tool, cmdName)
			toolCmd.AddCommand(cmd)
		}

		rootCmd.AddCommand(toolCmd)
	}

	return nil
}

// createToolCommand creates a cobra command that proxies to make
func createToolCommand(tool Tool, cmdName string) *cobra.Command {
	return &cobra.Command{
		Use:   cmdName,
		Short: fmt.Sprintf("Run %s-%s command", tool.FullName, cmdName),
		Run: func(cmd *cobra.Command, args []string) {
			runMakeCommand(tool.Path, tool.FullName, cmdName, args)
		},
	}
}

// runMakeCommand executes the corresponding make command
func runMakeCommand(toolPath, toolName, cmdName string, args []string) {
	makeTarget := fmt.Sprintf("%s-%s", toolName, cmdName)
	makeCmd := exec.Command("make", "-C", toolPath, makeTarget)
	
	// Add any additional args as make variables
	for _, arg := range args {
		if strings.Contains(arg, "=") {
			makeCmd.Args = append(makeCmd.Args, arg)
		}
	}
	
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr
	makeCmd.Stdin = os.Stdin

	if err := makeCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running make command: %v\n", err)
		os.Exit(1)
	}
}