# NimsForest Package Manager

**Bootstrap Event & Value-Driven Organization Workspaces**

NimsForest PM creates workspaces where organizations can explicitly optimize their coordination (organize) and value creation (productize) in an endless improvement cycle. This package manager bootstraps the complete workspace structure and orchestrates the nimsforest-organize and nimsforest-productize tools.

## Quick Start

### 1. Install NimsForest PM
```bash
git clone https://github.com/nimsforest/nimsforest-pm
cd nimsforest-pm
make install
```

### 2. Bootstrap Organization Workspace
```bash
# Create complete workspace structure
nimsforest-pm bootstrap my-org-workspace

# Or install tools individually  
nimsforest-pm install organize    # For organization coordination
nimsforest-pm install productize  # For product development
```

### 3. What You Get
```
my-org-workspace/
├── org-repository/              # Organization coordination (nimsforest-organize)
│   ├── docs/purpose/           # Vision, mission, goals, strategy
│   ├── docs/people/            # Teams, roles, skills, structure
│   ├── docs/processes/         # Workflows, procedures, methods
│   ├── docs/resources/         # Budget, tools, assets, constraints
│   ├── docs/knowledge/         # Decisions made, lessons learned
│   └── products/               # Git submodule links to products
└── products-workspace/         # Product development ecosystem
    ├── shared/                 # Common libraries, components, APIs
    ├── product-a-workspace/    # nimsforest-productize value streams
    │   ├── main/              # Main development branch
    │   └── feature-branches/  # Git worktrees for parallel work
    └── product-b-workspace/
        ├── main/
        └── feature-branches/
```

## The Full NimsForest Toolset

### Core Philosophy
Organizations are optimization engines where humans coordinate (organize) to create products that deliver user value (productize). Better coordination enables better products, which teaches better coordination—an endless improvement cycle.

### Tool Responsibilities

**nimsforest-pm** (this package manager):
- Bootstraps complete workspace architecture
- Orchestrates organize and productize tool installation
- Provides seamless integration with no dependencies beyond Unix tools

**nimsforest-organize**:
- Creates event-driven organizational coordination structure
- MECE (Mutually Exclusive, Collectively Exhaustive) documentation system
- Purpose → People → Processes → Resources → Knowledge architecture

**nimsforest-productize**:
- Generates complete value stream repositories
- Infrastructure as code (NixOS-style declarative systems)  
- Built-in metrics, feedback loops, and communication systems
- Communicate → Awareness → Usage → Feedback → Improve cycle

## Installation Patterns

### Complete Setup (Recommended)
```bash
# Bootstrap everything at once
nimsforest-pm bootstrap my-org-workspace
cd my-org-workspace

# Initialize organization coordination
cd org-repository/main
nimsforest-organize init

# Create your first product
cd ../../products-workspace
nimsforest-productize create my-first-product
cd my-first-product-workspace/main
nimsforest-productize init
```

### Individual Tool Installation
```bash
# Just organization coordination
nimsforest-pm install organize
nimsforest-organize init

# Just product development  
nimsforest-pm install productize
nimsforest-productize create my-product
```

### Legacy Integration
```bash
# Add to existing project as git submodule
git submodule add https://github.com/nimsforest/nimsforest-pm tools/nimsforest-pm
cd tools/nimsforest-pm
make legacy-install ORG_NAME=my-company
```

## Event-Driven Integration

### Organization → Product Events
- Organizational changes trigger events that flow to product development
- Strategy updates automatically sync to product roadmaps
- Resource changes update product capacity planning
- Team structure changes update product ownership

### Product → Organization Events  
- Product feedback generates events that improve organizational coordination
- User metrics inform organizational decision-making
- Product learnings update organizational knowledge base
- Value creation data drives coordination optimization

### Continuous Optimization Loop
1. Organization coordinates better → Products create more value
2. Products create more value → Organization learns coordination patterns  
3. Organization learns coordination patterns → Organization coordinates better
4. Endless cycle of improvement

## Package Manager Commands

### Core Commands
```bash
nimsforest-pm bootstrap <workspace-name>    # Create complete workspace
nimsforest-pm install organize             # Install organization tools
nimsforest-pm install productize           # Install product development tools
nimsforest-pm status                        # Check installation status
nimsforest-pm update                        # Update all installed tools
```

### Organization Commands (via nimsforest-organize)
```bash
nimsforest-organize init                    # Initialize org coordination structure
nimsforest-organize validate               # Validate MECE structure
nimsforest-organize events                 # Show active event streams
nimsforest-organize metrics                # Organization coordination metrics
```

### Product Commands (via nimsforest-productize)
```bash
nimsforest-productize create <product-name>  # Create new product workspace
nimsforest-productize init                   # Initialize value stream structure  
nimsforest-productize deploy                 # Deploy infrastructure as code
nimsforest-productize metrics                # Product value metrics
nimsforest-productize feedback               # User feedback analysis
```

### Legacy Support (Makefile-based)
```bash
make nimsforestpm-hello                     # System compatibility check
make nimsforestpm-create-organisation       # Legacy workspace creation
make nimsforestpm-install COMPONENTS=       # Legacy component installation
```

## Why NimsForest?

### For Organizations
- **Explicit Coordination**: Make invisible organizational patterns visible and improvable
- **Value-Driven**: Every decision anchored to user value creation
- **Event-Driven**: Changes trigger measurable responses across the system
- **Continuous Learning**: Product feedback improves organizational coordination

### For Developers
- **Infrastructure as Code**: Declarative, reproducible systems (NixOS-style)
- **Git Worktree Ready**: `/main/` structure supports advanced branching workflows
- **No Dependencies**: Pure Unix tools (make, bash) - works everywhere
- **Modular**: Install only what you need, when you need it

### For Teams
- **MECE Structure**: Mutually Exclusive, Collectively Exhaustive organization
- **Clear Ownership**: Every product has its own complete value stream
- **Automatic Integration**: Events flow between organization and products
- **Measurable Impact**: Track value creation from coordination to customer

Perfect for organizations ready to treat coordination like infrastructure as code—explicit, measurable, and continuously optimized.