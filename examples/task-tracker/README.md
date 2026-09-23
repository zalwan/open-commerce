# Example: Task Tracker

> **EXAMPLE ONLY.** This directory demonstrates how a concrete project fills the AI-Native Repository Standard v1.0 template.
> Nothing here is part of the core standard. The technology choices below are illustrative and impose no requirement on adopters.

A minimal task-tracking application used to show populated templates:

- Backend service + server-rendered frontend
- File-based relational database
- Single-host deployment
- Pipeline running tests on change

## What This Example Demonstrates

| Standard Artifact | Example File |
|-------------------|--------------|
| Populated technology inventory | [docs/system/technology.md](./docs/system/technology.md) |
| Populated architecture | [docs/system/architecture.md](./docs/system/architecture.md) |
| Populated component map | [docs/system/components.md](./docs/system/components.md) |
| Populated dependencies | [docs/system/dependencies.md](./docs/system/dependencies.md) |
| Populated constraints | [docs/system/constraints.md](./docs/system/constraints.md) |
| Project definition | [docs/system/project.md](./docs/system/project.md) |
| At least one ADR | [docs/decisions/0001-use-layered-architecture.md](./docs/decisions/0001-use-layered-architecture.md) |
| Adapted root agent file | [AGENTS.md](./AGENTS.md) |
| Local infrastructure guidance | [infra/AGENTS.md](./infra/AGENTS.md) |

## Example Stack (illustrative only)

- Language: Python 3.12 (example choice)
- Frontend: server-rendered pages from the backend (no separate framework — example choice)
- Database: SQLite file (example choice)
- Hosting: single virtual host with reverse proxy (example choice)
- CI: hosted pipeline running `pytest` (example choice)

If this stack does not match your project, ignore it. The core standard (`/AGENTS.md`, `/.ai/`, `/docs/system/*` templates, `/SPECIFICATION.md`) is technology-agnostic.
