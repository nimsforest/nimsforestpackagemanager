# NimsforestPM - Package Manager for Organizational Components

NimsForest enables the core organizational funnel: **inbound communication → issues → work → improvements → new content → outbound communication → inbound communication**. NimsforestPM helps you discover and install the right components to build this complete organizational intelligence cycle.

## Quick Setup

### 1. Add to your project
```bash
git submodule add https://github.com/nimsforest/nimsforestpackagemanager.git tools/nimsforestpm
```

### 2. Create your organizational workspace
```bash
cd nimsforestpm
make nimsforestpm-hello
make nimsforestpm-create-organisation ORG_NAME=my-company
```

### 3. Install components and start organizing
```bash
cd ../my-company
make nimsforestpm-install COMPONENTS=work,communication,organization
```

## Core Components

The four pillars of the organizational intelligence cycle:

1. **nimsforestcommunication** - Captures all inbound communication
2. **nimsforestwork** - Transforms communication into structured work items  
3. **nimsforestorganization** - Maintains organizational structure and identity
4. **nimsforestwebstack** - Enables outbound communication and web presence

## Commands

```bash
make nimsforestpm-hello                    # System compatibility check
make nimsforestpm-create-organisation      # Create organizational workspace
make nimsforestpm-install                  # Install specific components
make nimsforestpm-addtomainmake            # Add nimsforestpm to main Makefile
make nimsforestpm-test                     # Run comprehensive test suite
make nimsforestpm-lint                     # Validate installed components
```

## Component Installation

```bash
# Install specific components
make nimsforestpm-install COMPONENTS=work,communication
make nimsforestpm-install COMPONENTS=organization,webstack

# Install all components
make nimsforestpm-install COMPONENTS=all
```

## Integration

Each component follows the same pattern:
- Self-contained with README.md and MAKEFILE
- Git submodule ready
- Integrates via `make {component}-init`
- Works together to enable the complete organizational funnel

Perfect for organizations wanting to systematically capture, process, and act on all communication while maintaining complete organizational intelligence.