# Infrastructure Agent Guidance — Task Tracker (Example)

> **EXAMPLE ONLY.** Shows how an adopting project adapts `infra/AGENTS.md`. Not part of the core standard.

## Scope (Example)

- `infra/` in this example holds the single-host service definition and reverse-proxy configuration.
- No secrets are stored here. The SQLite file lives under `data/` on the host, not in this repository.

## Example Rules

1. Before changing proxy or service definitions, read `docs/system/dependencies.md` and the example `docs/system/architecture.md`.
2. Blast radius for this example: one demo host. Call out any change that would affect data (`data/tasks.db`) explicitly.
3. Validate by reviewing the rendered configuration and, where the example host supports it, running the platform's check/dry-run. If no check exists, state that explicitly.

## Deployment Notes (Example)

- Deploy: copy service files, restart the service, verify `GET /tasks` responds.
- Rollback: restore the previous service copy and restart.
- Backup `data/tasks.db` before each deploy.
