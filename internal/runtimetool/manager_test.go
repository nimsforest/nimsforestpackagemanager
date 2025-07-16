package runtimetool

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/nimsforest/nimsforestpackagemanager/internal/workspace"
)

func TestManager_BasicOperations(t *testing.T) {
	// Create a temporary workspace
	ws := workspace.NewWorkspace()
	
	// Add a sample tool
	tool := workspace.ToolEntry{
		Name:    "testtool",
		Mode:    "binary",
		Path:    "/usr/bin/echo",
		Version: "1.0.0",
	}
	ws.AddTool(tool)
	
	// Create manager
	manager := NewManager(ws)
	
	// Test GetTool
	runtimeTool, err := manager.GetTool("testtool")
	if err != nil {
		t.Fatalf("GetTool failed: %v", err)
	}
	
	if runtimeTool.Name() != "testtool" {
		t.Errorf("Expected tool name 'testtool', got %s", runtimeTool.Name())
	}
	
	// Test ListTools
	tools := manager.ListTools()
	if len(tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(tools))
	}
	
	// Test tool doesn't exist
	_, err = manager.GetTool("nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent tool")
	}
}

func TestManager_WorkspaceOperations(t *testing.T) {
	// Create a temporary workspace
	ws := workspace.NewWorkspace()
	manager := NewManager(ws)
	
	// Add tool via manager
	tool := workspace.ToolEntry{
		Name:    "newtool",
		Mode:    "clone",
		Path:    "./tools/newtool",
		Version: "2.0.0",
	}
	
	err := manager.AddTool(tool)
	if err != nil {
		t.Fatalf("AddTool failed: %v", err)
	}
	
	// Verify tool was added
	retrievedTool, err := manager.GetTool("newtool")
	if err != nil {
		t.Fatalf("GetTool failed after adding: %v", err)
	}
	
	if retrievedTool.Version() != "2.0.0" {
		t.Errorf("Expected version '2.0.0', got %s", retrievedTool.Version())
	}
	
	// Remove tool
	err = manager.RemoveTool("newtool")
	if err != nil {
		t.Fatalf("RemoveTool failed: %v", err)
	}
	
	// Verify tool was removed
	_, err = manager.GetTool("newtool")
	if err == nil {
		t.Error("Expected error after removing tool")
	}
}

func TestRuntimeTool_GetExecutablePath(t *testing.T) {
	ws := workspace.NewWorkspace()
	
	// Test binary mode
	binaryTool := workspace.ToolEntry{
		Name:    "binarytool",
		Mode:    "binary",
		Path:    "/usr/bin/echo",
		Version: "1.0.0",
	}
	
	rt := NewRuntimeTool(binaryTool, ws)
	path, err := rt.GetExecutablePath()
	if err != nil {
		t.Fatalf("GetExecutablePath failed: %v", err)
	}
	
	if path != "/usr/bin/echo" {
		t.Errorf("Expected path '/usr/bin/echo', got %s", path)
	}
	
	// Test clone mode
	cloneTool := workspace.ToolEntry{
		Name:    "clonetool",
		Mode:    "clone",
		Path:    "./tools/clonetool",
		Version: "1.0.0",
	}
	
	rt2 := NewRuntimeTool(cloneTool, ws)
	path2, err := rt2.GetExecutablePath()
	if err != nil {
		t.Fatalf("GetExecutablePath failed for clone mode: %v", err)
	}
	
	if path2 != "./tools/clonetool" {
		t.Errorf("Expected path './tools/clonetool', got %s", path2)
	}
}

func TestRuntimeTool_Validate(t *testing.T) {
	ws := workspace.NewWorkspace()
	
	// Test with existing binary
	tool := workspace.ToolEntry{
		Name:    "echo",
		Mode:    "binary",
		Path:    "/usr/bin/echo",
		Version: "1.0.0",
	}
	
	rt := NewRuntimeTool(tool, ws)
	err := rt.Validate(context.Background())
	if err != nil {
		t.Errorf("Validation failed for existing binary: %v", err)
	}
	
	// Test with non-existent binary
	badTool := workspace.ToolEntry{
		Name:    "nonexistent",
		Mode:    "binary",
		Path:    "/nonexistent/path",
		Version: "1.0.0",
	}
	
	rt2 := NewRuntimeTool(badTool, ws)
	err = rt2.Validate(context.Background())
	if err == nil {
		t.Error("Expected validation error for non-existent binary")
	}
}

func TestWorkspaceWithTools(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "workspace-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create a workspace file with tools
	workspaceContent := `nimsforest 1.0

tools (
    testtool binary /usr/bin/echo 1.0.0
    anothertool clone ./tools/another 2.0.0
)`
	
	workspaceFile := filepath.Join(tmpDir, "nimsforest.workspace")
	err = os.WriteFile(workspaceFile, []byte(workspaceContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write workspace file: %v", err)
	}
	
	// Load workspace
	ws, err := workspace.LoadWorkspace(workspaceFile)
	if err != nil {
		t.Fatalf("Failed to load workspace: %v", err)
	}
	
	// Verify tools were loaded
	tools := ws.GetInstalledTools()
	if len(tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(tools))
	}
	
	// Test manager with loaded workspace
	manager := NewManager(ws)
	runtimeTools := manager.ListTools()
	
	if len(runtimeTools) != 2 {
		t.Errorf("Expected 2 runtime tools, got %d", len(runtimeTools))
	}
	
	// Test getting specific tool
	tool, err := manager.GetTool("testtool")
	if err != nil {
		t.Fatalf("Failed to get testtool: %v", err)
	}
	
	if tool.Mode() != "binary" {
		t.Errorf("Expected mode 'binary', got %s", tool.Mode())
	}
}