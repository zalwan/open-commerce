# AGENTS.md — Task Tracker (Example)

> **EXAMPLE ONLY.** Shows how an adopting project adapts the root `AGENTS.md`.
> The authoritative standard is the repository root `AGENTS.md` and `SPECIFICATION.md`.

## Project Specifics

- Stack: Python 3.12 service + SQLite file (example choices — see `docs/system/technology.md`).
- Architecture: layered monolith (see `docs/system/architecture.md` and `docs/decisions/0001-use-layered-architecture.md`).
- Boundaries: `app/web/` never issues SQL; all persistence goes through `app/tasks/`.

## Operating Rules for This Example

1. Follow the root workflow (Request → Understand → Context Discovery → Architecture → Impact Analysis → Plan → Implement → Validate → Documentation → Final Review).
2. For task changes, read `docs/system/components.md` and the task service before modifying the web layer.
3. Respect forbidden changes in `docs/system/constraints.md` (no web-layer SQL, no new integrations without an ADR).
4. Validate with the example test command (`pytest`) and report results.

## Local Files

- `infra/AGENTS.md` applies to `infra/` in this example.
