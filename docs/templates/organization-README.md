# {{ORG_NAME}} Organization

Organizational workspace powered by NimsForest components.

## Why This Structure?

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

1. **Organization Repository** (`{{ORG_NAME}}-repository`): Core organizational structure
2. **Tools Repository** (`tools-repository`): Shared utilities and NimsForest components
3. **Product Repositories** (`product-repositories/`): Individual product development

This separation allows for:
- **Independent versioning**: Each component can evolve at its own pace
- **Flexible permissions**: Different access levels for different repositories
- **Git worktree support**: The `/main/` structure supports branching strategies

## Getting Started

### Install NimsForest Components

```bash
# Install core organizational intelligence cycle
make nimsforestpm-install COMPONENTS=work,communication,organization,webstack

# Add advanced folder management
make nimsforestpm-install COMPONENTS=folders

# Or install everything at once
make nimsforestpm-install COMPONENTS=all
```

### Add Your First Product

```bash
# Create a software product
make nimsforestpm-add-product PRODUCT_NAME=my-app PRODUCT_TYPE=software

# Or hardware product
make nimsforestpm-add-product PRODUCT_NAME=my-device PRODUCT_TYPE=hardware

# Or service product
make nimsforestpm-add-product PRODUCT_NAME=my-service PRODUCT_TYPE=service
```

### Validate Your Setup

```bash
# Check organizational structure and installed components
make nimsforestpm-lint
```

## Design Philosophy

This structure embodies several key principles:

1. **Game Engine Thinking**: Organizations are dynamic scenes with interacting entities
2. **Learning Systems**: Nims provide objective, transparent optimization
3. **Hierarchical Tools**: Each level (workspace, org, product) has its own tools
4. **Clean Separation**: Actors do things, assets are resources, tools enable work
5. **Git Worktree Ready**: Structure supports advanced branching workflows

## Next Steps for {{ORG_NAME}}

1. **Install components**: `make nimsforestpm-install COMPONENTS=all`
2. **Add your first product**: `make nimsforestpm-add-product PRODUCT_NAME=my-app PRODUCT_TYPE=software`
3. **Validate setup**: `make nimsforestpm-lint`
4. **Initialize components**: Run the respective `make {component}-init` commands for installed components

---

*Created with [NimsForest Package Manager](https://github.com/nimsforest/nimsforestpackagemanager)*