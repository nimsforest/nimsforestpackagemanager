package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ToolInfo represents information about a tool
type ToolInfo struct {
	Repository  string `json:"repository"`
	Description string `json:"description"`
}

// ToolRegistry represents the tools.json structure
type ToolRegistry struct {
	Tools   map[string]ToolInfo `json:"tools"`
	Version string              `json:"version"`
	Updated string              `json:"updated"`
}

var registry *ToolRegistry

// LoadRegistry loads the tools.json file
func LoadRegistry() (*ToolRegistry, error) {
	if registry != nil {
		return registry, nil
	}

	data, err := os.ReadFile("docs/tools.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read tools.json: %v", err)
	}

	var reg ToolRegistry
	if err := json.Unmarshal(data, &reg); err != nil {
		return nil, fmt.Errorf("failed to parse tools.json: %v", err)
	}

	registry = &reg
	return registry, nil
}

// ResolveToolRepository converts tool names to GitHub repository paths
func ResolveToolRepository(toolName string) (string, error) {
	// Handle full repository paths directly
	if strings.Contains(toolName, "/") {
		return toolName, nil
	}
	
	// Load registry and look up tool
	reg, err := LoadRegistry()
	if err != nil {
		return "", err
	}
	
	if tool, exists := reg.Tools[toolName]; exists {
		return tool.Repository, nil
	}
	
	return "", fmt.Errorf("unknown tool: %s. Available tools: %s", toolName, strings.Join(AvailableTools(), ", "))
}

// InstallTool installs a tool using go get and go install
func InstallTool(toolName string) error {
	repo, err := ResolveToolRepository(toolName)
	if err != nil {
		return err
	}
	
	fmt.Printf("Installing %s from %s...\n", toolName, repo)

	// Step 1: go get the tool
	cmd := exec.Command("go", "get", repo+"@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to get %s: %v", toolName, err)
	}

	// Step 2: go install the tool
	cmd = exec.Command("go", "install", repo+"@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install %s: %v", toolName, err)
	}

	fmt.Printf("✓ %s installed successfully!\n", toolName)
	fmt.Printf("Tool available as: %s\n", toolName)
	return nil
}

// UpdateTool updates a tool using go get -u and go install
func UpdateTool(toolName string) error {
	repo, err := ResolveToolRepository(toolName)
	if err != nil {
		return err
	}
	
	fmt.Printf("Updating %s from %s...\n", toolName, repo)

	// Step 1: go get -u the tool
	cmd := exec.Command("go", "get", "-u", repo+"@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update %s: %v", toolName, err)
	}

	// Step 2: go install the tool
	cmd = exec.Command("go", "install", repo+"@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install updated %s: %v", toolName, err)
	}

	fmt.Printf("✓ %s updated successfully!\n", toolName)
	return nil
}

// IsToolInstalled checks if a tool is installed in $GOPATH/bin
func IsToolInstalled(toolName string) bool {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		// Use default GOPATH
		home, err := os.UserHomeDir()
		if err != nil {
			return false
		}
		gopath = filepath.Join(home, "go")
	}

	binaryPath := filepath.Join(gopath, "bin", toolName)
	_, err := os.Stat(binaryPath)
	return err == nil
}

// AvailableTools returns a list of known nimsforest tools
func AvailableTools() []string {
	reg, err := LoadRegistry()
	if err != nil {
		return []string{} // Return empty if registry can't be loaded
	}
	
	tools := make([]string, 0, len(reg.Tools))
	for name := range reg.Tools {
		tools = append(tools, name)
	}
	return tools
}

// InstalledTools returns a list of installed nimsforest tools
func InstalledTools() []string {
	available := AvailableTools()
	installed := make([]string, 0)
	
	for _, tool := range available {
		if IsToolInstalled(tool) {
			installed = append(installed, tool)
		}
	}
	return installed
}

// GetToolInfo returns information about a specific tool
func GetToolInfo(toolName string) (ToolInfo, error) {
	reg, err := LoadRegistry()
	if err != nil {
		return ToolInfo{}, err
	}
	
	if tool, exists := reg.Tools[toolName]; exists {
		return tool, nil
	}
	
	return ToolInfo{}, fmt.Errorf("unknown tool: %s", toolName)
}