# Operations

> Template directory. Replace with project-specific guides when adopting this standard.
> Purpose: how to deploy, observe, and recover this system.
> These templates are intentionally generic. Do not assume a specific hosting, container, or cloud platform.

## Expected Guides

| Guide | Suggested File | Contents |
|-------|---------------|----------|
| Deployment | `deployment.md` | Environments, release steps, configuration, rollback |
| Monitoring | `monitoring.md` | Health signals, logs, metrics, alerts, dashboards |
| Troubleshooting | `troubleshooting.md` | Symptom → diagnosis → remediation tables, escalation |

Create only the guides the project needs. A small project may start with a single `README.md` section; a larger one may split into the files above.

## Deployment (`deployment.md` sketch)

```markdown
# Deployment

## Environments

- [Environment 1: purpose, access, configuration source.]

## Release Process

1. [Step 1.]
2. [Step 2.]

## Configuration

- [Required settings and where they are provided.]

## Rollback

- [How to revert a release.]
```

## Monitoring (`monitoring.md` sketch)

```markdown
# Monitoring

## Health Signals

- [Signal 1: where to observe, what healthy looks like.]

## Alerts

| Alert | Meaning | Response |
|-------|---------|----------|
| [Name] | [Cause] | [Action + link to troubleshooting] |
```

## Troubleshooting (`troubleshooting.md` sketch)

```markdown
# Troubleshooting

| Symptom | Likely Cause | Diagnosis | Remediation |
|---------|--------------|-----------|-------------|
| [Observed behavior] | [Hypothesis] | [Where to confirm] | [Fix + escalation] |
```

## Rules

- Describe only what actually exists. Do not document aspirational infrastructure.
- If infrastructure contains resources not documented anywhere, that is context drift — document them or report the discrepancy.
- Operational changes that affect behavior must also update `docs/system/` where applicable.
