# Components — Task Tracker (Example)

> **EXAMPLE ONLY.** Populated illustration of `docs/system/components.md`. Not part of the core standard.

## `Web layer`

- **Responsibility:** HTTP routing, input validation, page rendering. Does not contain task rules or SQL.
- **Source location:** `app/web/` (example path).
- **Dependencies:** Task service.
- **Consumers:** Browser.
- **Interfaces:** `GET /tasks`, `POST /tasks`, `POST /tasks/{id}/complete` (example routes).
- **Data stores:** None directly (via task service only).
- **Operational notes:** Stateless; safe to restart at any time.

---

## `Task service`

- **Responsibility:** Task creation rules and status transitions (`open` → `doing` → `done`). Owns validation of state changes.
- **Source location:** `app/tasks/service.py` (example path).
- **Dependencies:** Task store.
- **Consumers:** Web layer.
- **Interfaces:** `create_task(title)`, `start_task(id)`, `complete_task(id)`, `list_tasks(status)` (example functions).
- **Data stores:** Via task store only.
- **Operational notes:** No background work; failures surface as request errors.

---

## `Task store`

- **Responsibility:** SQLite persistence for tasks. Owns schema and queries; no business rules.
- **Source location:** `app/tasks/store.py` (example path).
- **Dependencies:** SQLite file at `data/tasks.db`.
- **Consumers:** Task service.
- **Interfaces:** `insert(task)`, `update_status(id, status)`, `fetch(filter)` (example functions).
- **Data stores:** `data/tasks.db` (single table `tasks`).
- **Operational notes:** File must be backed up with the host; concurrent writes serialize.
