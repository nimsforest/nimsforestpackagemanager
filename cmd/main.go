package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/nimsforest/nimsforestpackagemanager/internal/workspace"
	"github.com/nimsforest/nimsforestpackagemanager/internal/runtimetool"
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


// Tool represents an installed nimsforest tool
type Tool struct {
	Name        string   // e.g., "work"
	FullName    string   // e.g., "nimsforestwork"
	Path        string   // Path to the tool workspace
	Commands    []string // Available commands for this tool
	Description string   // Tool description
}

// Note: Tool discovery and validation is now handled by the Makefile
// The status command delegates to 'make nimsforestpm-lint' for consistent behavior

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

// addDynamicCommands discovers and adds dynamic subcommands for installed tools
func addDynamicCommands() error {
	// Try to find workspace file
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}
	
	// Load workspace if available
	ws, err := workspace.LoadWorkspaceFromDir(currentDir)
	if err != nil {
		// Not in a workspace, skip dynamic commands
		return nil
	}
	
	// Create runtime tool manager
	manager := runtimetool.NewManager(ws)
	
	// Get all installed tools
	tools := manager.ListTools()
	if len(tools) == 0 {
		return nil
	}
	
	// Add commands for each tool
	for _, tool := range tools {
		// Create tool command
		toolCmd := createDynamicToolCommand(tool, manager)
		rootCmd.AddCommand(toolCmd)
	}
	
	return nil
}

// createDynamicToolCommand creates a cobra command for a runtime tool
func createDynamicToolCommand(tool *runtimetool.RuntimeTool, manager *runtimetool.Manager) *cobra.Command {
	// Extract short name from full tool name (e.g., "nimsforestwork" -> "work")
	shortName := extractShortName(tool.Name())
	
	return &cobra.Command{
		Use:   shortName,
		Short: fmt.Sprintf("Run %s commands", tool.Name()),
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// First argument is the command, rest are arguments
			command := args[0]
			cmdArgs := args[1:]
			
			// Execute the command using the runtime tool manager
			return manager.ExecuteCommand(context.Background(), tool.Name(), command, cmdArgs)
		},
	}
}

// extractShortName extracts the short name from a full tool name
func extractShortName(fullName string) string {
	// Remove "nimsforest" prefix if present
	if strings.HasPrefix(fullName, "nimsforest") {
		return strings.TrimPrefix(fullName, "nimsforest")
	}
	return fullName
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