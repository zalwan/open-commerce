# AI-Native Repository Standard v1.0

A technology-agnostic repository standard that provides structured engineering context for humans and AI agents.

## What Is This?

AI-Native Repository Standard is a technology-agnostic repository standard designed to provide structured engineering context for humans and AI agents.

It defines:

- a reusable repository structure,
- an AI instruction layer (`.ai/`, `AGENTS.md`),
- system knowledge templates (`docs/system/`),
- an architecture decision system (`docs/decisions/`),
- task workflows and checklists,
- a context model with explicit source-of-truth and drift handling,
- machine-readable project metadata (`.ai/project.yaml`),
- a reference example (`examples/task-tracker/`).

The canonical definition is [`SPECIFICATION.md`](./SPECIFICATION.md). Agent behavior is defined in [`AGENTS.md`](./AGENTS.md).

## Why?

Modern AI coding agents are capable of modifying large codebases, but their effectiveness depends heavily on the quality and structure of available context.

This standard organizes that context.

Without it, agents infer architecture from filenames, invent missing decisions, and drift from intent. With it, agents discover the minimum sufficient context progressively, respect boundaries and constraints, and keep documentation current with implementation.

## Core Model

```text
Instructions
    ↓
.ai/

Knowledge
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

Decision History
    ↓
docs/decisions/
```

| Layer | Location | Role |
|-------|----------|------|
| Instructions | `.ai/`, `AGENTS.md` | How to operate. Never system-specific docs. |
| Knowledge | `docs/system/` | Intended architecture, components, dependencies, constraints. |
| Implementation | `app/`, `config/` | What actually exists. |
| Infrastructure | `infra/` | Environment and platform definitions. |
| Tests | `tests/` | Verification. |
| Decision History | `docs/decisions/` | Why important choices were made. Immutable. |

When layers conflict, the agent reports the conflict instead of silently choosing (see `AGENTS.md` §3, §9–§10).

## Design Goals

- technology agnostic
- AI-friendly
- human-friendly
- progressive context discovery
- explicit architecture
- explicit constraints
- traceable architectural decisions
- portable across AI coding agents

## Non-Goals

v1.0 explicitly does **not** include:

- AI model
- AI API
- RAG
- vector database
- CLI
- cloud platform
- deployment framework

These are out of scope for the standard itself. They may appear in adopting projects or future extensions, but the core standard does not require them.

## Repository Layout

```text
.
├── AGENTS.md
├── README.md
├── SPECIFICATION.md
├── LICENSE
├── .gitignore
├── .ai/
├── app/
├── infra/
├── config/
├── tests/
├── docs/
├── examples/task-tracker/
└── .github/workflows/
```

See [`.ai/context.md`](./.ai/context.md) for the navigation map.

## Adoption

1. Clone or copy this repository as a starting point.
2. Fill `docs/system/project.md` and `docs/system/technology.md` for your stack.
3. Declare your architecture in `docs/system/architecture.md` and `components.md`.
4. Record dependencies and constraints.
5. Create your first ADR in `docs/decisions/` when the first architectural decision is made.
6. Fill `.ai/project.yaml` metadata.
7. Define stack-specific rules in `docs/development/conventions.md` and `docs/development/setup.md`.
8. Review the example in `examples/task-tracker/` for a populated illustration (example stack only — not part of the core standard).

## Example

`examples/task-tracker/` demonstrates how a concrete project fills the generic template: populated system docs, one ADR, and adapted agent files. Its specific technology choices are **example-only** and impose no requirement on adopters.

## Versioning

Standard versioning is `MAJOR.MINOR.PATCH` (see `SPECIFICATION.md` §18). This release is **1.0.0**.

## License

MIT — see [`LICENSE`](./LICENSE).
