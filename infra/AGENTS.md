# Infrastructure Agent Guidance

> Scope: this file applies to `infra/` and, where noted, `config/`.
> It inherits from the root `AGENTS.md` and `.ai/` instructions. More specific instructions apply within their scope.
> This file is technology-agnostic. It must not assume a specific hosting, provisioning, or pipeline platform.

## Scope

- `infra/` holds infrastructure definitions and environment configuration templates for this project.
- `config/` holds application configuration. Application code lives in `app/`.
- Do not place application logic in `infra/`. Do not place infrastructure definitions in `app/`.

## Local Inheritance Model

```text
/
├── AGENTS.md            ← repository-wide rules (highest general authority)
├── app/
│   └── AGENTS.md        ← optional, only if app needs specific guidance
└── infra/
    └── AGENTS.md        ← this file: infrastructure-specific guidance
```

- Root `AGENTS.md` always applies.
- An `AGENTS.md` deeper in the tree adds scope-specific rules; it does not cancel parent rules.
- When parent and local instructions conflict, the more specific file governs the conflicted point, and the conflict must be reported in the change summary.
- Do not create local `AGENTS.md` files speculatively. The standard does not require multiple local files — create one only when a directory genuinely needs distinct guidance.

## Before Changing Infrastructure

1. Read this file, the relevant `infra/` files, `docs/system/dependencies.md`, and `docs/operations/`.
2. Assess: dependency impact, security impact, blast radius (environments, services, data), configuration impact.
3. Follow `.ai/checklists/infrastructure.md`.

## Safety Rules

- Never commit real secrets, credentials, or private access material.
- Prefer parameterized, environment-scoped configuration over hardcoded values.
- Changes with cross-environment or destructive potential (data loss, service replacement, network exposure) must be called out explicitly and validated with the project's applicable static check or dry-run before submission.
- Document deployment implications: ordering, migration, and rollback.

## Validation

Run the project's applicable infrastructure checks (static validation, policy check, dry-run, or plan) as defined by the project's stack. If no check is configured yet, state that explicitly — do not claim validation that did not occur.
