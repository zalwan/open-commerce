# Bugfix Checklist

Use for: Bugfix requests.

```text
[ ] Bug reproduced or behavior confirmed
[ ] Root cause identified
[ ] Relevant architecture inspected
[ ] Fix implemented
[ ] Regression test added/updated
[ ] Validation completed
[ ] Documentation impact checked
```

Notes:

- Prefer the smallest fix that addresses the root cause, not the symptoms.
- Check `docs/system/constraints.md` for forbidden changes before fixing.
- If documentation and implementation disagree about expected behavior, report the discrepancy per `AGENTS.md` instead of silently choosing one.
