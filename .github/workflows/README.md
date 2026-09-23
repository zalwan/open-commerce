# CI Workflows

> Template directory. This repository ships no active pipelines by default.
> Purpose: explain where automation lives once the project defines it.

## Conventions

- Workflow definitions live in `.github/workflows/` (or the project's chosen pipeline location, documented in `docs/system/technology.md`).
- Keep workflows minimal: build, test, and validate on change; deploy only via an explicit, documented mechanism.
- Every workflow must have a human-readable name and a documented trigger.
- Pipeline changes follow `.ai/checklists/infrastructure.md` (blast radius, rollback, validation).

## Adding the First Workflow

1. Document the pipeline choice in `docs/system/technology.md` (CI/CD section).
2. Add the workflow file(s) here.
3. Describe expected checks in `docs/development/conventions.md`.
4. Delete or replace this README once real workflows exist — do not keep stale placeholder docs alongside live pipelines.
