# 0001: Use Layered Architecture

> **EXAMPLE ONLY.** Example ADR demonstrating the decision-record format. Not part of the core standard.

- Status: Accepted
- Date: 2026-01-15
- Deciders: Example maintainer

## Context

The task tracker needs clear separation between HTTP handling, task rules, and persistence so that small changes do not leak across concerns. The system is small enough for a single deployable unit.

## Decision

Adopt a layered monolith: web layer → task service → task store. The web layer never touches SQL directly; all persistence goes through the task service.

## Alternatives Considered

- Single-module script with mixed concerns — rejected: would tangle rendering, rules, and SQL.
- Separate services per concern — rejected: operational overhead unjustified for a prototype with one data store and no external integrations.

## Consequences

- Positive: rule and persistence changes stay localized; testing the service does not require HTTP.
- Negative: in-process calls cannot scale independently; acceptable for example load.
- Risk: future real-time or multi-client needs may require revisiting this decision.

## Related

- [Architecture](../system/architecture.md)
- [Components](../system/components.md)
