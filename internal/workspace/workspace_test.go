package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWorkspace(t *testing.T) {
	ws := NewWorkspace()
	
	if ws.Version != "1.0" {
		t.Errorf("expected version '1.0', got '%s'", ws.Version)
	}
	
	if ws.Organization != "" {
		t.Errorf("expected empty organization, got '%s'", ws.Organization)
	}
	
	if len(ws.Products) != 0 {
		t.Errorf("expected empty products slice, got %d items", len(ws.Products))
	}
}

func TestParseWorkspace_ValidFormat(t *testing.T) {
	content := `nimsforest 1.0

organization ./acme-organization-workspace

products (
    ./products-workspace/nimsforestwork-workspace
    ./products-workspace/nimsforestcommunication-workspace  
    ./products-workspace/nimsforestwebstack-workspace
    ./products-workspace/acme-app-workspace
)`

	ws, err := ParseWorkspace(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if ws.Version != "1.0" {
		t.Errorf("expected version '1.0', got '%s'", ws.Version)
	}
	
	if ws.Organization != "./acme-organization-workspace" {
		t.Errorf("expected organization './acme-organization-workspace', got '%s'", ws.Organization)
	}
	
	expectedProducts := []string{
		"./products-workspace/nimsforestwork-workspace",
		"./products-workspace/nimsforestcommunication-workspace",
		"./products-workspace/nimsforestwebstack-workspace",
		"./products-workspace/acme-app-workspace",
	}
	
	if len(ws.Products) != len(expectedProducts) {
		t.Errorf("expected %d products, got %d", len(expectedProducts), len(ws.Products))
	}
	
	for i, expected := range expectedProducts {
		if i >= len(ws.Products) {
			t.Errorf("missing product at index %d: %s", i, expected)
			continue
		}
		if ws.Products[i] != expected {
			t.Errorf("expected product at index %d: '%s', got '%s'", i, expected, ws.Products[i])
		}
	}
}

func TestParseWorkspace_WithComments(t *testing.T) {
	content := `# This is a comment
nimsforest 1.0

# Organization comment
organization ./acme-organization-workspace

# Products comment
products (
    # Product 1 comment
    ./products-workspace/nimsforestwork-workspace
    ./products-workspace/nimsforestcommunication-workspace  
    # Product 2 comment
    ./products-workspace/nimsforestwebstack-workspace
    ./products-workspace/acme-app-workspace
)`

	ws, err := ParseWorkspace(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if ws.Version != "1.0" {
		t.Errorf("expected version '1.0', got '%s'", ws.Version)
	}
	
	if len(ws.Products) != 4 {
		t.Errorf("expected 4 products, got %d", len(ws.Products))
	}
}

func TestParseWorkspace_MinimalFormat(t *testing.T) {
	content := `nimsforest 1.0`

	ws, err := ParseWorkspace(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if ws.Version != "1.0" {
		t.Errorf("expected version '1.0', got '%s'", ws.Version)
	}
	
	if ws.Organization != "" {
		t.Errorf("expected empty organization, got '%s'", ws.Organization)
	}
	
	if len(ws.Products) != 0 {
		t.Errorf("expected empty products, got %d", len(ws.Products))
	}
}

func TestParseWorkspace_OrganizationOnly(t *testing.T) {
	content := `nimsforest 1.0

organization ./acme-organization-workspace`

	ws, err := ParseWorkspace(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if ws.Version != "1.0" {
		t.Errorf("expected version '1.0', got '%s'", ws.Version)
	}
	
	if ws.Organization != "./acme-organization-workspace" {
		t.Errorf("expected organization './acme-organization-workspace', got '%s'", ws.Organization)
	}
	
	if len(ws.Products) != 0 {
		t.Errorf("expected empty products, got %d", len(ws.Products))
	}
}

func TestParseWorkspace_ProductsOnly(t *testing.T) {
	content := `nimsforest 1.0

products (
    ./products-workspace/nimsforestwork-workspace
    ./products-workspace/nimsforestcommunication-workspace
)`

	ws, err := ParseWorkspace(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if ws.Version != "1.0" {
		t.Errorf("expected version '1.0', got '%s'", ws.Version)
	}
	
	if ws.Organization != "" {
		t.Errorf("expected empty organization, got '%s'", ws.Organization)
	}
	
	if len(ws.Products) != 2 {
		t.Errorf("expected 2 products, got %d", len(ws.Products))
	}
}

func TestParseWorkspace_InvalidFormat(t *testing.T) {
	testCases := []struct {
		name    string
		content string
	}{
		{
			name:    "empty content",
			content: "",
		},
		{
			name:    "missing version",
			content: "organization ./test",
		},
		{
			name:    "invalid version line",
			content: "invalid 1.0",
		},
		{
			name:    "invalid version format",
			content: "nimsforest invalid",
		},
		{
			name:    "invalid organization line",
			content: "nimsforest 1.0\norganization",
		},
		{
			name:    "unclosed products section",
			content: "nimsforest 1.0\nproducts (\n./test",
		},
		{
			name:    "invalid products start",
			content: "nimsforest 1.0\nproducts invalid",
		},
		{
			name:    "content after products section",
			content: "nimsforest 1.0\nproducts (\n./test\n)\nextra content",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseWorkspace(tc.content)
			if err == nil {
				t.Errorf("expected error for case '%s', but got none", tc.name)
			}
		})
	}
}

func TestWorkspaceString(t *testing.T) {
	ws := &Workspace{
		Version:      "1.0",
		Organization: "./acme-organization-workspace",
		Products: []string{
			"./products-workspace/nimsforestwork-workspace",
			"./products-workspace/nimsforestcommunication-workspace",
		},
	}
	
	result := ws.String()
	
	expectedLines := []string{
		"nimsforest 1.0",
		"",
		"organization ./acme-organization-workspace",
		"",
		"products (",
		"    ./products-workspace/nimsforestwork-workspace",
		"    ./products-workspace/nimsforestcommunication-workspace",
		")",
	}
	
	expected := strings.Join(expectedLines, "\n")
	
	if result != expected {
		t.Errorf("expected:\n%s\n\ngot:\n%s", expected, result)
	}
}

func TestWorkspaceAddProduct(t *testing.T) {
	ws := NewWorkspace()
	
	ws.AddProduct("./test-product")
	
	if len(ws.Products) != 1 {
		t.Errorf("expected 1 product, got %d", len(ws.Products))
	}
	
	if ws.Products[0] != "./test-product" {
		t.Errorf("expected product './test-product', got '%s'", ws.Products[0])
	}
	
	// Test duplicate addition
	ws.AddProduct("./test-product")
	
	if len(ws.Products) != 1 {
		t.Errorf("expected 1 product after duplicate addition, got %d", len(ws.Products))
	}
}

func TestWorkspaceRemoveProduct(t *testing.T) {
	ws := NewWorkspace()
	ws.AddProduct("./test-product-1")
	ws.AddProduct("./test-product-2")
	
	if len(ws.Products) != 2 {
		t.Errorf("expected 2 products, got %d", len(ws.Products))
	}
	
	ws.RemoveProduct("./test-product-1")
	
	if len(ws.Products) != 1 {
		t.Errorf("expected 1 product after removal, got %d", len(ws.Products))
	}
	
	if ws.Products[0] != "./test-product-2" {
		t.Errorf("expected remaining product './test-product-2', got '%s'", ws.Products[0])
	}
	
	// Test removal of non-existent product
	ws.RemoveProduct("./non-existent")
	
	if len(ws.Products) != 1 {
		t.Errorf("expected 1 product after removing non-existent, got %d", len(ws.Products))
	}
}

func TestLoadWorkspace(t *testing.T) {
	// Create a temporary file
	tmpDir, err := os.MkdirTemp("", "workspace_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	workspaceFile := filepath.Join(tmpDir, "nimsforest.workspace")
	content := `nimsforest 1.0

organization ./acme-organization-workspace

products (
    ./products-workspace/nimsforestwork-workspace
    ./products-workspace/nimsforestcommunication-workspace
)`
	
	err = os.WriteFile(workspaceFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	
	ws, err := LoadWorkspace(workspaceFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if ws.Version != "1.0" {
		t.Errorf("expected version '1.0', got '%s'", ws.Version)
	}
	
	if ws.FilePath != workspaceFile {
		t.Errorf("expected file path '%s', got '%s'", workspaceFile, ws.FilePath)
	}
	
	if len(ws.Products) != 2 {
		t.Errorf("expected 2 products, got %d", len(ws.Products))
	}
}

func TestLoadWorkspace_FileNotFound(t *testing.T) {
	_, err := LoadWorkspace("/non/existent/path/nimsforest.workspace")
	if err == nil {
		t.Error("expected error for non-existent file, but got none")
	}
}

func TestLoadWorkspace_EmptyPath(t *testing.T) {
	_, err := LoadWorkspace("")
	if err == nil {
		t.Error("expected error for empty path, but got none")
	}
}

func TestWorkspaceSave(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "workspace_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	ws := &Workspace{
		Version:      "1.0",
		Organization: "./acme-organization-workspace",
		Products: []string{
			"./products-workspace/nimsforestwork-workspace",
			"./products-workspace/nimsforestcommunication-workspace",
		},
	}
	
	workspaceFile := filepath.Join(tmpDir, "nimsforest.workspace")
	err = ws.Save(workspaceFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Verify file was created
	if _, err := os.Stat(workspaceFile); os.IsNotExist(err) {
		t.Error("workspace file was not created")
	}
	
	// Verify content
	content, err := os.ReadFile(workspaceFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}
	
	expected := ws.String()
	if string(content) != expected {
		t.Errorf("expected:\n%s\n\ngot:\n%s", expected, string(content))
	}
	
	// Verify FilePath was set
	if ws.FilePath != workspaceFile {
		t.Errorf("expected file path '%s', got '%s'", workspaceFile, ws.FilePath)
	}
}

func TestWorkspaceSave_EmptyPath(t *testing.T) {
	ws := NewWorkspace()
	err := ws.Save("")
	if err == nil {
		t.Error("expected error for empty path, but got none")
	}
}

func TestFindWorkspaceFile(t *testing.T) {
	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "workspace_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create nested directories
	subDir := filepath.Join(tmpDir, "sub", "deep")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("failed to create sub directories: %v", err)
	}
	
	// Create workspace file in root
	workspaceFile := filepath.Join(tmpDir, WorkspaceFileName)
	err = os.WriteFile(workspaceFile, []byte("nimsforest 1.0"), 0644)
	if err != nil {
		t.Fatalf("failed to create workspace file: %v", err)
	}
	
	// Search from deep directory
	foundFile, err := FindWorkspaceFile(subDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if foundFile != workspaceFile {
		t.Errorf("expected found file '%s', got '%s'", workspaceFile, foundFile)
	}
}

func TestFindWorkspaceFile_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "workspace_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	_, err = FindWorkspaceFile(tmpDir)
	if err == nil {
		t.Error("expected error when workspace file not found, but got none")
	}
}

func TestLoadWorkspaceFromDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "workspace_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create workspace file
	workspaceFile := filepath.Join(tmpDir, WorkspaceFileName)
	content := `nimsforest 1.0

organization ./acme-organization-workspace`
	
	err = os.WriteFile(workspaceFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	
	ws, err := LoadWorkspaceFromDir(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	if ws.Version != "1.0" {
		t.Errorf("expected version '1.0', got '%s'", ws.Version)
	}
	
	if ws.Organization != "./acme-organization-workspace" {
		t.Errorf("expected organization './acme-organization-workspace', got '%s'", ws.Organization)
	}
}

func TestWorkspaceValidate(t *testing.T) {
	t.Run("valid workspace", func(t *testing.T) {
		ws := &Workspace{
			Version:      "1.0",
			Organization: "",
			Products:     []string{},
		}
		
		err := ws.Validate()
		if err != nil {
			t.Errorf("unexpected error for valid workspace: %v", err)
		}
	})
	
	t.Run("empty version", func(t *testing.T) {
		ws := &Workspace{
			Version:      "",
			Organization: "",
			Products:     []string{},
		}
		
		err := ws.Validate()
		if err == nil {
			t.Error("expected error for empty version, but got none")
		}
	})
}

func TestWorkspaceGetAbsolutePaths(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "workspace_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create test directories
	orgDir := filepath.Join(tmpDir, "org")
	prodDir := filepath.Join(tmpDir, "prod")
	err = os.MkdirAll(orgDir, 0755)
	if err != nil {
		t.Fatalf("failed to create org dir: %v", err)
	}
	err = os.MkdirAll(prodDir, 0755)
	if err != nil {
		t.Fatalf("failed to create prod dir: %v", err)
	}
	
	ws := &Workspace{
		Version:      "1.0",
		Organization: "./org",
		Products:     []string{"./prod"},
		FilePath:     filepath.Join(tmpDir, "nimsforest.workspace"),
	}
	
	absOrg, absProducts, err := ws.GetAbsolutePaths()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	expectedOrg := filepath.Join(tmpDir, "org")
	if absOrg != expectedOrg {
		t.Errorf("expected organization path '%s', got '%s'", expectedOrg, absOrg)
	}
	
	if len(absProducts) != 1 {
		t.Errorf("expected 1 product path, got %d", len(absProducts))
	}
	
	expectedProd := filepath.Join(tmpDir, "prod")
	if absProducts[0] != expectedProd {
		t.Errorf("expected product path '%s', got '%s'", expectedProd, absProducts[0])
	}
}

func TestWorkspaceGetAbsolutePaths_NoFilePath(t *testing.T) {
	ws := &Workspace{
		Version:      "1.0",
		Organization: "./org",
		Products:     []string{"./prod"},
	}
	
	_, _, err := ws.GetAbsolutePaths()
	if err == nil {
		t.Error("expected error when FilePath is not set, but got none")
	}
}

func TestValidateWorkspaceFormat(t *testing.T) {
	t.Run("valid format", func(t *testing.T) {
		content := `nimsforest 1.0

organization ./test`
		
		err := ValidateWorkspaceFormat(content)
		if err != nil {
			t.Errorf("unexpected error for valid format: %v", err)
		}
	})
	
	t.Run("empty content", func(t *testing.T) {
		err := ValidateWorkspaceFormat("")
		if err == nil {
			t.Error("expected error for empty content, but got none")
		}
	})
	
	t.Run("missing version line", func(t *testing.T) {
		content := `organization ./test`
		
		err := ValidateWorkspaceFormat(content)
		if err == nil {
			t.Error("expected error for missing version line, but got none")
		}
	})
}

func TestNormalizeWorkspaceContent(t *testing.T) {
	input := `  nimsforest 1.0  

  organization ./test  

  products (  
./prod1
    ./prod2  
  )  `
	
	expected := `nimsforest 1.0

organization ./test

products (
    ./prod1
    ./prod2
)`
	
	result := NormalizeWorkspaceContent(input)
	if result != expected {
		t.Errorf("expected:\n%s\n\ngot:\n%s", expected, result)
	}
}