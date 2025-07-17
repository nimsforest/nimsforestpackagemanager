package main

import (
	"testing"

	"github.com/nimsforest/nimsforestpackagemanager/internal/registry"
)

func TestRegistryBasicFunctions(t *testing.T) {
	// Test that registry functions work
	available := registry.AvailableTools()
	if available == nil {
		t.Error("AvailableTools should return a slice, not nil")
	}

	installed := registry.InstalledTools()
	if installed == nil {
		t.Error("InstalledTools should return a slice, not nil")
	}

	// Test tool existence check
	isInstalled := registry.IsToolInstalled("non-existent-tool")
	if isInstalled {
		t.Error("Non-existent tool should not be reported as installed")
	}
}

func TestRootCommand(t *testing.T) {
	// Test that root command is properly configured
	if rootCmd.Use != "nimsforestpm" {
		t.Errorf("Expected rootCmd.Use to be 'nimsforestpm', got '%s'", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("rootCmd.Short should not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("rootCmd.Long should not be empty")
	}
}
