# Refactor Checklist

Use for: Refactor requests (behavior-preserving restructuring).

```text
[ ] Existing behavior understood
[ ] Invariants identified
[ ] Architecture impact assessed
[ ] Refactor scope defined
[ ] Refactor implemented
[ ] Existing behavior verified
[ ] Tests executed
[ ] No unrelated changes introduced
```

Notes:

- Define scope explicitly before starting; do not expand scope mid-refactor.
- Behavior must be verified by existing tests plus any new characterization tests, not by inspection alone.
- If the refactor crosses component boundaries in `docs/system/components.md`, treat it as an architecture-affecting change.
