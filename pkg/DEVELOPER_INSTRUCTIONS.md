# nimsforest Package Manager - Developer Instructions

## Overview

This document provides comprehensive instructions for developers working with the nimsforest package manager's `pkg` directory, specifically the `tool` package that provides the core interfaces and types for all nimsforest tools.

## Package Structure

```
pkg/
├── tool/                    # Core tool interface package
│   ├── doc.go              # Package documentation with examples
│   ├── interface.go        # Core interfaces (Tool, Manager, Registry, etc.)
│   ├── types.go            # Common types, enums, and structures
│   ├── base.go             # Base implementation for tools to embed
│   ├── registry.go         # Tool registry system
│   └── errors.go           # Standard error types
└── DEVELOPER_INSTRUCTIONS.md # This file
```

## Core Concepts

### 1. Tool Interface

The `Tool` interface is the foundation that all nimsforest tools must implement:

```go
type Tool interface {
    Name() string
    Version() string
    Description() string
    Commands() []Command
    Execute(ctx context.Context, commandName string, args []string) error
    Install(ctx context.Context, options InstallOptions) error
    Update(ctx context.Context, options UpdateOptions) error
    Uninstall(ctx context.Context, options UninstallOptions) error
    Status() ToolStatus
    Info() ToolInfo
    Validate(ctx context.Context) error
}
```

### 2. Installation Modes

Tools support three installation modes:

- **Binary**: Pre-compiled binaries for direct execution
- **Clone**: Full repository clone for development/modification
- **Submodule**: Git submodule integration for workspace inclusion

### 3. Optional Interfaces

Tools can implement additional interfaces for extended functionality:

- `Configurable`: Configuration management
- `Healthcheck`: Health monitoring
- `Updatable`: Advanced update capabilities
- `DependencyProvider`: Dependency management
- `Workspace`: Workspace-specific operations
- `Plugin`: Plugin system support

## Creating a New Tool

### Step 1: Basic Implementation

```go
package main

import (
    "context"
    "fmt"
    "github.com/nimsforest/nimsforestpackagemanager/pkg/tool"
)

type MyTool struct {
    *tool.BaseTool
}

func NewMyTool() *MyTool {
    base := tool.NewBaseTool("mytool", "1.0.0", "A sample tool")
    t := &MyTool{BaseTool: base}
    
    // Add commands
    t.AddCommand(tool.Command{
        Name:        "hello",
        Description: "Say hello",
        Handler:     t.handleHello,
    })
    
    return t
}

func (t *MyTool) handleHello(ctx context.Context, args []string) error {
    fmt.Println("Hello from mytool!")
    return nil
}
```

### Step 2: Registration

```go
func main() {
    mytool := NewMyTool()
    
    if err := tool.Register(mytool); err != nil {
        panic(err)
    }
    
    // Tool is now discoverable
    fmt.Printf("Registered tool: %s\n", mytool.Name())
}
```

### Step 3: Advanced Features

```go
// Add configuration support
func (t *MyTool) Configure(ctx context.Context, config tool.Config) error {
    // Custom configuration logic
    return t.BaseTool.Configure(ctx, config)
}

// Add health checks
func (t *MyTool) HealthCheck(ctx context.Context) tool.HealthCheck {
    // Custom health check logic
    return t.BaseTool.HealthCheck(ctx)
}

// Add dependency management
func (t *MyTool) Dependencies() []tool.Dependency {
    return []tool.Dependency{
        {
            Name:     "git",
            Version:  ">=2.0.0",
            Required: true,
            Type:     tool.DependencyTypeSystem,
        },
    }
}
```

## Testing Your Tool

### Unit Tests

```go
package main

import (
    "context"
    "testing"
    "github.com/nimsforest/nimsforestpackagemanager/pkg/tool"
)

func TestMyTool(t *testing.T) {
    mytool := NewMyTool()
    
    // Test basic properties
    if mytool.Name() != "mytool" {
        t.Errorf("Expected name 'mytool', got %s", mytool.Name())
    }
    
    // Test command execution
    err := mytool.Execute(context.Background(), "hello", []string{})
    if err != nil {
        t.Errorf("Command execution failed: %v", err)
    }
    
    // Test validation
    err = mytool.Validate(context.Background())
    if err != nil {
        t.Errorf("Tool validation failed: %v", err)
    }
}
```

### Integration Tests

```go
func TestToolRegistry(t *testing.T) {
    // Clear registry for clean test
    tool.Clear()
    
    mytool := NewMyTool()
    
    // Test registration
    err := tool.Register(mytool)
    if err != nil {
        t.Errorf("Registration failed: %v", err)
    }
    
    // Test retrieval
    retrieved, err := tool.Get("mytool")
    if err != nil {
        t.Errorf("Tool retrieval failed: %v", err)
    }
    
    if retrieved.Name() != "mytool" {
        t.Errorf("Retrieved wrong tool")
    }
}
```

## Best Practices

### 1. Error Handling

Always use the provided error types:

```go
// Good
return tool.NewCommandNotFoundError(t.Name(), commandName)

// Bad
return fmt.Errorf("command not found: %s", commandName)
```

### 2. Context Usage

Always respect context cancellation:

```go
func (t *MyTool) longRunningOperation(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Continue operation
        }
    }
}
```

### 3. Configuration Validation

Always validate configuration:

```go
func (t *MyTool) ValidateConfig(config tool.Config) error {
    if _, exists := config["required_field"]; !exists {
        return tool.NewValidationFailedError("required_field", nil, "field is required")
    }
    return nil
}
```

### 4. Command Design

Design commands to be composable and follow Unix principles:

```go
tool.Command{
    Name:        "list",
    Description: "List items",
    Usage:       "mytool list [--format json|table]",
    Handler:     t.handleList,
}
```

## Common Patterns

### 1. Configuration-Driven Commands

```go
func (t *MyTool) handleConfigurableCommand(ctx context.Context, args []string) error {
    config := t.GetConfig()
    
    if enabled, ok := config.GetBool("feature_enabled"); ok && enabled {
        return t.executeFeature(ctx, args)
    }
    
    return fmt.Errorf("feature not enabled")
}
```

### 2. Dependency Checking

```go
func (t *MyTool) Execute(ctx context.Context, commandName string, args []string) error {
    if err := t.CheckDependencies(ctx); err != nil {
        return fmt.Errorf("dependency check failed: %w", err)
    }
    
    return t.BaseTool.Execute(ctx, commandName, args)
}
```

### 3. Progressive Installation

```go
func (t *MyTool) Install(ctx context.Context, options tool.InstallOptions) error {
    // Pre-install checks
    if err := t.CheckDependencies(ctx); err != nil {
        return err
    }
    
    // Install dependencies if needed
    if !options.SkipDependencies {
        if err := t.InstallDependencies(ctx); err != nil {
            return err
        }
    }
    
    // Perform actual installation
    return t.BaseTool.Install(ctx, options)
}
```

## Debugging

### 1. Enable Debug Logging

```go
func (t *MyTool) Execute(ctx context.Context, commandName string, args []string) error {
    fmt.Printf("DEBUG: Executing command %s with args %v\n", commandName, args)
    return t.BaseTool.Execute(ctx, commandName, args)
}
```

### 2. Health Check Information

```go
func (t *MyTool) HealthCheck(ctx context.Context) tool.HealthCheck {
    health := t.BaseTool.HealthCheck(ctx)
    
    // Add custom health information
    health.Details["custom_status"] = t.getCustomStatus()
    
    return health
}
```

### 3. Registry Inspection

```go
func inspectRegistry() {
    stats := tool.Stats()
    fmt.Printf("Registry stats: %+v\n", stats)
    
    for _, tool := range tool.List() {
        fmt.Printf("Tool: %s, Status: %s\n", tool.Name(), tool.Status())
    }
}
```

## Performance Considerations

### 1. Lazy Loading

```go
type MyTool struct {
    *tool.BaseTool
    heavyResource *HeavyResource
}

func (t *MyTool) getHeavyResource() *HeavyResource {
    if t.heavyResource == nil {
        t.heavyResource = NewHeavyResource()
    }
    return t.heavyResource
}
```

### 2. Command Caching

```go
func (t *MyTool) Commands() []tool.Command {
    if t.cachedCommands == nil {
        t.cachedCommands = t.buildCommands()
    }
    return t.cachedCommands
}
```

### 3. Concurrent Operations

```go
func (t *MyTool) processItems(ctx context.Context, items []string) error {
    sem := make(chan struct{}, 10) // Limit concurrency
    
    for _, item := range items {
        select {
        case sem <- struct{}{}:
            go func(item string) {
                defer func() { <-sem }()
                t.processItem(ctx, item)
            }(item)
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    
    return nil
}
```

## Contributing

### 1. Code Style

- Follow Go conventions
- Use meaningful variable names
- Add comments for exported functions
- Keep functions small and focused

### 2. Testing

- Write unit tests for all public functions
- Test error conditions
- Use table-driven tests where appropriate
- Mock external dependencies

### 3. Documentation

- Update doc.go with new features
- Add examples for complex functionality
- Document breaking changes

### 4. Compatibility

- Maintain backward compatibility
- Use semantic versioning
- Deprecate features before removal

## Support

For questions or issues:

1. Check the documentation in `doc.go`
2. Review existing tools for examples
3. Run tests to verify functionality
4. Create issues for bugs or feature requests

## Version History

- v1.0.0: Initial release with core interfaces
- Future versions will be documented here

---

This package is the foundation of the nimsforest ecosystem. Changes should be made carefully and with consideration for all dependent tools.