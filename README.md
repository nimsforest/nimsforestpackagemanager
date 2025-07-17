# NimsForest Package Manager

**Simple Go-based tool manager for the NimsForest ecosystem**

A lightweight package manager that installs and manages NimsForest tools via `go get` and `go install`. No complex dependencies, no configuration files—just a simple wrapper around Go's native tooling.

## Installation

### Platform-Specific Installers
```bash
# Linux/macOS
curl -fsSL get.nimsforest.com/install.sh | sh

# Windows (PowerShell)
irm get.nimsforest.com/install.ps1 | iex
```

The package manager will guide you through installing Go if you don't have it.

## Quick Start

### 1. Create Organization Workspace
```bash
# Install workspace tool first
nimsforestpm install workspace

# Create organizational workspace structure
nimsforestworkspace create my-org
cd my-org-workspace
```

### 2. Install Tools
```bash
# Install individual tools
nimsforestpm install organize
nimsforestpm install work
nimsforestpm install communicate

# Or install all tools at once
nimsforestpm install all
```

### 3. Check Status
```bash
nimsforestpm status
```

## Available Tools

- **organize**: Organization coordination and structure management
- **work**: Work and task management
- **communicate**: Communication and collaboration tools
- **webstack**: Web development and deployment
- **productize**: Product development workflows
- **folders**: File and folder management utilities

## Commands

### Core Commands
```bash
nimsforestpm install <tool> [tool2] [tool3]       # Install tools
nimsforestpm install all                           # Install all tools
nimsforestpm update [tool]                         # Update tools (all if no tool specified)
nimsforestpm status                                # Show installation status
nimsforestpm hello                                 # System compatibility check
nimsforestpm hello --dev                           # Developer mode compatibility check
nimsforestpm validate <tool>                       # Validate tool installation
```

### Workspace Commands
```bash
nimsforestpm install workspace                     # Install workspace tool
nimsforestworkspace create <name>                  # Create workspace structure
```

### Installation Examples
```bash
# Install single tool
nimsforestpm install organize

# Install multiple tools
nimsforestpm install organize work communicate

# Install all available tools
nimsforestpm install all

# Install from full repository path
nimsforestpm install github.com/nimsforest/nimsforestorganize
```

## How It Works

1. **Tool Registry**: Tools are defined in `docs/tools.json` with repository mappings
2. **Go-based Installation**: Uses `go get` and `go install` to install tools to `$GOPATH/bin`
3. **No Configuration**: No workspace files or complex configuration needed
4. **Simple Management**: Tools are standard Go binaries in your PATH

## Workspace Structure

```
my-org-workspace/
├── my-org-organization-workspace/    # Organization coordination
│   └── main/                         # Main organization repo
│       └── README.md                 # Organization documentation
└── products-workspace/               # Product development area
```

## Tool Development

Tools are standard Go programs that can be installed via `go install`. To create a compatible tool:

1. Build as a Go binary
2. Implement standard commands (version, help, etc.)
3. Add to the tools registry for easy installation

See [pkg/tool/README.md](pkg/tool/README.md) for the tool interface specification.

### Simple Tool Example
```go
package main

import (
    "context"
    "fmt"
    "github.com/nimsforest/nimsforestpackagemanager/pkg/tool"
)

func main() {
    // Create a simple tool
    mytool := tool.NewSimpleTool("mytool", "1.0.0", "My awesome tool")
    
    // Add commands
    mytool.AddCommand("hello", func(ctx context.Context, args []string) error {
        fmt.Println("Hello from mytool!")
        return nil
    })
    
    // Handle standard main logic
    mytool.HandleMain()
}
```

## Development

### Build
```bash
# Build for development
task build

# Build release binaries for all platforms
task build-release

# Full release with checksums
task release
```

### Releases
Binary releases are automatically built and published via GitHub Actions when a new release is created. The workflow builds cross-platform binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64) 
- Windows (amd64)

To create a release:
1. Create a git tag: `git tag v1.0.0`
2. Push the tag: `git push origin v1.0.0`
3. Create a GitHub release from the tag
4. Binaries will be automatically built and attached to the release

### Test
```bash
# Run tests
task test

# Run integration tests
task test-integration

# Run all tests
task test-all
```

## Key Features

- **Zero Dependencies**: No complex setup or configuration files
- **Go Native**: Leverages Go's built-in package management
- **Cross-Platform**: Works on Linux, macOS, and Windows
- **Simple**: Just a wrapper around `go get` and `go install`
- **Extensible**: Easy to add new tools to the registry

## License

MIT License - see [LICENSE](LICENSE) for details.