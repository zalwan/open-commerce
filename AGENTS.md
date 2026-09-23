# AGENTS.md — AI-Native Repository Standard v1.0

> This file is the repository-wide operating model for AI coding agents and human contributors.
> It inherits nothing above it. Local `AGENTS.md` files (e.g. `infra/AGENTS.md`) add scope-specific rules.
> Behavior principles: `.ai/principles.md`. Task process: `.ai/workflow.md`. Navigation: `.ai/context.md`. Metadata: `.ai/project.yaml`.

## 1. What This Repository Standard Is

This repository follows the **AI-Native Repository Standard v1.0**: a technology-agnostic structure that provides structured engineering context for humans and AI agents.

The repository is organized as an **AI-Native Engineering Context**:

```text
Instructions
    ↓
.ai/

System Knowledge
    ↓
docs/

Implementation
    ↓
app/

Infrastructure
    ↓
infra/

Tests
    ↓
tests/

Architecture History
    ↓
docs/decisions/
```

- `.ai/` — how to operate (instructions, workflow, checklists, metadata). Never system-specific technical documentation.
- `docs/` — what the system is and why (system knowledge, decisions, operations, development).
- `app/` — application implementation.
- `infra/` — infrastructure definitions. Scope-specific rules in `infra/AGENTS.md`.
- `config/` — application configuration.
- `tests/` — test suites.

## 2. Repository Structure

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
│   ├── operations/
│   └── development/
│       ├── setup.md
│       └── conventions.md
├── examples/
└── .github/
    └── workflows/
```

Do not add directories without purpose. Every document must have a clear purpose. Reference implementation: `examples/task-tracker/`.

## 3. Source-of-Truth Rules

| Source | Role |
|--------|------|
| Implementation (`app/`, `infra/`, `config/`) | **Implementation truth** — what actually exists. |
| `docs/system/` | **Architecture intent** — what the system is designed to be. |
| `docs/decisions/` | **Decision history** — why important choices were made. Immutable history. |
| `.ai/` + `AGENTS.md` | **AI instructions** — how an agent should operate. |

When these conflict, **identify the conflict instead of silently overwriting one source with another**. A conflict that materially affects the task must be reported before proceeding with a change that depends on it.

## 4. Context Discovery Process

Use **progressive context discovery**: start narrow, expand only as needed. Do not read the entire repository by default.

Inspect in this order:

1. Current user request (authoritative intent).
2. Root `AGENTS.md` (this file).
3. Applicable local `AGENTS.md` (e.g. `infra/AGENTS.md` for infrastructure tasks).
4. `.ai/context.md` (navigation map).
5. `.ai/project.yaml` (metadata, when stack or path resolution matters).
6. Relevant `docs/system/*` — only the files relevant to the task.
7. Relevant decision records in `docs/decisions/`.
8. Relevant source code, configuration, infrastructure, and tests.

Relevant checklists in `.ai/checklists/` define task-type expectations (feature, bugfix, refactor, infrastructure).

## 5. Engineering Workflow

The default workflow for every task:

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

- **Request:** classify per `.ai/workflow.md` §1 (Feature, Bugfix, Refactor, Infrastructure, Documentation, Configuration, Security, Performance).
- **Understand:** restate the requirement in own terms; confirm scope.
- **Context Discovery:** per §4 above.
- **Architecture:** consult `docs/system/architecture.md`, `components.md`, and relevant ADRs; confirm boundaries.
- **Impact Analysis:** explicitly determine affected areas (Application, API, Database, Dependencies, Configuration, Infrastructure, Security, Tests, CI/CD, Documentation, Operations). Not every area must change.
- **Plan:** smallest correct change consistent with existing patterns and `docs/system/constraints.md`.
- **Implement:** follow `.ai/principles.md` and `docs/development/conventions.md`.
- **Validate:** run applicable tests, static analysis, formatting, type checking, infrastructure validation, security checks. Report what ran and results.
- **Documentation:** update only the documents the change affects (see §6).
- **Final Review:** re-read the diff; confirm no unrelated changes, validation complete, uncertainty reported.

## 6. Documentation Requirements

- System changes affecting architecture, dependencies, constraints, operations, or behavior **must** update the relevant documentation.
- Update only affected documents. Do not rewrite unrelated docs.
- Technology changes update `docs/system/technology.md` and `.ai/project.yaml`.
- Architectural changes update `docs/system/architecture.md` / `components.md` and, when decision-worthy, add a new ADR. Never silently rewrite an existing ADR; mark superseded records and link to successors.
- Documentation must evolve with the system. Stale documentation is context drift (see §9).

## 7. Safety Rules

- Human intent remains authoritative. Explicit current user instructions outrank all other context.
- Respect `docs/system/constraints.md`, including **Forbidden Changes**.
- Never commit real secrets, credentials, or private access material.
- Prefer parameterized, environment-scoped configuration over hardcoded values.
- Infrastructure changes with destructive or cross-environment potential require explicit blast-radius assessment and the project's applicable dry-run or static validation (see `infra/AGENTS.md` and `.ai/checklists/infrastructure.md`).
- Do not introduce AI provider integrations, model dependencies, vector stores, code-generation services, or deployment platforms as part of a routine task. v1.0 scope is the standard and templates only.

## 8. Validation Requirements

A change is not complete until relevant validation has been performed:

- Tests (new or updated where behavior changed; regression tests for bugfixes).
- Static analysis, formatting, type checking — as configured by the project.
- Infrastructure validation — applicable static check or dry-run for `infra/`/`config/` changes.
- Security checks — where the project defines them or the change touches trust boundaries.

If a check cannot be run, report it explicitly. Do not claim unperformed validation.

## 9. Handling of Uncertainty and Context Drift

**Explicit uncertainty is mandatory.** Never invent missing architecture or technology decisions.

Context drift occurs when:

```text
Documentation
    ≠
Implementation
    ≠
Infrastructure
```

Examples: architecture describes a removed component; technology lists an obsolete store; infrastructure contains undocumented resources; code introduces an unmapped component; documented dependencies differ from actual ones.

When drift or a contradiction is detected that materially affects the task:

```text
Architecture discrepancy detected.

Documentation states:
<documented behavior>

Implementation indicates:
<actual behavior>

This should be resolved or explicitly acknowledged before
making a change that depends on this information.
```

Report, do not silently choose. Proceed only after the discrepancy is resolved or explicitly acknowledged by the user.

## 10. Context Priority

When sources disagree, this is the precedence for *acting* (not permission to ignore contradictions):

```text
1. Explicit current user instruction
2. Applicable local AGENTS.md
3. Root AGENTS.md
4. Current implementation
5. Current system documentation
6. Architecture decisions
7. Other repository documentation
8. Assumptions
```

Assumptions are last and must be stated explicitly. A conflict between implementation and documentation must still be reported when it materially affects the task.

## 11. Local `AGENTS.md` Behavior

- A local `AGENTS.md` applies to files inside its directory scope.
- More specific instructions govern their scope; parent rules still apply unless directly conflicted.
- Inheritance example:

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

- This standard ships one local file (`infra/AGENTS.md`) as the canonical example. Do not create additional local files speculatively — only when a directory genuinely needs distinct guidance.

## 12. Compatibility and Portability

This standard is portable across AI coding agents and human-only workflows. No vendor-specific agent configuration is required. The core standard remains technology-agnostic; technology-specific guidance belongs in future optional profiles (see `SPECIFICATION.md` §15), not in this file.
