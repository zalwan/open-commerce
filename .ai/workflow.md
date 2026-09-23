# AI Engineering Workflow

> Standard: AI-Native Repository Standard v1.0
> This file defines the default task workflow. Root `AGENTS.md` defines the full agent operating model; this file defines the step-by-step process.

## 1. Request Classification

Classify each request as one or more of:

- Feature
- Bugfix
- Refactor
- Infrastructure
- Documentation
- Configuration
- Security
- Performance

Classification determines which checklist in `.ai/checklists/` applies and which impact areas require attention. A request may carry multiple labels (for example, a performance bugfix).

## 2. Context Discovery

Inspect context in this order. Stop expanding scope as soon as sufficient context is gathered (progressive disclosure).

1. Root `AGENTS.md`
2. Applicable local `AGENTS.md` (for example, `infra/AGENTS.md` when touching infrastructure)
3. `.ai/context.md` (navigation map)
4. Relevant `docs/system/*` (only the files relevant to the task)
5. Relevant decision records in `docs/decisions/`
6. Relevant source code in `app/`, `config/`, or `tests/`
7. Relevant infrastructure in `infra/`

Do not read the entire repository by default. Start narrow, then expand only when the task requires it.

## 3. Impact Analysis

For each task, determine which of the following areas are actually affected:

```text
Application
API
Database
Dependencies
Configuration
Infrastructure
Security
Tests
CI/CD
Documentation
Operations
```

Not every area must change. The agent must determine which areas are actually affected and state that determination explicitly before implementing.

For infrastructure tasks, also assess blast radius, dependency impact, and deployment implications (see `.ai/checklists/infrastructure.md`).

## 4. Implementation

- Prefer the smallest correct change.
- Follow `.ai/principles.md`: reuse existing patterns, preserve boundaries, minimize surface.
- Follow project conventions in `docs/development/conventions.md`.
- Respect constraints in `docs/system/constraints.md`, including forbidden changes.

## 5. Validation

Run all applicable checks for the change type:

- tests
- static analysis
- formatting
- type checking
- infrastructure validation
- security checks

If a check cannot be run (missing tooling, environment limitation), report it explicitly rather than claiming success.

## 6. Documentation

Update documentation only where the change affects documented system knowledge:

- Architecture or component change → `docs/system/architecture.md`, `docs/system/components.md`
- Technology change → `docs/system/technology.md`, `.ai/project.yaml`
- Dependency change → `docs/system/dependencies.md`
- Constraint change → `docs/system/constraints.md`
- New architectural decision → new file in `docs/decisions/`
- Operational change → `docs/operations/`

Do not update documentation speculatively. Documentation must evolve with the system, but every edit must trace to an actual change.

## 7. Final Review

Before finishing:

1. Re-read the final diff.
2. Confirm no unrelated changes were introduced.
3. Confirm validation results.
4. Confirm documentation impact was checked.
5. Confirm any uncertainty or discrepancy was reported, not hidden.
