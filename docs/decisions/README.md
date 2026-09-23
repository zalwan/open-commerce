# Architecture Decision Records (ADRs)

> Purpose: durable, traceable history of important architectural decisions (decision history).
> System intent lives in `docs/system/`; this directory explains *why* important choices were made.

## Format

Files are numbered sequentially:

```text
0001-decision-name.md
0002-another-decision.md
```

Use short kebab-case names. Numbers never repeat and are never reused.

## Template

Each ADR must contain:

```text
Status
Context
Decision
Alternatives Considered
Consequences
Related
```

Supported statuses:

```text
Proposed
Accepted
Deprecated
Superseded
```

Copy the template below for each new ADR:

```markdown
# 000X: [Title]

- Status: Proposed
- Date: YYYY-MM-DD
- Deciders: [Names/roles]

## Context

[What problem or force requires a decision.]

## Decision

[What was decided, stated plainly.]

## Alternatives Considered

- [Alternative 1 — why rejected.]
- [Alternative 2 — why rejected.]

## Consequences

- [Effect 1 — positive, negative, or risk.]
- [Effect 2.]

## Related

- [Links to docs/system/*, other ADRs, or issues.]
```

## Lifecycle Rules

1. Historical decisions must not be silently rewritten. Edit an ADR only for typos or to update its `Status` line.
2. When an architectural decision changes, create a new ADR and mark the old one as `Superseded` with a link to the new record.
3. `Deprecated` means no longer in force but not replaced by a single successor.
4. A new project starts with zero ADRs. Create the first ADR when the first architectural decision is made — do not pre-fill with speculative decisions.

## What Qualifies as an ADR

- Architecture style or structural pattern choice.
- Technology choice with system-wide impact.
- Boundary, integration, or data-ownership decision.
- Constraint or policy with long-term consequences.

Routine implementation details do not need ADRs.
