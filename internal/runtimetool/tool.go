package runtimetool

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nimsforest/nimsforestpackagemanager/internal/workspace"
)

// RuntimeTool represents a tool that can be executed at runtime
type RuntimeTool struct {
	entry     workspace.ToolEntry
	workspace *workspace.Workspace
}

// NewRuntimeTool creates a new runtime tool from a workspace entry
func NewRuntimeTool(entry workspace.ToolEntry, ws *workspace.Workspace) *RuntimeTool {
	return &RuntimeTool{
		entry:     entry,
		workspace: ws,
	}
}

// Name returns the tool name
func (t *RuntimeTool) Name() string {
	return t.entry.Name
}

// Version returns the tool version
func (t *RuntimeTool) Version() string {
	return t.entry.Version
}

// Mode returns the installation mode
func (t *RuntimeTool) Mode() string {
	return t.entry.Mode
}

// Path returns the installation path
func (t *RuntimeTool) Path() string {
	return t.entry.Path
}

// GetExecutablePath returns the path to the executable based on the mode
func (t *RuntimeTool) GetExecutablePath() (string, error) {
	switch t.entry.Mode {
	case "binary":
		// For binary mode, the path should point directly to the executable
		return t.entry.Path, nil
	case "clone", "submodule":
		// For clone/submodule mode, look for the executable in the workspace
		// Check for a Makefile or go.mod to determine how to build/find the executable
		if t.workspace != nil && t.workspace.FilePath != "" {
			workspaceDir := filepath.Dir(t.workspace.FilePath)
			toolPath := filepath.Join(workspaceDir, t.entry.Path)
			
			// Check if there's a binary with the tool name in the tool directory
			binaryPath := filepath.Join(toolPath, t.entry.Name)
			if _, err := os.Stat(binaryPath); err == nil {
				return binaryPath, nil
			}
			
			// Check if there's a bin directory
			binPath := filepath.Join(toolPath, "bin", t.entry.Name)
			if _, err := os.Stat(binPath); err == nil {
				return binPath, nil
			}
			
			// Return the tool path and let the caller handle execution
			return toolPath, nil
		}
		return t.entry.Path, nil
	default:
		return "", fmt.Errorf("unsupported tool mode: %s", t.entry.Mode)
	}
}

// GetCommands discovers available commands by calling the tool
func (t *RuntimeTool) GetCommands(ctx context.Context) ([]string, error) {
	execPath, err := t.GetExecutablePath()
	if err != nil {
		return nil, err
	}
	
	// Try common help flags to discover commands
	helpFlags := []string{"--help", "-h", "help"}
	
	for _, flag := range helpFlags {
		cmd := exec.CommandContext(ctx, execPath, flag)
		output, err := cmd.CombinedOutput()
		if err != nil {
			continue // Try next flag
		}
		
		// Parse output to extract commands
		commands := t.parseCommandsFromHelp(string(output))
		if len(commands) > 0 {
			return commands, nil
		}
	}
	
	// If no help available, try to run without arguments to see if it lists commands
	cmd := exec.CommandContext(ctx, execPath)
	output, err := cmd.CombinedOutput()
	if err == nil {
		commands := t.parseCommandsFromHelp(string(output))
		if len(commands) > 0 {
			return commands, nil
		}
	}
	
	return []string{}, nil
}

// parseCommandsFromHelp attempts to parse command names from help output
func (t *RuntimeTool) parseCommandsFromHelp(output string) []string {
	var commands []string
	lines := strings.Split(output, "\n")
	
	// Look for patterns like "Commands:" or "Available commands:"
	inCommandsSection := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Check if we're entering a commands section
		if strings.Contains(strings.ToLower(line), "command") && strings.Contains(line, ":") {
			inCommandsSection = true
			continue
		}
		
		// If we're in commands section, parse command names
		if inCommandsSection {
			// Stop if we hit an empty line or a line that doesn't look like a command
			if line == "" || (!strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t")) {
				break
			}
			
			// Extract command name (first word after whitespace)
			fields := strings.Fields(line)
			if len(fields) > 0 {
				command := fields[0]
				// Basic validation - command names shouldn't contain special characters
				if isValidCommandName(command) {
					commands = append(commands, command)
				}
			}
		}
	}
	
	return commands
}

// isValidCommandName checks if a string looks like a valid command name
func isValidCommandName(name string) bool {
	if name == "" {
		return false
	}
	
	// Command names should be alphanumeric with hyphens/underscores
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || 
			 (char >= 'A' && char <= 'Z') || 
			 (char >= '0' && char <= '9') || 
			 char == '-' || char == '_') {
			return false
		}
	}
	
	return true
}

// Execute runs a command with the given arguments
func (t *RuntimeTool) Execute(ctx context.Context, command string, args []string) error {
	execPath, err := t.GetExecutablePath()
	if err != nil {
		return err
	}
	
	// Build the full command with arguments
	cmdArgs := []string{command}
	cmdArgs = append(cmdArgs, args...)
	
	cmd := exec.CommandContext(ctx, execPath, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	return cmd.Run()
}

// Validate checks if the tool is properly installed and accessible
func (t *RuntimeTool) Validate(ctx context.Context) error {
	execPath, err := t.GetExecutablePath()
	if err != nil {
		return err
	}
	
	// Check if the executable/directory exists
	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		return fmt.Errorf("tool path does not exist: %s", execPath)
	}
	
	// For binary mode, check if it's executable
	if t.entry.Mode == "binary" {
		info, err := os.Stat(execPath)
		if err != nil {
			return err
		}
		
		// Check if it's executable (Unix permissions)
		if info.Mode()&0111 == 0 {
			return fmt.Errorf("tool binary is not executable: %s", execPath)
		}
	}
	
	return nil
}