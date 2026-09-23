# Infrastructure Checklist

Use for: Infrastructure and configuration changes (`infra/`, `config/`, pipelines).

```text
[ ] Infrastructure context inspected
[ ] Dependency impact assessed
[ ] Security impact assessed
[ ] Blast radius assessed
[ ] Configuration impact assessed
[ ] Validation performed
[ ] Deployment implications documented
[ ] Final diff reviewed
```

Notes:

- "Infrastructure context" includes `infra/AGENTS.md`, relevant `infra/` files, `docs/system/dependencies.md`, and `docs/operations/`.
- Assess blast radius before applying: which environments, services, and data stores are affected.
- Infrastructure validation means the applicable static check, policy check, or dry-run for the project's platform — see `infra/AGENTS.md`.
- Document deployment implications (ordering, migration, rollback) in the change summary or `docs/operations/` when applicable.
