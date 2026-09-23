# Architecture — Task Tracker (Example)

> **EXAMPLE ONLY.** Populated illustration of `docs/system/architecture.md`. Not part of the core standard.

## Architectural Style

- Style: `layered monolith` (example declaration).
- Key patterns: request handler → service → store.
- Decision references: [0001-use-layered-architecture](../decisions/0001-use-layered-architecture.md).

## System Overview

A single Python service renders task pages, enforces task rules, and persists tasks to a local relational file. One deployable unit, one data store, no external calls.

```text
browser --> web layer --> task service --> sqlite file
```

## Components

Summarized here; detail in [components.md](./components.md).

| Component | Responsibility | Location |
|-----------|----------------|----------|
| Web layer | HTTP routing and page rendering | `app/web/` (example path) |
| Task service | Task rules and state transitions | `app/tasks/` (example path) |
| Task store | Persistence to SQLite file | `app/tasks/store.py` (example path) |

## Communication Mechanisms

| From → To | Mechanism | Notes |
|-----------|-----------|-------|
| Browser → Web layer | Request-response over HTTP | Synchronous; HTML responses |
| Web layer → Task service | In-process function calls | Synchronous |
| Task service → Task store | In-process calls + SQL | Synchronous |

## Data Flows

1. Create task: form submit → web layer validates input → task service creates `open` task → store inserts row → list page re-rendered.
2. Complete task: action → service checks `open`/`doing` → store updates status → list page re-rendered.

## Boundaries

- Module boundaries: `web` never touches SQL directly; all persistence goes through the task service.
- Trust boundaries: input validation at the web layer; state-transition rules in the task service.
- External boundaries: none (no outbound integrations in this example).

## External Integrations

| Integration | Direction | Contract | Failure Handling |
|-------------|-----------|----------|------------------|
| None | — | — | — |

## Security Boundaries

- Unauthenticated local demo: acceptable for example prototype; production adoption would add authentication (see constraints).
- File permissions on `data/tasks.db` restrict access to the service user.

## Scalability Considerations

- Single-process service with file store: sufficient for example load; concurrent writes serialize at the database.
- Stateless request handling except for the database file.

## Failure Modes

| Failure | Impact | Mitigation |
|---------|--------|------------|
| Database file locked | Write fails with error page | Retry guidance; backup documented in operations |
| Process crash | Downtime until restart | Host process supervision restarts service |
