# AI-Native Repository Standard — Specification v1.0.0

> Canonical specification for v1.0. Normative unless marked non-normative (guidance, examples).
> Audience: repository designers, contributors, and AI coding agents operating under `AGENTS.md`.
> Companion files: `AGENTS.md` (operating model), `.ai/` (instructions), `docs/` (knowledge templates), `examples/task-tracker/` (reference implementation).

## 1. Purpose

The AI-Native Repository Standard defines a minimal, technology-agnostic repository structure and documentation model that provides structured engineering context for humans and AI coding agents.

The standard ensures that:

1. An agent can discover relevant context without reading the entire repository.
2. Architecture, components, dependencies, and constraints are explicit.
3. Architectural decisions are traceable over time.
4. Documentation, implementation, and infrastructure discrepancies are identifiable.
5. Human intent remains authoritative over inferred context.

Non-normative principle: **the repository should describe the system, not force the system to fit the repository.**

## 2. Scope

In scope for v1.0:

- Repository structure and file conventions.
- AI instruction layer (behavior, workflow, checklists, metadata schema).
- System knowledge templates (project, technology, architecture, components, dependencies, constraints).
- Architecture Decision Record (ADR) system.
- Operations and development documentation templates.
- Context discovery, source-of-truth, drift, and priority models.
- Local agent-file inheritance rules.
- Validation principles and conformance criteria.
- One reference example demonstrating template adoption.

Out of scope for v1.0 (explicit non-goals):

- Command-line tooling for the standard.
- AI provider integrations or model dependencies.
- Retrieval, vector stores, or code-generation services.
- Automatic architecture detection or dependency-graph generation.
- Deployment systems or platform bindings.
- Vendor-specific agent configuration (unless explicitly optional).

Future extensions are addressed in §19.

## 3. Terminology

| Term | Definition |
|------|------------|
| AI-Native Engineering Context | The combined structure of instructions (`.ai/`, `AGENTS.md`), knowledge (`docs/`), implementation (`app/`, `config/`), infrastructure (`infra/`), tests (`tests/`), and history (`docs/decisions/`). |
| Architecture Intent | The designed system as described in `docs/system/`. |
| Implementation Truth | What actually exists in `app/`, `infra/`, and `config/`. |
| Decision History | The immutable record of architectural decisions in `docs/decisions/`. |
| AI Instructions | Behavioral and procedural rules in `.ai/` and `AGENTS.md`. |
| Context Drift | A material divergence between documentation, implementation, and infrastructure (§11). |
| Progressive Context Discovery | Ordered, minimal expansion of inspected context (§9). |
| Local Agent File | An `AGENTS.md` inside a subdirectory that adds scope-specific rules (§12). |
| Technology Profile | A future optional extension adding stack-specific conventions (§15). Not part of v1.0. |

RFC-style keywords (MUST, SHOULD, MAY) apply where used.

## 4. Repository Structure

A conforming repository MUST include the following structure (additional project files are permitted; required paths MUST NOT be removed or repurposed):

```text
.
├── AGENTS.md
├── README.md
├── SPECIFICATION.md
├── LICENSE
├── .gitignore
├── .ai/
│   ├── principles.md
│   ├── workflow.md
│   ├── context.md
│   ├── project.yaml
│   └── checklists/
│       ├── feature.md
│       ├── bugfix.md
│       ├── refactor.md
│       └── infrastructure.md
├── app/
├── infra/
│   └── AGENTS.md
├── config/
├── tests/
├── docs/
│   ├── system/
│   │   ├── project.md
│   │   ├── technology.md
│   │   ├── architecture.md
│   │   ├── components.md
│   │   ├── dependencies.md
│   │   └── constraints.md
│   ├── decisions/
│   │   └── README.md
│   ├── operations/
│   │   └── README.md
│   └── development/
│       ├── setup.md
│       └── conventions.md
├── examples/
└── .github/
    └── workflows/
```

Rules:

1. `.ai/` MUST NOT contain system-specific technical documentation. It contains instructions and metadata only.
2. `docs/system/` describes the adopting project's actual system, not the standard itself.
3. `infra/AGENTS.md` is the canonical local-file example. Additional local files MAY exist only where a directory genuinely needs distinct guidance.
4. Empty implementation directories MUST contain a `.gitkeep` placeholder until populated.
5. No directory SHALL be added without a documented purpose.

## 5. AI Instruction Model

The instruction layer consists of:

- Root `AGENTS.md`: repository-wide operating model (authoritative for agent behavior within this repository).
- `.ai/principles.md`: behavior principles (understand before modifying; prefer existing patterns; minimize change surface; preserve boundaries; explicit uncertainty; verify; keep knowledge current).
- `.ai/workflow.md`: default task workflow (classify → discover → analyze impact → implement → validate → document → review).
- `.ai/context.md`: navigation map linking to system knowledge, decisions, operations, development, and metadata. It MUST NOT duplicate documentation content.
- `.ai/project.yaml`: machine-readable metadata per the schema below.
- `.ai/checklists/`: task-type expectations for feature, bugfix, refactor, and infrastructure work.

### 5.1 Metadata Schema

`.ai/project.yaml` MUST be valid YAML and MUST preserve the following top-level keys (values are project-declared; empty values are permitted for unadopted templates):

```yaml
standard:
  name: "AI-Native Repository Standard"
  version: "1.0"

project:
  name: ""
  status: ""
  maturity: ""

technology:
  languages: []
  frameworks: []
  runtimes: []
  databases: []
  infrastructure: []
  cicd: []
  observability: []
  external_services: []

architecture:
  style: ""
  patterns: []

documentation:
  system: "docs/system"
  decisions: "docs/decisions"
  operations: "docs/operations"
  development: "docs/development"

ai:
  instructions: ".ai"
```

The schema MUST remain technology-agnostic: keys are generic; values are free-form. The schema MUST NOT require any specific language, framework, store, or platform.

`AGENTS.md` and `.ai/workflow.md` MUST NOT contradict each other. `AGENTS.md` defines the operating model; `.ai/workflow.md` defines the process steps within it.

## 6. System Knowledge Model

`docs/system/` holds architecture intent for the adopting project. Each file has a defined purpose:

| File | Purpose |
|------|---------|
| `project.md` | Name, purpose, users, scope, non-goals, status, ownership. |
| `technology.md` | Languages, frameworks, runtimes, data stores, hosting, pipelines, observability, external services. Versions documented when relevant. |
| `architecture.md` | Declared style, overview, components summary, communication, data flows, boundaries, integrations, security boundaries, scalability, failure modes. |
| `components.md` | Per-component responsibility, location, dependencies, consumers, interfaces, data stores, operational notes. |
| `dependencies.md` | Component graph, runtime, external, and infrastructure dependencies. Markdown representation is sufficient; graph tooling is NOT required. |
| `constraints.md` | Technical, business, operational, security, compatibility, and performance constraints, plus explicit forbidden changes. |

Rules:

1. Templates MUST NOT prescribe an architecture style. The project MUST explicitly declare its own.
2. Templates MUST NOT require specific technologies. Illustrative stacks belong in `examples/` only.
3. `constraints.md` is binding on agents. Forbidden changes override convenience.

## 7. Architecture Documentation

Architecture documentation MUST:

1. Declare the architectural style explicitly (no default is assumed).
2. Summarize components in `architecture.md` and detail them in `components.md`.
3. Describe communication mechanisms, data flows, and boundaries (module, trust, external).
4. Reference decision records for significant choices.
5. Document failure modes and scalability considerations at the level of detail the system requires — no more.

Documentation SHOULD link to code and contracts rather than duplicating them.

## 8. ADR Model

`docs/decisions/` holds the immutable decision history.

### 8.1 File Convention

- Files are named `NNNN-kebab-case-title.md` with monotonically increasing numbers (e.g. `0001-...`, `0002-...`). Numbers are never reused.
- `docs/decisions/README.md` defines the system and template.

### 8.2 Record Structure

Each ADR MUST support:

```text
Status
Context
Decision
Alternatives Considered
Consequences
Related
```

Supported statuses:

```text
Proposed
Accepted
Deprecated
Superseded
```

### 8.3 Lifecycle

1. Historical records MUST NOT be silently rewritten. Amendments are limited to typographical corrections and status-line updates.
2. When a decision changes, a new ADR MUST be created and the prior record marked `Superseded` with a forward link.
3. `Deprecated` indicates no longer in force without a single successor.
4. Routine implementation details SHOULD NOT become ADRs. Architecture-wide, boundary, integration, data-ownership, and long-term policy decisions SHOULD.

## 9. Context Discovery

Agents MUST use progressive context discovery: start narrow, expand only as required.

Normative inspection order:

1. Explicit current user instruction.
2. Applicable local `AGENTS.md`.
3. Root `AGENTS.md`.
4. `.ai/context.md`.
5. `.ai/project.yaml` (when stack or path resolution matters).
6. Relevant `docs/system/*` (only task-relevant files).
7. Relevant decision records.
8. Relevant source, configuration, infrastructure, and tests.

Reading the entire repository by default is non-conforming. `.ai/context.md` functions as the navigation map for minimal sufficient context.

## 10. Source-of-Truth Model

| Source | Role |
|--------|------|
| Implementation (`app/`, `infra/`, `config/`) | Implementation truth — what actually exists. |
| `docs/system/` | Architecture intent — what the system is designed to be. |
| `docs/decisions/` | Decision history — why important choices were made. |
| `.ai/` + `AGENTS.md` | AI instructions — how an agent should operate. |

When sources conflict, the agent MUST identify the conflict rather than silently overwriting one source with another. A conflict that materially affects the task MUST be reported before proceeding with a dependent change.

## 11. Context Drift

Context drift is defined as:

```text
Documentation
    ≠
Implementation
    ≠
Infrastructure
```

Non-normative examples: architecture describes a removed component; technology lists an obsolete store; infrastructure contains undocumented resources; code introduces an unmapped component; documented dependencies differ from actual dependencies.

On detecting drift that materially affects the task, the agent MUST report using the discrepancy format:

```text
Architecture discrepancy detected.

Documentation states:
<documented behavior>

Implementation indicates:
<actual behavior>

This should be resolved or explicitly acknowledged before
making a change that depends on this information.
```

The agent MUST NOT silently choose one source. Work that depends on the disputed information MUST wait for resolution or explicit user acknowledgment.

## 12. Local `AGENTS.md`

1. A local `AGENTS.md` applies to files within its directory scope.
2. More specific instructions govern their scope; parent rules remain in force except where directly conflicted.
3. On conflict, the more specific file governs the conflicted point, and the conflict MUST be reported in the change summary.
4. The standard ships one local file (`infra/AGENTS.md`) as the canonical example. Additional local files MUST NOT be created speculatively.

Non-normative inheritance illustration:

```text
/
├── AGENTS.md
├── app/
│   ├── AGENTS.md
│   └── backend/
│       └── AGENTS.md
└── infra/
    └── AGENTS.md
```

## 13. Task Workflows

### 13.1 Classification

Every request MUST be classified as one or more of: Feature, Bugfix, Refactor, Infrastructure, Documentation, Configuration, Security, Performance.

### 13.2 Default Workflow

```text
Request
  ↓
Understand
  ↓
Context Discovery
  ↓
Architecture
  ↓
Impact Analysis
  ↓
Plan
  ↓
Implement
  ↓
Validate
  ↓
Documentation
  ↓
Final Review
```

### 13.3 Impact Analysis

The agent MUST explicitly determine affected areas among: Application, API, Database, Dependencies, Configuration, Infrastructure, Security, Tests, CI/CD, Documentation, Operations. Not every area must change.

### 13.4 Implementation and Checklists

- Implementation MUST prefer the smallest correct change consistent with existing patterns and constraints.
- Task-type checklists in `.ai/checklists/` (feature, bugfix, refactor, infrastructure) define completion expectations, including regression tests for bugfixes and blast-radius assessment for infrastructure.

## 14. Documentation Lifecycle

1. Documentation MUST evolve with the system. Stale documentation is drift (§11).
2. Only affected documents SHALL be updated. Unrelated rewrites are non-conforming.
3. Technology changes MUST update `docs/system/technology.md` and `.ai/project.yaml`.
4. Architectural changes MUST update `architecture.md` / `components.md` and, when decision-worthy, add a new ADR.
5. Existing ADRs MUST NOT be silently rewritten (§8.3).

## 15. Technology Profiles (Future — Non-Normative, Not Part of v1.0)

Profiles are a reserved extension point. They MAY in future extend the core standard with stack-specific conventions, commands, testing guidance, and local instructions.

Illustrative future layout — **example only, not part of v1.0, imposes no requirement**:

```text
profiles/
├── generic/
├── python/
├── node/
├── go/
├── rust/
└── ...
```

The core standard MUST remain technology-agnostic regardless of profiles. Profiles MUST NOT alter the §4 structure or §5 schema keys; they MAY add scoped conventions and checklists.

## 16. Validation Principles

A change is not complete until relevant validation has been performed: tests, static analysis, formatting, type checking, infrastructure validation (static check or dry-run), and security checks where applicable.

Conformance checks for the standard itself:

- **Structure:** all §4 required paths exist.
- **Documentation:** templates are internally consistent; every document has a clear purpose.
- **Links:** internal documentation links resolve.
- **Example:** `examples/task-tracker/` follows the standard with populated templates and at least one ADR.
- **YAML:** `.ai/project.yaml` is valid YAML and preserves §5.1 keys.
- **Instructions:** `AGENTS.md` and `.ai/workflow.md` are mutually consistent.
- **Scope:** no AI tooling, provider integration, or application code beyond the reference example has been introduced.
- **Technology agnosticism:** the core standard MUST NOT require specific technologies. The following are **illustrative example names only** — listed here solely to define the neutrality test, not as dependencies, recommendations, or requirements: illustrative language/framework/store/platform/provider names (for example, any specific web framework, UI library, relational store, container or orchestration platform, cloud provider, or model provider that an adopter might otherwise assume). Any such name appearing in the core standard outside `examples/` and outside an explicitly marked example block is a defect.

## 17. Compatibility Expectations

The standard MUST be usable with human-only workflows and with diverse AI coding agents without vendor-specific configuration. Optional vendor or product integrations MAY exist only as explicitly optional additions that do not alter conformance.

The following are **example agent categories only** (illustrative, non-exhaustive, imposing no requirement): general-purpose coding agents including command-line and editor-integrated assistants. Any specific product names mentioned in adoption guides are examples of compatible tooling, not required dependencies.

## 18. Versioning

The standard uses semantic-style versioning:

```text
MAJOR.MINOR.PATCH
```

Example: `1.0.0`.

- MAJOR: breaking changes to structure, schema keys, or normative rules.
- MINOR: backward-compatible additions (new optional templates, checklist items, metadata fields).
- PATCH: documentation corrections and clarifications with no normative change.

This document specifies version **1.0.0**. The `standard.version` key in `.ai/project.yaml` tracks the major.minor lineage (`1.0`).

## 19. Future Extension Points

Reserved, non-breaking directions (none part of v1.0):

1. Technology profiles (§15).
2. Optional validation tooling (structure linters, link checkers, YAML validators) — provided as opt-in, never required for conformance.
3. Optional reference implementations for additional system shapes.
4. Interoperability notes for specific agent platforms — explicitly optional.

Any extension MUST preserve technology agnosticism of the core and MUST NOT make v1.0-conforming repositories non-conforming.

---

## Appendix A — Conformance Checklist (Normative)

```text
[ ] SPECIFICATION.md exists and covers §§1–19
[ ] README.md explains adoption
[ ] AGENTS.md is complete
[ ] .ai/principles.md exists
[ ] .ai/workflow.md exists
[ ] .ai/context.md exists
[ ] .ai/project.yaml exists and is valid YAML
[ ] All system documentation templates exist
[ ] ADR system exists
[ ] Operations templates exist
[ ] Development templates exist
[ ] Task checklists exist
[ ] Local AGENTS.md example exists
[ ] Context drift is documented
[ ] Source-of-truth model is documented
[ ] Technology profiles documented as future extensions
[ ] Reference example exists and follows the standard
[ ] Internal links are valid
[ ] Core standard is technology agnostic
[ ] No unnecessary tooling was introduced
```
