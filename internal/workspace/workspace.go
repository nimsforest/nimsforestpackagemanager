package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ToolEntry represents a tool installation in the workspace
type ToolEntry struct {
	Name        string `json:"name"`
	Mode        string `json:"mode"`        // "binary", "clone", "submodule"
	Path        string `json:"path"`
	Version     string `json:"version"`
}

// Workspace represents a nimsforest workspace configuration
type Workspace struct {
	Version      string      `json:"version"`
	Organization string      `json:"organization"`
	Products     []string    `json:"products"`
	Tools        []ToolEntry `json:"tools"`
	FilePath     string      `json:"file_path"`
}

// WorkspaceFileName is the standard name for workspace files
const WorkspaceFileName = "nimsforest.workspace"

// NewWorkspace creates a new workspace with default values
func NewWorkspace() *Workspace {
	return &Workspace{
		Version:  "1.0",
		Products: make([]string, 0),
		Tools:    make([]ToolEntry, 0),
	}
}

// FindWorkspaceFile searches for a workspace file starting from the given directory
// and walking up the directory tree until it finds one or reaches the root
func FindWorkspaceFile(startDir string) (string, error) {
	if startDir == "" {
		var err error
		startDir, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current working directory: %w", err)
		}
	}

	// Make sure we have an absolute path
	absPath, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	currentDir := absPath
	for {
		workspaceFile := filepath.Join(currentDir, WorkspaceFileName)
		if _, err := os.Stat(workspaceFile); err == nil {
			return workspaceFile, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			// We've reached the root directory
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("workspace file not found in directory tree starting from %s", startDir)
}

// LoadWorkspace loads a workspace from a file
func LoadWorkspace(filePath string) (*Workspace, error) {
	if filePath == "" {
		return nil, fmt.Errorf("file path cannot be empty")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("workspace file does not exist: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read workspace file: %w", err)
	}

	workspace, err := ParseWorkspace(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse workspace file: %w", err)
	}

	workspace.FilePath = filePath
	return workspace, nil
}

// LoadWorkspaceFromDir loads a workspace by searching for the workspace file
// starting from the given directory
func LoadWorkspaceFromDir(dir string) (*Workspace, error) {
	workspaceFile, err := FindWorkspaceFile(dir)
	if err != nil {
		return nil, err
	}

	return LoadWorkspace(workspaceFile)
}

// Save saves the workspace to a file
func (w *Workspace) Save(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	content := w.String()
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write workspace file: %w", err)
	}

	w.FilePath = filePath
	return nil
}

// String returns the string representation of the workspace in nimsforest.workspace format
func (w *Workspace) String() string {
	var sb strings.Builder
	
	sb.WriteString(fmt.Sprintf("nimsforest %s", w.Version))
	
	if w.Organization != "" || len(w.Products) > 0 {
		sb.WriteString("\n")
	}
	
	if w.Organization != "" {
		sb.WriteString(fmt.Sprintf("\norganization %s", w.Organization))
		if len(w.Products) > 0 {
			sb.WriteString("\n")
		}
	}
	
	if len(w.Products) > 0 {
		sb.WriteString("\nproducts (\n")
		for _, product := range w.Products {
			sb.WriteString(fmt.Sprintf("    %s\n", product))
		}
		sb.WriteString(")")
	}
	
	if len(w.Tools) > 0 {
		if len(w.Products) > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString("\ntools (\n")
		for _, tool := range w.Tools {
			sb.WriteString(fmt.Sprintf("    %s %s %s %s\n", tool.Name, tool.Mode, tool.Path, tool.Version))
		}
		sb.WriteString(")")
	}
	
	return sb.String()
}

// AddProduct adds a product to the workspace
func (w *Workspace) AddProduct(productPath string) {
	// Check if product already exists
	for _, existing := range w.Products {
		if existing == productPath {
			return
		}
	}
	w.Products = append(w.Products, productPath)
}

// AddTool adds a tool to the workspace
func (w *Workspace) AddTool(tool ToolEntry) {
	// Check if tool already exists, if so update it
	for i, existing := range w.Tools {
		if existing.Name == tool.Name {
			w.Tools[i] = tool
			return
		}
	}
	w.Tools = append(w.Tools, tool)
}

// RemoveTool removes a tool from the workspace
func (w *Workspace) RemoveTool(toolName string) {
	for i, tool := range w.Tools {
		if tool.Name == toolName {
			w.Tools = append(w.Tools[:i], w.Tools[i+1:]...)
			return
		}
	}
}

// GetTool retrieves a tool by name
func (w *Workspace) GetTool(toolName string) (*ToolEntry, error) {
	for _, tool := range w.Tools {
		if tool.Name == toolName {
			return &tool, nil
		}
	}
	return nil, fmt.Errorf("tool %s not found in workspace", toolName)
}

// GetInstalledTools returns all tools in the workspace
func (w *Workspace) GetInstalledTools() []ToolEntry {
	return w.Tools
}

// RemoveProduct removes a product from the workspace
func (w *Workspace) RemoveProduct(productPath string) {
	for i, product := range w.Products {
		if product == productPath {
			w.Products = append(w.Products[:i], w.Products[i+1:]...)
			return
		}
	}
}

// Validate validates the workspace structure
func (w *Workspace) Validate() error {
	if w.Version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	// Check if organization path exists if specified
	if w.Organization != "" {
		if !filepath.IsAbs(w.Organization) {
			// If relative path, make it relative to workspace file directory
			if w.FilePath != "" {
				orgPath := filepath.Join(filepath.Dir(w.FilePath), w.Organization)
				if _, err := os.Stat(orgPath); os.IsNotExist(err) {
					return fmt.Errorf("organization path does not exist: %s", orgPath)
				}
			}
		} else {
			// Absolute path
			if _, err := os.Stat(w.Organization); os.IsNotExist(err) {
				return fmt.Errorf("organization path does not exist: %s", w.Organization)
			}
		}
	}

	// Check if product paths exist
	for _, product := range w.Products {
		if !filepath.IsAbs(product) {
			// If relative path, make it relative to workspace file directory
			if w.FilePath != "" {
				productPath := filepath.Join(filepath.Dir(w.FilePath), product)
				if _, err := os.Stat(productPath); os.IsNotExist(err) {
					return fmt.Errorf("product path does not exist: %s", productPath)
				}
			}
		} else {
			// Absolute path
			if _, err := os.Stat(product); os.IsNotExist(err) {
				return fmt.Errorf("product path does not exist: %s", product)
			}
		}
	}

	return nil
}

// GetAbsolutePaths returns absolute paths for organization and products
func (w *Workspace) GetAbsolutePaths() (string, []string, error) {
	if w.FilePath == "" {
		return "", nil, fmt.Errorf("workspace file path not set")
	}

	workspaceDir := filepath.Dir(w.FilePath)
	
	var absOrganization string
	if w.Organization != "" {
		if filepath.IsAbs(w.Organization) {
			absOrganization = w.Organization
		} else {
			absOrganization = filepath.Join(workspaceDir, w.Organization)
		}
	}

	absProducts := make([]string, len(w.Products))
	for i, product := range w.Products {
		if filepath.IsAbs(product) {
			absProducts[i] = product
		} else {
			absProducts[i] = filepath.Join(workspaceDir, product)
		}
	}

	return absOrganization, absProducts, nil
}