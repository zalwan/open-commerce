# Project Context Index

> This file is a navigation map, not a duplicate documentation source.
> For behavior rules, see `../AGENTS.md` and `.ai/principles.md`.
> For process, see `.ai/workflow.md`.
> For machine-readable metadata, see `.ai/project.yaml`.

Start with the root `AGENTS.md`, then use this index to locate the minimum context required for the task.

## System Knowledge (`docs/system/`)

| Topic | Location |
|-------|----------|
| Purpose, users, scope, status | [docs/system/project.md](../docs/system/project.md) |
| Languages, frameworks, data stores, hosting, pipelines | [docs/system/technology.md](../docs/system/technology.md) |
| Style, overview, flows, boundaries | [docs/system/architecture.md](../docs/system/architecture.md) |
| Responsibilities and interfaces per component | [docs/system/components.md](../docs/system/components.md) |
| Runtime, external, and infrastructure dependencies | [docs/system/dependencies.md](../docs/system/dependencies.md) |
| Technical, business, security, and forbidden changes | [docs/system/constraints.md](../docs/system/constraints.md) |

## Decision History

- [docs/decisions/](../docs/decisions/) — Architecture Decision Records (why the system is the way it is). Read the records relevant to the task; do not assume the latest record describes the whole system.

## Operations

- [docs/operations/](../docs/operations/) — Deployment, monitoring, and troubleshooting guides (generic templates until populated for this project).

## Development

- [docs/development/setup.md](../docs/development/setup.md) — Environment setup.
- [docs/development/conventions.md](../docs/development/conventions.md) — Engineering conventions.

## Instructions and Metadata

- [../AGENTS.md](../AGENTS.md) — Agent operating model and source-of-truth rules.
- [.ai/principles.md](./principles.md) — Behavior principles.
- [.ai/workflow.md](./workflow.md) — Task workflow.
- [.ai/project.yaml](./project.yaml) — Machine-readable project metadata.
- [.ai/checklists/](./checklists/) — Task-type checklists (feature, bugfix, refactor, infrastructure).

## Suggested Entry Points

- New feature → `docs/system/project.md` → `architecture.md` → `components.md` → relevant ADR → code.
- Bugfix → `components.md` → relevant code and tests → `constraints.md` → relevant ADR.
- Infrastructure change → `infra/AGENTS.md` → `docs/system/dependencies.md` → `docs/operations/` → `infra/`.
- Unfamiliar project → this index → `project.md` → `technology.md` → `architecture.md`, then expand only as needed.
