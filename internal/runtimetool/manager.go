package runtimetool

import (
	"context"
	"fmt"

	"github.com/nimsforest/nimsforestpackagemanager/internal/workspace"
)

// Manager manages runtime tools based on workspace configuration
type Manager struct {
	workspace *workspace.Workspace
}

// NewManager creates a new tool manager
func NewManager(ws *workspace.Workspace) *Manager {
	return &Manager{
		workspace: ws,
	}
}

// GetTool retrieves a tool by name
func (m *Manager) GetTool(name string) (*RuntimeTool, error) {
	entry, err := m.workspace.GetTool(name)
	if err != nil {
		return nil, err
	}
	
	return NewRuntimeTool(*entry, m.workspace), nil
}

// ListTools returns all available tools
func (m *Manager) ListTools() []*RuntimeTool {
	entries := m.workspace.GetInstalledTools()
	tools := make([]*RuntimeTool, len(entries))
	
	for i, entry := range entries {
		tools[i] = NewRuntimeTool(entry, m.workspace)
	}
	
	return tools
}

// ExecuteCommand executes a command on a specific tool
func (m *Manager) ExecuteCommand(ctx context.Context, toolName, command string, args []string) error {
	tool, err := m.GetTool(toolName)
	if err != nil {
		return err
	}
	
	return tool.Execute(ctx, command, args)
}

// ValidateAllTools validates all tools in the workspace
func (m *Manager) ValidateAllTools(ctx context.Context) error {
	tools := m.ListTools()
	
	for _, tool := range tools {
		if err := tool.Validate(ctx); err != nil {
			return fmt.Errorf("validation failed for tool %s: %w", tool.Name(), err)
		}
	}
	
	return nil
}

// GetToolCommands gets available commands for a specific tool
func (m *Manager) GetToolCommands(ctx context.Context, toolName string) ([]string, error) {
	tool, err := m.GetTool(toolName)
	if err != nil {
		return nil, err
	}
	
	return tool.GetCommands(ctx)
}

// AddTool adds a new tool to the workspace
func (m *Manager) AddTool(entry workspace.ToolEntry) error {
	m.workspace.AddTool(entry)
	return nil
}

// RemoveTool removes a tool from the workspace
func (m *Manager) RemoveTool(toolName string) error {
	m.workspace.RemoveTool(toolName)
	return nil
}

// SaveWorkspace saves the current workspace state
func (m *Manager) SaveWorkspace() error {
	if m.workspace.FilePath == "" {
		return fmt.Errorf("workspace file path not set")
	}
	
	return m.workspace.Save(m.workspace.FilePath)
}

// LoadWorkspace loads or reloads the workspace
func (m *Manager) LoadWorkspace(filePath string) error {
	ws, err := workspace.LoadWorkspace(filePath)
	if err != nil {
		return err
	}
	
	m.workspace = ws
	return nil
}

// GetWorkspace returns the underlying workspace
func (m *Manager) GetWorkspace() *workspace.Workspace {
	return m.workspace
}