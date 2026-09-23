# Constraints — Task Tracker (Example)

> **EXAMPLE ONLY.** Populated illustration of `docs/system/constraints.md`. Not part of the core standard.

## Technical Constraints

- Python 3.12 required (example constraint).
- SQLite file must remain the only data store for this example; no additional stores without a new ADR.

## Business Constraints

- Demonstration project only; no production data.

## Operational Constraints

- Single-host deployment; no multi-environment promotion in this example.
- Database backups via host file copy before each deploy.

## Security Constraints

- No authentication in the prototype; MUST NOT be exposed publicly without adding authentication.
- Never commit the database file with real data.

## Compatibility Constraints

- Supports current desktop browsers for server-rendered pages (example scope).

## Performance Constraints

- None defined beyond acceptable interactive latency for a demo.

## Forbidden Changes

- Do not bypass the task service to write SQL from the web layer.
- Do not add an external integration without updating `architecture.md`, `dependencies.md`, and adding an ADR.
