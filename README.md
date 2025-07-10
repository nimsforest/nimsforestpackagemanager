# NimsforestPM - Package Manager for Organizational Components

NimsForest enables the core organizational funnel: **inbound communication → issues → work → improvements → new content → outbound communication → inbound communication**. NimsforestPM helps you discover and install the right components to build this complete organizational intelligence cycle.

## Quick Setup

### 1. Add to your project
```bash
git submodule add https://github.com/nimsforest/nimsforestpm.git tools/nimsforestpm
```

### 2. Check compatibility and initialize
```bash
cd tools/nimsforestpm
make nimsforestpm-hello
make nimsforestpm-init
```

## Core Components

The four pillars of the organizational intelligence cycle:

1. **nimsforestcommunication** - Captures all inbound communication
2. **nimsforestwork** - Transforms communication into structured work items  
3. **nimsforestorganization** - Maintains organizational structure and identity
4. **nimsforestwebstack** - Enables outbound communication and web presence

## Commands

```bash
make nimsforestpm-hello        # System compatibility check
make nimsforestpm-init         # Initialize component discovery and setup
make nimsforestpm-lint         # Validate installed components
```

## Integration

Each component follows the same pattern:
- Self-contained with README.md and MAKEFILE
- Git submodule ready
- Integrates via `make {component}-init`
- Works together to enable the complete organizational funnel

Perfect for organizations wanting to systematically capture, process, and act on all communication while maintaining complete organizational intelligence.