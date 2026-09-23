# 0004: Terraform/OpenTofu IaC for GCP (planned, not yet applied)

- Status: Accepted
- Date: 2026-09-23
- Deciders: Maintainer (user request: GCP later, scaffold now)

## Context

Production deploy target is GCP (future). ClickOps consoles are not reviewable, team-deployable, or reproducible. The repo already has local Compose; a cloud path needs versioned infrastructure with secrets handled outside git. No GCP project, billing, or credentials exist yet — this decision covers the scaffold only, not a live environment.

## Decision

- IaC in `infra/terraform/`, plain HCL compatible with both HashiCorp Terraform (BSL) and OpenTofu (MPL): no vendor-specific blocks, `required_version >= 1.9`, providers `hashicorp/google ~> 6` and `hashicorp/random ~> 3`.
- Scope: VPC + subnet + Serverless VPC connector, Cloud SQL Postgres 16 (private IP), Artifact Registry, two Cloud Run services (api, web). DB password generated (`random_password`), stored in Secret Manager, injected as full `DATABASE_URL` secret.
- All identity/region/project parameterized; `project_id` and image references have no defaults. Remote state via GCS backend configured per environment at init time (default local state for review only).
- Nothing is applied by this change: zero live resources, zero blast radius. First apply requires a GCP project + billing + API enablement (documented in `infra/terraform/README.md`).

## Alternatives Considered

- ClickOps console / gcloud scripts — rejected: not reviewable, drifts silently, no state tracking.
- Deployment Manager — rejected: Google-deprecated in favor of Terraform.
- Pulumi — rejected: extra language runtime for a 4-resource footprint; HCL is enough.
- Terraform Cloud / Spacelift — rejected for now: vendor coupling before we even have a project; re-evaluate when CI-apply is needed.

## Consequences

- Plus: prod path is reviewable (`plan` in PR), reproducible per env, secrets never in git.
- Minus/risiko: HCL must be re-validated (`init` + `validate` + `plan`) once a real project exists; provider major upgrades need a follow-up. Compose remains the dev path — no change for contributors.
- Perlu: `infra/terraform/README.md` (prereqs, backend, plan/apply/teardown), dependency/technology doc rows marked `(planned)`, CI `plan` job when credentials exist (future).

## Related

- ADR-0001 (modular monolith — one api + one web service maps 1:1 to Cloud Run), ADR-0003 (Postgres — maps to Cloud SQL).
- `infra/terraform/`, `docs/system/dependencies.md`, `docs/system/technology.md`, `docs/operations/README.md`.
