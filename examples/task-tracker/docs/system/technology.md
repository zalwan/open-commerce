# Technology — Task Tracker (Example)

> **EXAMPLE ONLY.** Populated illustration of `docs/system/technology.md`. Technology choices below are example-only.

## Languages

| Language | Version | Purpose | Notes |
|----------|---------|---------|-------|
| Python | 3.12 | Backend + page rendering | Example choice |

## Frameworks and Libraries

| Name | Version | Purpose | Notes |
|------|---------|---------|-------|
| Standard-library HTTP server | 3.12 | Request handling | Example choice; no external web framework |
| `sqlite3` (standard library) | 3.12 | Database access | Example choice |

## Runtimes

| Runtime | Version | Purpose | Notes |
|---------|---------|---------|-------|
| Python interpreter | 3.12 | Runs backend service | Example choice |

## Data Stores

| Store | Type | Purpose | Notes |
|-------|------|---------|-------|
| `data/tasks.db` | Relational, file-based (SQLite) | Task persistence | Example choice; backed up with host volume |

## Infrastructure and Hosting

| Area | Choice | Notes |
|------|--------|-------|
| Hosting | Single virtual host (example) | One environment for demo |
| Networking | Reverse proxy → backend on localhost | Example topology |

## CI/CD

| Stage | Tool / Mechanism | Notes |
|-------|------------------|-------|
| Build | None (interpreted) | Example |
| Test | Hosted pipeline running `pytest` | Example |
| Deploy | Copy artifact + restart service | Example |

## Observability

| Area | Tool / Mechanism | Notes |
|------|------------------|-------|
| Logs | Standard output, collected by host | Example |
| Metrics | None | Example; acceptable for prototype |
| Traces | None | Example |

## External Services

| Service | Purpose | Integration | Notes |
|---------|---------|-------------|-------|
| None | — | — | Example has no external dependencies |
