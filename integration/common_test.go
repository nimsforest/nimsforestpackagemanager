//go:build integration

package integration

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// TestEnvironment provides a complete test environment for integration tests
type TestEnvironment struct {
	TempDir        string
	OriginalDir    string
	CLIPath        string
	WorkspaceDir   string
	OrganizationDir string
	ProductsDir    string
	T              *testing.T
}

// NewTestEnvironment creates a new test environment with workspace setup
func NewTestEnvironment(t *testing.T, orgName string) *TestEnvironment {
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	
	env := &TestEnvironment{
		TempDir:     tempDir,
		OriginalDir: originalDir,
		T:           t,
	}
	
	// Build CLI binary
	env.CLIPath = buildCLIBinary(t, tempDir)
	
	// Setup workspace
	env.setupWorkspace(orgName)
	
	return env
}

// Cleanup restores the original working directory
func (env *TestEnvironment) Cleanup() {
	os.Chdir(env.OriginalDir)
}

// setupWorkspace creates the workspace structure
func (env *TestEnvironment) setupWorkspace(orgName string) {
	env.WorkspaceDir = filepath.Join(env.TempDir, "test-workspace")
	env.OrganizationDir = filepath.Join(env.WorkspaceDir, orgName+"-organization-workspace")
	env.ProductsDir = filepath.Join(env.WorkspaceDir, "products-workspace")
	
	dirs := []string{env.OrganizationDir, env.ProductsDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			env.T.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}
	
	// Copy makefile for testing
	if err := env.copyMakefile(); err != nil {
		env.T.Logf("Warning: Could not copy makefile: %v", err)
	}
}

// copyMakefile copies the project makefile to the workspace
func (env *TestEnvironment) copyMakefile() error {
	_, currentFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(currentFile))
	makefileSrc := filepath.Join(projectRoot, "MAKEFILE.nimsforestpm")
	makefileDst := filepath.Join(env.WorkspaceDir, "MAKEFILE.nimsforestpm")
	
	return copyFile(makefileSrc, makefileDst)
}

// RunCLI runs the CLI binary with given arguments from the organization workspace
func (env *TestEnvironment) RunCLI(args ...string) (*CLIResult, error) {
	return env.RunCLIFromDir(env.OrganizationDir, args...)
}

// RunCLIFromDir runs the CLI binary from a specific directory
func (env *TestEnvironment) RunCLIFromDir(dir string, args ...string) (*CLIResult, error) {
	cmd := exec.Command(env.CLIPath, args...)
	cmd.Dir = dir
	
	output, err := cmd.CombinedOutput()
	
	return &CLIResult{
		Output:   string(output),
		Error:    err,
		ExitCode: cmd.ProcessState.ExitCode(),
	}, err
}

// RunCLIWithTimeout runs the CLI binary with a timeout
func (env *TestEnvironment) RunCLIWithTimeout(timeout time.Duration, args ...string) (*CLIResult, error) {
	cmd := exec.Command(env.CLIPath, args...)
	cmd.Dir = env.OrganizationDir
	
	done := make(chan error, 1)
	var output []byte
	
	go func() {
		var err error
		output, err = cmd.CombinedOutput()
		done <- err
	}()
	
	select {
	case err := <-done:
		return &CLIResult{
			Output:   string(output),
			Error:    err,
			ExitCode: cmd.ProcessState.ExitCode(),
		}, err
	case <-time.After(timeout):
		cmd.Process.Kill()
		return nil, fmt.Errorf("command timed out after %v", timeout)
	}
}

// InitGit initializes git in the products directory
func (env *TestEnvironment) InitGit() error {
	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = env.ProductsDir
	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git: %w", err)
	}
	
	// Set up git config for testing
	gitCmds := [][]string{
		{"git", "config", "user.name", "Test User"},
		{"git", "config", "user.email", "test@example.com"},
	}
	
	for _, cmdArgs := range gitCmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = env.ProductsDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git config failed: %w", err)
		}
	}
	
	return nil
}

// CreateMockComponent creates a mock component for testing
func (env *TestEnvironment) CreateMockComponent(name, fullName string, commands []string) error {
	componentDir := filepath.Join(env.ProductsDir, fullName+"-workspace", "main")
	if err := os.MkdirAll(componentDir, 0755); err != nil {
		return fmt.Errorf("failed to create component directory: %w", err)
	}
	
	// Create mock makefile
	makefilePath := filepath.Join(componentDir, "MAKEFILE."+fullName)
	var content strings.Builder
	content.WriteString("# Mock " + fullName + " makefile\n\n")
	
	for _, cmd := range commands {
		content.WriteString(fullName + "-" + cmd + ":\n")
		content.WriteString("\t@echo \"Mock " + name + " " + cmd + "\"\n\n")
	}
	
	return os.WriteFile(makefilePath, []byte(content.String()), 0644)
}

// CheckComponentInstalled verifies that a component is installed
func (env *TestEnvironment) CheckComponentInstalled(fullName string) bool {
	componentDir := filepath.Join(env.ProductsDir, fullName+"-workspace")
	_, err := os.Stat(componentDir)
	return err == nil
}

// GetInstalledComponents returns list of installed components
func (env *TestEnvironment) GetInstalledComponents() ([]string, error) {
	entries, err := os.ReadDir(env.ProductsDir)
	if err != nil {
		return nil, err
	}
	
	var components []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasSuffix(entry.Name(), "-workspace") {
			components = append(components, entry.Name())
		}
	}
	
	return components, nil
}

// CLIResult holds the result of running a CLI command
type CLIResult struct {
	Output   string
	Error    error
	ExitCode int
}

// Success returns true if the command succeeded
func (r *CLIResult) Success() bool {
	return r.Error == nil
}

// ContainsOutput checks if the output contains the given text
func (r *CLIResult) ContainsOutput(text string) bool {
	return strings.Contains(r.Output, text)
}

// ContainsOutputIgnoreCase checks if the output contains the given text (case insensitive)
func (r *CLIResult) ContainsOutputIgnoreCase(text string) bool {
	return strings.Contains(strings.ToLower(r.Output), strings.ToLower(text))
}

// HasError checks if the command failed
func (r *CLIResult) HasError() bool {
	return r.Error != nil
}

// ErrorContains checks if the error contains the given text
func (r *CLIResult) ErrorContains(text string) bool {
	return r.Error != nil && strings.Contains(r.Error.Error(), text)
}

// TestHelper provides utility functions for tests
type TestHelper struct {
	T *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{T: t}
}

// AssertSuccess asserts that the CLI result was successful
func (h *TestHelper) AssertSuccess(result *CLIResult, message string) {
	if !result.Success() {
		h.T.Fatalf("%s: command failed with error: %v\nOutput: %s", message, result.Error, result.Output)
	}
}

// AssertError asserts that the CLI result had an error
func (h *TestHelper) AssertError(result *CLIResult, message string) {
	if result.Success() {
		h.T.Fatalf("%s: expected error but command succeeded\nOutput: %s", message, result.Output)
	}
}

// AssertContains asserts that the output contains the expected text
func (h *TestHelper) AssertContains(result *CLIResult, expected, message string) {
	if !result.ContainsOutput(expected) {
		h.T.Fatalf("%s: expected output to contain '%s'\nActual output: %s", message, expected, result.Output)
	}
}

// AssertNotContains asserts that the output does not contain the text
func (h *TestHelper) AssertNotContains(result *CLIResult, unexpected, message string) {
	if result.ContainsOutput(unexpected) {
		h.T.Fatalf("%s: expected output not to contain '%s'\nActual output: %s", message, unexpected, result.Output)
	}
}

// AssertFileExists asserts that a file exists
func (h *TestHelper) AssertFileExists(path, message string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		h.T.Fatalf("%s: file %s does not exist", message, path)
	}
}

// AssertFileNotExists asserts that a file does not exist
func (h *TestHelper) AssertFileNotExists(path, message string) {
	if _, err := os.Stat(path); err == nil {
		h.T.Fatalf("%s: file %s should not exist", message, path)
	}
}

// AssertDirExists asserts that a directory exists
func (h *TestHelper) AssertDirExists(path, message string) {
	if stat, err := os.Stat(path); os.IsNotExist(err) || !stat.IsDir() {
		h.T.Fatalf("%s: directory %s does not exist", message, path)
	}
}

// ExpectSkip skips a test with a formatted message
func (h *TestHelper) ExpectSkip(result *CLIResult, condition string, args ...interface{}) {
	if !result.Success() {
		h.T.Skipf("Test skipped - %s: %v\nOutput: %s", fmt.Sprintf(condition, args...), result.Error, result.Output)
	}
}

// MockRepository provides utilities for creating mock git repositories
type MockRepository struct {
	Path string
	T    *testing.T
}

// NewMockRepository creates a new mock git repository
func NewMockRepository(t *testing.T, path string) *MockRepository {
	repo := &MockRepository{
		Path: path,
		T:    t,
	}
	
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Failed to create mock repository directory: %v", err)
	}
	
	repo.initGit()
	return repo
}

// initGit initializes git in the repository
func (r *MockRepository) initGit() {
	cmd := exec.Command("git", "init")
	cmd.Dir = r.Path
	if err := cmd.Run(); err != nil {
		r.T.Fatalf("Failed to initialize git repository: %v", err)
	}
	
	// Set up git config
	gitCmds := [][]string{
		{"git", "config", "user.name", "Test User"},
		{"git", "config", "user.email", "test@example.com"},
	}
	
	for _, cmdArgs := range gitCmds {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = r.Path
		if err := cmd.Run(); err != nil {
			r.T.Fatalf("Git config failed: %v", err)
		}
	}
}

// AddFile adds a file to the repository
func (r *MockRepository) AddFile(filename, content string) {
	filePath := filepath.Join(r.Path, filename)
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		r.T.Fatalf("Failed to write file %s: %v", filename, err)
	}
}

// Commit commits all changes
func (r *MockRepository) Commit(message string) {
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = r.Path
	if err := cmd.Run(); err != nil {
		r.T.Fatalf("Failed to add files: %v", err)
	}
	
	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = r.Path
	if err := cmd.Run(); err != nil {
		r.T.Fatalf("Failed to commit: %v", err)
	}
}

// buildCLIBinary builds the CLI binary for testing (shared function)
func buildCLIBinary(t *testing.T, tempDir string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	cmdDir := filepath.Dir(filepath.Dir(currentFile))
	
	binaryPath := filepath.Join(tempDir, "nimsforestpm-test")
	
	cmd := exec.Command("go", "build", "-o", binaryPath, "./cmd")
	cmd.Dir = cmdDir
	
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI binary: %v\nOutput: %s", err, string(output))
	}
	
	return binaryPath
}

// copyFile copies a file from src to dst (shared function)
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// WaitForFile waits for a file to be created (with timeout)
func WaitForFile(path string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("file %s not created within timeout", path)
}

// WaitForDir waits for a directory to be created (with timeout)
func WaitForDir(path string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if stat, err := os.Stat(path); err == nil && stat.IsDir() {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("directory %s not created within timeout", path)
}

// CleanupTempDirs removes all temporary directories (call in defer)
func CleanupTempDirs(dirs ...string) {
	for _, dir := range dirs {
		os.RemoveAll(dir)
	}
}

// GetAvailableTools returns the list of available tools for testing
func GetAvailableTools() []string {
	return []string{"work", "organize", "communicate", "productize", "folders", "webstack"}
}

// GetToolFullName returns the full name for a tool
func GetToolFullName(tool string) string {
	toolMap := map[string]string{
		"work":        "nimsforestwork",
		"organize":    "nimsforestorganize",
		"communicate": "nimsforestcommunication",
		"productize":  "nimsforestproductize",
		"folders":     "nimsforestfolders",
		"webstack":    "nimsforestwebstack",
	}
	return toolMap[tool]
}