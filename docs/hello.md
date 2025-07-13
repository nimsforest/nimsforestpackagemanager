# Hello, NimsForest

**Event & Value-Driven Organization System**

## The Foundation

Organizations are social constructs where humans organize to achieve their goals. Goals are achieved through work, which manifests as products (including services) that create value for users.

**Examples:**
- A **startup** organizes to build profitable software → creates a mobile app that saves users time
- A **hospital** organizes to improve patient outcomes → delivers medical services that heal people  
- A **school** organizes to educate students → provides learning experiences that develop skills
- A **nonprofit** organizes to solve social problems → offers programs that improve communities

In essence: humans organize to achieve goals, then continuously optimize to achieve those goals better. This optimization process manifests in products.

The better the organization coordinates (organize), the better products they create (productize), which teaches them how to coordinate even better—an endless cycle of improvement.

## The Core Insight

**Organizations are optimization engines. This optimization manifests in products that create user value.**

Organizations exist to continuously **organize** themselves better so they can **productize** more effectively, which teaches them how to organize even better, enabling superior productizing—an infinite loop of improvement.

## Two Active Verbs, One System

**Organize** → **Productize** → **Organize** → **Productize**

- **Organize**: How humans coordinate to achieve goals more effectively
- **Productize**: How that coordination manifests as products that create user value

## Quick Start

```bash
# Clone and install nimsforest package manager
git clone https://github.com/nimsforest/nimsforest-pm
cd nimsforest-pm
make install

# Bootstrap a new organization workspace
nimsforest-pm bootstrap my-org-workspace

# Or install tools individually
nimsforest-pm install organize  # For organization architecture
nimsforest-pm install productize # For product development
```

## What You Get

A workspace structure that makes organizational coordination explicit and value creation measurable:

```
my-org-workspace/
├── org-repository/              # The organization coordination system
│   ├── docs/                   # nimsforest-organize goes here
│   └── products/               # Git submodule links to product repositories
└── products-workspace/         # The product development ecosystem
    ├── shared/                 # Common libraries, components, APIs
    ├── product-a-workspace/    # nimsforest-productize goes here
    │   ├── main/              # Main development branch
    │   └── feature-branches/  # Git worktrees for parallel development
    └── product-b-workspace/    # nimsforest-productize goes here
        ├── main/
        └── feature-branches/
```

**Important**: Products include everything that creates user value - custom software, off-the-shelf tools with internal configurations, accounting systems with organizational agreements, or services with defined processes. If it delivers value to users (internal or external), it's a product.

## The Tools

### nimsforest-organize
Creates event-driven organizational architecture with MECE structure for human coordination.

### nimsforest-productize  
Generates complete value stream repositories with integrated infrastructure, metrics, and feedback loops.

### nimsforest-pm
Make and bash-based package manager that orchestrates both tools and provides seamless integration. No dependencies beyond standard Unix tools.

---

# Chapter 1: The Philosophy

## Why This Matters

Most organizations are invisible to themselves. Coordination happens in meetings, decisions live in people's heads, and value creation is disconnected from organizational learning. This creates inefficiency, misalignment, and missed opportunities.

## The Solution

Make the invisible visible by treating organizational coordination like infrastructure as code:
- **Event-driven**: Changes trigger measurable responses
- **Value-focused**: Every decision anchored to user value creation  
- **Continuously optimized**: Feedback loops drive improvement

## The Result

Organizations that operate like well-architected software systems—responsive, measurable, and continuously improving.

---

# Chapter 2: Organization Architecture

## The MECE Structure

```
organize/
├── purpose/     # Why: Vision, mission, goals, strategy
├── people/      # Who: Teams, roles, skills, structure  
├── processes/   # How: Workflows, procedures, methods
├── resources/   # What: Budget, tools, assets, constraints
└── knowledge/   # Learning: Decisions made, lessons learned
```

## Event-Driven Coordination

- Organizational changes trigger events that flow to product development
- Product feedback generates events that improve organizational coordination
- Real-time optimization through continuous event streams

## Value Alignment

Every folder serves the ultimate goal of creating user value more efficiently and effectively.

---

# Chapter 3: Product Development

## Value Stream Structure

Each product is a complete value delivery system:

```
product-name/
├── src/           # Application code
├── infra/         # Declarative infrastructure (NixOS style)
├── metrics/       # Value measurement and tracking
├── events/        # Event integration with organization
├── feedback/      # User feedback collection and analysis
└── communication/ # How product value is communicated
```

## The Product Loop

**Communicate** → **Awareness** → **Usage** → **Feedback** → **Improve**

Built into every product repository to ensure continuous value optimization.

## Infrastructure as Code

Following NixOS principles, all infrastructure is declarative, reproducible, and lives with the product code.

---

# Chapter 4: Integration & Usage

## Getting Started

1. **Bootstrap**: `nimsforest-pm bootstrap my-org-workspace` creates complete system
2. **Define Purpose**: Start with value propositions in `/org-repository/docs/purpose/`
3. **Create Products**: Use `nimsforestproductize-init` in `/products-workspace/` for each value stream
4. **Establish Events**: Define how organization and products communicate
5. **Measure Everything**: Track value creation at every level

## Daily Usage

- Organizational changes update `/organize/` documentation
- Product development happens in `/productize/` value streams  
- Events flow between organization and products automatically
- Metrics aggregate from products to organizational dashboards

## Evolution

The system grows with your organization:
- New value streams become new product repositories
- Organizational learnings update coordination patterns
- Event flows optimize based on real usage data

---

# Chapter 5: Advanced Patterns

## Multi-Team Organizations

- Each team can have its own `/organize/` documentation
- Shared `/productize/` workspace with clear value stream ownership
- Event flows coordinate between teams automatically

## Platform Teams

- Platform services live in `/productize/platform-services/`
- Internal tools follow same value stream structure as customer products
- Infrastructure teams productize their coordination as reusable services

## Scaling

- Repository federation for large organizations
- Event aggregation across multiple organization repositories
- Shared templates and patterns in `/productize/shared/`

## Measurement

- Value metrics roll up from products to organizational KPIs
- Event analytics show coordination effectiveness
- Feedback velocity indicates organizational learning speed

---

*Read time: 4 minutes*

*Organizations are optimization engines. This optimization manifests in products that create user value. NimsForest makes this process explicit, measurable, and continuously improving.*

---

## FAQ

**Q: Can you have organizations within organizations?**

A: Organizations are defined by their leadership, not ownership. If a group has its own leadership making decisions about goals and coordination, it's a separate organization. A holding company and its subsidiaries are multiple organizations that may coordinate through events, but each optimizes independently under its own leadership.

**Q: How does this scale to large enterprises or government?**

A: Any group with unified leadership and shared goals can use this structure. Whether it's 3 people in a startup, 300 people in a division, or 300,000 people in a government agency - if they have common leadership making coordination decisions, they're one organization that can benefit from explicit organize → productize structure.

**Q: What about matrix organizations or shared services?**

A: Shared services are products that create value for internal users. Matrix reporting creates event flows between organizations. The key is identifying who has decision-making authority (leadership) for specific goals and coordination patterns.
