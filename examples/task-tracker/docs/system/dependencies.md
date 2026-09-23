# Dependencies — Task Tracker (Example)

> **EXAMPLE ONLY.** Populated illustration of `docs/system/dependencies.md`. Not part of the core standard.

## Component Dependency Graph

```text
browser --> web layer --> task service --> task store --> data/tasks.db
```

## Runtime Dependencies

| Dependent | Depends On | Type | Notes |
|-----------|------------|------|-------|
| Web layer | Task service | internal | In-process calls |
| Task service | Task store | internal | In-process calls |
| Task store | SQLite file | platform | File-based relational store (example) |
| Test suite | `pytest` | third-party | Example test runner |

## External Dependencies

| Dependency | Provider | Purpose | Contract | Fallback |
|------------|----------|---------|----------|----------|
| None | — | — | — | — |

## Infrastructure Dependencies

| Dependent | Infrastructure | Notes |
|-----------|----------------|-------|
| Backend service | Single virtual host (example) | Demo environment only |
| Backend service | Reverse proxy (example) | Forwards HTTP to localhost |
| Pipeline | Hosted test runner (example) | Runs test suite on change |
