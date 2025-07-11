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
cd ../my-company-workspace/my-company-repository/main
make nimsforestpm-install COMPONENTS=work,communication,organization
```

## Why This Organizational Structure?

This organizational structure is inspired by **game engine architecture** - treating your organization as a dynamic scene where different entities interact to create value.

We also drew inspiration from **Pixar's USD (Universal Scene Description)** for hierarchical organization and **git worktree** patterns for flexible development.

## Organizational Structure

### actors/
All entities that can perform actions in your organizational scene:

- **nims/**: Intelligent advisory actors that learn optimal patterns through reinforcement learning. These are the smart shadows that provide transparent, objective guidance to help work flow better.

- **humans/**: People in the organization - employees, stakeholders, customers. The decision makers and creative force.

- **machines/**: Physical systems that perform work:
  - **mobile/**: Drones, robots, vehicles - systems that can move around
  - **fixed/**: Servers, ASML machines, production equipment - stationary systems

### assets/
Resources and files that actors use:

- **documentation/**: Knowledge and processes
- **data/**: Information and datasets
- **media/**: Images, videos, presentations
- **templates/**: Reusable patterns and structures

### tools/
Capabilities and utilities that enable work:

- **shared/**: Tools shared across the organization (via tools-repository)
- **org-specific/**: Tools specific to this organization

### products/
What the organization builds and delivers:

- Each product has its own repository with the same actor/asset/tool structure
- Products are linked as git submodules for version control
- Can be software, hardware, or services

## Workspace Architecture

The workspace follows a **three-repository pattern**:

1. **Organization Repository**: Core organizational structure
2. **Tools Repository**: Shared utilities and NimsForest components
3. **Product Repositories**: Individual product development

This separation allows for:
- **Independent versioning**: Each component can evolve at its own pace
- **Flexible permissions**: Different access levels for different repositories
- **Git worktree support**: The `/main/` structure supports branching strategies

## Core Components

The organizational intelligence cycle components:

1. **nimsforestcommunication** - Captures all inbound communication
2. **nimsforestwork** - Transforms communication into structured work items  
3. **nimsforestorganization** - Maintains organizational structure and identity
4. **nimsforestwebstack** - Enables outbound communication and web presence
5. **nimsforestfolders** - Advanced folder management system

## Commands

```bash
make nimsforestpm-hello                           # System compatibility check
make nimsforestpm-create-organisation ORG_NAME=  # Create organizational workspace
make nimsforestpm-add-product PRODUCT_NAME=      # Add product repository
make nimsforestpm-install COMPONENTS=            # Install specific components
make nimsforestpm-install-component COMPONENT=   # Install individual component
make nimsforestpm-install-folders                # Install folder management
make nimsforestpm-addtomainmake                  # Add nimsforestpm to main Makefile
make nimsforestpm-lint                           # Validate installed components
```

## Component Installation

```bash
# Install specific components
make nimsforestpm-install COMPONENTS=work,communication
make nimsforestpm-install COMPONENTS=organization,webstack,folders

# Install individual component
make nimsforestpm-install-component COMPONENT=nimsforestfolders

# Install all components
make nimsforestpm-install COMPONENTS=all
```

## Product Management

```bash
# Create software product
make nimsforestpm-add-product PRODUCT_NAME=my-app PRODUCT_TYPE=software

# Create hardware product
make nimsforestpm-add-product PRODUCT_NAME=my-device PRODUCT_TYPE=hardware

# Create service product
make nimsforestpm-add-product PRODUCT_NAME=my-service PRODUCT_TYPE=service
```

## Design Philosophy

This structure embodies several key principles:

1. **Game Engine Thinking**: Organizations are dynamic scenes with interacting entities
2. **Learning Systems**: Nims provide objective, transparent optimization
3. **Hierarchical Tools**: Each level (workspace, org, product) has its own tools
4. **Clean Separation**: Actors do things, assets are resources, tools enable work
5. **Git Worktree Ready**: Structure supports advanced branching workflows

## Integration

Each component follows the same pattern:
- Self-contained with README.md and MAKEFILE
- Git submodule ready
- Integrates via `make {component}-init`
- Works together to enable the complete organizational funnel

Perfect for organizations wanting to systematically capture, process, and act on all communication while maintaining complete organizational intelligence.