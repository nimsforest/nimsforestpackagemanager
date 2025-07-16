package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nimsforest/nimsforestpackagemanager/pkg/tool"
)

// ExampleTool demonstrates how to create a basic tool
type ExampleTool struct {
	*tool.BaseTool
}

// NewExampleTool creates a new example tool
func NewExampleTool() *ExampleTool {
	base := tool.NewBaseTool("example", "1.0.0", "An example tool for demonstration")
	t := &ExampleTool{BaseTool: base}

	// Add commands
	t.AddCommand(tool.Command{
		Name:        "hello",
		Description: "Say hello",
		Usage:       "example hello [name]",
		Handler:     t.handleHello,
	})

	t.AddCommand(tool.Command{
		Name:        "version",
		Description: "Show version",
		Usage:       "example version",
		Handler:     t.handleVersion,
	})

	t.AddCommand(tool.Command{
		Name:        "status",
		Description: "Show tool status",
		Usage:       "example status",
		Handler:     t.handleStatus,
	})

	// Add dependencies
	t.AddDependency(tool.Dependency{
		Name:     "git",
		Version:  ">=2.0.0",
		Required: true,
		Type:     tool.DependencyTypeSystem,
	})

	return t
}

func (t *ExampleTool) handleHello(ctx context.Context, args []string) error {
	name := "World"
	if len(args) > 0 {
		name = args[0]
	}
	fmt.Printf("Hello, %s!\n", name)
	return nil
}

func (t *ExampleTool) handleVersion(ctx context.Context, args []string) error {
	fmt.Printf("%s version %s\n", t.Name(), t.Version())
	return nil
}

func (t *ExampleTool) handleStatus(ctx context.Context, args []string) error {
	info := t.Info()
	fmt.Printf("Tool: %s\n", info.Name)
	fmt.Printf("Version: %s\n", info.Version)
	fmt.Printf("Description: %s\n", info.Description)
	fmt.Printf("Status: %s\n", info.Status)
	fmt.Printf("Install Path: %s\n", info.InstallPath)
	fmt.Printf("Install Mode: %s\n", info.InstallMode)
	fmt.Printf("Commands: %d\n", len(t.Commands()))
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: example <command> [args...]")
		fmt.Println("Available commands:")
		fmt.Println("  hello [name]  - Say hello")
		fmt.Println("  version       - Show version")
		fmt.Println("  status        - Show tool status")
		os.Exit(1)
	}

	// Create the example tool
	exampleTool := NewExampleTool()

	// Install the tool (mark as installed for demo)
	err := exampleTool.Install(context.Background(), tool.InstallOptions{
		Mode: tool.InstallModeBinary,
		Path: "/tmp/example-tool",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to install tool: %v\n", err)
		os.Exit(1)
	}

	// Execute the command
	command := os.Args[1]
	args := os.Args[2:]

	err = exampleTool.Execute(context.Background(), command, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Command failed: %v\n", err)
		os.Exit(1)
	}
}